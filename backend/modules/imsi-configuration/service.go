package imsiconfiguration

import (
	"checkpoint/.gen/checkpoint/public/model"
	. "checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/utils"
	"fmt"
	"log"
	"strings"
	"time"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func GetImsiConfigurations(data GetImsiConfigurationsData) ([]ImsiConfigurationResponse, int, error) {
	dbClient := db.GetPrimaryClient()

	conditions := pg.Bool(true).
		AND(ImsiConfiguration.DeletedAt.IS_NULL()).
		AND(ImsiConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId)))

	if data.Label != "" {
		conditions = conditions.AND(ImsiConfiguration.Label.EQ(pg.NewEnumValue(data.Label)))
	}

	if len(data.Tags) != 0 {
		conditions = conditions.AND(pg.RawBool("imsi_configuration.tags @> array[string_to_array(#tags,'~^~')]", pg.RawArgs{"#tags": strings.Join(data.Tags, "~^~")}))
	}

	if data.Search != "" {
		conditions = conditions.AND(ImsiConfiguration.Imsi.LIKE(pg.String(data.Search)))
	}

	if data.Mcc != "" {
		conditions = conditions.AND(ImsiConfiguration.Mcc.EQ(pg.String(data.Mcc)))
	}

	if data.Mnc != "" {
		conditions = conditions.AND(ImsiConfiguration.Mnc.EQ(pg.String(data.Mnc)))
	}

	getImsiConfigurationsStmt := ImsiConfiguration.
		SELECT(ImsiConfiguration.AllColumns).
		FROM(ImsiConfiguration).
		WHERE(conditions).
		LIMIT(data.Pagination.Limit).
		OFFSET(data.Pagination.Skip)

	imsiConfigurations := []model.ImsiConfiguration{}
	fmt.Println(getImsiConfigurationsStmt.DebugSql())
	err := getImsiConfigurationsStmt.Query(dbClient, &imsiConfigurations)

	if err != nil {
		log.Println("get-imsi-configurations-error", err.Error())
		return nil, 500, utils.InternalServerError
	}

	imsiConfigurationsResponse := lo.Map(imsiConfigurations, func(item model.ImsiConfiguration, index int) ImsiConfigurationResponse {
		return ImsiConfigurationResponse{
			ID:                item.ID,
			ProjectId:         item.ProjectId,
			Imsi:              item.Imsi,
			CreatedBy:         item.CreatedBy,
			UpdatedBy:         item.UpdatedBy,
			CreatedAt:         item.CreatedAt,
			UpdatedAt:         item.UpdatedAt,
			Label:             item.Label,
			Priority:          item.Priority,
			StationLocationId: item.StationLocationId,
			Mcc:               item.Mcc,
			Mnc:               item.Mnc,
			Tags:              db.ConvertArrayDbStringToArrayString(item.Tags),
		}
	})

	return imsiConfigurationsResponse, 200, nil
}

func GetImsiConfigurationById(data GetImsiConfigurationByIdData) (*ImsiConfigurationResponse, int, error) {
	dbClient := db.GetPrimaryClient()
	getImsiStmt := ImsiConfiguration.
		SELECT(ImsiConfiguration.AllColumns).
		FROM(ImsiConfiguration).
		WHERE(ImsiConfiguration.ID.EQ(pg.UUID(data.ID)).
			AND(ImsiConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId))).
			AND(ImsiConfiguration.DeletedAt.IS_NULL()))

	imsiConfiguration := model.ImsiConfiguration{}

	err := getImsiStmt.Query(dbClient, &imsiConfiguration)
	if err != nil && db.HasNoRow(err) {
		return nil, 403, utils.ForbiddenOperation
	}

	if err != nil {
		return nil, 500, utils.InternalServerError
	}

	return &ImsiConfigurationResponse{
		ID:        imsiConfiguration.ID,
		ProjectId: imsiConfiguration.ProjectId,
		Imsi:      imsiConfiguration.Imsi,
		CreatedBy: imsiConfiguration.CreatedBy,
		CreatedAt: imsiConfiguration.CreatedAt,
		UpdatedBy: imsiConfiguration.UpdatedBy,
		UpdatedAt: imsiConfiguration.UpdatedAt,
		Label:     imsiConfiguration.Label,
		Priority:  imsiConfiguration.Priority,
		Tags:      db.ConvertArrayDbStringToArrayString(imsiConfiguration.Tags),
	}, 200, nil
}

