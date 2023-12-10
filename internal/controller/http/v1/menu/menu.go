package menu

import (
	"errors"
	"net/http"
	"reflect"
	"xs/internal/pkg"
	menu_repo "xs/internal/repository/postgres/menu"
	"xs/internal/service/request"
	"xs/internal/service/response"
	"xs/internal/service/slug"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	post Menu
}

func NewController(post Menu) *Controller {
	return &Controller{post}
}

func (mc Controller) CreateMenu(c *gin.Context) {
	var data menu_repo.CreateMenuRequest

	if err := request.BindFunc(c, &data, "Title", "Content", "Type"); err != nil {
		response.RespondError(c, err)
		return
	}

	if data.Type == "EXTRA" && data.ParentId == nil {
		c.JSON(http.StatusBadRequest, response.StatusOk{
			Message: "parent_id is required!",
			Status:  false,
		})
		return
	}

	if data.Type == "MAIN" {
		data.ParentId = nil
	}

	// Dereference the pointer to access the map
	title := data.Title
	// Retrieve the title value for the "uz" key
	uzTitle := (*title)["uz"]
	// Convert the title to slug
	data.Slug = slug.Make(&uzTitle)

	detail, er := mc.post.MenuCreate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (mc Controller) GetMenuList(c *gin.Context) {
	var filter menu_repo.Filter

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

	if len(fieldErrors) > 0 {
		response.RespondError(c, &pkg.Error{
			Err:    errors.New("invalid query"),
			Fields: fieldErrors,
			Status: http.StatusBadRequest,
		})
		return
	}

	list, count, er := mc.post.MenuGetAll(c, filter)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, map[string]interface{}{
		"results": list,
		"count":   count,
	})
}

/*
func (mc Controller) UpdateMenu(c *gin.Context) {
	var data menu_repo.UpdateMenuRequest
	if err := request.BindFunc(c, &data, "Title", "Content", "Status"); err != nil {
		response.RespondError(c, err)
		return
	}

	data.Id = c.Param("id")

	er := mc.post.MenuUpdate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}

func (mc Controller) DeleteMenu(c *gin.Context) {

	Id := c.Param("id")

	er := mc.post.MenuDelete(c, Id)
	if er != nil {
		response.RespondError(c, er)
		return
	}
	response.RespondNoData(c)
}
*/
