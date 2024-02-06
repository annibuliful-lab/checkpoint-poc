package imsiconfiguration

import (
	"checkpoint/.gen/checkpoint/public/model"
	table "checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/utils"
	"log"
	"strings"
	"time"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/samber/lo"
)

type ImsiConfigurationService struct{}

func (ImsiConfigurationService) FindMany(data GetImsiConfigurationsData) ([]Imsiconfiguration, int, error) {
	dbClient := db.GetPrimaryClient()

	conditions := pg.Bool(true).
		AND(table.ImsiConfiguration.DeletedAt.IS_NULL()).
		AND(table.ImsiConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId)))

	if data.Label != "" {
		conditions = conditions.AND(table.ImsiConfiguration.PermittedLabel.EQ(pg.NewEnumValue(data.Label)))
	}

	if len(data.Tags) != 0 {
		conditions = conditions.AND(pg.RawBool("imsi_configuration.tags @> array[string_to_array(#tags,'~^~')]", pg.RawArgs{"#tags": strings.Join(data.Tags, "~^~")}))
	}

	if data.Search != "" {
		conditions = conditions.AND(table.ImsiConfiguration.Imsi.LIKE(pg.String(data.Search)))
	}

	if data.Mcc != "" {
		conditions = conditions.AND(table.ImsiConfiguration.Mcc.EQ(pg.String(data.Mcc)))
	}

	if data.Mnc != "" {
		conditions = conditions.AND(table.ImsiConfiguration.Mnc.EQ(pg.String(data.Mnc)))
	}

	getImsiConfigurationsStmt := table.ImsiConfiguration.
		SELECT(table.ImsiConfiguration.AllColumns).
		FROM(table.ImsiConfiguration).
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
		return Imsiconfiguration{
			ID:                graphql.ID(item.ID.String()),
			ProjectId:         graphql.ID(item.ProjectId.String()),
			Imsi:              item.Imsi,
			CreatedBy:         item.CreatedBy,
			UpdatedBy:         item.UpdatedBy,
			CreatedAt:         item.CreatedAt,
			UpdatedAt:         item.UpdatedAt,
			Label:             item.PermittedLabel,
			Priority:          item.Priority,
			StationLocationId: item.StationLocationId,
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
		return nil, 403, utils.ForbiddenOperation
	}

	if err != nil {
		return nil, 500, utils.InternalServerError
	}

	return &Imsiconfiguration{
		ID:        graphql.ID(imsiConfiguration.ID.String()),
		ProjectId: graphql.ID(imsiConfiguration.ProjectId.String()),
		Imsi:      imsiConfiguration.Imsi,
		CreatedBy: imsiConfiguration.CreatedBy,
		CreatedAt: imsiConfiguration.CreatedAt,
		UpdatedBy: imsiConfiguration.UpdatedBy,
		UpdatedAt: imsiConfiguration.UpdatedAt,
		Label:     imsiConfiguration.PermittedLabel,
		Priority:  imsiConfiguration.Priority,
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
	mcc, mnc, err := utils.ExtractMCCMNC(data.Imsi)

	if err != nil {
		return nil, 400, err
	}
	now := time.Now()

	updateImsiStmt := table.ImsiConfiguration.
		UPDATE(table.ImsiConfiguration.Imsi, table.ImsiConfiguration.Priority, table.ImsiConfiguration.PermittedLabel, table.ImsiConfiguration.Mcc, table.ImsiConfiguration.Mnc, table.ImsiConfiguration.UpdatedAt, table.ImsiConfiguration.UpdatedBy).
		MODEL(model.ImsiConfiguration{
			Imsi:           data.Imsi,
			Priority:       model.BlacklistPriority(data.Priority),
			PermittedLabel: model.DevicePermittedLabel(data.Label),
			Mcc:            mcc,
			Mnc:            mnc,
			UpdatedBy:      &data.UpdatedBy,
			UpdatedAt:      &now,
		}).
		WHERE(table.ImsiConfiguration.ID.EQ(pg.UUID(data.ID)).
			AND(table.ImsiConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId))).
			AND(table.ImsiConfiguration.DeletedAt.IS_NULL())).
		RETURNING(table.ImsiConfiguration.AllColumns)

	imsiConfiguration := model.ImsiConfiguration{}
	err = updateImsiStmt.Query(dbClient, &imsiConfiguration)

	if err != nil && db.HasNoRow(err) {
		return nil, 403, utils.ForbiddenOperation
	}

	if err != nil && db.InvalidInput(err) {
		log.Println("invalid-update-imsi-configuraiton-error", err.Error())
		return nil, 400, err
	}

	if err != nil {
		log.Println("insert-imsi-configuraiton-error", err.Error())
		return nil, 500, utils.InternalServerError
	}

	return &Imsiconfiguration{
		ID:        graphql.ID(imsiConfiguration.ID.String()),
		ProjectId: graphql.ID(imsiConfiguration.ProjectId.String()),
		Imsi:      imsiConfiguration.Imsi,
		CreatedBy: imsiConfiguration.CreatedBy,
		CreatedAt: imsiConfiguration.CreatedAt,
		UpdatedBy: imsiConfiguration.UpdatedBy,
		UpdatedAt: imsiConfiguration.UpdatedAt,
		Label:     imsiConfiguration.PermittedLabel,
		Priority:  imsiConfiguration.Priority,
	}, 200, nil
}

