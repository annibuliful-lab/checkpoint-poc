package vehicletargetconfigurationimage

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"context"
	"log"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/graph-gophers/dataloader"
	"github.com/samber/lo"
)

type VehicleTargetConfigurationImageService struct{}

// func (VehicleTargetConfigurationImageService) UpsertStatement() pg.Statement {

// }
func (VehicleTargetConfigurationImageService) Dataloader() *dataloader.Loader {

	return dataloader.NewBatchedLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		dbClient := db.GetPrimaryClient()
		var ids []pg.Expression

		for _, id := range keys {
			ids = append(ids, pg.UUID(uuid.MustParse(id.String())))
		}

		vehicleTargetImages := []model.VehicleTargetConfigurationImage{}

		getVehicleTargetTagsStmt := table.VehicleTargetConfigurationImage.
			SELECT(table.VehicleTargetConfigurationImage.AllColumns).
			WHERE(table.VehicleTargetConfigurationImage.VehicleTargetConfigurationId.IN(ids...))

		err := getVehicleTargetTagsStmt.Query(dbClient, &vehicleTargetImages)

		if err != nil {
			log.Println("dataloader-vehicle-target-configuration-image-error", err.Error())
			return nil
		}

		var results []*dataloader.Result

		for _, key := range keys {
			filtered := lo.Filter(vehicleTargetImages, func(item model.VehicleTargetConfigurationImage, index int) bool {
				return item.VehicleTargetConfigurationId == uuid.MustParse(key.String())
			})

			results = append(results, &dataloader.Result{Data: filtered})
		}

		return results
	}, dataloader.WithClearCacheOnBatch())
}
