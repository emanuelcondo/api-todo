package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emanuelcondo/api-todo/models"
	"github.com/emanuelcondo/api-todo/utils"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	json.NewDecoder(r.Body).Decode(user)

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Println(err)
		err := map[string]string{
			"message": "An error occurred when encrypting password.",
		}
		utils.Reply(w, r, http.StatusInternalServerError, err)
	}

	user.Password = string(pass)
	createdUser, err := userRepository.Create(user)

	if err != nil {
		err := map[string]string{
			"message": "An error occurred when creating the user.",
		}
		utils.Reply(w, r, http.StatusInternalServerError, err)
	} else {
		result := map[string]interface{}{
			"user": createdUser,
		}
		utils.Reply(w, r, http.StatusCreated, result)
	}
}
