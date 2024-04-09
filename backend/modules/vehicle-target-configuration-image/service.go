package vehicletargetconfigurationimage

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/modules/upload"
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

type VehicleTargetConfigurationImageService struct{}

func (VehicleTargetConfigurationImageService) Delete(ctx context.Context, data DeleteVehicleTargetConfigurationImageData) (*VehicleTargetConfigurationImage, error) {
	dbClient := db.GetPrimaryClient()
	now := time.Now()

	updateStmt := table.VehicleTargetConfigurationImage.
		UPDATE(
			table.VehicleTargetConfigurationImage.DeletedBy,
			table.VehicleTargetConfigurationImage.DeletedAt,
		).
		MODEL(model.VehicleTargetConfigurationImage{
			DeletedBy: &data.DeletedBy,
			DeletedAt: &now,
		}).
		WHERE(table.VehicleTargetConfigurationImage.ID.EQ(pg.UUID(data.Id))).
		RETURNING(table.VehicleTargetConfigurationImage.AllColumns)

	image := model.VehicleTargetConfigurationImage{}

	err := updateStmt.Query(dbClient, &image)

	if err != nil {
		log.Println("delete-vehicle-target-configuration-image-error", err.Error())
		return nil, utils.InternalServerError
	}

	return transformToGraphql(ctx, image), nil
}
func (VehicleTargetConfigurationImageService) Update(ctx context.Context, data UpdateVehicleTargetConfigurationImageData) (*VehicleTargetConfigurationImage, error) {

	dbClient := db.GetPrimaryClient()
	now := time.Now()

	updateStmt := table.VehicleTargetConfigurationImage.
		UPDATE(
			table.VehicleTargetConfigurationImage.Type,
			table.VehicleTargetConfigurationImage.S3key,
			table.VehicleTargetConfigurationImage.UpdatedBy,
			table.VehicleTargetConfigurationImage.UpdatedAt,
		).
		MODEL(model.VehicleTargetConfigurationImage{
			Type:      model.ImageType(data.Type.String()),
			S3key:     data.S3Key,
			UpdatedBy: &data.UpdatedBy,
			UpdatedAt: &now,
		}).
		WHERE(
			table.VehicleTargetConfigurationImage.ID.EQ(pg.UUID(data.Id)).
				AND(table.VehicleTargetConfigurationImage.DeletedAt.IS_NULL()),
		).
		RETURNING(table.VehicleTargetConfigurationImage.AllColumns)

	image := model.VehicleTargetConfigurationImage{}
	err := updateStmt.Query(dbClient, &image)

	if err != nil && db.HasNoRow(err) {
		return nil, utils.IdNotfound
	}

	if err != nil {
		log.Println("update-vehicle-target-configuration-image-error", err.Error())
		return nil, utils.InternalServerError
	}

	return transformToGraphql(ctx, image), nil
}

func (VehicleTargetConfigurationImageService) Create(ctx context.Context, data CreateVehicleTargetConfigurationImageData) (*VehicleTargetConfigurationImage, error) {
	dbClient := db.GetPrimaryClient()

	insertStmt := table.VehicleTargetConfigurationImage.
		INSERT(
			table.VehicleTargetConfigurationImage.ID,
			table.VehicleTargetConfigurationImage.VehicleTargetConfigurationId,
			table.VehicleTargetConfigurationImage.Type,
			table.VehicleTargetConfigurationImage.S3key,
			table.VehicleTargetConfigurationImage.CreatedBy,
		).
		MODEL(model.VehicleTargetConfigurationImage{
			ID:                           uuid.New(),
			VehicleTargetConfigurationId: data.VehicleTargetConfigurationId,
			Type:                         model.ImageType(data.Type.String()),
			S3key:                        data.S3Key,
			CreatedBy:                    data.CreatedBy,
		}).
		RETURNING(table.VehicleTargetConfigurationImage.AllColumns)

	image := model.VehicleTargetConfigurationImage{}

	err := insertStmt.Query(dbClient, &image)

	if err != nil {
		log.Println("insert-vehicle-target-configuration-image-error", err.Error())
		return nil, utils.InternalServerError
	}

	return transformToGraphql(ctx, image), nil
}

func transformToGraphql(ctx context.Context, data model.VehicleTargetConfigurationImage) *VehicleTargetConfigurationImage {
	url, err := upload.SignedUrl(ctx, data.S3key)

	if err != nil {
		log.Println("get-vehicle-target-image-signed-url-error", err.Error())
	}

	return &VehicleTargetConfigurationImage{
		Id:                           graphql.ID(data.ID.String()),
		VehicleTargetConfigurationId: graphql.ID(data.VehicleTargetConfigurationId.String()),
		S3Key:                        data.S3key,
		Url:                          *url,
	}
}

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
			WHERE(
				table.VehicleTargetConfigurationImage.VehicleTargetConfigurationId.IN(ids...).
					AND(table.VehicleTargetConfigurationImage.DeletedAt.IS_NULL()),
			)

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
