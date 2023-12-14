package router

import (
	"xs/internal/auth"
	auth_controller "xs/internal/controller/http/v1/auth"
	menu_controller "xs/internal/controller/http/v1/menu"
	menu_file_controller "xs/internal/controller/http/v1/menu_file"
	post_file_controller "xs/internal/controller/http/v1/post_file"

	hashtag_controller "xs/internal/controller/http/v1/hashtag"
	post_hashtag_controller "xs/internal/controller/http/v1/post_hashtag"

	post_controller "xs/internal/controller/http/v1/post"
	user_controller "xs/internal/controller/http/v1/user"

	"xs/internal/pkg/repository/postgres"
	menu_repo "xs/internal/repository/postgres/menu"
	menu_file_repo "xs/internal/repository/postgres/menu_file"
	post_file_repo "xs/internal/repository/postgres/post_file"

	hashtag_repo "xs/internal/repository/postgres/hashtag"
	post_hashtag_repo "xs/internal/repository/postgres/post_hashtag"

	post_repo "xs/internal/repository/postgres/post"
	user_repo "xs/internal/repository/postgres/user"

	"github.com/gin-gonic/gin"
)

type Auth interface {
	HasPermission(roles ...string) gin.HandlerFunc
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
	menuFileRepo := menu_file_repo.NewRepository(r.postgresDB)
	postFileRepo := post_file_repo.NewRepository(r.postgresDB)
	hashtagRepo := hashtag_repo.NewRepository(r.postgresDB)
	postHashtagRepo := post_hashtag_repo.NewRepository(r.postgresDB)

	//controller
	authController := auth_controller.NewController(userRepo, r.auth)

	userController := user_controller.NewController(userRepo)
	postController := post_controller.NewController(postRepo)
	menuController := menu_controller.NewController(menuRepo)
	menuFileController := menu_file_controller.NewController(menuFileRepo)
	postFileController := post_file_controller.NewController(postFileRepo)
	hashtagController := hashtag_controller.NewController(hashtagRepo)
	postHashtagController := post_hashtag_controller.NewController(postHashtagRepo)

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
	router.GET("/api/v1/admin/menu/list", r.auth.HasPermission("ADMIN"), menuController.GetMenuList)
	router.PUT("/api/v1/admin/menu/:id", r.auth.HasPermission("ADMIN"), menuController.UpdateMenu)
	router.DELETE("/api/v1/admin/menu/:id", r.auth.HasPermission("ADMIN"), menuController.DeleteMenu)

	// #MENU_FILE
	router.POST("api/v1/admin/menu-file/create", r.auth.HasPermission("ADMIN"), menuFileController.CreateMenuFile)
	router.GET("/api/v1/admin/menu-file/list", r.auth.HasPermission("ADMIN"), menuFileController.GetMenuFileList)
	router.GET("/api/v1/admin/menu-file/:id", r.auth.HasPermission("ADMIN"), menuFileController.GetMenuFileById)
	router.PUT("/api/v1/admin/menu-file/:id", r.auth.HasPermission("ADMIN"), menuFileController.UpdateMenuFile)
	router.DELETE("/api/v1/admin/menu-file/:id", r.auth.HasPermission("ADMIN"), menuFileController.DeleteMenuFile)

	// #POST_FILE
	router.POST("api/v1/admin/post-file/create", r.auth.HasPermission("ADMIN"), postFileController.CreatePostFile)
	router.GET("/api/v1/admin/post-file/list", r.auth.HasPermission("ADMIN"), postFileController.GetPostFileList)
	router.GET("/api/v1/admin/post-file/:id", r.auth.HasPermission("ADMIN"), postFileController.GetPostFileById)
	router.PUT("/api/v1/admin/post-file/:id", r.auth.HasPermission("ADMIN"), postFileController.UpdatePostFile)
	router.DELETE("/api/v1/admin/post-file/:id", r.auth.HasPermission("ADMIN"), postFileController.DeletePostFile)

	// #HASHTAG
	router.POST("api/v1/admin/hashtag/create", r.auth.HasPermission("ADMIN"), hashtagController.CreateHashtag)
	router.GET("/api/v1/admin/hashtag/list", r.auth.HasPermission("ADMIN"), hashtagController.GetHashtagList)
	router.GET("/api/v1/admin/hashtag/:id", r.auth.HasPermission("ADMIN"), hashtagController.GetHashtagById)
	router.PUT("/api/v1/admin/hashtag/:id", r.auth.HasPermission("ADMIN"), hashtagController.UpdateHashtag)
	router.DELETE("/api/v1/admin/hashtag/:id", r.auth.HasPermission("ADMIN"), hashtagController.DeleteHashtag)

	// #POST_HASHTAG
	router.POST("api/v1/admin/post-hashtag/create", r.auth.HasPermission("ADMIN"), postHashtagController.CreatePostHashtag)
	router.GET("/api/v1/admin/post-hashtag/list/:post_id", r.auth.HasPermission("ADMIN"), postHashtagController.GetPostHashtagList)
	router.DELETE("/api/v1/admin/post-hashtag/:id", r.auth.HasPermission("ADMIN"), postHashtagController.DeletePostHashtag)

	return router.Run(port)
}
