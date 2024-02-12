package imsiconfiguration

import (
	"checkpoint/.gen/checkpoint/public/model"
	table "checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/utils/graphql_utils"
	"context"
	"errors"

	"checkpoint/utils"
	"log"
	"time"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/graph-gophers/dataloader"
	"github.com/graph-gophers/graphql-go"
	"github.com/samber/lo"
)

type ImsiConfigurationService struct{}

func (ImsiConfigurationService) FindMany(data GetImsiConfigurationsData) ([]Imsiconfiguration, int, error) {
	dbClient := db.GetPrimaryClient()

	conditions := pg.Bool(true).
		AND(table.ImsiConfiguration.DeletedAt.IS_NULL()).
		AND(table.ImsiConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId))).
		AND(table.ImsiConfiguration.StationLocationId.EQ(pg.UUID(data.StationLocationId)))

	var fromConditions pg.ReadableTable = table.ImsiConfiguration

	if data.PermittedLabel != nil {
		conditions = conditions.AND(table.ImsiConfiguration.PermittedLabel.EQ(pg.NewEnumValue(*data.PermittedLabel)))
	}

	if data.BlacklistPriority != nil {
		conditions = conditions.AND(table.ImsiConfiguration.BlacklistPriority.EQ(pg.NewEnumValue(*data.BlacklistPriority)))
	}

	if data.Tags != nil && len(*data.Tags) != 0 {
		var tagItems []pg.Expression

		for _, tag := range *data.Tags {
			tagItems = append(tagItems, pg.String(tag))
		}

		fromConditions = fromConditions.
			INNER_JOIN(table.ImsiConfigurationTag, table.ImsiConfigurationTag.ImsiConfigurationId.EQ(table.ImsiConfiguration.ID)).
			INNER_JOIN(table.Tag, table.ImsiConfigurationTag.TagId.EQ(table.Tag.ID))
		conditions = conditions.AND(table.Tag.Title.IN(tagItems...))
	}

	if data.Search != nil {
		conditions = conditions.AND(table.ImsiConfiguration.Imsi.LIKE(pg.String(*data.Search)))
	}

	if data.Mcc != nil {
		conditions = conditions.AND(table.ImsiConfiguration.Mcc.EQ(pg.String(*data.Mcc)))
	}

	if data.Mnc != nil {
		conditions = conditions.AND(table.ImsiConfiguration.Mnc.EQ(pg.String(*data.Mnc)))
	}

	getImsiConfigurationsStmt := table.ImsiConfiguration.
		SELECT(table.ImsiConfiguration.AllColumns).
		FROM(fromConditions).
		WHERE(conditions).
		LIMIT(data.Pagination.Limit).
		OFFSET(data.Pagination.Skip)

	imsiConfigurations := []model.ImsiConfiguration{}

	err := getImsiConfigurationsStmt.Query(dbClient, &imsiConfigurations)

	if err != nil {
		log.Println("get-imsi-configurations-error", err.Error())
		return nil, 500, utils.InternalServerError
	}

	imsiConfigurationsResponse := lo.Map(imsiConfigurations, func(item model.ImsiConfiguration, index int) Imsiconfiguration {
		var updatedBy graphql.NullID
		if item.UpdatedBy != nil {
			updatedBy = graphql_utils.ConvertStringToNullID(item.UpdatedBy)
		}

		var updatedAt graphql.NullTime
		if item.UpdatedAt != nil {
			updatedAt = graphql.NullTime{Value: &graphql.Time{Time: *item.UpdatedAt}}
		}
		return Imsiconfiguration{
			ID:                graphql.ID(item.ID.String()),
			ProjectId:         graphql.ID(item.ProjectId.String()),
			Imsi:              item.Imsi,
			CreatedBy:         graphql.ID(item.CreatedBy),
			UpdatedBy:         &updatedBy,
			CreatedAt:         graphql.Time{Time: item.CreatedAt},
			UpdatedAt:         &updatedAt,
			PermittedLabel:    string(item.PermittedLabel),
			BlacklistPriority: item.BlacklistPriority.String(),
			StationLocationId: graphql.ID(item.StationLocationId.String()),
			Mcc:               item.Mcc,
			Mnc:               item.Mnc,
		}
	})

	return imsiConfigurationsResponse, 200, nil
}

