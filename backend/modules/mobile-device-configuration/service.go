package mobiledeviceconfiguration

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/gql/enum"
	imeiconfiguration "checkpoint/modules/imei-configuration"
	imsiconfiguration "checkpoint/modules/imsi-configuration"
	tagUtils "checkpoint/modules/tag"
	utils "checkpoint/utils"
	"checkpoint/utils/graphql_utils"
	"context"
	"log"
	"time"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/samber/lo"
)

type MobileDeviceConfigurationService struct{}

func (MobileDeviceConfigurationService) FindMany(data GetMobileDeviceConfigurationsData) ([]MobileDeviceConfiguration, error) {
	dbClient := db.GetPrimaryClient()
	conditions := pg.Bool(true).
		AND(table.MobileDeviceConfiguration.StationLocationId.EQ(pg.UUID(data.StationLocationId))).
		AND(table.MobileDeviceConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId))).
		AND(table.MobileDeviceConfiguration.DeletedAt.IS_NULL())

	var fromConditions pg.ReadableTable = table.MobileDeviceConfiguration

	if data.BlacklistPriority != nil {
		conditions = conditions.AND(table.MobileDeviceConfiguration.BlacklistPriority.EQ(pg.NewEnumValue(*data.BlacklistPriority)))
	}

	if data.PermittedLabel != nil {
		conditions = conditions.AND(table.MobileDeviceConfiguration.PermittedLabel.EQ(pg.NewEnumValue(*data.PermittedLabel)))
	}

	if data.Search != nil {
		conditions = conditions.AND(table.MobileDeviceConfiguration.Title.LIKE(pg.String(*data.Search)))
	}

	if data.Tags != nil && len(*data.Tags) != 0 {
		var tagItems []pg.Expression

		for _, tag := range *data.Tags {
			tagItems = append(tagItems, pg.String(tag))
		}

		fromConditions = fromConditions.
			INNER_JOIN(table.MobileDeviceConfigurationTag, table.MobileDeviceConfigurationTag.MobileDeviceConfigurationId.EQ(table.MobileDeviceConfigurationTag.ID)).
			INNER_JOIN(table.Tag, table.MobileDeviceConfigurationTag.TagId.EQ(table.Tag.ID))
		conditions = conditions.AND(table.Tag.Title.IN(tagItems...))

	}

	getMobileDevicesStmt := pg.
		SELECT(table.MobileDeviceConfiguration.AllColumns).
		FROM(fromConditions).
		WHERE(conditions).
		LIMIT(data.pagination.Limit).
		OFFSET(data.pagination.Skip)

	mobileDeviceConfigurations := []model.MobileDeviceConfiguration{}

	err := getMobileDevicesStmt.Query(dbClient, &mobileDeviceConfigurations)

	if err != nil && db.HasNoRow(err) {
		return []MobileDeviceConfiguration{}, nil
	}

	if err != nil {
		log.Println("get-mobile-devices-error", err.Error())
		return nil, utils.InternalServerError
	}

	result := lo.Map(mobileDeviceConfigurations, func(mobileDeviceConfiguration model.MobileDeviceConfiguration, index int) MobileDeviceConfiguration {

		return *transformToGraphql(mobileDeviceConfiguration)
	})

	return result, nil
}

func (MobileDeviceConfigurationService) FindById(data GetMobileDeviceConfigurationData) (*MobileDeviceConfiguration, error) {
	dbClient := db.GetPrimaryClient()
	getByIdStmt := pg.
		SELECT(table.MobileDeviceConfiguration.AllColumns).
		WHERE(table.MobileDeviceConfiguration.ID.EQ(pg.UUID(data.ID)).
			AND(table.MobileDeviceConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId))).
			AND(table.MobileDeviceConfiguration.DeletedAt.IS_NULL()),
		).LIMIT(1)

	mobileDeviceConfiguration := model.MobileDeviceConfiguration{}

	err := getByIdStmt.Query(dbClient, &mobileDeviceConfiguration)

	if err != nil && db.HasNoRow(err) {
		return nil, utils.Notfound
	}

	if err != nil {
		return nil, utils.InternalServerError
	}

	return transformToGraphql(mobileDeviceConfiguration), nil
}

