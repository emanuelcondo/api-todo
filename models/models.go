package models

import (
	//"github.com/jinzhu/gorm"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type User struct {
	//	gorm.Model
	ID        	uint   	`json:"id" gorm:"column:id;primary_key"`
	FirstName 	string 	`json:"first_name" gorm:"column:first_name;not null"`
	LastName  	string 	`json:"last_name" gorm:"column:last_name;not null"`
	Email     	string 	`json:"email" gorm:"column:email;not null;unique_index"`
	Password  	string 	`json:"password" gorm:"column:password;not null"`
	Role		string	`json:"role" gorm:"not null"`
}

type Todo struct {
	//gorm.Model
	ID          uint      `json:"id" gorm:"column:id;primary_key"`
	Title       string    `json:"title" gorm:"column:title;not null"`
	Description string    `json:"description" gorm:"column:description;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP"`
}

type Credentials struct {
	Email		string	`json:"email"`
	Password	string	`json:"password""`
}

type  JWTPayload struct {
	User	string	`json:"user"`
	jwt.StandardClaims
}

type SessionResponse struct {
	AccessToken		string		`json:"access_token"`
	Type			string		`json:"type"`
	Expires			time.Time	`json:"expires"`
	Role			string		`json:"role"`
}
