package stationlocation

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	tagUtils "checkpoint/modules/tag"
	"checkpoint/utils"
	"context"
	"log"
	"time"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/samber/lo"
)

type StationLocationService struct{}

func (StationLocationService) FindMany(data GetStationLocationsData) ([]*StationLocation, error) {
	dbClient := db.GetPrimaryClient()
	conditions := table.StationLocation.DeletedAt.IS_NULL().
		AND(table.StationLocation.ProjectId.EQ(pg.UUID(data.ProjectId)))

	var fromConditions pg.ReadableTable = table.StationLocation

	if data.Search != nil {
		conditions = conditions.
			AND(table.StationLocation.Title.LIKE(pg.String(*data.Search)))
	}

	if data.Tags != nil && len(*data.Tags) != 0 {
		var tagItems []pg.Expression

		for _, tag := range *data.Tags {
			tagItems = append(tagItems, pg.String(tag))
		}

		fromConditions = fromConditions.
			INNER_JOIN(table.StationLocationTag, table.StationLocationTag.StationLocationId.EQ(table.StationLocation.ID)).
			INNER_JOIN(table.Tag, table.StationLocationTag.TagId.EQ(table.Tag.ID))
		conditions = conditions.AND(table.Tag.Title.IN(tagItems...))
	}

	getStationLocationsStmt := pg.
		SELECT(table.StationLocation.AllColumns).
		FROM(fromConditions).
		WHERE(conditions).
		LIMIT(data.Limit).
		OFFSET(data.Skip)

	stationLocations := []model.StationLocation{}
	err := getStationLocationsStmt.Query(dbClient, &stationLocations)

	if err != nil {
		log.Println("get-stations-location-error", err.Error())
		return nil, utils.InternalServerError
	}

	return lo.Map(stationLocations, func(item model.StationLocation, index int) *StationLocation {
		return transformToGraphql(item)
	}), nil

}

func (StationLocationService) FindById(data GetStationLocationByIdData) (*StationLocation, error) {
	dbClient := db.GetPrimaryClient()
	getStationLocationStmt := pg.
		SELECT(table.StationLocation.AllColumns).
		FROM(table.StationLocation).
		WHERE(table.StationLocation.ID.EQ(pg.UUID(data.Id)).
			AND(table.StationLocation.ProjectId.EQ(pg.UUID(data.ProjectId))).
			AND(table.StationLocation.DeletedAt.IS_NULL()),
		)

	stationLocation := model.StationLocation{}
	err := getStationLocationStmt.Query(dbClient, &stationLocation)

	if err != nil && db.HasNoRow(err) {
		return nil, utils.ForbiddenOperation
	}

	if err != nil {
		log.Println("get-by-id-station-location-error", err.Error())
		return nil, utils.InternalServerError
	}

	return transformToGraphql(stationLocation), nil
}

func (StationLocationService) Update(data UpdateStationLocationData) (*StationLocation, error) {
	dbClient := db.GetPrimaryClient()
	ctx := context.Background()
	tx, err := dbClient.Begin()
	if err != nil {
		log.Println("update-station-location-transaction-error", err.Error())
		return nil, utils.InternalServerError
	}

	now := time.Now()
	var dataFieldsToUpdate = model.StationLocation{
		UpdatedAt: &now,
		UpdatedBy: &data.UpdatedBy,
	}

	var columnsToUpdate pg.ColumnList

	columnsToUpdate = append(
		columnsToUpdate,
		table.StationLocation.UpdatedAt,
		table.StationLocation.UpdatedBy,
	)

	if data.Department != nil {
		columnsToUpdate = append(columnsToUpdate, table.StationLocation.Department)
		dataFieldsToUpdate.Department = *data.Department
	}

	if data.Description != nil {
		columnsToUpdate = append(columnsToUpdate, table.StationLocation.Description)
		dataFieldsToUpdate.Description = data.Description
	}

	if data.Latitude != nil {
		columnsToUpdate = append(columnsToUpdate, table.StationLocation.Latitude)
		dataFieldsToUpdate.Latitude = *data.Latitude
	}

	if data.Longitude != nil {
		columnsToUpdate = append(columnsToUpdate, table.StationLocation.Longitude)
		dataFieldsToUpdate.Longitude = *data.Longitude
	}

	if data.Title != nil {
		columnsToUpdate = append(columnsToUpdate, table.StationLocation.Title)
		dataFieldsToUpdate.Title = *data.Title
	}

	if data.Remark != nil {
		columnsToUpdate = append(columnsToUpdate, table.StationLocation.Remark)
		dataFieldsToUpdate.Remark = data.Remark
	}

	updateStationLocationStmt := table.StationLocation.
		UPDATE(columnsToUpdate).
		MODEL(dataFieldsToUpdate).
		WHERE(table.StationLocation.ID.EQ(pg.UUID(data.Id)).
			AND(table.StationLocation.ProjectId.EQ(pg.UUID(data.ProjectId))),
		).
		RETURNING(table.StationLocation.AllColumns)

	stationLocation := model.StationLocation{}
	err = updateStationLocationStmt.QueryContext(ctx, tx, &stationLocation)

	if err != nil && db.HasNoRow(err) {
		return nil, utils.ForbiddenOperation
	}

	if err != nil {
		log.Println("update-station-location-error", err.Error())
		return nil, utils.InternalServerError
	}

	if data.Tags != nil && len(*data.Tags) != 0 {
		deleteAllTagsStmt := table.StationLocationTag.
			DELETE().
			WHERE(table.StationLocationTag.StationLocationId.
				EQ(pg.UUID(stationLocation.ID)),
			)

		_, err = deleteAllTagsStmt.ExecContext(ctx, tx)

		if err != nil {
			tx.Rollback()
			log.Println("delete-all-station-location-tag-error", err.Error())
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
				log.Println("upsert-station-location-tag-error", err.Error())
				return nil, utils.InternalServerError
			}

			insertStationLocationTagStmt := table.StationLocationTag.
				INSERT(table.StationLocationTag.ID, table.StationLocationTag.StationLocationId, table.StationLocationTag.TagId, table.StationLocationTag.CreatedBy).
				MODEL(model.StationLocationTag{
					ID:                uuid.New(),
					StationLocationId: stationLocation.ID,
					TagId:             tagResult.ID,
					CreatedBy:         data.UpdatedBy,
				})

			_, err = insertStationLocationTagStmt.ExecContext(ctx, tx)

			if err != nil {
				tx.Rollback()
				log.Println("insert-station-location-tag-error", err.Error())
				return nil, utils.InternalServerError
			}
		}
	}

	tx.Commit()

	return transformToGraphql(stationLocation), nil
}

