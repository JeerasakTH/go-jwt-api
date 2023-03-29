package model

type User struct {
	ID       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username" binding:"required"`
	Password string `db:"password" json:"password" binding:"required"`
	Fullname string `db:"fullname" json:"fullname" binding:"required"`
	Avartar  string `db:"avartar" json:"avartar" binding:"required"`
}
