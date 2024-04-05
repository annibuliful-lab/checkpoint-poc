package imeiconfiguration

import (
	"checkpoint/.gen/checkpoint/public/model"
	table "checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/gql/enum"
	tagUtils "checkpoint/modules/tag"
	"checkpoint/utils"
	"checkpoint/utils/graphql_utils"
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/graph-gophers/dataloader"
	"github.com/graph-gophers/graphql-go"
	"github.com/samber/lo"

	"github.com/google/uuid"
)

type ImeiConfigurationService struct{}

func (ImeiConfigurationService) Upsert(data UpsertImeiConfigurationData) (*ImeiConfiguration, error) {
	dbClient := db.GetPrimaryClient()
	ctx := context.Background()
	tx, err := dbClient.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
	})

	if err != nil {
		log.Println("init-upsert-imei-transaction-error", err.Error())
		return nil, utils.InternalServerError
	}

	var columnToUpsert pg.ColumnList
	columnToUpsert = append(columnToUpsert,
		table.ImeiConfiguration.Imei,
		table.ImeiConfiguration.ID,
		table.ImeiConfiguration.ProjectId,
		table.ImeiConfiguration.StationLocationId,
		table.ImeiConfiguration.CreatedBy,
	)

	modelToUpsert := model.ImeiConfiguration{
		ID:                uuid.New(),
		ProjectId:         data.ProjectId,
		Imei:              data.Imei,
		CreatedBy:         data.UpdatedBy,
		StationLocationId: data.StationLocationId,
	}

	if data.BlacklistPriority != nil {
		columnToUpsert = append(columnToUpsert, table.ImeiConfiguration.BlacklistPriority)
		modelToUpsert.BlacklistPriority = model.BlacklistPriority(*data.BlacklistPriority)
	}

	if data.PermittedLabel != nil {
		columnToUpsert = append(columnToUpsert, table.ImeiConfiguration.PermittedLabel)
		modelToUpsert.PermittedLabel = model.DevicePermittedLabel(*data.PermittedLabel)
	}

	insertImeiStmt := table.ImeiConfiguration.
		INSERT(columnToUpsert).
		MODEL(modelToUpsert).
		ON_CONFLICT(table.ImeiConfiguration.Imei, table.ImeiConfiguration.ProjectId).
		DO_UPDATE(pg.SET(table.ImeiConfiguration.Imei.SET(pg.String(data.Imei)))).
		RETURNING(table.ImeiConfiguration.AllColumns)

	imeiConfiguration := model.ImeiConfiguration{}

	err = insertImeiStmt.QueryContext(ctx, tx, &imeiConfiguration)

	if err != nil {
		tx.Rollback()
		log.Println("upsert-imei-configuration-error", err.Error())
		return nil, utils.InternalServerError
	}

	if data.Tags != nil && len(*data.Tags) != 0 {
		deleteAllImeiTagsStmt := table.ImeiConfigurationTag.
			DELETE().
			WHERE(table.ImeiConfigurationTag.ImeiConfigurationId.EQ(pg.UUID(imeiConfiguration.ID)))
		_, err := deleteAllImeiTagsStmt.ExecContext(ctx, tx)
		if err != nil {
			tx.Rollback()
			log.Println("delete-imei-configuration-tag-error", err.Error())
			return nil, utils.InternalServerError
		}

		for _, tag := range *data.Tags {
			upsertTagStmt := table.Tag.
				INSERT(table.Tag.ID, table.Tag.ProjectId, table.Tag.Title, table.Tag.CreatedBy, table.Tag.CreatedAt).
				MODEL(model.Tag{
					ID:        uuid.New(),
					ProjectId: data.ProjectId,
					Title:     tag,
					CreatedBy: data.UpdatedBy,
					CreatedAt: time.Now(),
				}).
				ON_CONFLICT(table.Tag.Title, table.Tag.ProjectId).
				DO_UPDATE(pg.SET(table.Tag.ID.SET(table.Tag.EXCLUDED.ID))).
				RETURNING(table.Tag.AllColumns)
			tagResult := model.Tag{}

			err := upsertTagStmt.QueryContext(ctx, tx, &tagResult)

			if err != nil {
				tx.Rollback()
				log.Println("upsert-imei-configuration-tag-error", err.Error())
				return nil, utils.InternalServerError
			}

			insertImsiTagStmt := table.ImeiConfigurationTag.
				INSERT(table.ImeiConfigurationTag.ID, table.ImeiConfigurationTag.ImeiConfigurationId, table.ImeiConfigurationTag.TagId, table.ImeiConfigurationTag.CreatedBy).
				MODEL(model.ImeiConfigurationTag{
					ID:                  uuid.New(),
					ImeiConfigurationId: imeiConfiguration.ID,
					TagId:               tagResult.ID,
					CreatedBy:           data.UpdatedBy,
				})

			_, err = insertImsiTagStmt.ExecContext(ctx, tx)

			if err != nil {
				tx.Rollback()
				log.Println("insert-imei-configuration-tag-error", err.Error())
				return nil, utils.InternalServerError
			}

		}
	}

	tx.Commit()

	var updatedBy graphql.NullID
	if imeiConfiguration.UpdatedBy != nil {
		updatedBy = graphql_utils.ConvertStringToNullID(imeiConfiguration.UpdatedBy)
	}

	var updatedAt graphql.NullTime
	if imeiConfiguration.UpdatedAt != nil {
		updatedAt = graphql.NullTime{Value: &graphql.Time{Time: *imeiConfiguration.UpdatedAt}}
	}

	return &ImeiConfiguration{
		ID:                graphql.ID(imeiConfiguration.ID.String()),
		ProjectId:         graphql.ID(imeiConfiguration.ProjectId.String()),
		StationLocationId: graphql.ID(imeiConfiguration.StationLocationId.String()),
		Imei:              imeiConfiguration.Imei,
		PermittedLabel:    enum.GetDevicePermittedLabel(imeiConfiguration.PermittedLabel.String()),
		BlacklistPriority: enum.GetBlacklistPriority(imeiConfiguration.BlacklistPriority.String()),
		CreatedBy:         imeiConfiguration.CreatedBy,
		CreatedAt:         graphql.Time{Time: imeiConfiguration.CreatedAt},
		UpdatedBy:         &updatedBy,
		UpdatedAt:         &updatedAt,
	}, nil
}

