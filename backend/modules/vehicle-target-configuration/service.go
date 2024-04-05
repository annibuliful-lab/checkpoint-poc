package vehicletargetconfiguration

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/gql/enum"
	tagUtils "checkpoint/modules/tag"
	vehicleproperty "checkpoint/modules/vehicle-property"
	"checkpoint/utils"
	"context"
	"database/sql"
	"log"
	"time"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/samber/lo"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

type VehicleTargetConfigurationService struct{}

func (VehicleTargetConfigurationService) FindMany(data GetVehicleTargetsConfigurationData) ([]*VehicleTargetConfiguration, error) {
	dbClient := db.GetPrimaryClient()

	conditions := table.VehicleTargetConfiguration.DeletedAt.IS_NULL().
		AND(table.VehicleTargetConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId)))

	var fromConditions pg.ReadableTable = table.VehicleTargetConfiguration

	if data.StationLocationId != nil {
		conditions = conditions.AND(
			table.VehicleTargetConfiguration.StationLocationId.EQ(
				pg.UUID(data.StationLocationId),
			),
		)
	}

	if data.BlacklistPriority != nil {
		conditions = conditions.
			AND(table.VehicleTargetConfiguration.BlacklistPriority.
				EQ(pg.NewEnumValue(string(*data.BlacklistPriority))),
			)
	}

	if data.PermittedLabel != nil {
		conditions = conditions.
			AND(table.VehicleTargetConfiguration.PermittedLabel.
				EQ(pg.NewEnumValue(string(*data.PermittedLabel))),
			)
	}

	if data.Search != nil {
		conditions = conditions.
			AND(
				table.VehicleTargetConfiguration.Prefix.LIKE(pg.String(*data.Search)).
					OR(table.VehicleTargetConfiguration.Number.LIKE(pg.String(*data.Search))).
					OR(table.VehicleTargetConfiguration.Province.LIKE(pg.String(*data.Search))).
					OR(table.VehicleTargetConfiguration.Country.LIKE(pg.String(*data.Search))),
			)
	}

	if data.Type != nil {
		conditions = conditions.AND(table.VehicleTargetConfiguration.Type.EQ(pg.String(*data.Type)))
	}

	if data.Tags != nil && len(*data.Tags) != 0 {
		var tagItems []pg.Expression

		for _, tag := range *data.Tags {
			tagItems = append(tagItems, pg.String(tag))
		}

		fromConditions = fromConditions.
			INNER_JOIN(table.VehicleTargetConfigurationTag, table.VehicleTargetConfigurationTag.VehicleTargetConfigurationId.EQ(table.VehicleTargetConfiguration.ID)).
			INNER_JOIN(table.Tag, table.VehicleTargetConfigurationTag.TagId.EQ(table.Tag.ID))
		conditions = conditions.AND(table.Tag.Title.IN(tagItems...))
	}

	getTargetsStmt := table.VehicleTargetConfiguration.
		SELECT(table.VehicleTargetConfiguration.AllColumns).
		FROM(fromConditions).
		WHERE(conditions).
		LIMIT(data.Limit).
		OFFSET(data.Skip)

	vehicleTargets := []model.VehicleTargetConfiguration{}

	err := getTargetsStmt.Query(dbClient, &vehicleTargets)

	if err != nil {
		log.Println("get-vehicle-targets-configuration-error", err.Error())
		return nil, utils.InternalServerError
	}

	return lo.Map(vehicleTargets, func(item model.VehicleTargetConfiguration, index int) *VehicleTargetConfiguration {
		return transformToGraphql(item)
	}), nil

}

func (VehicleTargetConfigurationService) FindById(data GetVehicleTargetConfigurationData) (*VehicleTargetConfiguration, error) {
	dbClient := db.GetPrimaryClient()

	getStmt := table.VehicleTargetConfiguration.
		SELECT(table.VehicleTargetConfiguration.AllColumns).
		WHERE(
			table.VehicleTargetConfiguration.ID.EQ(pg.UUID(data.ID)).
				AND(table.VehicleTargetConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId))).
				AND(table.VehicleTargetConfiguration.DeletedAt.IS_NULL()),
		)

	vehicleTarget := model.VehicleTargetConfiguration{}

	err := getStmt.Query(dbClient, &vehicleTarget)

	if err != nil && db.HasNoRow(err) {
		return nil, utils.Notfound
	}

	if err != nil {
		log.Println("get-vehicle-target-configuration-by-id-error", err.Error())
		return nil, utils.InternalServerError
	}

	return transformToGraphql(vehicleTarget), nil

}

