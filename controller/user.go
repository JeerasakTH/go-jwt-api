package controller

import (
	"fmt"
	"os"
	"time"

	db "github.com/JeerasakTH/go-jwt-api/database"
	"github.com/JeerasakTH/go-jwt-api/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Register(resource *db.PostgresDB) func(c *gin.Context) {
	type Body struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Fullname string `json:"fullname" binding:"required"`
		Avartar  string `json:"avartar" binding:"required"`
	}
	return func(c *gin.Context) {

		var body Body
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		user := model.User{}
		q := "select id, username from users where username=$1"
		_ = resource.DB.Get(&user, q, body.Username)
		// if err != nil {
		// 	c.JSON(500, gin.H{
		// 		"message": err.Error(),
		// 		"data":    "Get Error",
		// 	})
		// 	return
		// }
		if user.ID > 0 {
			c.JSON(500, gin.H{
				"message": "error",
				"data":    "User exists",
			})
			return
		}

		tx, err := resource.DB.Begin()
		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
				"data":    "Begin Error",
			})
			return
		}
		encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 0)

		var lastInsertID int64
		q = "insert into users (username, password, fullname, avartar) values ($1, $2, $3, $4) RETURNING id"
		err = tx.QueryRow(q, body.Username, encryptedPassword, body.Fullname, body.Avartar).Scan(&lastInsertID)
		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
				"data":    "Insert Error",
			})
			return
		}

		// ผ่านทั้งหมดให้คอมมิท
		err = tx.Commit()
		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
				"data":    "Commit error",
			})
			return
		}

		user = model.User{
			ID:       int(lastInsertID),
			Username: body.Username,
			Password: body.Password,
			Fullname: body.Fullname,
			Avartar:  body.Avartar,
		}

		c.JSON(200, gin.H{
			"status": "ok",
			"data":   user,
		})
	}
}

func Login(resource *db.PostgresDB) func(c *gin.Context) {
	type Body struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	return func(c *gin.Context) {
		var body Body
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		user := model.User{}
		q := "select id, username, password from users where username=$1"
		err := resource.DB.Get(&user, q, body.Username)
		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
				"data":    "User does not exist",
			})
			return
		}

		hmacSampleSecret := []byte(os.Getenv("JWT_SECRET_KEY"))
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userID": user.ID,
			"exp":    time.Now().Add(time.Second * 30).Unix(),
		})
		tokenString, err := token.SignedString(hmacSampleSecret)

		fmt.Println(tokenString, err)

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
		if err == nil {
			c.JSON(200, gin.H{
				"message": "ok",
				"data":    "Login successful",
				"token":   tokenString,
			})
			return
		} else {
			c.JSON(500, gin.H{
				"message": err.Error(),
				"data":    "Login fialed",
			})
			return
		}
	}
}

func GetAllUsers(resource *db.PostgresDB) func(c *gin.Context) {
	return func(c *gin.Context) {
		users := []model.User{}
		query := "select * from users"
		err := resource.DB.Select(&users, query)
		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
				"data":    "Select error",
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "success",
			"data":    users,
		})
	}
}

func Profile(resource *db.PostgresDB) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID := c.MustGet("userID")
		user := model.User{}
		q := "select * from users where id=$1"
		err := resource.DB.Get(&user, q, userID)
		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
				"data":    "Get Error",
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "success",
			"data":    user,
		})
	}
}

// func GetAllCouponDetail(resource *db.PostgresDB) func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		query := "select id, name, start_date, end_date, status, coupon_count, coupon_type, reward from coupon_detail"
// 		var couponDetail []model.CouponDetail
// 		err := resource.DB.Select(&couponDetail, query)
// 		if err != nil {
// 			c.JSON(500, gin.H{
// 				"message": err.Error(),
// 				"data":    "Select error",
// 			})
// 			return
// 		}
// 		c.JSON(200, gin.H{
// 			"message": "success",
// 			"data":    couponDetail,
// 		})
// 	}
// }

// func GetCouponDetailByID(resource *db.PostgresDB) func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		couponID := c.Query("id")
// 		query := "select id, name, start_date, end_date, status, coupon_count, coupon_type, reward from coupon_detail where id=$1"
// 		couponDetail := model.CouponDetail{}
// 		err := resource.DB.Get(&couponDetail, query, couponID)
// 		if err != nil {
// 			c.JSON(500, gin.H{
// 				"message": err.Error(),
// 				"data":    "Insert Error",
// 			})
// 			return
// 		}

