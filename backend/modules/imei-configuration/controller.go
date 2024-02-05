package imeiconfiguration

import (
	"checkpoint/jwt"
	"checkpoint/utils"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

func DeleteImeiConfigurationController(ctx iris.Context) {
	headers := utils.GetAuthenticationHeaders(ctx)
	token := utils.GetAuthToken(headers.Authorization)
	id, err := uuid.Parse(ctx.Params().Get("id"))

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	payload, _ := jwt.VerifyToken(token)

	code, err := DeleteImeiConfiguration(DeleteImeiConfigurationData{
		ID:        id,
		ProjectId: uuid.MustParse(headers.ProjectId),
		DeletedBy: payload.AccountId.String(),
	})

	if err != nil {
		ctx.StatusCode(code)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.StatusCode(code)
	ctx.JSON(iris.Map{
		"message": "deleted",
		"data":    nil,
	})
}

func GetImeiConfigurationsController(ctx iris.Context) {
	headers := utils.GetAuthenticationHeaders(ctx)
	var query struct {
		Tags   []string `url:"tags"`
		Search string   `url:"search"`
		Label  string   `url:"label"`
		Limit  int64    `url:"limit"`
		Skip   int64    `url:"skip"`
	}

	err := ctx.ReadQuery(&query)

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	imeiconfigurations, code, err := GetImeiConfigurations(GetImeiConfigurationsData{
		ProjectId: uuid.MustParse(headers.ProjectId),
		Tags:      query.Tags,
		Label:     query.Label,
		Search:    query.Search,
		Pagination: utils.OffsetPagination{
			Limit: query.Limit,
			Skip:  query.Skip,
		},
	})

	if err != nil {
		ctx.StatusCode(code)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.StatusCode(200)
	ctx.JSON(iris.Map{
		"data": imeiconfigurations,
	})

}

func GetImeiConfigurationByIdController(ctx iris.Context) {
	headers := utils.GetAuthenticationHeaders(ctx)
	id, err := uuid.Parse(ctx.Params().Get("id"))

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	imeiconfiguration, code, err := GetImeiConfiguration(GetImeiConfigurationData{
		ID:        id,
		ProjectId: uuid.MustParse(headers.ProjectId),
	})

	if err != nil {
		ctx.StatusCode(code)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.StatusCode(200)
	ctx.JSON(iris.Map{
		"data": imeiconfiguration,
	})
}

func UpdateImeiConfigurationController(ctx iris.Context) {
	headers := utils.GetAuthenticationHeaders(ctx)
	token := utils.GetAuthToken(headers.Authorization)
	id, err := uuid.Parse(ctx.Params().Get("id"))

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	var data struct {
		Imei      string   `json:"imei"`
		Priority  string   `json:"priority"`
		CreatedBy string   `json:"createdBy"`
		Label     string   `json:"label"`
		Tags      []string `json:"tags"`
	}

	err = ctx.ReadJSON(&data)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	payload, _ := jwt.VerifyToken(token)

	imeiconfiguration, code, err := UpdateImeiConfiguration(UpdateImeiConfigurationData{
		ID:        id,
		Imei:      data.Imei,
		UpdatedBy: payload.AccountId.String(),
		ProjectId: uuid.MustParse(headers.ProjectId),
		Priority:  data.Priority,
		Label:     data.Label,
		Tags:      data.Tags,
	})

	if err != nil {
		ctx.StatusCode(code)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.StatusCode(200)
	ctx.JSON(iris.Map{
		"message": "updated",
		"data":    imeiconfiguration,
	})
}

func CreateImeiConfigurationController(ctx iris.Context) {
	headers := utils.GetAuthenticationHeaders(ctx)
	token := utils.GetAuthToken(headers.Authorization)
	var data struct {
		Imei              string    `json:"imei"`
		Priority          string    `json:"priority"`
		StationLocationId uuid.UUID `json:"stationLocationId"`
		CreatedBy         string    `json:"createdBy"`
		Label             string    `json:"label"`
		Tags              []string  `json:"tags"`
	}

	err := ctx.ReadJSON(&data)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	payload, _ := jwt.VerifyToken(token)

	imeiconfiguration, code, err := CreateImeiConfiguration(CreateImeiConfigurationData{
		Imei:              data.Imei,
		CreatedBy:         payload.AccountId.String(),
		ProjectId:         uuid.MustParse(headers.ProjectId),
		Priority:          data.Priority,
		Label:             data.Label,
		Tags:              data.Tags,
		StationLocationId: data.StationLocationId,
	})

	if err != nil {
		ctx.StatusCode(code)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.StatusCode(201)
	ctx.JSON(iris.Map{
		"message": "created",
		"data":    imeiconfiguration,
	})
}