func (MobileDeviceConfigurationService) Update(data UpdateMobileDeviceConfigurationData) (*MobileDeviceConfiguration, error) {
	dbClient := db.GetPrimaryClient()
	var columnsToUpdate pg.ColumnList

	columnsToUpdate = append(columnsToUpdate,
		table.MobileDeviceConfiguration.UpdatedAt,
		table.MobileDeviceConfiguration.UpdatedBy,
	)

	now := time.Now()

	var dataFieldsToUpdate = model.MobileDeviceConfiguration{
		UpdatedAt: &now,
		UpdatedBy: &data.UpdatedBy,
	}

	if data.Msisdn != nil {
		dataFieldsToUpdate.Msisdn = data.Msisdn
		columnsToUpdate = append(columnsToUpdate, table.MobileDeviceConfiguration.Msisdn)
	}

	if data.PermittedLabel != nil {
		dataFieldsToUpdate.PermittedLabel = model.DevicePermittedLabel(*data.PermittedLabel)
		columnsToUpdate = append(columnsToUpdate, table.MobileDeviceConfiguration.PermittedLabel)
	}

	if data.BlacklistPriority != nil {
		dataFieldsToUpdate.BlacklistPriority = model.BlacklistPriority(*data.BlacklistPriority)
		columnsToUpdate = append(columnsToUpdate, table.MobileDeviceConfiguration.BlacklistPriority)
	}

	if data.Imsi != nil {
		if !utils.ValidateIMSI(*data.Imsi) {
			return nil, utils.ErrInvalidIMSI
		}

		imsiConfiguration, err := imsiService.Upsert(imsiconfiguration.UpsertImsiConfigurationData{
			UpdatedBy:         data.UpdatedBy,
			ProjectId:         data.ProjectId,
			StationLocationId: data.StationId,
			Imsi:              *data.Imsi,
		})

		if err != nil {
			return nil, err
		}

		dataFieldsToUpdate.ReferenceImeiConfigurationId = uuid.MustParse(string(imsiConfiguration.ID))
		columnsToUpdate = append(columnsToUpdate, table.MobileDeviceConfiguration.ReferenceImeiConfigurationId)
	}

	if data.Imei != nil {
		if !utils.ValidateIMEI(*data.Imei) {
			return nil, utils.ErrInvalidIMEI
		}

		imsiConfiguration, err := imeiService.Upsert(imeiconfiguration.UpsertImeiConfigurationData{
			UpdatedBy:         data.UpdatedBy,
			ProjectId:         data.ProjectId,
			StationLocationId: data.StationId,
			Imei:              *data.Imei,
		})

		if err != nil {
			return nil, err
		}

		dataFieldsToUpdate.ReferenceImsiConfigurationId = uuid.MustParse(string(imsiConfiguration.ID))
		columnsToUpdate = append(columnsToUpdate, table.MobileDeviceConfiguration.ReferenceImsiConfigurationId)
	}

	ctx := context.Background()
	tx, err := dbClient.Begin()

	if err != nil {
		log.Println("init-insert-mobile-configuration-error", err.Error())
		return nil, utils.InternalServerError
	}

	mobileDeviceConfiguration := model.MobileDeviceConfiguration{}
	updateMobileDeviceConfigurationStmt := table.MobileDeviceConfiguration.
		UPDATE(columnsToUpdate).
		MODEL(dataFieldsToUpdate).
		WHERE(table.MobileDeviceConfiguration.ID.EQ(pg.UUID(data.ID)).
			AND(table.MobileDeviceConfiguration.DeletedAt.IS_NULL()).
			AND(table.MobileDeviceConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId)))).
		RETURNING(table.MobileDeviceConfiguration.AllColumns)

	err = updateMobileDeviceConfigurationStmt.QueryContext(ctx, tx, &mobileDeviceConfiguration)

	if err != nil {
		log.Println("insert-mobile-configuration-error", err.Error())
		return nil, utils.InternalServerError
	}

	if data.Tags != nil && len(*data.Tags) != 0 {
		deleteAllTagsStmt := table.MobileDeviceConfigurationTag.
			DELETE().
			WHERE(table.MobileDeviceConfigurationTag.MobileDeviceConfigurationId.
				EQ(pg.UUID(mobileDeviceConfiguration.ID)),
			)

		_, err = deleteAllTagsStmt.ExecContext(ctx, tx)

		if err != nil {
			tx.Rollback()
			log.Println("delete-all-mobile-configuration-tag-error", err.Error())
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
				log.Println("upsert-mobile-device-configuration-tag-error", err.Error())
				return nil, utils.InternalServerError
			}

			insertImsiTagStmt := table.MobileDeviceConfigurationTag.
				INSERT(table.MobileDeviceConfiguration.ID, table.MobileDeviceConfigurationTag.MobileDeviceConfigurationId, table.MobileDeviceConfigurationTag.TagId, table.MobileDeviceConfigurationTag.CreatedBy).
				MODEL(model.MobileDeviceConfigurationTag{
					ID:                          uuid.New(),
					MobileDeviceConfigurationId: mobileDeviceConfiguration.ID,
					TagId:                       tagResult.ID,
					CreatedBy:                   data.UpdatedBy,
				})
			_, err = insertImsiTagStmt.ExecContext(ctx, tx)

			if err != nil {
				tx.Rollback()
				log.Println("insert-mobile-device-configuration-tag-error", err.Error())
				return nil, utils.InternalServerError
			}

		}
	}

	tx.Commit()

	return transformToGraphql(mobileDeviceConfiguration), nil
}

