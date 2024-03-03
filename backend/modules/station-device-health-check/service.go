package stationdevicehealthcheck

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/gql/enum"
	"checkpoint/utils"
	"errors"
	"log"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/samber/lo"
)

type StationDeviceHealthCheckActivityService struct{}

func (StationDeviceHealthCheckActivityService) Create(data CreateStationDeviceHealthCheckActivityData) (*StationDeviceHealthCheckActivity, error) {
	dbClient := db.GetPrimaryClient()

	var columnsToInsert pg.ColumnList

	columnsToInsert = append(
		columnsToInsert,
		table.StationDeviceHealthCheckActivity.ID,
		table.StationDeviceHealthCheckActivity.Status,
		table.StationDeviceHealthCheckActivity.StationDeviceId,
		table.StationDeviceHealthCheckActivity.ActivityTime,
	)

	fieldsToInsert := model.StationDeviceHealthCheckActivity{
		ID:              uuid.New(),
		StationDeviceId: data.StationDeviceId,
		Status:          model.DeviceStatus(data.Status.String()),
		ActivityTime:    data.ActivityTime,
	}

	if data.Issue != nil {
		columnsToInsert = append(columnsToInsert, table.StationDeviceHealthCheckActivity.Issue)
		fieldsToInsert.Issue = data.Issue
	}

	createActivityStmt := table.StationDeviceHealthCheckActivity.
		INSERT(columnsToInsert).
		MODEL(fieldsToInsert).
		RETURNING(table.StationDeviceHealthCheckActivity.AllColumns)

	healthCheckActivity := model.StationDeviceHealthCheckActivity{}

	err := createActivityStmt.Query(dbClient, &healthCheckActivity)

	if err != nil && db.IsInvalidForeignKey(err) {
		return nil, errors.New("invalid stationDeviceId")
	}

	if err != nil {
		log.Println("create-station-device-health-check-activity-error", err.Error())
		return nil, utils.InternalServerError
	}

	return transformToGraphql(healthCheckActivity), nil
}

func (StationDeviceHealthCheckActivityService) FindMany(data GetStationDeviceHealthCheckActivitiesData) ([]*StationDeviceHealthCheckActivity, error) {
	dbClient := db.GetPrimaryClient()

	conditions := table.StationDeviceHealthCheckActivity.StationDeviceId.EQ(pg.UUID(data.StationDeviceId))

	if data.Status != nil {
		conditions = conditions.AND(
			table.StationDeviceHealthCheckActivity.Status.
				EQ(pg.
					NewEnumValue(data.Status.String()),
				),
		)
	}

	if data.StartDate != nil {
		conditions = conditions.AND(
			table.StationDeviceHealthCheckActivity.ActivityTime.
				GT_EQ(
					pg.TimestampzT(*data.StartDate),
				),
		)
	}

	if data.EndDate != nil {
		conditions = conditions.AND(
			table.StationDeviceHealthCheckActivity.ActivityTime.
				LT_EQ(
					pg.TimestampzT(*data.EndDate),
				),
		)
	}

	getActivitiesStmt := pg.
		SELECT(table.StationDeviceHealthCheckActivity.AllColumns).
		FROM(table.StationDeviceHealthCheckActivity).
		WHERE(conditions).
		LIMIT(data.Limit).
		OFFSET(data.Skip).
		ORDER_BY(table.StationDeviceHealthCheckActivity.ActivityTime.DESC())

	activities := []model.StationDeviceHealthCheckActivity{}

	err := getActivitiesStmt.Query(dbClient, &activities)

	if err != nil {
		log.Println("create-station-device-health-check-activity-error", err.Error())
		return nil, utils.InternalServerError
	}

	return lo.Map(activities, func(item model.StationDeviceHealthCheckActivity, index int) *StationDeviceHealthCheckActivity {
		return transformToGraphql(item)
	}), nil
}

func transformToGraphql(item model.StationDeviceHealthCheckActivity) *StationDeviceHealthCheckActivity {
	return &StationDeviceHealthCheckActivity{
		ID:              graphql.ID(item.ID.String()),
		StationDeviceId: graphql.ID(item.StationDeviceId.String()),
		Issue:           item.Issue,
		Status:          enum.GetDeviceStatus(item.Status.String()),
		ActivityTime:    graphql.Time{Time: item.ActivityTime},
	}
}
