package imsiconfiguration

import (
	"checkpoint/jwt"
	"checkpoint/utils"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

func GetImsiConfigurationsController(ctx iris.Context) {
	headers := utils.GetAuthenticationHeaders(ctx)

	var query struct {
		Tags   []string `url:"tags"`
		Search string   `url:"search"`
		Label  string   `url:"label"`
		Mcc    string   `url:"mcc"`
		Mnc    string   `url:"mnc"`
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

	imsiconfigurations, code, err := GetImsiConfigurations(GetImsiConfigurationsData{
		Search:    query.Search,
		Mnc:       query.Mnc,
		Mcc:       query.Mcc,
		Label:     query.Label,
		Tags:      query.Tags,
		ProjectId: uuid.MustParse(headers.ProjectId),
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
		"data": imsiconfigurations,
	})

}

func GetImsiConfigurationController(ctx iris.Context) {
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

	imsiConfiguration, code, err := GetImsiConfigurationById(GetImsiConfigurationByIdData{
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

	ctx.StatusCode(201)
	ctx.JSON(iris.Map{
		"message": "result",
		"data":    imsiConfiguration,
	})
}

func DeleteImsiConfigurationByIdController(ctx iris.Context) {
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

	code, err := DeleteImsiConfigurationById(DeleteImsiConfigurationData{
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

func UpdateImsiConfigurationController(ctx iris.Context) {
	headers := utils.GetAuthenticationHeaders(ctx)
	token := utils.GetAuthToken(headers.Authorization)

	var data struct {
		Imsi              string    `json:"imsi"`
		Label             string    `json:"label"`
		Priority          string    `json:"priority"`
		StationLocationId uuid.UUID `json:"stationLocationId"`
		Tags              []string  `json:"tags"`
	}
	id, err := uuid.Parse(ctx.Params().Get("id"))

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
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

	imsiConfiguration, code, err := UpdateImsiConfiguration(UpdateImsiConfigurationData{
		ID:        id,
		Imsi:      data.Imsi,
		Label:     data.Label,
		ProjectId: uuid.MustParse(headers.ProjectId),
		Tags:      data.Tags,
		UpdatedBy: payload.AccountId.String(),
		Priority:  data.Priority,
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
		"message": "updated",
		"data":    imsiConfiguration,
	})
}

func CreateImsiConfigurationController(ctx iris.Context) {
	headers := utils.GetAuthenticationHeaders(ctx)
	token := utils.GetAuthToken(headers.Authorization)

	var data struct {
		Imsi              string    `json:"imsi"`
		Label             string    `json:"label"`
		Priority          string    `json:"priority"`
		StationLocationId uuid.UUID `json:"stationLocationId"`
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

	imsiConfiguration, code, err := CreateImsiConfiguration(CreateImsiConfigurationData{
		Imsi:              data.Imsi,
		Label:             data.Label,
		ProjectId:         uuid.MustParse(headers.ProjectId),
		StationLocationId: data.StationLocationId,
		Tags:              data.Tags,
		CreatedBy:         payload.AccountId.String(),
		Priority:          data.Priority,
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
		"data":    imsiConfiguration,
	})
}
