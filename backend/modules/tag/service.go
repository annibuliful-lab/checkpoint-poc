package tag

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/.gen/checkpoint/public/table"
	"checkpoint/db"
	"checkpoint/utils"
	"log"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

type TagService struct{}

func (TagService) FindByImsiConfigurationId(imsiConfigurationId uuid.UUID) (*[]model.Tag, error) {
	dbClient := db.GetPrimaryClient()
	tags := []model.Tag{}
	getByImsiConfigurationStmt := table.Tag.
		SELECT(table.Tag.AllColumns).
		FROM(table.Tag.
			INNER_JOIN(table.ImsiConfigurationTag, table.ImsiConfigurationTag.TagId.EQ(table.Tag.ID))).
		WHERE(table.ImsiConfigurationTag.ImsiConfigurationId.EQ(pg.UUID(imsiConfigurationId)))

	err := getByImsiConfigurationStmt.Query(dbClient, &tags)
	if err != nil {
		log.Println("select-imsi-configuration-tags-error", err.Error())
		return nil, utils.InternalServerError
	}

	return &tags, nil
}

func (TagService) FindByImeiConfigurationId(imeiConfigurationId uuid.UUID) (*[]model.Tag, error) {
	dbClient := db.GetPrimaryClient()
	tags := []model.Tag{}
	getByImsiConfigurationStmt := table.Tag.
		SELECT(table.Tag.AllColumns).
		FROM(table.Tag.
			INNER_JOIN(table.ImeiConfigurationTag, table.ImeiConfigurationTag.TagId.EQ(table.Tag.ID))).
		WHERE(table.ImeiConfigurationTag.ImeiConfigurationId.EQ(pg.UUID(imeiConfigurationId)))

	err := getByImsiConfigurationStmt.Query(dbClient, &tags)
	if err != nil {
		log.Println("select-imsi-configuration-tags-error", err.Error())
		return nil, utils.InternalServerError
	}

	return &tags, nil
}