func (ImeiConfigurationService) Delete(data DeleteImeiConfigurationData) (int, error) {
	dbClient := db.GetPrimaryClient()
	now := time.Now()
	softDeleteStmt := table.ImeiConfiguration.
		UPDATE(table.ImeiConfiguration.DeletedBy, table.ImeiConfiguration.DeletedAt).
		MODEL(model.ImeiConfiguration{
			DeletedAt: &now,
			DeletedBy: &data.DeletedBy,
		}).
		WHERE(table.ImeiConfiguration.ID.EQ(pg.UUID(data.ID)).
			AND(table.ImeiConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId))))

	_, err := softDeleteStmt.Exec(dbClient)

	if err != nil && db.HasNoRow(err) {
		return 403, utils.ForbiddenOperation
	}

	if err != nil {
		log.Println("delete-imei-configuration", err.Error())
		return 500, utils.InternalServerError
	}

	return 200, nil
}

func (ImeiConfigurationService) FindMany(data GetImeiConfigurationsData) ([]ImeiConfiguration, int, error) {
	dbClient := db.GetPrimaryClient()
	conditions := pg.Bool(true).
		AND(table.ImeiConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId))).
		AND(table.ImeiConfiguration.DeletedAt.IS_NULL()).
		AND(table.ImeiConfiguration.StationLocationId.EQ(pg.UUID(data.StationLocationId)))

	var fromConditions pg.ReadableTable = table.ImeiConfiguration

	if data.PermittedLabel != nil {
		conditions = conditions.AND(table.ImeiConfiguration.PermittedLabel.EQ(pg.NewEnumValue(*data.PermittedLabel)))
	}

	if data.BlacklistPriority != nil {
		conditions = conditions.AND(table.ImeiConfiguration.BlacklistPriority.EQ(pg.NewEnumValue(*data.BlacklistPriority)))
	}

	if data.Search != nil {
		conditions = conditions.AND(table.ImeiConfiguration.Imei.LIKE(pg.String(*data.Search)))
	}

	if data.Tags != nil && len(*data.Tags) != 0 {
		var tagItems []pg.Expression

		for _, tag := range *data.Tags {
			tagItems = append(tagItems, pg.String(tag))
		}

		fromConditions = fromConditions.
			INNER_JOIN(table.ImeiConfigurationTag, table.ImeiConfigurationTag.ImeiConfigurationId.EQ(table.ImeiConfiguration.ID)).
			INNER_JOIN(table.Tag, table.ImeiConfigurationTag.TagId.EQ(table.Tag.ID))

		conditions = conditions.AND(table.Tag.Title.IN(tagItems...))
	}

	getImeisStmt := table.ImeiConfiguration.SELECT(table.ImeiConfiguration.AllColumns).
		FROM(fromConditions).
		WHERE(conditions).
		LIMIT(data.Pagination.Limit).
		OFFSET(data.Pagination.Skip)

	imeiConfigurations := []model.ImeiConfiguration{}

	err := getImeisStmt.Query(dbClient, &imeiConfigurations)

	if err != nil {
		log.Println("get-imei-configurations-error", err.Error())
		return nil, 500, utils.InternalServerError
	}

	imeiConfigurationsResponse := lo.Map(imeiConfigurations, func(item model.ImeiConfiguration, index int) ImeiConfiguration {
		var updatedBy graphql.NullID
		if item.UpdatedBy != nil {
			updatedBy = graphql_utils.ConvertStringToNullID(item.UpdatedBy)
		}

		var updatedAt graphql.NullTime
		if item.UpdatedAt != nil {
			updatedAt = graphql.NullTime{Value: &graphql.Time{Time: *item.UpdatedAt}}
		}

		return ImeiConfiguration{
			ID:                graphql.ID(item.ID.String()),
			ProjectId:         graphql.ID(item.ProjectId.String()),
			StationLocationId: graphql.ID(item.StationLocationId.String()),
			Imei:              item.Imei,
			BlacklistPriority: enum.GetBlacklistPriority(item.BlacklistPriority.String()),
			PermittedLabel:    enum.GetDevicePermittedLabel(item.PermittedLabel.String()),
			CreatedBy:         item.CreatedBy,
			CreatedAt:         graphql.Time{Time: item.CreatedAt},
			UpdatedBy:         &updatedBy,
			UpdatedAt:         &updatedAt,
		}
	})

	return imeiConfigurationsResponse, 200, nil
}

