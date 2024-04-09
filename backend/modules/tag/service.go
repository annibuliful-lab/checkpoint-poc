package tag

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/utils"
	"context"
	"log"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	dataloader "github.com/graph-gophers/dataloader"
	"github.com/graph-gophers/graphql-go"
	"github.com/samber/lo"
)

type TagService struct{}

func (TagService) FindMany(data GetTagsInput) (*[]Tag, error) {
	dbClient := db.GetPrimaryClient()
	conditions := pg.Bool(true)

	if data.Search != nil {
		conditions = conditions.AND(table.Tag.Title.LIKE(pg.String(*data.Search)))
	}

	modelTags := []model.Tag{}
	getTagsStmt := table.Tag.
		SELECT(table.Tag.AllColumns).
		FROM(table.Tag).
		WHERE(conditions).
		LIMIT(data.Limit).
		OFFSET(data.Skip)

	err := getTagsStmt.Query(dbClient, modelTags)
	if err != nil {
		return nil, err
	}

	tags := lo.Map(modelTags, func(item model.Tag, index int) Tag {
		return Tag{
			Id:        graphql.ID(item.ID.String()),
			Title:     item.Title,
			ProjectId: graphql.ID(item.ProjectId.String()),
		}
	})

	return &tags, nil
}

func (TagService) FindByImsiConfigurationId(imsiConfigurationId uuid.UUID) (*[]Tag, error) {
	dbClient := db.GetPrimaryClient()
	modelTags := []model.Tag{}
	getByImsiConfigurationStmt := table.Tag.
		SELECT(table.Tag.AllColumns).
		FROM(table.Tag.
			INNER_JOIN(table.ImsiConfigurationTag, table.ImsiConfigurationTag.TagId.EQ(table.Tag.ID))).
		WHERE(table.ImsiConfigurationTag.ImsiConfigurationId.EQ(pg.UUID(imsiConfigurationId)))

	err := getByImsiConfigurationStmt.Query(dbClient, &modelTags)
	if err != nil {
		log.Println("select-imsi-configuration-tags-error", err.Error())
		return nil, utils.InternalServerError
	}

	tags := lo.Map(modelTags, func(item model.Tag, index int) Tag {
		return Tag{
			Id:        graphql.ID(item.ID.String()),
			Title:     item.Title,
			ProjectId: graphql.ID(item.ProjectId.String()),
		}
	})

	return &tags, nil
}

func (TagService) FindByImeiConfigurationId(imeiConfigurationId uuid.UUID) (*[]Tag, error) {
	dbClient := db.GetPrimaryClient()
	modelTags := []model.Tag{}
	getByImsiConfigurationStmt := table.Tag.
		SELECT(table.Tag.AllColumns).
		FROM(table.Tag.
			INNER_JOIN(table.ImeiConfigurationTag, table.ImeiConfigurationTag.TagId.EQ(table.Tag.ID))).
		WHERE(table.ImeiConfigurationTag.ImeiConfigurationId.EQ(pg.UUID(imeiConfigurationId)))

	err := getByImsiConfigurationStmt.Query(dbClient, &modelTags)
	if err != nil {
		log.Println("select-imsi-configuration-tags-error", err.Error())
		return nil, utils.InternalServerError
	}

	tags := lo.Map(modelTags, func(item model.Tag, index int) Tag {
		return Tag{
			Id:        graphql.ID(item.ID.String()),
			Title:     item.Title,
			ProjectId: graphql.ID(item.ProjectId.String()),
		}
	})

	return &tags, nil
}

func (TagService) ImeiConfigurationDataloader() *dataloader.Loader {

	return dataloader.NewBatchedLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		dbClient := db.GetPrimaryClient()
		var ids []pg.Expression

		for _, id := range keys {
			ids = append(ids, pg.UUID(uuid.MustParse(id.String())))
		}

		var imeiTags = []ImeiTag{}

		getByImsiConfigurationStmt := table.Tag.
			SELECT(table.Tag.AllColumns, table.ImeiConfigurationTag.AllColumns).
			FROM(table.Tag.
				INNER_JOIN(table.ImeiConfigurationTag, table.ImeiConfigurationTag.TagId.EQ(table.Tag.ID))).
			WHERE(table.ImeiConfigurationTag.ImeiConfigurationId.IN(ids...))

		err := getByImsiConfigurationStmt.Query(dbClient, &imeiTags)

		if err != nil {
			log.Println("dataloader-imei-configuration-tags-error", err.Error())
			return nil
		}

		var results []*dataloader.Result

		for _, key := range keys {
			filtered := lo.Filter(imeiTags, func(item ImeiTag, index int) bool {
				return item.ImeiConfigurationId == uuid.MustParse(key.String())
			})

			results = append(results, &dataloader.Result{Data: filtered})
		}

		return results
	}, dataloader.WithClearCacheOnBatch())
}

