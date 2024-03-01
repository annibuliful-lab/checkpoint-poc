package stationofficer

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/utils"
	"context"
	"log"
	"time"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/graph-gophers/dataloader"
	"github.com/graph-gophers/graphql-go"
	"github.com/samber/lo"
)

type StationOfficerService struct{}

func (StationOfficerService) FindMany(data GetStationOfficersData) ([]*StationOfficer, error) {
	dbClient := db.GetPrimaryClient()

	conditions := table.StationOfficer.StationLocationId.EQ(pg.UUID(data.StationLocationId)).
		AND(table.StationOfficer.DeletedAt.IS_NULL())

	if data.Search != nil {
		conditions = conditions.AND(
			table.StationOfficer.Firstname.LIKE(pg.String(*data.Search)).
				OR(table.StationOfficer.Lastname.LIKE(pg.String(*data.Search)).
					OR(table.StationOfficer.Msisdn.LIKE(pg.String(*data.Search))),
				),
		)
	}

	getStationOfficersStmt := table.StationOfficer.
		SELECT(table.StationOfficer.AllColumns).
		WHERE(conditions).
		LIMIT(data.Limit).
		OFFSET(data.Skip)

	stationOfficers := []model.StationOfficer{}

	err := getStationOfficersStmt.Query(dbClient, &stationOfficers)

	if err != nil {
		log.Println("get-station-officers-error", err.Error())
		return nil, utils.InternalServerError
	}

	return lo.Map(stationOfficers, func(item model.StationOfficer, index int) *StationOfficer {
		return &StationOfficer{
			ID:                graphql.ID(item.ID.String()),
			StationLocationId: graphql.ID(item.StationLocationId.String()),
			Msisdn:            item.Msisdn,
			Firstname:         item.Firstname,
			Lastname:          item.Lastname,
		}
	}), nil
}

func (StationOfficerService) FindById(data GetStationOfficerData) (*StationOfficer, error) {
	dbClient := db.GetPrimaryClient()
	getByIdStmt := table.StationOfficer.
		SELECT(table.StationOfficer.AllColumns).
		WHERE(
			table.StationOfficer.ID.EQ(pg.UUID(data.ID)).
				AND(table.StationOfficer.DeletedAt.IS_NULL()),
		)

	stationOfficer := model.StationOfficer{}

	err := getByIdStmt.Query(dbClient, &stationOfficer)

	if err != nil && db.HasNoRow(err) {
		return nil, utils.ForbiddenOperation
	}

	if err != nil {
		log.Println("get-id-station-officer-error", err.Error())
		return nil, utils.InternalServerError
	}

	return transformToGraphql(stationOfficer), nil
}

func (StationOfficerService) Delete(data DeleteStationOfficerData) (*utils.DeleteOperation, error) {
	dbClient := db.GetPrimaryClient()
	now := time.Now()
	softDeleteStmt := table.StationOfficer.
		UPDATE(table.StationOfficer.DeletedAt, table.StationOfficer.DeletedBy).
		MODEL(model.StationOfficer{
			DeletedAt: &now,
			DeletedBy: &data.DeletedBy,
		}).
		WHERE(table.StationOfficer.ID.EQ(pg.UUID(data.ID)))

	result, err := softDeleteStmt.Exec(dbClient)

	if err != nil {
		log.Println("delete-station-officer-error", err.Error())
		return nil, utils.InternalServerError
	}

	affectedRow, err := result.RowsAffected()

	if affectedRow == 0 {
		return nil, utils.ForbiddenOperation
	}

	if err != nil {
		log.Println("delete-station-officer-error", err.Error())
		return nil, utils.InternalServerError
	}

	return &utils.DeleteOperation{
		Success: true,
	}, nil
}

