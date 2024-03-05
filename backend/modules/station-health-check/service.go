package stationhealthcheck

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/gql/enum"
	"checkpoint/utils"
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/samber/lo"
)

type StationHealthCheckActivityService struct{}

func (StationHealthCheckActivityService) Update(data UpdateStationHealthCheckActivityData) (*StationLocationHealthCheckActivity, error) {
	dbClient := db.GetPrimaryClient()

	now := time.Now()
	fieldsToUpdate := model.StationLocationHealthCheckActivity{
		UpdatedAt: &now,
		UpdatedBy: &data.UpdatedBy,
	}

	var columnsToUpdate pg.ColumnList
	columnsToUpdate = append(columnsToUpdate,
		table.StationLocationHealthCheckActivity.UpdatedAt,
		table.StationLocationHealthCheckActivity.UpdatedBy,
	)

	if data.StationStatus != nil {
		columnsToUpdate = append(columnsToUpdate, table.StationLocationHealthCheckActivity.StationStatus)
		fieldsToUpdate.StationStatus = model.StationStatus(*data.StationStatus)
	}

	if data.EndDatetime != nil {
		columnsToUpdate = append(columnsToUpdate, table.StationLocationHealthCheckActivity.EndDatetime)
		fieldsToUpdate.EndDatetime = data.EndDatetime
	}

	if data.StartDatetime != nil {
		columnsToUpdate = append(columnsToUpdate, table.StationLocationHealthCheckActivity.StartDatetime)
		fieldsToUpdate.StartDatetime = *data.StartDatetime
	}

	updateActivityStmt := table.StationLocationHealthCheckActivity.
		UPDATE(columnsToUpdate).
		MODEL(fieldsToUpdate).
		WHERE(table.StationLocationHealthCheckActivity.ID.EQ(pg.UUID(data.ID))).
		RETURNING(table.StationLocationHealthCheckActivity.AllColumns)

	updatedActivity := model.StationLocationHealthCheckActivity{}

	err := updateActivityStmt.Query(dbClient, &updatedActivity)

	if err != nil && db.HasNoRow(err) {
		return nil, utils.Notfound
	}

	if err != nil {
		log.Println("update-station-location-health-check-error", err.Error())
		return nil, utils.InternalServerError
	}

	return transformToGraphql(updatedActivity), nil
}

func (StationHealthCheckActivityService) FindMany(data GetStationHealthCheckActivitiesData) ([]*StationLocationHealthCheckActivity, error) {
	dbClient := db.GetPrimaryClient()

	conditions := table.StationLocationHealthCheckActivity.StationId.EQ(pg.UUID(data.StationId))

	if data.StartDatetime != nil {
		conditions = conditions.AND(
			table.StationLocationHealthCheckActivity.StartDatetime.
				GT_EQ(
					pg.TimestampzT(*data.StartDatetime),
				),
		)
	}

	if data.EndDatetime != nil {
		conditions = conditions.AND(
			table.StationLocationHealthCheckActivity.EndDatetime.
				LT_EQ(
					pg.TimestampzT(*data.EndDatetime),
				),
		)
	}

	if data.StationStatus != nil {
		conditions = conditions.
			AND(
				table.StationLocationHealthCheckActivity.StationStatus.EQ(
					pg.NewEnumValue(string(*data.StationStatus)),
				),
			)
	}

	getActitivitiesStmt := pg.
		SELECT(table.StationLocationHealthCheckActivity.AllColumns).
		FROM(table.StationLocationHealthCheckActivity).
		WHERE(conditions).
		LIMIT(data.Limit).
		OFFSET(data.Skip).
		ORDER_BY(table.StationLocationHealthCheckActivity.StartDatetime.DESC())

	activities := []model.StationLocationHealthCheckActivity{}

	err := getActitivitiesStmt.Query(dbClient, &activities)

	if err != nil {
		log.Println("get-station-location-health-check-error", err.Error())
		return nil, utils.InternalServerError
	}

	return lo.Map(activities, func(item model.StationLocationHealthCheckActivity, index int) *StationLocationHealthCheckActivity {
		return transformToGraphql(item)
	}), nil
}

func (StationHealthCheckActivityService) Create(data CreateStationHealthCheckActivityData) (*StationLocationHealthCheckActivity, error) {
	dbClient := db.GetPrimaryClient()
	ctx := context.Background()
	tx, err := dbClient.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
	})

	if err != nil {
		log.Println("init-tx-create-station-health-check-activity-error", err.Error())
		return nil, utils.InternalServerError
	}

	createStationHealthCheckActivityStmt := table.StationLocationHealthCheckActivity.
		INSERT(
			table.StationLocationHealthCheckActivity.ID,
			table.StationLocationHealthCheckActivity.CreatedBy,
			table.StationLocationHealthCheckActivity.StationId,
			table.StationLocationHealthCheckActivity.StationStatus,
			table.StationLocationHealthCheckActivity.EndDatetime,
			table.StationLocationHealthCheckActivity.StartDatetime,
		).
		MODEL(model.StationLocationHealthCheckActivity{
			ID:            uuid.New(),
			CreatedBy:     data.CreatedBy,
			StationId:     data.StationId,
			StationStatus: data.StationStatus,
			StartDatetime: data.StartDatetime,
			EndDatetime:   data.EndDatetime,
		}).
		RETURNING(table.StationLocationHealthCheckActivity.AllColumns)

	createdActivity := model.StationLocationHealthCheckActivity{}

	err = createStationHealthCheckActivityStmt.QueryContext(ctx, tx, &createdActivity)

	if err != nil && db.IsInvalidForeignKey(err) {
		tx.Rollback()
		return nil, errors.New("invalid stationId")
	}

	if err != nil {
		tx.Rollback()
		log.Println("create-station-location-health-check-activity-error", err.Error())
		return nil, utils.InternalServerError
	}

	now := time.Now()
	updateStationLocationCurrentHealthCheckStmt := table.StationLocation.
		UPDATE(
			table.StationLocation.UpdatedAt,
			table.StationLocation.UpdatedBy,
			table.StationLocation.CurrentHealthCheckId,
		).
		MODEL(model.StationLocation{
			UpdatedAt:            &now,
			UpdatedBy:            &data.CreatedBy,
			CurrentHealthCheckId: &createdActivity.ID,
		}).
		WHERE(table.StationLocation.ID.EQ(pg.UUID(data.StationId)))

	_, err = updateStationLocationCurrentHealthCheckStmt.ExecContext(ctx, tx)

	if err != nil {
		tx.Rollback()
		log.Println("update-station-location-current-health-check-id-error", err.Error())
		return nil, utils.InternalServerError
	}

	tx.Commit()

	return transformToGraphql(createdActivity), nil
}

func transformToGraphql(item model.StationLocationHealthCheckActivity) *StationLocationHealthCheckActivity {

	result := StationLocationHealthCheckActivity{
		ID:            graphql.ID(item.ID.String()),
		StationId:     graphql.ID(item.StationId.String()),
		StationStatus: enum.GetStationStatus(item.StationStatus.String()),
		StartDatetime: graphql.Time{Time: item.StartDatetime},
	}

	if item.EndDatetime != nil {
		result.EndDatetime = &graphql.Time{Time: *item.EndDatetime}
	}

	return &result
}
