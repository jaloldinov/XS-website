package menu_file

import (
	"xs/internal/pkg"
	menu_file_repo "xs/internal/repository/postgres/menu_file"
	"xs/internal/service/file"
	"xs/internal/service/request"
	"xs/internal/service/response"

	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	menu_file MenuFile
}

func NewController(user MenuFile) *Controller {
	return &Controller{user}
}

func (mfc Controller) CreateMenuFile(c *gin.Context) {
	var data menu_file_repo.CreateMenuFileRequest

	if err := request.BindFunc(c, &data, "File", "Type", "MenuId"); err != nil {
		response.RespondError(c, err)
		return
	}

	var folder string
	if data.Type == "FILE" {
		folder = "menu/file"
	} else if data.Type == "IMAGE" {
		folder = "menu/images"
	}

	fileLink, err := file.NewService().Upload(c, data.File, folder)
	if err != nil {
		response.RespondError(c, err)
		return
	}
	data.FileLink = &fileLink

	detail, er := mfc.menu_file.MenuFileCreate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (mfc Controller) GetMenuFileById(c *gin.Context) {
	idParam := c.Param("id")

	detail, er := mfc.menu_file.MenuFileGetById(c, idParam)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (mfc Controller) GetMenuFileList(c *gin.Context) {
	var filter menu_file_repo.Filter
	fieldErrors := make([]pkg.FieldError, 0)

	limit, err := request.GetQuery(c, reflect.Int, "limit")
	if err != nil {
		fieldErrors = append(fieldErrors, *err)
	} else if value, ok := limit.(*int); ok {
		filter.Limit = value
	}

	offset, err := request.GetQuery(c, reflect.Int, "offset")
	if err != nil {
		fieldErrors = append(fieldErrors, *err)
	} else if value, ok := offset.(*int); ok {
		filter.Offset = value
	}

	fileType, err := request.GetQuery(c, reflect.String, "type")
	if err != nil {
		fieldErrors = append(fieldErrors, *err)
	} else if value, ok := fileType.(*string); ok {
		filter.Type = value
	}

	menuId, err := request.GetQuery(c, reflect.String, "menu_id")
	if err != nil {
		fieldErrors = append(fieldErrors, *err)
	} else if value, ok := menuId.(*string); ok {
		filter.MenuId = value
	}

	if len(fieldErrors) > 0 {
		response.RespondError(c, &pkg.Error{
			Err:    errors.New("invalid query"),
			Fields: fieldErrors,
			Status: http.StatusBadRequest,
		})
		return
	}

	list, count, er := mfc.menu_file.MenuFileGetAll(c, filter)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, map[string]interface{}{
		"results": list,
		"count":   count,
	})
}

func (mfc Controller) UpdateMenuFile(c *gin.Context) {
	var data menu_file_repo.UpdateMenuFileRequest
	if err := request.BindFunc(c, &data); err != nil {
		response.RespondError(c, err)
		return
	}

	data.Id = c.Param("id")

	var folder string
	if data.Type == "FILE" {
		folder = "menu/file"
	} else if data.Type == "IMAGE" {
		folder = "menu/images"
	}

	fileLink, err := file.NewService().Upload(c, data.File, folder)
	if err != nil {
		response.RespondError(c, err)
		return
	}
	data.FileLink = &fileLink
	er := mfc.menu_file.MenuFileUpdate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}

func (mfc Controller) DeleteMenuFile(c *gin.Context) {

	Id := c.Param("id")

	er := mfc.menu_file.MenuFileDelete(c, Id)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}
