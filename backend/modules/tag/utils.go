package tag

import (
	"checkpoint/.gen/checkpoint/public/model"
	"checkpoint/.gen/checkpoint/public/table"
	"time"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

func UpsertStatement(data UpsertTagData) pg.Statement {
	return table.Tag.
		INSERT(table.Tag.ID, table.Tag.ProjectId, table.Tag.Title, table.Tag.CreatedBy, table.Tag.CreatedAt).
		MODEL(model.Tag{
			ID:        uuid.New(),
			ProjectId: uuid.MustParse(data.ProjectId),
			Title:     data.Tag,
			CreatedBy: data.CreatedBy,
			CreatedAt: time.Now(),
		}).
		ON_CONFLICT(table.Tag.Title, table.Tag.ProjectId).
		DO_UPDATE(pg.SET(table.Tag.ID.SET(table.Tag.EXCLUDED.ID))).
		RETURNING(table.Tag.AllColumns)
}