func (MobileDeviceConfigurationService) Delete(data DeleteMobileDeviceConfigurationData) (*utils.DeleteOperation, error) {
	dbClient := db.GetPrimaryClient()
	now := time.Now()
	softDeleteMobileDeviceStmt := table.MobileDeviceConfiguration.
		UPDATE(table.MobileDeviceConfiguration.DeletedBy, table.MobileDeviceConfiguration.DeletedAt).
		MODEL(model.MobileDeviceConfiguration{
			DeletedAt: &now,
			DeletedBy: &data.DeletedBy,
		}).
		WHERE(table.MobileDeviceConfiguration.ID.EQ(pg.UUID(data.ID)).
			AND(table.MobileDeviceConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId))).
			AND(table.MobileDeviceConfiguration.DeletedAt.IS_NULL()))

	mobileDeviceRowsAffected, err := softDeleteMobileDeviceStmt.Exec(dbClient)

	if err != nil {
		return nil, utils.InternalServerError
	}

	rowAffected, err := mobileDeviceRowsAffected.RowsAffected()
	if err != nil {
		return nil, utils.InternalServerError
	}

	if rowAffected == 0 {
		return nil, utils.ForbiddenOperation
	}

	return &utils.DeleteOperation{
		Success: true,
	}, nil

}