func (StationLocationService) Delete(data DeleteStationLocationData) error {
	dbClient := db.GetPrimaryClient()
	now := time.Now()

	softDeleteStationLocationStmt := table.StationLocation.
		UPDATE(table.StationLocation.DeletedAt, table.StationLocation.DeletedBy).
		MODEL(model.StationLocation{
			DeletedAt: &now,
			DeletedBy: &data.DeletedBy,
		}).
		WHERE(table.StationLocation.ID.EQ(pg.UUID(data.ID)).
			AND(table.StationLocation.ProjectId.EQ(pg.UUID(data.ProjectId))),
		)

	result, err := softDeleteStationLocationStmt.Exec(dbClient)

	if err != nil {
		log.Println("delete-station-location-error", err.Error())
		return utils.InternalServerError
	}

	affectedRow, err := result.RowsAffected()

	if affectedRow == 0 {
		return utils.Notfound
	}

	if err != nil {
		log.Println("delete-station-location-error", err.Error())
		return utils.InternalServerError
	}

	return nil
}

func (StationLocationService) Create(data CreateStationLocationData) (*StationLocation, error) {
	dbClient := db.GetPrimaryClient()
	ctx := context.Background()
	tx, err := dbClient.Begin()
	if err != nil {
		log.Println("create-station-location-transaction-error", err.Error())
		return nil, utils.InternalServerError
	}

	createStationLocationStmt := table.StationLocation.INSERT(
		table.StationLocation.ID,
		table.StationLocation.ProjectId,
		table.StationLocation.CreatedBy,
		table.StationLocation.Department,
		table.StationLocation.Latitude,
		table.StationLocation.Longitude,
		table.StationLocation.Remark,
		table.StationLocation.Title,
	).MODEL(model.StationLocation{
		ID:          uuid.New(),
		ProjectId:   uuid.MustParse(data.ProjectId),
		CreatedBy:   data.CreatedBy,
		Department:  data.Department,
		Description: data.Description,
		Latitude:    data.Latitude,
		Longitude:   data.Longitude,
		Remark:      data.Remark,
		Title:       data.Title,
	}).
		RETURNING(table.StationLocation.AllColumns)

	stationLocation := model.StationLocation{}
	err = createStationLocationStmt.QueryContext(ctx, tx, &stationLocation)

	if err != nil {
		tx.Rollback()
		log.Println("create-station-location-error", err.Error())
		return nil, utils.InternalServerError
	}

	if data.Tags != nil && len(*data.Tags) != 0 {
		for _, tag := range *data.Tags {
			upsertTagStmt := tagUtils.UpsertStatement(tagUtils.UpsertTagData{
				Tag:       tag,
				ProjectId: data.ProjectId,
				CreatedBy: data.CreatedBy,
			})

			tagResult := model.Tag{}

			err := upsertTagStmt.QueryContext(ctx, tx, &tagResult)

			if err != nil {
				tx.Rollback()
				log.Println("upsert-station-location-tag-error", err.Error())
				return nil, utils.InternalServerError
			}

			insertStationLocationTagStmt := table.StationLocationTag.
				INSERT(table.StationLocationTag.ID, table.StationLocationTag.StationLocationId, table.StationLocationTag.TagId, table.StationLocationTag.CreatedBy).
				MODEL(model.StationLocationTag{
					ID:                uuid.New(),
					StationLocationId: stationLocation.ID,
					TagId:             tagResult.ID,
					CreatedBy:         data.CreatedBy,
				})

			_, err = insertStationLocationTagStmt.ExecContext(ctx, tx)

			if err != nil {
				tx.Rollback()
				log.Println("insert-station-location-tag-error", err.Error())
				return nil, utils.InternalServerError
			}
		}
	}

	tx.Commit()

	return transformToGraphql(stationLocation), nil
}

func transformToGraphql(stationLocation model.StationLocation) *StationLocation {
	return &StationLocation{
		ID:          graphql.ID(stationLocation.ID.String()),
		ProjectId:   graphql.ID(stationLocation.ProjectId.String()),
		Title:       stationLocation.Title,
		Department:  stationLocation.Department,
		Description: stationLocation.Description,
		Latitude:    stationLocation.Latitude,
		Longitude:   stationLocation.Longitude,
		Remark:      stationLocation.Remark,
	}
}
