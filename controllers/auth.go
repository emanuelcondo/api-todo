package controllers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"

	"github.com/emanuelcondo/api-todo/models"
	"github.com/emanuelcondo/api-todo/utils"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	credentials := &models.Credentials{}
	json.NewDecoder(r.Body).Decode(credentials)

	user, err := userRepository.FindByEmail(credentials.Email)
	if err != nil {
		err := map[string]string{
			"message": "An error occurred when looking for the user.",
		}
		utils.Reply(w, r, http.StatusInternalServerError, err)
		return
	} else if user == nil {
		err := map[string]string{
			"message": "Invalid credentials.",
		}
		utils.Reply(w, r, http.StatusUnauthorized, err)
		return
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		err := map[string]string{
			"message": "Invalid credentials.",
		}
		utils.Reply(w, r, http.StatusUnauthorized, err)
		return
	}

	expirationTime := time.Now().Add(1440 * time.Minute)
	payload := &models.JWTPayload{
		User:           user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	accessToken, err := token.SignedString([]byte("1234567890987654321"))

	if err != nil {
		err := map[string]string{
			"message": "An error occurred when generating token for user.",
		}
		utils.Reply(w, r, http.StatusInternalServerError, err)
	} else {
		result := models.SessionResponse{
			AccessToken: accessToken,
			Type:        "Bearer",
			Expires:     expirationTime,
			Role:        user.Role,
		}
		utils.Reply(w, r, http.StatusCreated, result)
	}
}
