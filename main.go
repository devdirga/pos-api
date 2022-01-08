package main

import (
	"fmt"
	"myapp/handler"
	"myapp/helper"

	"myapp/model"

	// "myapp/model/user"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	configuration := helper.GetConfig()
	gormParameters := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", configuration.DbHost, configuration.DbPort, configuration.DbName, configuration.DbUsername, configuration.DbPassword)
	gormDB, err := gorm.Open("postgres", gormParameters)
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	h := &handler.Handler{GormDB: gormDB}
	helper.GormDB = gormDB
	// Migrate the schema (tables): User, Authentication
	helper.GormDB.AutoMigrate(&helper.User{})
	helper.GormDB.AutoMigrate(&helper.Authentication{})
	helper.GormDB.AutoMigrate(&model.Student{})
	helper.GormDB.AutoMigrate(&handler.Area{})

	helper.GormDB.Model(&helper.Authentication{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	echoFramework := echo.New()
	echoFramework.Use(middleware.Logger()) // log
	echoFramework.Use(middleware.CORS())   // CORS from Any Origin, Any Method

	echoGroupUseJWT := echoFramework.Group("/api/v1")
	echoGroupUseJWT.Use(middleware.JWT([]byte(configuration.EncryptionKey)))

	echoGroupNoJWT := echoFramework.Group("/api/v1")

	// /api/v1/users : logged in users
	echoGroupUseJWT.POST("/users/logout", model.Logout)
	echoGroupUseJWT.GET("/users/gets", model.GetUser)
	echoGroupUseJWT.POST("/students/create", h.Create)
	// echoGroupUseJWT.GET("/area/insert", h.InsertArea(20, 20, "persegi"))

	// /api/v1/users : public accessible
	echoGroupNoJWT.POST("/users", model.CreateUser)
	echoGroupNoJWT.POST("/users/login", model.Login)

	defer helper.GormDB.Close()
	echoFramework.Logger.Fatal(echoFramework.Start(":1323"))

}