func (StationOfficerService) Update(data UpdateStationOfficerData) (*StationOfficer, error) {
	dbClient := db.GetPrimaryClient()
	var columnsToUpdate pg.ColumnList

	columnsToUpdate = append(columnsToUpdate,
		table.StationOfficer.UpdatedAt,
		table.StationOfficer.UpdatedBy,
	)
	now := time.Now()

	dataFieldsToUpdate := model.StationOfficer{
		UpdatedAt: &now,
		UpdatedBy: &data.UpdatedBy,
	}

	if data.Firstname != nil {
		columnsToUpdate = append(columnsToUpdate, table.StationOfficer.Firstname)
		dataFieldsToUpdate.Firstname = *data.Firstname
	}

	if data.Lastname != nil {
		columnsToUpdate = append(columnsToUpdate, table.StationOfficer.Lastname)
		dataFieldsToUpdate.Firstname = *data.Lastname

	}

	if data.Msisdn != nil {
		columnsToUpdate = append(columnsToUpdate, table.StationOfficer.Msisdn)
		dataFieldsToUpdate.Firstname = *data.Lastname
	}

	updateStmt := table.StationOfficer.
		UPDATE(columnsToUpdate).
		MODEL(dataFieldsToUpdate).
		WHERE(
			table.StationOfficer.ID.EQ(pg.UUID(data.ID)).
				AND(table.StationOfficer.DeletedAt.IS_NULL()),
		).
		RETURNING(table.StationOfficer.AllColumns)

	stationOfficer := model.StationOfficer{}

	err := updateStmt.Query(dbClient, &stationOfficer)

	if err != nil && db.HasNoRow(err) {
		return nil, utils.ForbiddenOperation
	}

	if err != nil {
		log.Println("update-station-officer-error", err.Error())
		return nil, utils.InternalServerError
	}

	return transformToGraphql(stationOfficer), nil
}

func (StationOfficerService) Create(data CreateStationOfficerData) (*StationOfficer, error) {
	dbClient := db.GetPrimaryClient()

	insertStationOfficerStmt := table.StationOfficer.
		INSERT(table.StationOfficer.ID,
			table.StationOfficer.StationLocationId,
			table.StationOfficer.Firstname,
			table.StationOfficer.Lastname,
			table.StationOfficer.Msisdn,
			table.StationOfficer.CreatedBy,
		).
		MODEL(model.StationOfficer{
			ID:                uuid.New(),
			StationLocationId: data.StationLocationId,
			Firstname:         data.Firstname,
			Lastname:          data.Lastname,
			Msisdn:            data.Msisdn,
			CreatedBy:         data.CreatedBy,
		}).
		RETURNING(table.StationOfficer.AllColumns)

	stationOfficer := model.StationOfficer{}

	err := insertStationOfficerStmt.Query(dbClient, &stationOfficer)

	if err != nil {
		log.Println("insert-station-officer-error", err.Error())
		return nil, utils.InternalServerError
	}

	return transformToGraphql(stationOfficer), nil

}

func (StationOfficerService) StationLocationOfficerDataloader() *dataloader.Loader {

	return dataloader.NewBatchedLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		dbClient := db.GetPrimaryClient()
		var ids []pg.Expression

		for _, id := range keys {
			ids = append(ids, pg.UUID(uuid.MustParse(id.String())))
		}

		var stationOfficers = []model.StationOfficer{}

		getStationOfficersStmt := table.StationOfficer.
			SELECT(table.StationOfficer.AllColumns).
			FROM(table.StationOfficer).
			WHERE(table.StationOfficer.StationLocationId.IN(ids...))

		err := getStationOfficersStmt.Query(dbClient, &stationOfficers)

		if err != nil {
			log.Println("dataloader-station-officer-error", err.Error())
			return nil
		}

		var results []*dataloader.Result

		for _, key := range keys {
			filtered := lo.Filter(stationOfficers, func(item model.StationOfficer, index int) bool {
				return item.StationLocationId.String() == key.String()
			})

			results = append(results, &dataloader.Result{Data: filtered})
		}

		return results
	}, dataloader.WithClearCacheOnBatch())
}

func transformToGraphql(data model.StationOfficer) *StationOfficer {
	return &StationOfficer{
		ID:                graphql.ID(data.ID.String()),
		StationLocationId: graphql.ID(data.StationLocationId.String()),
		Firstname:         data.Firstname,
		Lastname:          data.Lastname,
		Msisdn:            data.Msisdn,
	}
}
