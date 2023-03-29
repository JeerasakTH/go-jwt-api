package main

import (
	"fmt"

	"github.com/JeerasakTH/go-jwt-api/database"
	"github.com/JeerasakTH/go-jwt-api/router"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "golang.org/x/crypto/bcrypt"
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		fmt.Println("ConnectDB error:")
		panic(err)
	}
	r := gin.Default()
	router.Router(r, db)
	r.Use(cors.Default())
	r.Run(":8080")
}
