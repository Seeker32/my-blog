package main

import (
	"net/http"

	"github.com/Seeker32/my-blog/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dsn := "bloguser:U2Fcz%(brg9d6!LH@tcp(127.0.0.1:13306)/blog?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&model.User{}, &model.Category{}, &model.Tag{}, &model.Article{})
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	apiGroup := router.Group("/api/v1")
	apiGroup.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	_ = router.Run(":8080")
}