// 		c.JSON(200, gin.H{
// 			"message": "success",
// 			"data":    couponDetail,
// 		})
// 	}
// }

// func CreateCoupon(resource *db.PostgresDB) func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		tx, err := resource.DB.Begin()
// 		if err != nil {
// 			c.JSON(500, gin.H{
// 				"message": err.Error(),
// 				"data":    "Begin Error",
// 			})
// 			return
// 		}

// 		couponDetail := model.CouponDetail{}
// 		if err := c.ShouldBind(&couponDetail); err != nil {
// 			c.JSON(500, gin.H{
// 				"message": err.Error(),
// 				"data":    "ShouldBind Error",
// 			})
// 			return
// 		}

// 		q := "insert into coupon_detail (id, name, start_date, end_date, status, coupon_count, coupon_type, reward) values (:id, :name, :start_date, :end_date, :status, :coupon_count, :coupon_type, :reward)"
// 		result, err := resource.DB.NamedExec(q, &couponDetail)
// 		if err != nil {
// 			c.JSON(500, gin.H{
// 				"message": err.Error(),
// 				"data":    "Insert Error",
// 			})
// 			return
// 		}

// 		// Check ความเปลี่ยนแปลง
// 		affected, err := result.RowsAffected()
// 		if err != nil {
// 			// ถ้าพังตรงไหนก็โรลแบคกลับ
// 			tx.Rollback()
// 			return
// 		}

// 		if affected <= 0 {
// 			c.JSON(500, gin.H{
// 				"message": err.Error(),
// 				"data":    "Cannot insert",
// 			})
// 			return
// 		}

// 		// ผ่านทั้งหมดให้คอมมิท
// 		err = tx.Commit()
// 		if err != nil {
// 			c.JSON(500, gin.H{
// 				"message": err.Error(),
// 				"data":    "Commit error",
// 			})
// 			return
// 		}

// 		c.JSON(200, gin.H{
// 			"message": "success",
// 			"data":    couponDetail,
// 		})
// 	}
// }

// func UpdateCouponByID(resource *db.PostgresDB) func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		name := c.Query("name")
// 		id := c.Query("id")
// 		tx, err := resource.DB.Begin()
// 		if err != nil {
// 			c.JSON(500, gin.H{
// 				"message": err.Error(),
// 				"data":    "Begin Error",
// 			})
// 			return
// 		}

// 		query := "update coupon_detail set name=$1 where id=$2"
// 		result, err := resource.DB.Exec(query, name, id)
// 		if err != nil {
// 			c.JSON(500, gin.H{
// 				"message": err.Error(),
// 				"data":    "Insert Error",
// 			})
// 			return
// 		}

// 		// Check ความเปลี่ยนแปลง
// 		affected, err := result.RowsAffected()
// 		if err != nil {
// 			// ถ้าพังตรงไหนก็โรลแบคกลับ
// 			tx.Rollback()
// 			return
// 		}

// 		if affected <= 0 {
// 			c.JSON(500, gin.H{
// 				"message": err.Error(),
// 				"data":    "Cannot insert",
// 			})
// 			return
// 		}

// 		// ผ่านทั้งหมดให้คอมมิท
// 		err = tx.Commit()
// 		if err != nil {
// 			c.JSON(500, gin.H{
// 				"message": err.Error(),
// 				"data":    "Commit error",
// 			})
// 			return
// 		}

// 		c.JSON(200, gin.H{
// 			"message": "success",
// 			"data":    map[string]string{"name": name},
// 		})
// 	}
// }

// func DeleteCoupon(resource *db.PostgresDB) func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		id := c.Query("id")
// 		query := "delete from coupon_detail where id=$1"
// 		result, err := resource.DB.Exec(query, id)
// 		if err != nil {
// 			c.JSON(500, gin.H{
// 				"message": err.Error(),
// 				"data":    "Delete Error",
// 			})
// 			return
// 		}

// 		// Check ความเปลี่ยนแปลง
// 		affected, err := result.RowsAffected()
// 		if affected <= 0 {
// 			c.JSON(500, gin.H{
// 				"message": err.Error(),
// 				"data":    "Cannot Delete",
// 			})
// 			return
// 		}

// 		c.JSON(200, gin.H{
// 			"message": "success",
// 			"data":    nil,
// 		})
// 	}
// }
