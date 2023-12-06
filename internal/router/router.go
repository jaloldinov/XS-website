package router

import (
	"github.com/gin-gonic/gin"
)

type Auth interface {
	HasPermission(roles ...string) gin.HandlerFunc
}

type User interface {
	CreateUser(*gin.Context)
	GetUserById(*gin.Context)
	GetUserList(*gin.Context)
	UpdateUser(*gin.Context)
	DeleteUser(*gin.Context)
	ResetUserPassword(*gin.Context)
}

type AuthController interface {
	SignIn(c *gin.Context)
}

type Router struct {
	auth           Auth
	authController AuthController
	user           User
}

func New(auth Auth, user User, authController AuthController) *Router {
	return &Router{
		auth:           auth,
		user:           user,
		authController: authController,
	}
}

func (r *Router) Init(port string) error {
	router := gin.Default()

	// #AUTH
	router.POST("/api/v1/admin/sign-in", r.authController.SignIn)

	// #USER
	router.POST("api/v1/admin/user/create", r.auth.HasPermission("ADMIN"), r.user.CreateUser)
	router.GET("/api/v1/admin/user/:id", r.auth.HasPermission("ADMIN"), r.user.GetUserById)
	router.GET("/api/v1/admin/user/list", r.auth.HasPermission("ADMIN"), r.user.GetUserList)
	router.PUT("/api/v1/admin/user/:id", r.auth.HasPermission("ADMIN"), r.user.UpdateUser)
	router.DELETE("/api/v1/admin/user/:id", r.auth.HasPermission("ADMIN"), r.user.DeleteUser)

	router.POST("/api/v1/admin/user/reset/password", r.auth.HasPermission("ADMIN"), r.user.ResetUserPassword)

	return router.Run(port)
}
