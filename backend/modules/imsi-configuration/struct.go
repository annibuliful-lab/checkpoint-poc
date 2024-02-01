package imsiconfiguration

import (
	"checkpoint/utils"
	"time"

	"github.com/google/uuid"
)

type ImsiConfigurationResponse struct {
	ID        uuid.UUID  `json:"id"`
	ProjectId uuid.UUID  `json:"projectId"`
	Imsi      string     `json:"imsi"`
	CreatedBy string     `json:"createdBy"`
	UpdatedBy *string    `json:"updatedBy"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	Label     *string    `json:"label"`
}

type CreateImsiConfigurationData struct {
	ProjectId uuid.UUID `json:"projectId"`
	Imsi      string    `json:"imsi"`
	CreatedBy string    `json:"createdBy"`
	Label     *string   `json:"label"`
}

type UpdateImsiConfigurationData struct {
	ID        uuid.UUID `json:"id"`
	Imsi      string    `json:"imsi"`
	ProjectId uuid.UUID `json:"projectId"`
	UpdatedBy string    `json:"updatedBy"`
	Label     *string   `json:"label"`
}

type GetImsiConfigurationsData struct {
	Search     string    `json:"search"`
	ProjectId  uuid.UUID `json:"projectId"`
	Label      *string   `json:"label"`
	Pagination utils.OffsetPagination
}

type GetImsiConfigurationByIdData struct {
	ID        uuid.UUID `json:"id"`
	ProjectId uuid.UUID `json:"projectId"`
}

type DeleteImsiConfigurationData struct {
	ID        uuid.UUID `json:"id"`
	ProjectId uuid.UUID `json:"projectId"`
	UpdatedBy string    `json:"updatedBy"`
}
