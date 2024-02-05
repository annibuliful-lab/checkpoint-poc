package imeiconfiguration

import (
	"checkpoint/.gen/checkpoint/public/model"
	. "checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/utils"
	"log"
	"strings"
	"time"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/samber/lo"

	"github.com/google/uuid"
)

func DeleteImeiConfiguration(data DeleteImeiConfigurationData) (int, error) {
	dbClient := db.GetPrimaryClient()
	now := time.Now()
	softDeleteStmt := ImeiConfiguration.
		UPDATE(ImeiConfiguration.DeletedBy, ImeiConfiguration.DeletedAt).
		MODEL(model.ImeiConfiguration{
			DeletedAt: &now,
			DeletedBy: &data.DeletedBy,
		}).
		WHERE(ImeiConfiguration.ID.EQ(pg.UUID(data.ID)).
			AND(ImeiConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId))))

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

func GetImeiConfigurations(data GetImeiConfigurationsData) ([]ImeiConfigurationResponse, int, error) {
	dbClient := db.GetPrimaryClient()
	conditions := pg.Bool(true).
		AND(ImeiConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId))).
		AND(ImeiConfiguration.DeletedAt.IS_NULL())

	if data.Label != "" {
		conditions = conditions.AND(ImeiConfiguration.Label.EQ(pg.NewEnumValue(data.Label)))
	}

	if data.Search != "" {
		conditions = conditions.AND(ImeiConfiguration.Imei.LIKE(pg.String(data.Search)))
	}

	if len(data.Tags) != 0 {
		conditions = conditions.AND(pg.RawBool("imei_configuration.tags @> array[string_to_array(#tags,'~^~')]", pg.RawArgs{"#tags": strings.Join(data.Tags, "~^~")}))
	}

	getImeisStmt := ImeiConfiguration.SELECT(ImeiConfiguration.AllColumns).
		WHERE(conditions).
		LIMIT(data.Pagination.Limit).
		OFFSET(data.Pagination.Skip)

	imeiConfigurations := []model.ImeiConfiguration{}

	err := getImeisStmt.Query(dbClient, &imeiConfigurations)

	if err != nil {
		log.Println("get-imei-configurations-error", err.Error())
		return nil, 500, utils.InternalServerError
	}

	imeiConfigurationsResponse := lo.Map(imeiConfigurations, func(item model.ImeiConfiguration, index int) ImeiConfigurationResponse {
		return ImeiConfigurationResponse{
			ID:                item.ID,
			ProjectId:         item.ProjectId,
			StationLocationId: item.StationLocationId,
			Imei:              item.Imei,
			Priority:          item.Priority.String(),
			Label:             item.Label.String(),
			Tags:              db.ConvertArrayDbStringToArrayString(item.Tags),
			CreatedBy:         item.CreatedBy,
			CreatedAt:         item.CreatedAt,
			UpdatedBy:         item.UpdatedBy,
			UpdatedAt:         item.UpdatedAt,
		}
	})

	return imeiConfigurationsResponse, 200, nil
}

