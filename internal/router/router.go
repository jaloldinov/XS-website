package router

import (
	"xs/internal/auth"
	auth_controller "xs/internal/controller/http/v1/auth"
	menu_controller "xs/internal/controller/http/v1/menu"
	post_controller "xs/internal/controller/http/v1/post"
	user_controller "xs/internal/controller/http/v1/user"

	"xs/internal/pkg/repository/postgres"
	menu_repo "xs/internal/repository/postgres/menu"
	post_repo "xs/internal/repository/postgres/post"
	user_repo "xs/internal/repository/postgres/user"

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

type Post interface {
	CreatePost(*gin.Context)
	GetPostById(*gin.Context)
	GetPostList(*gin.Context)
	UpdatePost(*gin.Context)
}

type Menu interface {
	CreatePost(*gin.Context)
	// GetPostById(*gin.Context)
	// GetPostList(*gin.Context)
	// UpdatePost(*gin.Context)
}

type AuthController interface {
	SignIn(c *gin.Context)
}

type Router struct {
	postgresDB *postgres.Database
	auth       *auth.Auth
}

func New(auth *auth.Auth, postgresDB *postgres.Database) *Router {
	return &Router{
		auth:       auth,
		postgresDB: postgresDB,
	}
}

func (r *Router) Init(port string) error {
	router := gin.Default()

	//repository
	userRepo := user_repo.NewRepository(r.postgresDB)
	postRepo := post_repo.NewRepository(r.postgresDB)
	menuRepo := menu_repo.NewRepository(r.postgresDB)

	// mediaRepo := media_repo.NewRepository(postgresDB, fileService)

	//controller
	userController := user_controller.NewController(userRepo)
	postController := post_controller.NewController(postRepo)
	menuController := menu_controller.NewController(menuRepo)

	authController := auth_controller.NewController(userRepo, r.auth)

	// #AUTH
	router.POST("/api/v1/admin/sign-in", authController.SignIn)

	// # ADMIN USER
	router.POST("api/v1/admin/user/create", r.auth.HasPermission("ADMIN"), userController.CreateUser)
	router.GET("/api/v1/admin/user/:id", r.auth.HasPermission("ADMIN"), userController.GetUserById)
	router.GET("/api/v1/admin/user/list", r.auth.HasPermission("ADMIN"), userController.GetUserList)
	router.PUT("/api/v1/admin/user/:id", r.auth.HasPermission("ADMIN"), userController.UpdateUser)
	router.DELETE("/api/v1/admin/user/:id", r.auth.HasPermission("ADMIN"), userController.DeleteUser)
	router.POST("/api/v1/admin/user/reset/password", r.auth.HasPermission("ADMIN"), userController.ResetUserPassword)

	// #POST
	router.POST("api/v1/admin/post/create", r.auth.HasPermission("ADMIN"), postController.CreatePost)
	router.GET("/api/v1/admin/post/:id", r.auth.HasPermission("ADMIN"), postController.PostGetById)
	router.GET("/api/v1/admin/post/list", r.auth.HasPermission("ADMIN"), postController.GetPostList)
	router.PUT("/api/v1/admin/post/:id", r.auth.HasPermission("ADMIN"), postController.UpdatePost)
	router.DELETE("/api/v1/admin/post/:id", r.auth.HasPermission("ADMIN"), postController.DeletePost)

	// #MENU
	router.POST("api/v1/admin/menu/create", r.auth.HasPermission("ADMIN"), menuController.CreateMenu)

	return router.Run(port)
}