func (ImsiConfigurationService) FindById(data GetImsiConfigurationByIdData) (*Imsiconfiguration, int, error) {
	dbClient := db.GetPrimaryClient()
	getImsiStmt := table.ImsiConfiguration.
		SELECT(table.ImsiConfiguration.AllColumns).
		FROM(table.ImsiConfiguration).
		WHERE(table.ImsiConfiguration.ID.EQ(pg.UUID(data.ID)).
			AND(table.ImsiConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId))).
			AND(table.ImsiConfiguration.DeletedAt.IS_NULL()))

	imsiConfiguration := model.ImsiConfiguration{}

	err := getImsiStmt.Query(dbClient, &imsiConfiguration)
	if err != nil && db.HasNoRow(err) {
		return nil, 404, errors.New("imsi id not found")
	}

	if err != nil {
		return nil, 500, utils.InternalServerError
	}

	var updatedBy graphql.NullID
	if imsiConfiguration.UpdatedBy != nil {

		updatedBy = graphql_utils.ConvertStringToNullID(imsiConfiguration.UpdatedBy)
	}

	var updatedAt graphql.NullTime
	if imsiConfiguration.UpdatedAt != nil {
		updatedAt = graphql.NullTime{Value: &graphql.Time{Time: *imsiConfiguration.UpdatedAt}}
	}

	return &Imsiconfiguration{
		ID:                graphql.ID(imsiConfiguration.ID.String()),
		ProjectId:         graphql.ID(imsiConfiguration.ProjectId.String()),
		Imsi:              imsiConfiguration.Imsi,
		CreatedBy:         graphql.ID(imsiConfiguration.CreatedBy),
		UpdatedBy:         &updatedBy,
		CreatedAt:         graphql.Time{Time: imsiConfiguration.CreatedAt},
		UpdatedAt:         &updatedAt,
		PermittedLabel:    imsiConfiguration.PermittedLabel.String(),
		BlacklistPriority: imsiConfiguration.BlacklistPriority.String(),
		StationLocationId: graphql.ID(imsiConfiguration.StationLocationId.String()),
		Mcc:               imsiConfiguration.Mcc,
		Mnc:               imsiConfiguration.Mnc,
	}, 200, nil
}

func (ImsiConfigurationService) Delete(data DeleteImsiConfigurationData) (int, error) {
	dbClient := db.GetPrimaryClient()
	now := time.Now()
	deleteImsiStmt := table.ImsiConfiguration.
		UPDATE(table.ImsiConfiguration.DeletedAt, table.ImsiConfiguration.DeletedBy).
		MODEL(model.ImsiConfiguration{
			DeletedAt: &now,
			DeletedBy: &data.DeletedBy,
		}).
		WHERE(table.ImsiConfiguration.ID.EQ(pg.UUID(data.ID)).
			AND(table.ImsiConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId))))

	_, err := deleteImsiStmt.Exec(dbClient)
	if err != nil && db.HasNoRow(err) {
		return 403, utils.ForbiddenOperation
	}

	return 200, nil
}