func DeleteImsiConfigurationById(data DeleteImsiConfigurationData) (int, error) {
	dbClient := db.GetPrimaryClient()
	now := time.Now()
	deleteImsiStmt := ImsiConfiguration.
		UPDATE(ImsiConfiguration.DeletedAt, ImsiConfiguration.DeletedBy).
		MODEL(model.ImsiConfiguration{
			DeletedAt: &now,
			DeletedBy: &data.DeletedBy,
		}).
		WHERE(ImsiConfiguration.ID.EQ(pg.UUID(data.ID)).
			AND(ImsiConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId))))

	_, err := deleteImsiStmt.Exec(dbClient)
	if err != nil && db.HasNoRow(err) {
		return 403, utils.ForbiddenOperation
	}

	return 200, nil
}

func UpdateImsiConfiguration(data UpdateImsiConfigurationData) (*ImsiConfigurationResponse, int, error) {
	dbClient := db.GetPrimaryClient()
	mcc, mnc, err := utils.ExtractMCCMNC(data.Imsi)

	if err != nil {
		return nil, 400, err
	}
	now := time.Now()

	updateImsiStmt := ImsiConfiguration.
		UPDATE(ImsiConfiguration.Imsi, ImsiConfiguration.Priority, ImsiConfiguration.Label, ImsiConfiguration.Mcc, ImsiConfiguration.Mnc, ImsiConfiguration.Tags, ImsiConfiguration.UpdatedAt, ImsiConfiguration.UpdatedBy).
		MODEL(model.ImsiConfiguration{
			Imsi:      data.Imsi,
			Priority:  model.BlacklistPriority(data.Priority),
			Label:     model.DevicePermittedLabel(data.Label),
			Mcc:       mcc,
			Mnc:       mnc,
			UpdatedBy: &data.UpdatedBy,
			UpdatedAt: &now,
			Tags:      db.ConvertArrayStringToInput(data.Tags),
		}).
		WHERE(ImsiConfiguration.ID.EQ(pg.UUID(data.ID)).
			AND(ImsiConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId)))).
		RETURNING(ImsiConfiguration.AllColumns)

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

	return &ImsiConfigurationResponse{
		ID:        imsiConfiguration.ID,
		ProjectId: imsiConfiguration.ProjectId,
		Imsi:      imsiConfiguration.Imsi,
		CreatedBy: imsiConfiguration.CreatedBy,
		CreatedAt: imsiConfiguration.CreatedAt,
		UpdatedBy: imsiConfiguration.UpdatedBy,
		UpdatedAt: imsiConfiguration.UpdatedAt,
		Label:     imsiConfiguration.Label,
		Priority:  imsiConfiguration.Priority,
		Tags:      db.ConvertArrayDbStringToArrayString(imsiConfiguration.Tags),
	}, 200, nil
}

func CreateImsiConfiguration(data CreateImsiConfigurationData) (*ImsiConfigurationResponse, int, error) {
	dbClient := db.GetPrimaryClient()
	mcc, mnc, err := utils.ExtractMCCMNC(data.Imsi)

	if err != nil {
		return nil, 400, err
	}

	insertImsiStmt := ImsiConfiguration.
		INSERT(ImsiConfiguration.ID, ImsiConfiguration.Imsi, ImsiConfiguration.Priority, ImsiConfiguration.StationLocationId, ImsiConfiguration.Label, ImsiConfiguration.CreatedBy, ImsiConfiguration.ProjectId, ImsiConfiguration.Mcc, ImsiConfiguration.Mnc, ImsiConfiguration.Tags).
		MODEL(model.ImsiConfiguration{
			ID:                uuid.New(),
			Imsi:              data.Imsi,
			Priority:          model.BlacklistPriority(data.Priority),
			StationLocationId: data.StationLocationId,
			Label:             model.DevicePermittedLabel(data.Label),
			CreatedBy:         data.CreatedBy,
			ProjectId:         data.ProjectId,
			Mcc:               mcc,
			Mnc:               mnc,
			Tags:              db.ConvertArrayStringToInput(data.Tags),
		}).
		RETURNING(ImsiConfiguration.AllColumns)

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

	return &ImsiConfigurationResponse{
		ID:                imsiConfiguration.ID,
		ProjectId:         imsiConfiguration.ProjectId,
		StationLocationId: imsiConfiguration.StationLocationId,
		Mcc:               imsiConfiguration.Mcc,
		Mnc:               imsiConfiguration.Mnc,
		Imsi:              imsiConfiguration.Imsi,
		CreatedBy:         imsiConfiguration.CreatedBy,
		CreatedAt:         imsiConfiguration.CreatedAt,
		UpdatedBy:         imsiConfiguration.UpdatedBy,
		UpdatedAt:         imsiConfiguration.UpdatedAt,
		Label:             imsiConfiguration.Label,
		Priority:          imsiConfiguration.Priority,
		Tags:              db.ConvertArrayDbStringToArrayString(imsiConfiguration.Tags),
	}, 201, nil
}