func (ImsiConfigurationService) Create(data CreateImsiConfigurationData) (*Imsiconfiguration, int, error) {
	dbClient := db.GetPrimaryClient()
	mcc, mnc, err := utils.ExtractMCCMNC(data.Imsi)

	if err != nil {
		return nil, 400, err
	}

	insertImsiStmt := table.ImsiConfiguration.
		INSERT(table.ImsiConfiguration.ID, table.ImsiConfiguration.Imsi, table.ImsiConfiguration.Priority, table.ImsiConfiguration.StationLocationId, table.ImsiConfiguration.PermittedLabel, table.ImsiConfiguration.CreatedBy, table.ImsiConfiguration.ProjectId, table.ImsiConfiguration.Mcc, table.ImsiConfiguration.Mnc).
		MODEL(model.ImsiConfiguration{
			ID:                uuid.New(),
			Imsi:              data.Imsi,
			Priority:          model.BlacklistPriority(data.Priority),
			StationLocationId: data.StationLocationId,
			PermittedLabel:    model.DevicePermittedLabel(data.PermittedLabel),
			CreatedBy:         data.CreatedBy,
			ProjectId:         data.ProjectId,
			Mcc:               mcc,
			Mnc:               mnc,
		}).
		RETURNING(table.ImsiConfiguration.AllColumns)

	imsiConfiguration := model.ImsiConfiguration{}
	err = insertImsiStmt.Query(dbClient, &imsiConfiguration)

	if err != nil && db.InvalidInput(err) {
		log.Println("invalid-insert-imsi-configuraiton-error", err.Error())
		return nil, 400, err
	}

	if err != nil {
		log.Println("insert-imsi-configuraiton-error", err.Error())
		return nil, 500, utils.InternalServerError
	}

	return &Imsiconfiguration{
		ID:                graphql.ID(imsiConfiguration.ID.String()),
		ProjectId:         graphql.ID(imsiConfiguration.ProjectId.String()),
		StationLocationId: imsiConfiguration.StationLocationId,
		Mcc:               imsiConfiguration.Mcc,
		Mnc:               imsiConfiguration.Mnc,
		Imsi:              imsiConfiguration.Imsi,
		CreatedBy:         imsiConfiguration.CreatedBy,
		CreatedAt:         imsiConfiguration.CreatedAt,
		UpdatedBy:         imsiConfiguration.UpdatedBy,
		UpdatedAt:         imsiConfiguration.UpdatedAt,
		Label:             imsiConfiguration.PermittedLabel,
		Priority:          imsiConfiguration.Priority,
	}, 201, nil
}
