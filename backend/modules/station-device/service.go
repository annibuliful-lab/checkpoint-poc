package stationdevice

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/utils"
	"context"
	"errors"
	"log"
	"time"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/graph-gophers/dataloader"
	"github.com/graph-gophers/graphql-go"
	"github.com/samber/lo"
)

type StationDeviceService struct{}

func (StationDeviceService) StationLocationDataloader() *dataloader.Loader {
	dbClient := db.GetPrimaryClient()

	return dataloader.NewBatchedLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var stationLocationIds []pg.Expression

		for _, id := range keys {
			stationLocationIds = append(stationLocationIds, pg.UUID(uuid.MustParse(id.String())))
		}

		var stationDevices = []model.StationDevice{}

		getStationDevicesStmt := table.Tag.
			SELECT(table.StationDevice.AllColumns).
			FROM(table.StationDevice).
			WHERE(table.StationDevice.StationLocationId.IN(stationLocationIds...))

		err := getStationDevicesStmt.Query(dbClient, &stationDevices)

		if err != nil {
			log.Println("dataloader-station-location-error", err.Error())
			return nil
		}

		var results []*dataloader.Result

		for _, key := range keys {
			filtered := lo.Filter(stationDevices, func(item model.StationDevice, index int) bool {
				return item.StationLocationId == uuid.MustParse(key.String())
			})

			results = append(results, &dataloader.Result{Data: filtered})
		}

		return results
	}, dataloader.WithClearCacheOnBatch())
}

func (StationDeviceService) FindMany(data GetStationDevicesData) ([]*StationDevice, error) {
	dbClient := db.GetPrimaryClient()

	conditions := pg.Bool(true).AND(table.StationDevice.DeletedAt.IS_NULL())

	if data.Search != nil {
		conditions = conditions.AND(table.StationDevice.Title.LIKE(pg.String(*data.Search)))
	}

	getStationDevicesStmt := table.StationDevice.
		SELECT(table.StationDevice.AllColumns).
		WHERE(conditions).
		LIMIT(data.Limit).
		OFFSET(data.Skip)

	stationDevices := []model.StationDevice{}

	err := getStationDevicesStmt.Query(dbClient, &stationDevices)

	if err != nil {
		log.Println("get-station-devices-error", err.Error())
		return nil, utils.InternalServerError
	}

	return lo.Map(stationDevices, func(item model.StationDevice, index int) *StationDevice {
		return transformToGraphql(item)
	}), nil
}

func (StationDeviceService) FindById(data GetStationDeviceData) (*StationDevice, error) {
	dbClient := db.GetPrimaryClient()
	getByIdStmt := table.StationDevice.
		SELECT(table.StationDevice.AllColumns).
		WHERE(
			table.StationDevice.ID.EQ(pg.UUID(data.ID)).
				AND(table.StationDevice.DeletedAt.IS_NULL()),
		)

	stationDevice := model.StationDevice{}
	err := getByIdStmt.Query(dbClient, &stationDevice)

	if err != nil && db.HasNoRow(err) {
		return nil, utils.Notfound
	}

	if err != nil {
		log.Println("get-by-station-device-error", err.Error())
		return nil, utils.InternalServerError
	}

	return transformToGraphql(stationDevice), nil
}