func (ImeiConfigurationService) FindById(data GetImeiConfigurationData) (*ImeiConfiguration, int, error) {
	dbClient := db.GetPrimaryClient()

	getImeiStmt := table.ImeiConfiguration.
		SELECT(table.ImeiConfiguration.AllColumns).
		WHERE(table.ImeiConfiguration.ID.EQ(pg.UUID(data.ID)).
			AND(table.ImeiConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId))).
			AND(table.ImeiConfiguration.DeletedAt.IS_NULL()))

	imeiConfiguration := model.ImeiConfiguration{}

	err := getImeiStmt.Query(dbClient, &imeiConfiguration)

	if err != nil && db.HasNoRow(err) {
		return nil, 404, errors.New("imei id not found")
	}

	if err != nil {
		log.Println("insert-imei-configuration", err.Error())
		return nil, 500, utils.InternalServerError
	}

	var updatedBy graphql.NullID
	if imeiConfiguration.UpdatedBy != nil {
		updatedBy = graphql_utils.ConvertStringToNullID(imeiConfiguration.UpdatedBy)
	}

	var updatedAt graphql.NullTime
	if imeiConfiguration.UpdatedAt != nil {
		updatedAt = graphql.NullTime{Value: &graphql.Time{Time: *imeiConfiguration.UpdatedAt}}
	}

	return &ImeiConfiguration{
		ID:                graphql.ID(imeiConfiguration.ID.String()),
		ProjectId:         graphql.ID(imeiConfiguration.ProjectId.String()),
		StationLocationId: graphql.ID(imeiConfiguration.StationLocationId.String()),
		Imei:              imeiConfiguration.Imei,
		PermittedLabel:    enum.GetDevicePermittedLabel(imeiConfiguration.PermittedLabel.String()),
		BlacklistPriority: enum.GetBlacklistPriority(imeiConfiguration.BlacklistPriority.String()),
		CreatedBy:         imeiConfiguration.CreatedBy,
		CreatedAt:         graphql.Time{Time: imeiConfiguration.CreatedAt},
		UpdatedBy:         &updatedBy,
		UpdatedAt:         &updatedAt,
	}, 200, nil
}

