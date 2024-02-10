package imeiconfiguration

import (
	"checkpoint/.gen/checkpoint/public/model"
	table "checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/utils"
	"checkpoint/utils/graphql_utils"
	"context"
	"log"
	"time"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/graph-gophers/graphql-go"
	"github.com/samber/lo"

	"github.com/google/uuid"
)

type ImeiConfigurationService struct{}

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

	var fromCondition pg.ReadableTable = table.ImeiConfiguration

	if data.Label != nil {
		conditions = conditions.AND(table.ImeiConfiguration.PermittedLabel.EQ(pg.NewEnumValue(*data.Label)))
	}

	if data.Search != nil {
		conditions = conditions.AND(table.ImeiConfiguration.Imei.LIKE(pg.String(*data.Search)))
	}

	if data.Tags != nil {
		var tagItems []pg.Expression

		for _, tag := range *data.Tags {
			tagItems = append(tagItems, pg.String(tag))
		}

		fromCondition = fromCondition.
			INNER_JOIN(table.ImeiConfigurationTag, table.ImeiConfigurationTag.ImeiConfigurationId.EQ(table.ImeiConfiguration.ID)).
			INNER_JOIN(table.Tag, table.ImeiConfigurationTag.TagId.EQ(table.Tag.ID))

		conditions = conditions.AND(table.Tag.Title.IN(tagItems...))
	}

	getImeisStmt := table.ImeiConfiguration.SELECT(table.ImeiConfiguration.AllColumns).
		FROM(fromCondition).
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
			Priority:          item.Priority.String(),
			PermittedLabel:    item.PermittedLabel.String(),
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
		return nil, 403, utils.ForbiddenOperation
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
		PermittedLabel:    imeiConfiguration.PermittedLabel.String(),
		Priority:          imeiConfiguration.Priority.String(),
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
		UPDATE(table.ImeiConfiguration.Imei, table.ImeiConfiguration.UpdatedBy, table.ImeiConfiguration.PermittedLabel, table.ImeiConfiguration.Priority, table.ImeiConfiguration.UpdatedAt).
		MODEL(model.ImeiConfiguration{
			Imei:           data.Imei,
			UpdatedBy:      &data.UpdatedBy,
			UpdatedAt:      &now,
			PermittedLabel: model.DevicePermittedLabel(data.Label),
			Priority:       model.BlacklistPriority(data.Priority),
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

	if data.Tags != nil {
		deleteAllImeiTagsStmt := table.ImeiConfigurationTag.
			DELETE().
			WHERE(table.ImeiConfigurationTag.ImeiConfigurationId.EQ(pg.UUID(imeiConfiguration.ID)))
		_, err := deleteAllImeiTagsStmt.ExecContext(ctx, tx)
		if err != nil {
			tx.Rollback()
			log.Println("delete-imei-configuraiton-tag-error", err.Error())
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
				log.Println("upsert-imei-configuraiton-tag-error", err.Error())
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
		PermittedLabel:    imeiConfiguration.PermittedLabel.String(),
		Priority:          imeiConfiguration.Priority.String(),
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
		INSERT(table.ImeiConfiguration.Imei, table.ImeiConfiguration.ID, table.ImeiConfiguration.ProjectId, table.ImeiConfiguration.StationLocationId, table.ImeiConfiguration.CreatedBy, table.ImeiConfiguration.PermittedLabel, table.ImeiConfiguration.Priority).
		MODEL(model.ImeiConfiguration{
			ID:                uuid.New(),
			ProjectId:         data.ProjectId,
			Imei:              data.Imei,
			CreatedBy:         data.CreatedBy,
			PermittedLabel:    model.DevicePermittedLabel(data.Label),
			Priority:          model.BlacklistPriority(data.Priority),
			StationLocationId: data.StationLocationId,
		}).RETURNING(table.ImeiConfiguration.AllColumns)

	imeiConfiguration := model.ImeiConfiguration{}

	err = insertImeiStmt.QueryContext(ctx, tx, &imeiConfiguration)

	if err != nil {
		tx.Rollback()
		log.Println("insert-imei-configuration", err.Error())
		return nil, 500, utils.InternalServerError
	}

	if data.Tags != nil {
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
				log.Println("upsert-imei-configuraiton-tag-error", err.Error())
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
		PermittedLabel:    imeiConfiguration.PermittedLabel.String(),
		Priority:          imeiConfiguration.Priority.String(),
		CreatedBy:         imeiConfiguration.CreatedBy,
		CreatedAt:         graphql.Time{Time: imeiConfiguration.CreatedAt},
		UpdatedBy:         &updatedBy,
		UpdatedAt:         &updatedAt,
	}, 201, nil
}
