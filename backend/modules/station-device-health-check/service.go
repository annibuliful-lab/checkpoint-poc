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
		table.StationLocationDeviceHealthCheckActivity.ID,
		table.StationLocationDeviceHealthCheckActivity.Status,
		table.StationLocationDeviceHealthCheckActivity.StationDeviceId,
		table.StationLocationDeviceHealthCheckActivity.ActivityTime,
	)

	fieldsToInsert := model.StationLocationDeviceHealthCheckActivity{
		ID:              uuid.New(),
		StationDeviceId: data.StationDeviceId,
		Status:          model.DeviceStatus(data.Status.String()),
		ActivityTime:    data.ActivityTime,
	}

	if data.Issue != nil {
		columnsToInsert = append(columnsToInsert, table.StationLocationDeviceHealthCheckActivity.Issue)
		fieldsToInsert.Issue = data.Issue
	}

	createActivityStmt := table.StationLocationDeviceHealthCheckActivity.
		INSERT(columnsToInsert).
		MODEL(fieldsToInsert).
		RETURNING(table.StationLocationDeviceHealthCheckActivity.AllColumns)

	healthCheckActivity := model.StationLocationDeviceHealthCheckActivity{}

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

	conditions := table.StationLocationDeviceHealthCheckActivity.StationDeviceId.EQ(pg.UUID(data.StationDeviceId))

	if data.Status != nil {
		conditions = conditions.AND(
			table.StationLocationDeviceHealthCheckActivity.Status.
				EQ(pg.
					NewEnumValue(data.Status.String()),
				),
		)
	}

	if data.StartDate != nil {
		conditions = conditions.AND(
			table.StationLocationDeviceHealthCheckActivity.ActivityTime.
				GT_EQ(
					pg.TimestampzT(*data.StartDate),
				),
		)
	}

	if data.EndDate != nil {
		conditions = conditions.AND(
			table.StationLocationDeviceHealthCheckActivity.ActivityTime.
				LT_EQ(
					pg.TimestampzT(*data.EndDate),
				),
		)
	}

	getActivitiesStmt := pg.
		SELECT(table.StationLocationDeviceHealthCheckActivity.AllColumns).
		FROM(table.StationLocationDeviceHealthCheckActivity).
		WHERE(conditions).
		LIMIT(data.Limit).
		OFFSET(data.Skip).
		ORDER_BY(table.StationLocationDeviceHealthCheckActivity.ActivityTime.DESC())

	activities := []model.StationLocationDeviceHealthCheckActivity{}

	err := getActivitiesStmt.Query(dbClient, &activities)

	if err != nil {
		log.Println("create-station-device-health-check-activity-error", err.Error())
		return nil, utils.InternalServerError
	}

	return lo.Map(activities, func(item model.StationLocationDeviceHealthCheckActivity, index int) *StationDeviceHealthCheckActivity {
		return transformToGraphql(item)
	}), nil
}

func transformToGraphql(item model.StationLocationDeviceHealthCheckActivity) *StationDeviceHealthCheckActivity {
	return &StationDeviceHealthCheckActivity{
		ID:              graphql.ID(item.ID.String()),
		StationDeviceId: graphql.ID(item.StationDeviceId.String()),
		Issue:           item.Issue,
		Status:          enum.GetDeviceStatus(item.Status.String()),
		ActivityTime:    graphql.Time{Time: item.ActivityTime},
	}
}
