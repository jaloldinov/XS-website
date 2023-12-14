package post

import (
	"errors"
	"net/http"
	"reflect"
	"xs/internal/pkg"
	post_repo "xs/internal/repository/postgres/post"
	"xs/internal/service/request"
	"xs/internal/service/response"
	"xs/internal/service/slug"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	post Post
}

func NewController(post Post) *Controller {
	return &Controller{post}
}

func (pc Controller) CreatePost(c *gin.Context) {
	var data post_repo.CreatePostRequest

	if err := request.BindFunc(c, &data, "Title", "Content", "MenuId"); err != nil {
		response.RespondError(c, err)
		return
	}

	isStatic, err := pc.post.IsMenuStatic(c, *data.MenuId)
	if err != nil {
		response.RespondError(c, err)
		return
	}
	if isStatic {
		c.JSON(http.StatusNotAcceptable, response.StatusBadRequest{
			Error: "menu is static! you can't create a post on that menu",
		})
		return
	}
	// Dereference the pointer to access the map
	title := data.Title
	// Retrieve the title value for the "uz" key
	uzTitle := (title)["uz"]
	// Convert the title to slug
	data.Slug = slug.Make(&uzTitle)

	detail, er := pc.post.PostCreate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (pc Controller) PostGetById(c *gin.Context) {
	idParam := c.Param("id")

	detail, er := pc.post.PostGetById(c, idParam)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (pc Controller) GetPostList(c *gin.Context) {
	var filter post_repo.Filter

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

	order, err := request.GetQuery(c, reflect.String, "order")
	if err != nil {
		fieldErrors = append(fieldErrors, *err)
	} else if value, ok := order.(*string); ok {
		filter.Order = value
	}

	title, err := request.GetQuery(c, reflect.String, "title")
	if err != nil {
		fieldErrors = append(fieldErrors, *err)
	} else if value, ok := title.(*string); ok {
		filter.Title = value
	}

	if len(fieldErrors) > 0 {
		response.RespondError(c, &pkg.Error{
			Err:    errors.New("invalid query"),
			Fields: fieldErrors,
			Status: http.StatusBadRequest,
		})
		return
	}

	list, count, er := pc.post.PostGetAll(c, filter)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, map[string]interface{}{
		"results": list,
		"count":   count,
	})
}

func (pc Controller) UpdatePost(c *gin.Context) {
	var data post_repo.UpdatePostRequest
	if err := request.BindFunc(c, &data, "Title", "Content", "Status"); err != nil {
		response.RespondError(c, err)
		return
	}

	data.Id = c.Param("id")

	er := pc.post.PostUpdate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}

func (pc Controller) DeletePost(c *gin.Context) {

	Id := c.Param("id")

	er := pc.post.PostDelete(c, Id)
	if er != nil {
		response.RespondError(c, er)
		return
	}
	response.RespondNoData(c)
}