func (ImsiConfigurationService) Update(data UpdateImsiConfigurationData) (*Imsiconfiguration, int, error) {
	dbClient := db.GetPrimaryClient()
	ctx := context.Background()
	tx, err := dbClient.Begin()

	if err != nil {
		log.Println("init-db-tx", err.Error())
		return nil, 500, utils.InternalServerError
	}

	mcc, mnc, err := utils.ExtractMCCMNC(*data.Imsi)

	if err != nil {
		return nil, 400, err
	}

	var columnsToUpdate pg.ColumnList
	columnsToUpdate = append(columnsToUpdate, table.ImsiConfiguration.Imsi)

	if data.PermittedLabel != nil {
		columnsToUpdate = append(columnsToUpdate, table.ImsiConfiguration.PermittedLabel)
	}

	if data.BlacklistPriority != nil {
		columnsToUpdate = append(columnsToUpdate, table.ImsiConfiguration.BlacklistPriority)
	}

	now := time.Now()

	updateImsiStmt := table.ImsiConfiguration.
		UPDATE(columnsToUpdate).
		MODEL(model.ImsiConfiguration{
			Imsi:              *data.Imsi,
			BlacklistPriority: model.BlacklistPriority(*data.BlacklistPriority),
			PermittedLabel:    model.DevicePermittedLabel(*data.PermittedLabel),
			Mcc:               mcc,
			Mnc:               mnc,
			UpdatedBy:         &data.UpdatedBy,
			UpdatedAt:         &now,
		}).
		WHERE(table.ImsiConfiguration.ID.EQ(pg.UUID(data.ID)).
			AND(table.ImsiConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId))).
			AND(table.ImsiConfiguration.DeletedAt.IS_NULL())).
		RETURNING(table.ImsiConfiguration.AllColumns)

	imsiConfiguration := model.ImsiConfiguration{}
	err = updateImsiStmt.Query(tx, &imsiConfiguration)

	if err != nil && db.HasNoRow(err) {
		tx.Rollback()
		return nil, 403, utils.ForbiddenOperation
	}

	if err != nil && db.InvalidInput(err) {
		tx.Rollback()
		log.Println("invalid-update-imsi-configuration-error", err.Error())
		return nil, 400, err
	}

	if err != nil {
		tx.Rollback()
		log.Println("insert-imsi-configuration-error", err.Error())
		return nil, 500, utils.InternalServerError
	}

	if data.Tags != nil && len(*data.Tags) != 0 {
		deleteAllImsiTagStmt := table.ImsiConfigurationTag.
			DELETE().
			WHERE(table.ImsiConfigurationTag.ImsiConfigurationId.EQ(pg.UUID(data.ID)))

		_, err = deleteAllImsiTagStmt.ExecContext(ctx, tx)

		if err != nil {
			tx.Rollback()
			log.Println("delete-all-imsi-configuration-error", err.Error())
			return nil, 500, utils.InternalServerError
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
				log.Println("upsert-imsi-configuration-tag-error", err.Error())
				return nil, 500, utils.InternalServerError
			}

			insertImsiTagStmt := table.ImsiConfigurationTag.
				INSERT(table.ImsiConfigurationTag.ID, table.ImsiConfigurationTag.ImsiConfigurationId, table.ImsiConfigurationTag.TagId, table.ImsiConfigurationTag.CreatedBy).
				MODEL(model.ImsiConfigurationTag{
					ID:                  uuid.New(),
					ImsiConfigurationId: imsiConfiguration.ID,
					TagId:               tagResult.ID,
					CreatedBy:           data.UpdatedBy,
				})
			_, err = insertImsiTagStmt.ExecContext(ctx, tx)

			if err != nil {
				tx.Rollback()
				log.Println("insert-imsi-configuration-tag-error", err.Error())
				return nil, 500, utils.InternalServerError
			}

		}
	}

	var updatedBy graphql.NullID
	if imsiConfiguration.UpdatedBy != nil {
		updatedBy = graphql_utils.ConvertStringToNullID(imsiConfiguration.UpdatedBy)
	}

	var updatedAt graphql.NullTime
	if imsiConfiguration.UpdatedAt != nil {
		updatedAt = graphql.NullTime{Value: &graphql.Time{Time: *imsiConfiguration.UpdatedAt}}
	}

	tx.Commit()

	return &Imsiconfiguration{
		ID:                graphql.ID(imsiConfiguration.ID.String()),
		ProjectId:         graphql.ID(imsiConfiguration.ProjectId.String()),
		Imsi:              imsiConfiguration.Imsi,
		CreatedBy:         graphql.ID(imsiConfiguration.CreatedBy),
		UpdatedBy:         &updatedBy,
		CreatedAt:         graphql.Time{Time: imsiConfiguration.CreatedAt},
		UpdatedAt:         &updatedAt,
		PermittedLabel:    imsiConfiguration.PermittedLabel.String(),
		BlacklistPriority: imsiConfiguration.BlacklistPriority.String(),
		StationLocationId: graphql.ID(imsiConfiguration.StationLocationId.String()),
		Mcc:               imsiConfiguration.Mcc,
		Mnc:               imsiConfiguration.Mnc,
	}, 200, nil
}

