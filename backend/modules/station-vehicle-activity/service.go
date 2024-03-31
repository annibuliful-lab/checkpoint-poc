package stationvehicleactivity

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/gql/enum"
	vehicleproperty "checkpoint/modules/vehicle-property"
	"checkpoint/utils"
	"context"
	"database/sql"
	"errors"
	"log"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
)

type StationVehicleActivityService struct{}

func (StationVehicleActivityService) Create(data CreateStationVehicleActivityData) (*StationVehicleActivity, error) {
	dbClient := db.GetPrimaryClient()
	ctx := context.Background()
	tx, err := dbClient.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
	})

	if err != nil {
		log.Println("init-tx-create-vehicle-activity-error", err.Error())
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
		log.Println("tx-upsert-vehicle-activity-brand-error", err.Error())
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
		log.Println("tx-upsert-vehicle-activity-color-error", err.Error())
		return nil, utils.InternalServerError
	}

	upsertModelStmt := vehicleproperty.UpsertStatement(vehicleproperty.UpsertVehiclePropertyData{
		Property:  data.Model,
		Type:      model.PropertyType_VehicleModel,
		ProjectId: data.ProjectId,
	})

	_, err = upsertModelStmt.ExecContext(ctx, tx)
	if err != nil {
		tx.Rollback()
		log.Println("tx-upsert-vehicle-activity-model-error", err.Error())
		return nil, utils.InternalServerError
	}

	stationVehicleActivity := model.StationVehicleActivity{}

	var columnsToInsert pg.ColumnList
	columnsToInsert = append(columnsToInsert,
		table.StationVehicleActivity.ID,
		table.StationVehicleActivity.ProjectId,
		table.StationVehicleActivity.StationLocationId,
		table.StationVehicleActivity.Brand,
		table.StationVehicleActivity.BrandType,
		table.StationVehicleActivity.Color,
		table.StationVehicleActivity.ColorType,
		table.StationVehicleActivity.Model,
		table.StationVehicleActivity.ModelType,
		table.StationVehicleActivity.CreatedBy,
	)

	fieldsToInsert := model.StationVehicleActivity{
		ID:                uuid.New(),
		ProjectId:         data.ProjectId,
		StationLocationId: data.StationLocationId,
		Brand:             data.Brand,
		BrandType:         model.PropertyType_VehicleBrand,
		Color:             data.Color,
		ColorType:         model.PropertyType_VehicleColor,
		Model:             data.Model,
		ModelType:         model.PropertyType_VehicleModel,
		CreatedBy:         data.CreatedBy,
	}

	if data.Status != nil {
		columnsToInsert = append(columnsToInsert, table.StationVehicleActivity.Status)
		fieldsToInsert.Status = model.RemarkState(*data.Status)
	}

	createStationVehicleStmt := table.StationVehicleActivity.
		INSERT(columnsToInsert).
		MODEL(fieldsToInsert).
		RETURNING(table.StationVehicleActivity.AllColumns)

	err = createStationVehicleStmt.QueryContext(ctx, tx, &stationVehicleActivity)

	if err != nil && db.IsInvalidForeignKey(err) {
		tx.Rollback()
		log.Println("create-station-vehicle-invalid-foreign-key", err.Error())
		return nil, errors.New("projectId or stationLocationId are invalid")
	}

	if err != nil {
		tx.Rollback()
		log.Println("create-station-vehicle-error", err.Error())
		return nil, utils.InternalServerError
	}

	tx.Commit()

	return transformToGraphql(stationVehicleActivity), nil
}

func transformToGraphql(data model.StationVehicleActivity) *StationVehicleActivity {
	return &StationVehicleActivity{
		ID:                graphql.ID(data.ID.String()),
		StationLocationId: graphql.ID(data.StationLocationId.String()),
		ProjectId:         graphql.ID(data.ProjectId.String()),
		Brand:             data.Brand,
		Color:             data.Color,
		Model:             data.Model,
		Status:            enum.GetRemarkState(data.Status.String()),
	}
}
