package post

import (
	post_repo "xs/internal/repository/postgres/post"
	"xs/internal/service/request"
	"xs/internal/service/response"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	post Post
}

func NewController(post Post) *Controller {
	return &Controller{post}
}

func (uc Controller) CreatePost(c *gin.Context) {
	var data post_repo.CreatePostRequest

	if err := request.BindFunc(c, &data, "Title", "Content"); err != nil {
		response.RespondError(c, err)
		return
	}

	detail, er := uc.post.PostCreate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

/*
func (uc Controller) GetPostById(c *gin.Context) {
	idParam := c.Param("id")

	detail, er := uc.post.PostGetById(c, idParam)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (uc Controller) GetPostList(c *gin.Context) {
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

	search, err := request.GetQuery(c, reflect.String, "search")
	if err != nil {
		fieldErrors = append(fieldErrors, *err)
	} else if value, ok := search.(*string); ok {
		filter.Postname = value
	}

	if len(fieldErrors) > 0 {
		response.RespondError(c, &pkg.Error{
			Err:    errors.New("invalid query"),
			Fields: fieldErrors,
			Status: http.StatusBadRequest,
		})
		return
	}

	list, count, er := uc.post.PostGetAll(c, filter)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, map[string]interface{}{
		"results": list,
		"count":   count,
	})
}

func (uc Controller) UpdatePost(c *gin.Context) {
	var data post_repo.UpdatePostRequest
	if err := request.BindFunc(c, &data, "postname", "password", "role"); err != nil {
		response.RespondError(c, err)
		return
	}

	data.Id = c.Param("id")
	createdBy, _ := c.Keys["post_id"].(string)
	data.UpdatedBy = &createdBy

	avatarLink, err := file.NewService().Upload(c, data.Avatar, "avatar")
	if err != nil {
		fmt.Errorf("avatar file uploading error: %v", err)
		return
	}
	data.AvatarLink = &avatarLink

	er := uc.post.PostUpdate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}

func (uc Controller) DeletePost(c *gin.Context) {
	var req post.DeletePostRequest
	req.Id = c.Param("id")
	deletedBy, _ := c.Keys["post_id"].(string)
	req.DeletedBy = deletedBy

	er := uc.post.PostDelete(c, req)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}

func (uc Controller) ResetPostPassword(c *gin.Context) {
	var data post.UpdatePasswordRequest

	if err := request.BindFunc(c, &data, "id", "new_password"); err != nil {
		response.RespondError(c, err)
		return
	}

	updatedBy, _ := c.Keys["post_id"].(string)
	data.UpdatedBy = &updatedBy

	er := uc.post.PostUpdatePassword(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}
*/