func (VehicleTargetConfigurationService) Delete(data DeleteVehicleTargetConfigurationData) (*utils.DeleteOperation, error) {
	dbClient := db.GetPrimaryClient()
	now := time.Now()

	deleteStmt := table.VehicleTargetConfiguration.
		UPDATE(
			table.VehicleTargetConfiguration.DeletedAt,
			table.VehicleTargetConfiguration.DeletedBy,
		).
		MODEL(model.VehicleTargetConfiguration{
			DeletedAt: &now,
			DeletedBy: &data.DeletedBy,
		}).
		WHERE(
			table.VehicleTargetConfiguration.ID.EQ(pg.UUID(data.ID)).
				AND(table.VehicleTargetConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId))),
		)

	result, err := deleteStmt.Exec(dbClient)

	if db.HasNoAffectedRow(result) {
		return nil, utils.Notfound
	}

	if err != nil {
		log.Println("init-tx-update-vehicle-target-configuration-error", err.Error())
		return nil, utils.InternalServerError
	}

	return &utils.DeleteOperation{
		Success: true,
	}, nil
}

func (VehicleTargetConfigurationService) Update(data UpdateVehicleTargetConfigurationData) (*VehicleTargetConfiguration, error) {
	dbClient := db.GetPrimaryClient()
	ctx := context.Background()
	tx, err := dbClient.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
	})

	if err != nil {
		log.Println("init-tx-update-vehicle-target-configuration-error", err.Error())
		return nil, utils.InternalServerError
	}

	now := time.Now()
	var columnsToUpdate pg.ColumnList
	fieldsToUpdate := model.VehicleTargetConfiguration{
		UpdatedBy: &data.UpdatedBy,
		UpdatedAt: &now,
	}

	if data.BlacklistPriority != nil {
		columnsToUpdate = append(columnsToUpdate, table.VehicleTargetConfiguration.BlacklistPriority)
		fieldsToUpdate.BlacklistPriority = *data.BlacklistPriority
	}

	if data.Country != nil {
		columnsToUpdate = append(columnsToUpdate, table.VehicleTargetConfiguration.Country)
		fieldsToUpdate.Country = data.Country
	}

	if data.Number != nil {
		columnsToUpdate = append(columnsToUpdate, table.VehicleTargetConfiguration.Number)
		fieldsToUpdate.Number = *data.Number
	}

	if data.PermittedLabel != nil {
		columnsToUpdate = append(columnsToUpdate, table.VehicleTargetConfiguration.PermittedLabel)
		fieldsToUpdate.PermittedLabel = *data.PermittedLabel
	}

	if data.Prefix != nil {
		columnsToUpdate = append(columnsToUpdate, table.VehicleTargetConfiguration.Prefix)
		fieldsToUpdate.Prefix = *data.Prefix
	}

	if data.Province != nil {
		columnsToUpdate = append(columnsToUpdate, table.VehicleTargetConfiguration.Province)
		fieldsToUpdate.Province = *data.Province
	}

	if data.Type != nil {
		upsertTypeStmt := vehicleproperty.UpsertStatement(vehicleproperty.UpsertVehiclePropertyData{
			Property:  *data.Type,
			Type:      model.PropertyType_VehicleType,
			ProjectId: data.ProjectId,
		})

		_, err := upsertTypeStmt.ExecContext(ctx, tx)
		if err != nil {
			tx.Rollback()
			log.Println("tx-upsert-vehicle-target-configuration-type-error", err.Error())
			return nil, utils.InternalServerError
		}

		columnsToUpdate = append(columnsToUpdate, table.VehicleTargetConfiguration.Type)
		fieldsToUpdate.Type = *data.Type
	}

	if data.Color != nil {
		upsertColorStmt := vehicleproperty.UpsertStatement(vehicleproperty.UpsertVehiclePropertyData{
			Property:  *data.Color,
			Type:      model.PropertyType_VehicleColor,
			ProjectId: data.ProjectId,
		})

		_, err = upsertColorStmt.ExecContext(ctx, tx)
		if err != nil {
			tx.Rollback()
			log.Println("tx-upsert-vehicle-target-configuration-color-error", err.Error())
			return nil, utils.InternalServerError
		}
		columnsToUpdate = append(columnsToUpdate, table.VehicleTargetConfiguration.Color)
		fieldsToUpdate.Color = *data.Color
	}

	if data.Brand != nil {
		upsertBrandStmt := vehicleproperty.UpsertStatement(vehicleproperty.UpsertVehiclePropertyData{
			Property:  *data.Brand,
			Type:      model.PropertyType_VehicleBrand,
			ProjectId: data.ProjectId,
		})

		_, err = upsertBrandStmt.ExecContext(ctx, tx)
		if err != nil {
			tx.Rollback()
			log.Println("tx-upsert-vehicle-target-configuration-brand-error", err.Error())
			return nil, utils.InternalServerError
		}

		columnsToUpdate = append(columnsToUpdate, table.VehicleTargetConfiguration.Brand)
		fieldsToUpdate.Brand = *data.Brand
	}

	updateVehicleTargetStmt := table.VehicleTargetConfiguration.
		UPDATE(columnsToUpdate).
		MODEL(fieldsToUpdate).
		WHERE(
			table.VehicleTargetConfiguration.ID.EQ(pg.UUID(data.ID)).
				AND(table.VehicleTargetConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId))).
				AND(table.VehicleTargetConfiguration.DeletedAt.IS_NULL()),
		).
		RETURNING(table.VehicleTargetConfiguration.AllColumns)

	vehicleTarget := model.VehicleTargetConfiguration{}

	err = updateVehicleTargetStmt.QueryContext(ctx, tx, &vehicleTarget)

	if err != nil && db.HasNoRow(err) {
		return nil, utils.Notfound
	}

	if err != nil {
		tx.Rollback()
		log.Println("update-vehicle-target-configuration-error", err.Error())
		return nil, utils.InternalServerError
	}

	if data.Tags != nil && len(*data.Tags) != 0 {

		deleteTagsStmt := table.VehicleTargetConfigurationTag.DELETE().WHERE(table.VehicleTargetConfigurationTag.VehicleTargetConfigurationId.EQ(pg.UUID(data.ID)))
		_, err = deleteTagsStmt.ExecContext(ctx, tx)

		if err != nil {
			tx.Rollback()
			log.Println("delete-vehicle-target-configuration-tags-error", err.Error())
			return nil, utils.InternalServerError
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
				log.Println("upsert-vehicle-target-configuration-tag-error", err.Error())
				return nil, utils.InternalServerError
			}

			insertVehicleTargetStmt := table.VehicleTargetConfigurationTag.
				INSERT(
					table.VehicleTargetConfigurationTag.ID,
					table.VehicleTargetConfigurationTag.VehicleTargetConfigurationId,
					table.VehicleTargetConfigurationTag.TagId,
					table.VehicleTargetConfigurationTag.CreatedBy,
				).
				MODEL(model.VehicleTargetConfigurationTag{
					ID:                           uuid.New(),
					VehicleTargetConfigurationId: vehicleTarget.ID,
					TagId:                        tagResult.ID,
					CreatedBy:                    data.UpdatedBy,
				})

			_, err = insertVehicleTargetStmt.ExecContext(ctx, tx)

			if err != nil {
				tx.Rollback()
				log.Println("insert-vehicle-target-configuration-tag-error", err.Error())
				return nil, utils.InternalServerError
			}
		}
	}

	tx.Commit()

	return transformToGraphql(vehicleTarget), nil

}

