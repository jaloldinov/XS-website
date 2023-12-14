package hashtag

import (
	"xs/internal/pkg"
	hashtag_repo "xs/internal/repository/postgres/hashtag"
	"xs/internal/service/request"
	"xs/internal/service/response"

	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	hashtag Hashtag
}

func NewController(hashtag Hashtag) *Controller {
	return &Controller{hashtag}
}

func (hc Controller) CreateHashtag(c *gin.Context) {
	var data hashtag_repo.CreateHashtagRequest

	if err := request.BindFunc(c, &data, "Name"); err != nil {
		response.RespondError(c, err)
		return
	}

	detail, er := hc.hashtag.HashtagCreate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (hc Controller) GetHashtagById(c *gin.Context) {
	idParam := c.Param("id")

	detail, er := hc.hashtag.HashtagGetById(c, idParam)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (hc Controller) GetHashtagList(c *gin.Context) {
	var filter hashtag_repo.Filter
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

	list, count, er := hc.hashtag.HashtagGetAll(c, filter)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, map[string]interface{}{
		"results": list,
		"count":   count,
	})
}

func (hc Controller) UpdateHashtag(c *gin.Context) {
	var data hashtag_repo.UpdateHashtagRequest
	if err := request.BindFunc(c, &data, "Name"); err != nil {
		response.RespondError(c, err)
		return
	}

	data.Id = c.Param("id")

	er := hc.hashtag.HashtagUpdate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}

func (hc Controller) DeleteHashtag(c *gin.Context) {

	Id := c.Param("id")

	er := hc.hashtag.HashtagDelete(c, Id)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}
