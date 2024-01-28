package project

import (
	"checkpoint/jwt"
	"checkpoint/utils"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

func CreateProjectController(ctx iris.Context) {
	headers := utils.GetAuthenticationHeaders(ctx)

	payload, _ := jwt.VerifyToken(headers.Token)

	var data struct {
		Title string `json:"title"`
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

	project, code, err := CreateProject(CreateProjectData{
		Title:     data.Title,
		AccountId: payload.AccountId.String(),
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
		"data": iris.Map{
			"id":    project.ID,
			"title": project.Title,
		},
	})
}

func UpdateProjectController(ctx iris.Context) {
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

	payload, _ := jwt.VerifyToken(headers.Token)

	match := VerifyOwner(VerifyProjectAccountData{
		ID:        id,
		AccountId: payload.AccountId,
	})

	if !match {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(iris.Map{
			"message": utils.ForbiddenOperation.Error(),
			"data":    nil,
		})
		return
	}

	data := UpdateProjectData{}

	err = ctx.ReadJSON(&data)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	project, code, err := UpdateProject(UpdateProjectData{
		ID:    id,
		Title: data.Title,
	})

	if err != nil {
		ctx.StatusCode(code)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(iris.Map{
		"data": project,
	})
}

func GetProjectByIdController(ctx iris.Context) {
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

	payload, _ := jwt.VerifyToken(headers.Token)

	match := VerifyAccount(VerifyProjectAccountData{
		ID:        id,
		AccountId: payload.AccountId,
	})

	if !match {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(iris.Map{
			"message": utils.ForbiddenOperation.Error(),
			"data":    nil,
		})
		return
	}

	project, code, err := GetProjectById(GetProjectData{
		ID: id,
	})

	if err != nil {
		ctx.StatusCode(code)
		ctx.JSON(iris.Map{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(iris.Map{
		"data": project,
	})
}

func DeleteProjectByIdController(ctx iris.Context) {
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

	payload, _ := jwt.VerifyToken(headers.Token)

	match := VerifyOwner(VerifyProjectAccountData{
		ID:        id,
		AccountId: payload.AccountId,
	})

	if !match {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(iris.Map{
			"message": utils.ForbiddenOperation.Error(),
			"data":    nil,
		})
		return
	}
	accountId := payload.AccountId.String()

	code, err := DeleteProjectById(DeleteProjectData{
		ID:        id,
		AccountId: &accountId,
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
		"message": "success",
	})
}
