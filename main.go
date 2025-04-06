package main

import (
	"fmt"
	"go_effective/controller"
	_ "go_effective/docs"
	"go_effective/model"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Effective Mobile Test API
//	@version		1.0
//	@description	Cервис, который получает по API ФИО, из открытых API обогащает ответ наиболее вероятными возрастом, полом и национальностью и сохраняет данные в БД. По запросу выдает инфу о найденных людях.

func main() {
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := &model.DB{}
	if err := db.Init(); err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	defer db.Close()

	r := gin.Default()

	c := controller.NewController(db.Conn)

	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.GET(":id", c.ShowUser)
			users.GET("", c.ListUsers)
			users.POST("", c.AddUser)
			users.DELETE(":id", c.DeleteUser)
			users.PATCH(":id", c.UpdateUser)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}