func (TagService) ImsiConfigurationDataloader() *dataloader.Loader {

	return dataloader.NewBatchedLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		dbClient := db.GetPrimaryClient()
		var ids []pg.Expression

		for _, id := range keys {
			ids = append(ids, pg.UUID(uuid.MustParse(id.String())))
		}

		var imsiTags = []ImsiTag{}

		getByImsiConfigurationStmt := table.Tag.
			SELECT(table.Tag.AllColumns, table.ImsiConfigurationTag.AllColumns).
			FROM(table.Tag.
				INNER_JOIN(table.ImsiConfigurationTag, table.ImsiConfigurationTag.TagId.EQ(table.Tag.ID))).
			WHERE(table.ImsiConfigurationTag.ImsiConfigurationId.IN(ids...))

		err := getByImsiConfigurationStmt.Query(dbClient, &imsiTags)

		if err != nil {
			log.Println("dataloader-imsi-configuration-tags-error", err.Error())
			return nil
		}

		var results []*dataloader.Result

		for _, key := range keys {
			filtered := lo.Filter(imsiTags, func(item ImsiTag, index int) bool {
				return item.ImsiConfigurationId == uuid.MustParse(key.String())
			})

			results = append(results, &dataloader.Result{Data: filtered})
		}

		return results
	}, dataloader.WithClearCacheOnBatch())
}

func (TagService) StationLocationTagDataloader() *dataloader.Loader {
	return dataloader.NewBatchedLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		dbClient := db.GetPrimaryClient()
		var ids []pg.Expression

		for _, id := range keys {
			ids = append(ids, pg.UUID(uuid.MustParse(id.String())))
		}

		stationLocationTags := []StationLocationTag{}

		getStationLocationTagsStmt := table.Tag.
			SELECT(table.Tag.AllColumns, table.StationLocationTag.AllColumns).
			FROM(table.Tag.
				INNER_JOIN(table.StationLocationTag, table.StationLocationTag.TagId.EQ(table.Tag.ID)),
			).WHERE(table.StationLocationTag.StationLocationId.IN(ids...))

		err := getStationLocationTagsStmt.Query(dbClient, &stationLocationTags)

		if err != nil {
			log.Println("dataloader-mobile-configuration-tags-error", err.Error())
			return nil
		}

		var results []*dataloader.Result

		for _, key := range keys {
			filtered := lo.Filter(stationLocationTags, func(item StationLocationTag, index int) bool {
				return item.StationLocationId == uuid.MustParse(key.String())
			})

			results = append(results, &dataloader.Result{Data: filtered})
		}

		return results
	}, dataloader.WithClearCacheOnBatch())
}

func (TagService) MobileDeviceConfigurationTagDataloader() *dataloader.Loader {

	return dataloader.NewBatchedLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		dbClient := db.GetPrimaryClient()
		var ids []pg.Expression

		for _, id := range keys {
			ids = append(ids, pg.UUID(uuid.MustParse(id.String())))
		}

		mobileTags := []MobileDeviceTag{}

		getByMobileTagsStmt := table.Tag.
			SELECT(table.Tag.AllColumns, table.MobileDeviceConfigurationTag.AllColumns).
			FROM(table.Tag.
				INNER_JOIN(table.MobileDeviceConfigurationTag, table.MobileDeviceConfigurationTag.TagId.EQ(table.Tag.ID))).
			WHERE(table.MobileDeviceConfigurationTag.MobileDeviceConfigurationId.IN(ids...))

		err := getByMobileTagsStmt.Query(dbClient, &mobileTags)

		if err != nil {
			log.Println("dataloader-mobile-configuration-tags-error", err.Error())
			return nil
		}

		var results []*dataloader.Result

		for _, key := range keys {
			filtered := lo.Filter(mobileTags, func(item MobileDeviceTag, index int) bool {
				return item.MobileDeviceConfigurationId == uuid.MustParse(key.String())
			})

			results = append(results, &dataloader.Result{Data: filtered})
		}

		return results
	}, dataloader.WithClearCacheOnBatch())
}

func (TagService) VehicleTargetConfigurationTags() *dataloader.Loader {

	return dataloader.NewBatchedLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		dbClient := db.GetPrimaryClient()
		var ids []pg.Expression

		for _, id := range keys {
			ids = append(ids, pg.UUID(uuid.MustParse(id.String())))
		}

		vehicleTargetTags := []VehicleTargetConfigurationTag{}

		getVehicleTargetTagsStmt := table.Tag.
			SELECT(table.Tag.AllColumns, table.VehicleTargetConfigurationTag.AllColumns).
			FROM(table.Tag.
				INNER_JOIN(table.VehicleTargetConfigurationTag, table.VehicleTargetConfigurationTag.TagId.EQ(table.Tag.ID))).
			WHERE(table.VehicleTargetConfigurationTag.VehicleTargetConfigurationId.IN(ids...))

		err := getVehicleTargetTagsStmt.Query(dbClient, &vehicleTargetTags)

		if err != nil {
			log.Println("dataloader-vehicle-target-configuration-tag-error", err.Error())
			return nil
		}

		var results []*dataloader.Result

		for _, key := range keys {
			filtered := lo.Filter(vehicleTargetTags, func(item VehicleTargetConfigurationTag, index int) bool {
				return item.VehicleTargetConfigurationId == uuid.MustParse(key.String())
			})

			results = append(results, &dataloader.Result{Data: filtered})
		}

		return results
	}, dataloader.WithClearCacheOnBatch())
}