func GetImeiConfiguration(data GetImeiConfigurationData) (*ImeiConfigurationResponse, int, error) {
	dbClient := db.GetPrimaryClient()

	getImeiStmt := ImeiConfiguration.
		SELECT(ImeiConfiguration.AllColumns).
		WHERE(ImeiConfiguration.ID.EQ(pg.UUID(data.ID)).
			AND(ImeiConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId))).
			AND(ImeiConfiguration.DeletedAt.IS_NULL()))

	imeiConfiguration := model.ImeiConfiguration{}

	err := getImeiStmt.Query(dbClient, &imeiConfiguration)

	if err != nil && db.HasNoRow(err) {
		return nil, 403, utils.ForbiddenOperation
	}

	if err != nil {
		log.Println("insert-imei-configuration", err.Error())
		return nil, 500, utils.InternalServerError
	}

	return &ImeiConfigurationResponse{
		ID:                imeiConfiguration.ID,
		ProjectId:         imeiConfiguration.ProjectId,
		StationLocationId: imeiConfiguration.StationLocationId,
		Imei:              imeiConfiguration.Imei,
		Label:             imeiConfiguration.Label.String(),
		Priority:          imeiConfiguration.Priority.String(),
		Tags:              db.ConvertArrayDbStringToArrayString(imeiConfiguration.Tags),
		CreatedBy:         imeiConfiguration.CreatedBy,
		CreatedAt:         imeiConfiguration.CreatedAt,
		UpdatedBy:         imeiConfiguration.UpdatedBy,
		UpdatedAt:         imeiConfiguration.UpdatedAt,
	}, 200, nil
}
func UpdateImeiConfiguration(data UpdateImeiConfigurationData) (*ImeiConfigurationResponse, int, error) {
	dbClient := db.GetPrimaryClient()

	now := time.Now()
	updateImeiStmt := ImeiConfiguration.
		UPDATE(ImeiConfiguration.Tags, ImeiConfiguration.Imei, ImeiConfiguration.UpdatedBy, ImeiConfiguration.Label, ImeiConfiguration.Priority, ImeiConfiguration.UpdatedAt).
		MODEL(model.ImeiConfiguration{
			Imei:      data.Imei,
			UpdatedBy: &data.UpdatedBy,
			UpdatedAt: &now,
			Label:     model.DevicePermittedLabel(data.Label),
			Priority:  model.BlacklistPriority(data.Priority),
			Tags:      db.ConvertArrayStringToInput(data.Tags),
		}).
		RETURNING(ImeiConfiguration.AllColumns).
		WHERE(ImeiConfiguration.ID.EQ(pg.UUID(data.ID)).
			AND(ImeiConfiguration.ProjectId.EQ(pg.UUID(data.ProjectId))).
			AND(ImeiConfiguration.DeletedAt.IS_NULL()))

	imeiConfiguration := model.ImeiConfiguration{}

	err := updateImeiStmt.Query(dbClient, &imeiConfiguration)

	if err != nil && db.HasNoRow(err) {
		return nil, 403, utils.ForbiddenOperation
	}

	if err != nil {
		log.Println("insert-imei-configuration", err.Error())
		return nil, 500, utils.InternalServerError
	}

	return &ImeiConfigurationResponse{
		ID:                imeiConfiguration.ID,
		ProjectId:         imeiConfiguration.ProjectId,
		StationLocationId: imeiConfiguration.StationLocationId,
		Imei:              imeiConfiguration.Imei,
		Label:             imeiConfiguration.Label.String(),
		Priority:          imeiConfiguration.Priority.String(),
		Tags:              db.ConvertArrayDbStringToArrayString(imeiConfiguration.Tags),
		CreatedBy:         imeiConfiguration.CreatedBy,
		CreatedAt:         imeiConfiguration.CreatedAt,
		UpdatedBy:         imeiConfiguration.UpdatedBy,
		UpdatedAt:         imeiConfiguration.UpdatedAt,
	}, 200, nil
}

func CreateImeiConfiguration(data CreateImeiConfigurationData) (*ImeiConfigurationResponse, int, error) {
	dbClient := db.GetPrimaryClient()

	insertImeiStmt := ImeiConfiguration.
		INSERT(ImeiConfiguration.Tags, ImeiConfiguration.Imei, ImeiConfiguration.ID, ImeiConfiguration.ProjectId, ImeiConfiguration.StationLocationId, ImeiConfiguration.CreatedBy, ImeiConfiguration.Label, ImeiConfiguration.Priority).
		MODEL(model.ImeiConfiguration{
			ID:                uuid.New(),
			ProjectId:         data.ProjectId,
			Imei:              data.Imei,
			CreatedBy:         data.CreatedBy,
			Label:             model.DevicePermittedLabel(data.Label),
			Priority:          model.BlacklistPriority(data.Priority),
			StationLocationId: data.StationLocationId,
			Tags:              db.ConvertArrayStringToInput(data.Tags),
		}).RETURNING(ImeiConfiguration.AllColumns)

	imeiConfiguration := model.ImeiConfiguration{}

	err := insertImeiStmt.Query(dbClient, &imeiConfiguration)

	if err != nil {
		log.Println("insert-imei-configuration", err.Error())
		return nil, 500, utils.InternalServerError
	}

	return &ImeiConfigurationResponse{
		ID:                imeiConfiguration.ID,
		ProjectId:         imeiConfiguration.ProjectId,
		StationLocationId: imeiConfiguration.StationLocationId,
		Imei:              imeiConfiguration.Imei,
		Label:             imeiConfiguration.Label.String(),
		Priority:          imeiConfiguration.Priority.String(),
		Tags:              db.ConvertArrayDbStringToArrayString(imeiConfiguration.Tags),
		CreatedBy:         imeiConfiguration.CreatedBy,
		CreatedAt:         imeiConfiguration.CreatedAt,
		UpdatedBy:         imeiConfiguration.UpdatedBy,
		UpdatedAt:         imeiConfiguration.UpdatedAt,
	}, 201, nil
}