func (StationDeviceService) Update(data UpdateStationDeviceData) (*StationDevice, error) {
	dbClient := db.GetPrimaryClient()
	now := time.Now()

	var columnsToUpdate pg.ColumnList
	columnsToUpdate = append(columnsToUpdate,
		table.StationDevice.UpdatedAt,
		table.StationDevice.UpdatedBy,
	)

	fieldsToUpdate := model.StationDevice{
		UpdatedAt: &now,
		UpdatedBy: &data.UpdatedBy,
	}

	if data.Title != nil {
		columnsToUpdate = append(columnsToUpdate, table.StationDevice.Title)
		fieldsToUpdate.Title = *data.Title
	}

	if data.HardwareVersion != nil {
		columnsToUpdate = append(columnsToUpdate, table.StationDevice.HardwareVersion)
		fieldsToUpdate.HardwareVersion = data.HardwareVersion
	}

	if data.SoftwareVersion != nil {
		columnsToUpdate = append(columnsToUpdate, table.StationDevice.SoftwareVersion)
		fieldsToUpdate.SoftwareVersion = data.SoftwareVersion
	}

	updateStationDeviceStmt := table.StationDevice.
		UPDATE(columnsToUpdate).
		MODEL(fieldsToUpdate).
		WHERE(
			table.StationDevice.ID.EQ(pg.UUID(data.ID)).
				AND(table.StationDevice.DeletedAt.IS_NULL()),
		).
		RETURNING(table.StationDevice.AllColumns)

	stationDevice := model.StationDevice{}
	err := updateStationDeviceStmt.Query(dbClient, &stationDevice)

	if err != nil && db.HasNoRow(err) {
		return nil, utils.ForbiddenOperation
	}

	if err != nil {
		log.Println("update-station-device-error", err.Error())
		return nil, utils.InternalServerError
	}

	return transformToGraphql(stationDevice), nil
}

func (StationDeviceService) Delete(data DeleteStationDeviceData) (*utils.DeleteOperation, error) {
	dbClient := db.GetPrimaryClient()
	now := time.Now()

	softDeleteStmt := table.StationDevice.
		UPDATE(table.StationDevice.DeletedAt, table.StationDevice.DeletedBy).
		MODEL(model.StationDevice{
			DeletedAt: &now,
			DeletedBy: &data.DeletedBy,
		}).
		WHERE(table.StationDevice.ID.EQ(pg.UUID(data.ID)))

	result, err := softDeleteStmt.Exec(dbClient)

	if err != nil {
		log.Println("delete-station-device-error", err.Error())
		return nil, utils.InternalServerError
	}

	affectedRow, err := result.RowsAffected()

	if affectedRow == 0 {
		return nil, utils.ForbiddenOperation
	}

	if err != nil {
		log.Println("delete-station-device-error", err.Error())
		return nil, utils.InternalServerError
	}

	return &utils.DeleteOperation{
		Success: true,
	}, nil
}

func (StationDeviceService) Create(data CreateStationDeviceData) (*StationDevice, error) {
	dbClient := db.GetPrimaryClient()
	var columnsToInsert pg.ColumnList
	columnsToInsert = append(columnsToInsert,
		table.StationDevice.StationLocationId,
		table.StationDevice.ID,
		table.StationDevice.Title,
		table.StationDevice.CreatedBy,
		table.StationDevice.CreatedAt,
	)

	now := time.Now()
	fieldsToInsert := model.StationDevice{
		ID:                uuid.New(),
		StationLocationId: data.StationLocationId,
		Title:             data.Title,
		CreatedAt:         now,
	}

	if data.HardwareVersion != nil {
		columnsToInsert = append(columnsToInsert, table.StationDevice.HardwareVersion)
		fieldsToInsert.HardwareVersion = data.HardwareVersion
	}

	if data.SoftwareVersion != nil {
		columnsToInsert = append(columnsToInsert, table.StationDevice.SoftwareVersion)
		fieldsToInsert.SoftwareVersion = data.SoftwareVersion
	}

	createStationDeviceStmt := table.StationDevice.
		INSERT(columnsToInsert).
		MODEL(fieldsToInsert).
		RETURNING(table.StationDevice.AllColumns)

	stationDevice := model.StationDevice{}

	err := createStationDeviceStmt.Query(dbClient, &stationDevice)

	if err != nil && db.IsInvalidForeignKey(err) {
		return nil, errors.New("invalid stationLocationId")
	}

	if err != nil {
		log.Println("insert-station-device-error", err.Error())
		return nil, utils.InternalServerError
	}

	return transformToGraphql(stationDevice), nil
}

func transformToGraphql(data model.StationDevice) *StationDevice {
	return &StationDevice{
		ID:                graphql.ID(data.ID.String()),
		StationLocationId: graphql.ID(data.StationLocationId.String()),
		Title:             data.Title,
		HardwareVersion:   data.HardwareVersion,
		SoftwareVersion:   data.SoftwareVersion,
	}
}
