package router

import (
	"github.com/JeerasakTH/go-jwt-api/controller"
	db "github.com/JeerasakTH/go-jwt-api/database"
	"github.com/JeerasakTH/go-jwt-api/middleware"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine, resource *db.PostgresDB) {
	r.POST("/register", controller.Register(resource))
	r.POST("/login", controller.Login(resource))
	r.GET("/users", middleware.JWTAuthen(), controller.GetAllUsers(resource))
	r.GET("/user", middleware.JWTAuthen(), controller.Profile(resource))
}