func (ImeiConfigurationService) Update(data UpdateImeiConfigurationData) (*ImeiConfiguration, int, error) {
	dbClient := db.GetPrimaryClient()
	ctx := context.Background()
	tx, err := dbClient.Begin()
	if err != nil {
		log.Println("init-update-imei-transaction-error", err.Error())
		return nil, 500, utils.InternalServerError
	}

	now := time.Now()
	updateImeiStmt := table.ImeiConfiguration.
		UPDATE(table.ImeiConfiguration.Imei, table.ImeiConfiguration.UpdatedBy, table.ImeiConfiguration.PermittedLabel, table.ImeiConfiguration.BlacklistPriority, table.ImeiConfiguration.UpdatedAt).
		MODEL(model.ImeiConfiguration{
			Imei:              data.Imei,
			UpdatedBy:         &data.UpdatedBy,
			UpdatedAt:         &now,
			PermittedLabel:    model.DevicePermittedLabel(data.PermittedLabel),
			BlacklistPriority: model.BlacklistPriority(data.BlacklistPriority),
		}).
		RETURNING(table.ImeiConfiguration.AllColumns).
		WHERE(table.ImeiConfiguration.ID.EQ(pg.UUID(data.ID)).
			AND(table.ImeiConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId))).
			AND(table.ImeiConfiguration.DeletedAt.IS_NULL()))

	imeiConfiguration := model.ImeiConfiguration{}

	err = updateImeiStmt.QueryContext(ctx, tx, &imeiConfiguration)

	if err != nil && db.HasNoRow(err) {
		return nil, 403, utils.ForbiddenOperation
	}

	if err != nil {
		log.Println("insert-imei-configuration", err.Error())
		return nil, 500, utils.InternalServerError
	}

	if data.Tags != nil && len(*data.Tags) != 0 {
		deleteAllImeiTagsStmt := table.ImeiConfigurationTag.
			DELETE().
			WHERE(table.ImeiConfigurationTag.ImeiConfigurationId.EQ(pg.UUID(imeiConfiguration.ID)))
		_, err := deleteAllImeiTagsStmt.ExecContext(ctx, tx)
		if err != nil {
			tx.Rollback()
			log.Println("delete-imei-configuration-tag-error", err.Error())
			return nil, 500, utils.InternalServerError
		}

		for _, tag := range *data.Tags {
			upsertTagStmt := tagUtils.UpsertStatement(tagUtils.UpsertTagData{
				Tag:       tag,
				ProjectId: data.ProjectId.String(),
				CreatedBy: data.UpdatedBy,
			})
			tagResult := model.Tag{}

			err := upsertTagStmt.QueryContext(ctx, tx, &tagResult)

			if err != nil {
				tx.Rollback()
				log.Println("upsert-imei-configuration-tag-error", err.Error())
				return nil, 500, utils.InternalServerError
			}

			insertImsiTagStmt := table.ImeiConfigurationTag.
				INSERT(table.ImeiConfigurationTag.ID, table.ImeiConfigurationTag.ImeiConfigurationId, table.ImeiConfigurationTag.TagId, table.ImeiConfigurationTag.CreatedBy).
				MODEL(model.ImeiConfigurationTag{
					ID:                  uuid.New(),
					ImeiConfigurationId: imeiConfiguration.ID,
					TagId:               tagResult.ID,
					CreatedBy:           data.UpdatedBy,
				})

			_, err = insertImsiTagStmt.ExecContext(ctx, tx)

			if err != nil {
				tx.Rollback()
				log.Println("insert-imei-configuration-tag-error", err.Error())
				return nil, 500, utils.InternalServerError
			}

		}
	}

	tx.Commit()

	var updatedBy graphql.NullID
	if imeiConfiguration.UpdatedBy != nil {
		updatedBy = graphql_utils.ConvertStringToNullID(imeiConfiguration.UpdatedBy)
	}

	var updatedAt graphql.NullTime
	if imeiConfiguration.UpdatedAt != nil {
		updatedAt = graphql.NullTime{Value: &graphql.Time{Time: *imeiConfiguration.UpdatedAt}}
	}

	return &ImeiConfiguration{
		ID:                graphql.ID(imeiConfiguration.ID.String()),
		ProjectId:         graphql.ID(imeiConfiguration.ProjectId.String()),
		StationLocationId: graphql.ID(imeiConfiguration.StationLocationId.String()),
		Imei:              imeiConfiguration.Imei,
		PermittedLabel:    enum.GetDevicePermittedLabel(imeiConfiguration.PermittedLabel.String()),
		BlacklistPriority: enum.GetBlacklistPriority(imeiConfiguration.BlacklistPriority.String()),
		CreatedBy:         imeiConfiguration.CreatedBy,
		CreatedAt:         graphql.Time{Time: imeiConfiguration.CreatedAt},
		UpdatedBy:         &updatedBy,
		UpdatedAt:         &updatedAt,
	}, 200, nil
}

