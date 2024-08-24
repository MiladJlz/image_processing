package main

import (
	"context"
	"database/sql"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"image_processing/api"
	"image_processing/config"
	"image_processing/db"
	"log"
)

func main() {

	appConfig := config.GetConfig()
	conn, err := sql.Open("postgres", appConfig.GetPostgresConn())
	if err != nil {
		log.Fatal(err)
	}
	rdb := redis.NewClient(appConfig.GetRedisOptions())
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}

	var (
		userStore    = db.NewPostgresUserStore(conn)
		linkStore    = db.NewRedisLinkStore(rdb)
		imageStore   = db.NewPostgresImageStore(conn)
		userHandler  = api.NewUserHandler(userStore)
		imageHandler = api.NewImageHandler(imageStore, linkStore)
	)

	e := echo.New()

	e.POST("/register", userHandler.HandleRegister)
	e.GET("/image/:link", imageHandler.HandleImage)

	jwtConfig := appConfig.GetJwtConfig()

	apiv1 := e.Group("/api/v1")
	apiv1.Use(echojwt.WithConfig(jwtConfig))

	apiv1.POST("/upload", imageHandler.HandleUpload)
	apiv1.POST("/generate/:id", imageHandler.HandleGenerate)

	apiv1.GET("/transform/:transform/image/:id", imageHandler.HandleTransform)
	apiv1.GET("/filter/:filter/image/:id", imageHandler.HandleFilter)

	e.Logger.Fatal(e.Start(appConfig.Port))
}
