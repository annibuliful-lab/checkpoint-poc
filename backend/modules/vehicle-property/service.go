package vehicleproperty

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/utils"
	"log"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/samber/lo"
)

type VehiclePropertyService struct{}

func (VehiclePropertyService) FindMany(data GetVehiclePropertiesData) ([]*VehicleProperty, error) {
	dbClient := db.GetPrimaryClient()

	conditions := table.VehicleProperty.Type.EQ(pg.NewEnumValue(string(data.Type))).
		AND(table.VehicleProperty.ProjectId.EQ(pg.UUID(data.ProjectId)))

	if data.Search != nil {
		conditions = conditions.AND(table.VehicleProperty.Property.LIKE(pg.String(*data.Search)))
	}

	vehicleProperties := []model.VehicleProperty{}

	getVehiclePropertiesStmt := table.VehicleProperty.
		SELECT(table.VehicleProperty.AllColumns).
		WHERE(conditions).
		LIMIT(data.Limit).
		OFFSET(data.Skip)

	err := getVehiclePropertiesStmt.Query(dbClient, &vehicleProperties)
	if err != nil {
		log.Println("get-vehicle-properties-error", err.Error())
		return nil, utils.InternalServerError
	}

	return lo.Map(vehicleProperties, func(item model.VehicleProperty, index int) *VehicleProperty {
		return transformToGraphql(item)
	}), nil
}

func UpsertStatement(data UpsertVehiclePropertyData) pg.Statement {
	return table.VehicleProperty.
		INSERT(table.VehicleProperty.Property, table.VehicleProperty.Type, table.VehicleProperty.ProjectId).
		MODEL(model.VehicleProperty{
			Property:  data.Property,
			Type:      data.Type,
			ProjectId: data.ProjectId,
		}).
		ON_CONFLICT(
			table.VehicleProperty.Property,
			table.VehicleProperty.Type,
			table.VehicleProperty.ProjectId,
		).
		DO_NOTHING().
		RETURNING(table.VehicleProperty.AllColumns)
}

func transformToGraphql(data model.VehicleProperty) *VehicleProperty {
	return &VehicleProperty{
		Property: data.Property,
		Type:     model.PropertyType(data.Type.String()),
	}
}
