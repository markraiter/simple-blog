package app

import (
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
	"github.com/markraiter/simple-blog/cmd/migrate"
	"github.com/markraiter/simple-blog/internal/initializers"
	"github.com/markraiter/simple-blog/pkg/auth"
	"github.com/markraiter/simple-blog/pkg/handlers"
	"github.com/markraiter/simple-blog/pkg/middlewares"
	_ "github.com/markraiter/simple-blog/docs"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	migrate.Migrate()
}

// @title Simple Blog Project REST API Swagger Example
// @version 1.0
// @description This is a simple blog project for practicing go REST API.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func Start() {
	e := echo.New()

	// Serve the Swagger docs
	swagHandler := echoSwagger.WrapHandler
	e.GET("/swagger/*", swagHandler)

	//////////////////// Groups ////////////////////

	authGroup := e.Group("/api")
	authGroup.Use(middlewares.JWTMiddleware)
	postGroup := authGroup.Group("/v1/posts")
	commentsGroup := authGroup.Group("/v1/comments")

	//////////////////// ENDPOINTS ////////////////////

	// Registration
	e.POST("/registration", auth.Register(initializers.DB))

	// Authentification
	e.POST("/login", auth.Login(initializers.DB))

	// Operations with posts
	postGroup.GET("", handlers.GetPosts(initializers.DB))
	postGroup.GET("/:id", handlers.GetPostByID(initializers.DB))
	postGroup.POST("", handlers.CreatePost(initializers.DB))
	postGroup.PUT("/:id", handlers.UpdatePost(initializers.DB))
	postGroup.DELETE("/:id", handlers.DeletePost(initializers.DB))

	// Operations with comments
	commentsGroup.GET("", handlers.GetComments(initializers.DB))
	commentsGroup.GET("/:id", handlers.GetCommentByID(initializers.DB))
	commentsGroup.POST("", handlers.CreateComment(initializers.DB))
	commentsGroup.PUT("/:id", handlers.UpdateComment(initializers.DB))
	commentsGroup.DELETE("/:id", handlers.DeleteComment(initializers.DB))

	e.Logger.Fatal(e.Start(":8080"))
}
