package post_hashtag

import (
	post_hashtag_repo "xs/internal/repository/postgres/post_hashtag"
	"xs/internal/service/request"
	"xs/internal/service/response"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	post_hashtag PostHashtag
}

func NewController(post_hashtag PostHashtag) *Controller {
	return &Controller{post_hashtag}
}

func (hc Controller) CreatePostHashtag(c *gin.Context) {
	var data post_hashtag_repo.CreatePostHashtagRequest

	if err := request.BindFunc(c, &data, "PostId", "HashtagId"); err != nil {
		response.RespondError(c, err)
		return
	}

	detail, er := hc.post_hashtag.PostHashtagCreate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (hc Controller) GetPostHashtagList(c *gin.Context) {
	idParam := c.Param("post_id")

	list, count, er := hc.post_hashtag.PostHashtagGetAll(c, idParam)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, map[string]interface{}{
		"results": list,
		"count":   count,
	})
}

func (hc Controller) DeletePostHashtag(c *gin.Context) {

	Id := c.Param("id")

	er := hc.post_hashtag.PostHashtagDelete(c, Id)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}