func (VehicleTargetConfigurationService) Create(data CreateVehicleTargetConfigurationData) (*VehicleTargetConfiguration, error) {
	dbClient := db.GetPrimaryClient()
	ctx := context.Background()
	tx, err := dbClient.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
	})

	if err != nil {
		log.Println("init-tx-create-vehicle-target-configuration-error", err.Error())
		return nil, utils.InternalServerError
	}

	upsertColorStmt := vehicleproperty.UpsertStatement(vehicleproperty.UpsertVehiclePropertyData{
		Property:  data.Color,
		Type:      model.PropertyType_VehicleColor,
		ProjectId: data.ProjectId,
	})

	_, err = upsertColorStmt.ExecContext(ctx, tx)
	if err != nil {
		tx.Rollback()
		log.Println("tx-upsert-vehicle-target-configuration-color-error", err.Error())
		return nil, utils.InternalServerError
	}

	upsertTypeStmt := vehicleproperty.UpsertStatement(vehicleproperty.UpsertVehiclePropertyData{
		Property:  data.Type,
		Type:      model.PropertyType_VehicleType,
		ProjectId: data.ProjectId,
	})

	_, err = upsertTypeStmt.ExecContext(ctx, tx)
	if err != nil {
		tx.Rollback()
		log.Println("tx-upsert-vehicle-target-configuration-type-error", err.Error())
		return nil, utils.InternalServerError
	}

	upsertBrandStmt := vehicleproperty.UpsertStatement(vehicleproperty.UpsertVehiclePropertyData{
		Property:  data.Brand,
		Type:      model.PropertyType_VehicleBrand,
		ProjectId: data.ProjectId,
	})

	_, err = upsertBrandStmt.ExecContext(ctx, tx)
	if err != nil {
		tx.Rollback()
		log.Println("tx-upsert-vehicle-target-configuration-brand-error", err.Error())
		return nil, utils.InternalServerError
	}

	createVehicleTargetStmt := table.VehicleTargetConfiguration.
		INSERT(
			table.VehicleTargetConfiguration.ID,
			table.VehicleTargetConfiguration.StationLocationId,
			table.VehicleTargetConfiguration.ProjectId,
			table.VehicleTargetConfiguration.Prefix,
			table.VehicleTargetConfiguration.Number,
			table.VehicleTargetConfiguration.Province,
			table.VehicleTargetConfiguration.Country,
			table.VehicleTargetConfiguration.PermittedLabel,
			table.VehicleTargetConfiguration.BlacklistPriority,
			table.VehicleTargetConfiguration.CreatedBy,
			table.VehicleTargetConfiguration.Brand,
			table.VehicleTargetConfiguration.BrandType,
			table.VehicleTargetConfiguration.Color,
			table.VehicleTargetConfiguration.ColorType,
			table.VehicleTargetConfiguration.Type,
			table.VehicleTargetConfiguration.TypeType,
		).
		MODEL(model.VehicleTargetConfiguration{
			ID:                uuid.New(),
			StationLocationId: data.StationLocationId,
			ProjectId:         data.ProjectId,
			Prefix:            data.Prefix,
			Number:            data.Number,
			Province:          data.Province,
			Country:           data.Country,
			Type:              data.Type,
			TypeType:          model.PropertyType_VehicleType,
			ColorType:         model.PropertyType_VehicleColor,
			Color:             data.Color,
			BrandType:         model.PropertyType_VehicleBrand,
			Brand:             data.Brand,
			PermittedLabel:    data.PermittedLabel,
			BlacklistPriority: data.BlacklistPriority,
			CreatedBy:         data.CreatedBy,
		}).
		RETURNING(table.VehicleTargetConfiguration.AllColumns)

	vehicleTarget := model.VehicleTargetConfiguration{}

	err = createVehicleTargetStmt.QueryContext(ctx, tx, &vehicleTarget)

	if err != nil {
		tx.Rollback()
		log.Println("create-vehicle-target-configuration-error", err.Error())
		return nil, utils.InternalServerError
	}

	if data.Tags != nil && len(*data.Tags) != 0 {
		for _, tag := range *data.Tags {
			upsertTagStmt := tagUtils.UpsertStatement(tagUtils.UpsertTagData{
				Tag:       tag,
				ProjectId: data.ProjectId.String(),
				CreatedBy: data.CreatedBy,
			})

			tagResult := model.Tag{}

			err := upsertTagStmt.QueryContext(ctx, tx, &tagResult)

			if err != nil {
				tx.Rollback()
				log.Println("upsert-vehicle-target-configuration-tag-error", err.Error())
				return nil, utils.InternalServerError
			}

			insertVehicleTargetStmt := table.VehicleTargetConfigurationTag.
				INSERT(
					table.VehicleTargetConfigurationTag.ID,
					table.VehicleTargetConfigurationTag.VehicleTargetConfigurationId,
					table.VehicleTargetConfigurationTag.TagId,
					table.VehicleTargetConfigurationTag.CreatedBy,
				).
				MODEL(model.VehicleTargetConfigurationTag{
					ID:                           uuid.New(),
					VehicleTargetConfigurationId: vehicleTarget.ID,
					TagId:                        tagResult.ID,
					CreatedBy:                    data.CreatedBy,
				})

			_, err = insertVehicleTargetStmt.ExecContext(ctx, tx)

			if err != nil {
				tx.Rollback()
				log.Println("insert-vehicle-target-configuration-tag-error", err.Error())
				return nil, utils.InternalServerError
			}
		}
	}

	tx.Commit()

	return transformToGraphql(vehicleTarget), nil
}

func transformToGraphql(item model.VehicleTargetConfiguration) *VehicleTargetConfiguration {
	return &VehicleTargetConfiguration{
		ID:                graphql.ID(item.ID.String()),
		StationLocationId: graphql.ID(item.StationLocationId.String()),
		ProjectId:         graphql.ID(item.ProjectId.String()),
		Color:             item.Color,
		Brand:             item.Brand,
		Prefix:            item.Prefix,
		Number:            item.Number,
		Province:          item.Province,
		Type:              item.Type,
		Country:           item.Country,
		PermittedLabel:    enum.GetDevicePermittedLabel(item.PermittedLabel.String()),
		BlacklistPriority: enum.GetBlacklistPriority(item.BlacklistPriority.String()),
	}
}