func (MobileDeviceConfigurationService) Create(data CreateMobileDeviceConfigurationData) (*MobileDeviceConfiguration, int, error) {
	dbClient := db.GetPrimaryClient()

	tx, err := dbClient.Begin()
	if err != nil {
		log.Println("insert-mobile-device-configuration-error", err.Error())
		return nil, 500, utils.InternalServerError
	}

	ctx := context.Background()
	var columnsToInsert pg.ColumnList
	columnsToInsert = append(columnsToInsert,
		table.MobileDeviceConfiguration.Title,
		table.MobileDeviceConfiguration.ID,
		table.MobileDeviceConfiguration.BlacklistPriority,
		table.MobileDeviceConfiguration.PermittedLabel,
		table.MobileDeviceConfiguration.ProjectId,
		table.MobileDeviceConfiguration.CreatedBy,
		table.MobileDeviceConfiguration.ReferenceImeiConfigurationId,
		table.MobileDeviceConfiguration.ReferenceImsiConfigurationId,
		table.MobileDeviceConfiguration.StationLocationId,
	)

	if data.Msisdn != nil {
		columnsToInsert = append(columnsToInsert,
			table.MobileDeviceConfiguration.Msisdn,
		)
	}

	imsiConfiguration, err := imsiService.Upsert(imsiconfiguration.UpsertImsiConfigurationData{
		UpdatedBy:         data.CreatedBy,
		ProjectId:         data.ProjectId,
		StationLocationId: data.StationLocationId,
		Imsi:              data.Imsi,
		BlacklistPriority: data.BlacklistPriority,
		PermittedLabel:    data.PermittedLabel,
		Tags:              data.Tags,
	})

	if err != nil {
		return nil, 500, err
	}

	imeiConfiguration, err := imeiService.Upsert(imeiconfiguration.UpsertImeiConfigurationData{
		Imei:              data.Imei,
		StationLocationId: data.StationLocationId,
		UpdatedBy:         data.CreatedBy,
		ProjectId:         data.ProjectId,
		BlacklistPriority: data.BlacklistPriority,
		PermittedLabel:    data.PermittedLabel,
		Tags:              data.Tags,
	})

	if err != nil {
		return nil, 500, err
	}

	mobileDeviceConfiguration := model.MobileDeviceConfiguration{}
	insertMobileDeviceConfigurationStmt := table.MobileDeviceConfiguration.
		INSERT(columnsToInsert).
		MODEL(model.MobileDeviceConfiguration{
			ID:                           uuid.New(),
			ProjectId:                    data.ProjectId,
			CreatedBy:                    data.CreatedBy,
			Msisdn:                       data.Msisdn,
			BlacklistPriority:            data.BlacklistPriority,
			PermittedLabel:               data.PermittedLabel,
			Title:                        data.Title,
			ReferenceImeiConfigurationId: uuid.MustParse(string(imeiConfiguration.ID)),
			ReferenceImsiConfigurationId: uuid.MustParse(string(imsiConfiguration.ID)),
			StationLocationId:            data.StationLocationId,
		}).
		RETURNING(table.MobileDeviceConfiguration.AllColumns)

	err = insertMobileDeviceConfigurationStmt.QueryContext(ctx, tx, &mobileDeviceConfiguration)

	if err != nil {
		log.Println("insert-mobile-device-error", err.Error())
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
				log.Println("upsert-mobile-device-configuration-tag-error", err.Error())
				return nil, 500, utils.InternalServerError
			}

			insertImsiTagStmt := table.MobileDeviceConfigurationTag.
				INSERT(table.MobileDeviceConfiguration.ID, table.MobileDeviceConfigurationTag.MobileDeviceConfigurationId, table.MobileDeviceConfigurationTag.TagId, table.MobileDeviceConfigurationTag.CreatedBy).
				MODEL(model.MobileDeviceConfigurationTag{
					ID:                          uuid.New(),
					MobileDeviceConfigurationId: mobileDeviceConfiguration.ID,
					TagId:                       tagResult.ID,
					CreatedBy:                   data.CreatedBy,
				})
			_, err = insertImsiTagStmt.ExecContext(ctx, tx)

			if err != nil {
				tx.Rollback()
				log.Println("insert-mobile-device-configuration-tag-error", err.Error())
				return nil, 500, utils.InternalServerError
			}

		}
	}

	tx.Commit()

	return transformToGraphql(mobileDeviceConfiguration), 200, nil
}

func transformToGraphql(mobileDeviceConfiguration model.MobileDeviceConfiguration) *MobileDeviceConfiguration {
	var updatedBy graphql.NullID
	if mobileDeviceConfiguration.UpdatedBy != nil {

		updatedBy = graphql_utils.ConvertStringToNullID(mobileDeviceConfiguration.UpdatedBy)
	}

	var updatedAt graphql.NullTime
	if mobileDeviceConfiguration.UpdatedAt != nil {
		updatedAt = graphql.NullTime{Value: &graphql.Time{Time: *mobileDeviceConfiguration.UpdatedAt}}
	}

	return &MobileDeviceConfiguration{
		ID:                           graphql.ID(mobileDeviceConfiguration.ID.String()),
		ProjectId:                    graphql.ID(mobileDeviceConfiguration.ProjectId.String()),
		ReferenceImsiConfigurationId: graphql.ID(mobileDeviceConfiguration.ReferenceImsiConfigurationId.String()),
		ReferenceImeiConfigurationId: graphql.ID(mobileDeviceConfiguration.ReferenceImeiConfigurationId.String()),
		Title:                        mobileDeviceConfiguration.Title,
		Msisdn:                       mobileDeviceConfiguration.Msisdn,
		PermittedLabel:               enum.GetDevicePermittedLabel(mobileDeviceConfiguration.PermittedLabel.String()),
		BlacklistPriority:            enum.GetBlacklistPriority(mobileDeviceConfiguration.BlacklistPriority.String()),
		CreatedBy:                    graphql.ID(mobileDeviceConfiguration.CreatedBy),
		CreatedAt:                    graphql.Time{Time: mobileDeviceConfiguration.CreatedAt},
		UpdatedBy:                    &updatedBy,
		UpdatedAt:                    &updatedAt,
	}
}
