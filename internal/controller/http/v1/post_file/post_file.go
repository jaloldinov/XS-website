package post_file

import (
	"xs/internal/pkg"
	post_file_repo "xs/internal/repository/postgres/post_file"
	"xs/internal/service/file"
	"xs/internal/service/request"
	"xs/internal/service/response"

	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	post_file PostFile
}

func NewController(user PostFile) *Controller {
	return &Controller{user}
}

func (mfc Controller) CreatePostFile(c *gin.Context) {
	var data post_file_repo.CreatePostFileRequest

	if err := request.BindFunc(c, &data, "File", "Type", "PostId"); err != nil {
		response.RespondError(c, err)
		return
	}

	var folder string
	if data.Type == "FILE" {
		folder = "post/file"
	} else if data.Type == "IMAGE" {
		folder = "post/images"
	}

	fileLink, err := file.NewService().Upload(c, data.File, folder)
	if err != nil {
		response.RespondError(c, err)
		return
	}
	data.FileLink = &fileLink

	detail, er := mfc.post_file.PostFileCreate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (mfc Controller) GetPostFileById(c *gin.Context) {
	idParam := c.Param("id")

	detail, er := mfc.post_file.PostFileGetById(c, idParam)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (mfc Controller) GetPostFileList(c *gin.Context) {
	var filter post_file_repo.Filter
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

	postId, err := request.GetQuery(c, reflect.String, "post_id")
	if err != nil {
		fieldErrors = append(fieldErrors, *err)
	} else if value, ok := postId.(*string); ok {
		filter.PostId = value
	}

	if len(fieldErrors) > 0 {
		response.RespondError(c, &pkg.Error{
			Err:    errors.New("invalid query"),
			Fields: fieldErrors,
			Status: http.StatusBadRequest,
		})
		return
	}

	list, count, er := mfc.post_file.PostFileGetAll(c, filter)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, map[string]interface{}{
		"results": list,
		"count":   count,
	})
}

func (mfc Controller) UpdatePostFile(c *gin.Context) {
	var data post_file_repo.UpdatePostFileRequest
	if err := request.BindFunc(c, &data); err != nil {
		response.RespondError(c, err)
		return
	}

	data.Id = c.Param("id")

	var folder string
	if data.Type == "FILE" {
		folder = "post/file"
	} else if data.Type == "IMAGE" {
		folder = "post/images"
	}

	fileLink, err := file.NewService().Upload(c, data.File, folder)
	if err != nil {
		response.RespondError(c, err)
		return
	}
	data.FileLink = &fileLink
	er := mfc.post_file.PostFileUpdate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}

func (mfc Controller) DeletePostFile(c *gin.Context) {

	Id := c.Param("id")

	er := mfc.post_file.PostFileDelete(c, Id)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}