func (ImeiConfigurationService) Create(data CreateImeiConfigurationData) (*ImeiConfiguration, int, error) {
	dbClient := db.GetPrimaryClient()
	ctx := context.Background()
	tx, err := dbClient.Begin()
	if err != nil {
		log.Println("init-create-imei-transaction-error", err.Error())
		return nil, 500, utils.InternalServerError
	}

	insertImeiStmt := table.ImeiConfiguration.
		INSERT(table.ImeiConfiguration.Imei, table.ImeiConfiguration.ID, table.ImeiConfiguration.ProjectId, table.ImeiConfiguration.StationLocationId, table.ImeiConfiguration.CreatedBy, table.ImeiConfiguration.PermittedLabel, table.ImeiConfiguration.BlacklistPriority).
		MODEL(model.ImeiConfiguration{
			ID:                uuid.New(),
			ProjectId:         data.ProjectId,
			Imei:              data.Imei,
			CreatedBy:         data.CreatedBy,
			PermittedLabel:    model.DevicePermittedLabel(data.PermittedLabel),
			BlacklistPriority: model.BlacklistPriority(data.BlacklistPriority),
			StationLocationId: data.StationLocationId,
		}).RETURNING(table.ImeiConfiguration.AllColumns)

	imeiConfiguration := model.ImeiConfiguration{}

	err = insertImeiStmt.QueryContext(ctx, tx, &imeiConfiguration)

	if err != nil {
		tx.Rollback()
		log.Println("insert-imei-configuration", err.Error())
		return nil, 500, utils.InternalServerError
	}

	if data.Tags != nil && len(*data.Tags) != 0 {
		for _, tag := range *data.Tags {
			upsertTagStmt := table.Tag.
				INSERT(table.Tag.ID, table.Tag.ProjectId, table.Tag.Title, table.Tag.CreatedBy, table.Tag.CreatedAt).
				MODEL(model.Tag{
					ID:        uuid.New(),
					ProjectId: data.ProjectId,
					Title:     tag,
					CreatedBy: data.CreatedBy,
					CreatedAt: time.Now(),
				}).
				ON_CONFLICT(table.Tag.Title, table.Tag.ProjectId).
				DO_UPDATE(pg.SET(table.Tag.ID.SET(table.Tag.EXCLUDED.ID))).
				RETURNING(table.Tag.AllColumns)
			tagResult := model.Tag{}

			err := upsertTagStmt.QueryContext(ctx, tx, &tagResult)

			if err != nil {
				tx.Rollback()
				log.Println("upsert-imei-configuration-tag-error", err.Error())
				return nil, 500, utils.InternalServerError
			}

			insertImsiTagStmt := table.ImeiConfigurationTag.
				INSERT(table.ImeiConfigurationTag.ID, table.ImeiConfigurationTag.ImeiConfigurationId, table.ImeiConfigurationTag.TagId, table.ImeiConfigurationTag.CreatedBy).
				MODEL(model.ImeiConfigurationTag{
					ID:                  uuid.New(),
					ImeiConfigurationId: imeiConfiguration.ID,
					TagId:               tagResult.ID,
					CreatedBy:           data.CreatedBy,
				})

			_, err = insertImsiTagStmt.ExecContext(ctx, tx)

			if err != nil {
				tx.Rollback()
				log.Println("insert-imei-configuration-tag-error", err.Error())
				return nil, 500, utils.InternalServerError
			}

		}
	}

	tx.Commit()

	var updatedBy graphql.NullID
	if imeiConfiguration.UpdatedBy != nil {
		updatedBy = graphql_utils.ConvertStringToNullID(imeiConfiguration.UpdatedBy)
	}

	var updatedAt graphql.NullTime
	if imeiConfiguration.UpdatedAt != nil {
		updatedAt = graphql.NullTime{Value: &graphql.Time{Time: *imeiConfiguration.UpdatedAt}}
	}

	return &ImeiConfiguration{
		ID:                graphql.ID(imeiConfiguration.ID.String()),
		ProjectId:         graphql.ID(imeiConfiguration.ProjectId.String()),
		StationLocationId: graphql.ID(imeiConfiguration.StationLocationId.String()),
		Imei:              imeiConfiguration.Imei,
		PermittedLabel:    enum.GetDevicePermittedLabel(imeiConfiguration.PermittedLabel.String()),
		BlacklistPriority: enum.GetBlacklistPriority(imeiConfiguration.BlacklistPriority.String()),
		CreatedBy:         imeiConfiguration.CreatedBy,
		CreatedAt:         graphql.Time{Time: imeiConfiguration.CreatedAt},
		UpdatedBy:         &updatedBy,
		UpdatedAt:         &updatedAt,
	}, 201, nil
}