func (ImsiConfigurationService) Create(data CreateImsiConfigurationData) (*Imsiconfiguration, int, error) {
	dbClient := db.GetPrimaryClient()
	ctx := context.Background()

	mcc, mnc, err := utils.ExtractMCCMNC(data.Imsi)

	if err != nil {
		return nil, 400, err
	}

	tx, err := dbClient.Begin()

	if err != nil {
		return nil, 500, err
	}

	insertImsiStmt := table.ImsiConfiguration.
		INSERT(table.ImsiConfiguration.ID, table.ImsiConfiguration.Imsi, table.ImsiConfiguration.BlacklistPriority, table.ImsiConfiguration.StationLocationId, table.ImsiConfiguration.PermittedLabel, table.ImsiConfiguration.CreatedBy, table.ImsiConfiguration.ProjectId, table.ImsiConfiguration.Mcc, table.ImsiConfiguration.Mnc).
		MODEL(model.ImsiConfiguration{
			ID:                uuid.New(),
			Imsi:              data.Imsi,
			BlacklistPriority: model.BlacklistPriority(data.BlacklistPriority),
			StationLocationId: data.StationLocationId,
			PermittedLabel:    model.DevicePermittedLabel(data.PermittedLabel),
			CreatedBy:         data.CreatedBy,
			ProjectId:         data.ProjectId,
			Mcc:               mcc,
			Mnc:               mnc,
		}).
		RETURNING(table.ImsiConfiguration.AllColumns)

	imsiConfiguration := model.ImsiConfiguration{}
	err = insertImsiStmt.QueryContext(ctx, tx, &imsiConfiguration)

	if err != nil && db.InvalidInput(err) {
		tx.Rollback()
		log.Println("invalid-insert-imsi-configuration-error", err.Error())
		return nil, 400, err
	}

	if err != nil {
		tx.Rollback()
		log.Println("insert-imsi-configuration-error", err.Error())
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
				log.Println("upsert-imsi-configuration-tag-error", err.Error())
				return nil, 500, utils.InternalServerError
			}

			insertImsiTagStmt := table.ImsiConfigurationTag.
				INSERT(table.ImsiConfigurationTag.ID, table.ImsiConfigurationTag.ImsiConfigurationId, table.ImsiConfigurationTag.TagId, table.ImsiConfigurationTag.CreatedBy).
				MODEL(model.ImsiConfigurationTag{
					ID:                  uuid.New(),
					ImsiConfigurationId: imsiConfiguration.ID,
					TagId:               tagResult.ID,
					CreatedBy:           data.CreatedBy,
				})
			_, err = insertImsiTagStmt.ExecContext(ctx, tx)

			if err != nil {
				tx.Rollback()
				log.Println("insert-imsi-configuration-tag-error", err.Error())
				return nil, 500, utils.InternalServerError
			}

		}
	}
	err = tx.Commit()

	if err != nil {
		log.Println("commit-imsi-configuration-error", err.Error())
		return nil, 500, utils.InternalServerError
	}

	var updatedBy graphql.NullID
	if imsiConfiguration.UpdatedBy != nil {
		updatedBy = graphql_utils.ConvertStringToNullID(imsiConfiguration.UpdatedBy)
	}

	var updatedAt graphql.NullTime
	if imsiConfiguration.UpdatedAt != nil {
		updatedAt = graphql.NullTime{Value: &graphql.Time{Time: *imsiConfiguration.UpdatedAt}}
	}

	return &Imsiconfiguration{
		ID:                graphql.ID(imsiConfiguration.ID.String()),
		ProjectId:         graphql.ID(imsiConfiguration.ProjectId.String()),
		Imsi:              imsiConfiguration.Imsi,
		CreatedBy:         graphql.ID(imsiConfiguration.CreatedBy),
		UpdatedBy:         &updatedBy,
		CreatedAt:         graphql.Time{Time: imsiConfiguration.CreatedAt},
		UpdatedAt:         &updatedAt,
		PermittedLabel:    imsiConfiguration.PermittedLabel.String(),
		BlacklistPriority: imsiConfiguration.BlacklistPriority.String(),
		StationLocationId: graphql.ID(imsiConfiguration.StationLocationId.String()),
		Mcc:               imsiConfiguration.Mcc,
		Mnc:               imsiConfiguration.Mnc,
	}, 201, nil
}

func (ImsiConfigurationService) FindByIds(keys []uuid.UUID) ([]Imsiconfiguration, int, error) {
	dbClient := db.GetPrimaryClient()
	var ids []pg.Expression

	for _, id := range keys {
		ids = append(ids, pg.UUID(id))
	}

	getImsiConfigurationsStmt := table.ImsiConfiguration.
		SELECT(table.ImsiConfiguration.AllColumns).
		FROM(table.ImsiConfiguration).
		WHERE(table.ImsiConfiguration.ID.IN(ids...))

	imsiConfigurations := []model.ImsiConfiguration{}

	err := getImsiConfigurationsStmt.Query(dbClient, &imsiConfigurations)

	if err != nil {
		log.Println("get-imsi-configurations-error", err.Error())
		return nil, 500, utils.InternalServerError
	}

	imsiConfigurationsResponse := lo.Map(imsiConfigurations, func(item model.ImsiConfiguration, index int) Imsiconfiguration {
		var updatedBy graphql.NullID
		if item.UpdatedBy != nil {
			updatedBy = graphql_utils.ConvertStringToNullID(item.UpdatedBy)
		}

		var updatedAt graphql.NullTime
		if item.UpdatedAt != nil {
			updatedAt = graphql.NullTime{Value: &graphql.Time{Time: *item.UpdatedAt}}
		}
		return Imsiconfiguration{
			ID:                graphql.ID(item.ID.String()),
			ProjectId:         graphql.ID(item.ProjectId.String()),
			Imsi:              item.Imsi,
			CreatedBy:         graphql.ID(item.CreatedBy),
			UpdatedBy:         &updatedBy,
			CreatedAt:         graphql.Time{Time: item.CreatedAt},
			UpdatedAt:         &updatedAt,
			PermittedLabel:    string(item.PermittedLabel),
			BlacklistPriority: item.BlacklistPriority.String(),
			StationLocationId: graphql.ID(item.StationLocationId.String()),
			Mcc:               item.Mcc,
			Mnc:               item.Mnc,
		}
	})

	return imsiConfigurationsResponse, 200, nil
}

func (service ImsiConfigurationService) Dataloader() *dataloader.Loader {
	return dataloader.NewBatchedLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var ids []uuid.UUID
		for _, key := range keys {
			ids = append(ids, uuid.MustParse(key.String()))
		}

		imsiConfigurations, _, _ := service.FindByIds(ids)

		var results []*dataloader.Result

		for _, key := range keys {
			imsiConfiguration, match := lo.Find(imsiConfigurations, func(item Imsiconfiguration) bool {
				return string(item.ID) == string(graphql.ID(key.String()))
			})

			if match {
				results = append(results, &dataloader.Result{Data: imsiConfiguration})

			}

		}

		return results
	}, dataloader.WithClearCacheOnBatch())
}