func (ImeiConfigurationService) FindByIds(keys []uuid.UUID) ([]ImeiConfiguration, int, error) {
	dbClient := db.GetPrimaryClient()

	var ids []pg.Expression

	for _, id := range keys {
		ids = append(ids, pg.UUID(id))
	}

	getImeisStmt := table.ImeiConfiguration.SELECT(table.ImeiConfiguration.AllColumns).
		FROM(table.ImeiConfiguration).
		WHERE(table.ImeiConfiguration.ID.IN(ids...))
	imeiConfigurations := []model.ImeiConfiguration{}

	err := getImeisStmt.Query(dbClient, &imeiConfigurations)
	if err != nil {
		log.Println("get-imei-configuration-ids-error", err.Error())
		return nil, 500, utils.InternalServerError
	}

	imeiConfigurationsResponse := lo.Map(imeiConfigurations, func(item model.ImeiConfiguration, index int) ImeiConfiguration {
		var updatedBy graphql.NullID
		if item.UpdatedBy != nil {
			updatedBy = graphql_utils.ConvertStringToNullID(item.UpdatedBy)
		}

		var updatedAt graphql.NullTime
		if item.UpdatedAt != nil {
			updatedAt = graphql.NullTime{Value: &graphql.Time{Time: *item.UpdatedAt}}
		}

		return ImeiConfiguration{
			ID:                graphql.ID(item.ID.String()),
			ProjectId:         graphql.ID(item.ProjectId.String()),
			StationLocationId: graphql.ID(item.StationLocationId.String()),
			Imei:              item.Imei,
			BlacklistPriority: enum.GetBlacklistPriority(item.BlacklistPriority.String()),
			PermittedLabel:    enum.GetDevicePermittedLabel(item.PermittedLabel.String()),
			CreatedBy:         item.CreatedBy,
			CreatedAt:         graphql.Time{Time: item.CreatedAt},
			UpdatedBy:         &updatedBy,
			UpdatedAt:         &updatedAt,
		}
	})

	return imeiConfigurationsResponse, 200, nil

}

func (service ImeiConfigurationService) Dataloader() *dataloader.Loader {
	return dataloader.NewBatchedLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var ids []uuid.UUID
		for _, key := range keys {
			ids = append(ids, uuid.MustParse(key.String()))
		}

		imeiConfigurations, _, _ := service.FindByIds(ids)

		var results []*dataloader.Result

		for _, key := range keys {
			imsiConfiguration, match := lo.Find(imeiConfigurations, func(item ImeiConfiguration) bool {
				return string(item.ID) == string(graphql.ID(key.String()))
			})

			if match {
				results = append(results, &dataloader.Result{Data: imsiConfiguration})

			}

		}

		return results
	}, dataloader.WithClearCacheOnBatch())
}
