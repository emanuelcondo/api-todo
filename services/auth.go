package services

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/emanuelcondo/api-todo/models"
	"github.com/emanuelcondo/api-todo/repositories"
	"github.com/emanuelcondo/api-todo/utils"
	"net/http"
	"strings"
)

type Role int

const (
	ROLE_ADMIN	Role = 3
	ROLE_EDITOR	Role = 2
	ROLE_VIEWER	Role = 1
)

var mapRoles = map[string]Role{
	"ADMIN": ROLE_ADMIN,
	"EDITOR": ROLE_EDITOR,
	"VIEWER": ROLE_VIEWER,
}

var userRepository repositories.UserRepository

func Authenticate() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			authorizationSplitted := strings.Split(authorization, " ")

			if len(authorizationSplitted) != 2 {
				resp := map[string]string {
					"message": "Field 'Authorization' contains an invalid value.",
				}
				utils.Reply(w, r, http.StatusUnauthorized, resp)
				return
			}

			accessToken := authorizationSplitted[1]
			payload := &models.JWTPayload{}
			tkn, err := jwt.ParseWithClaims(accessToken, payload, func(token *jwt.Token) (interface{}, error) {
				return []byte("1234567890987654321"), nil
			})

			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					resp := map[string]string{
						"message": "Token: invalid signature.",
					}
					utils.Reply(w, r, http.StatusUnauthorized, resp)
				} else {
					resp := map[string]string{
						"message": "An error occurred when validating token.",
					}
					utils.Reply(w, r, http.StatusUnauthorized, resp)
				}
				return
			}
			if !tkn.Valid {
				resp := map[string]string {
					"message": "Invalid token.",
				}
				utils.Reply(w, r, http.StatusUnauthorized, resp)
				return
			}

			user, err := userRepository.FindByEmail(payload.User)
			if err != nil {
				err := map[string]string{
					"message": "An error occurred when looking for the user.",
				}
				utils.Reply(w, r, http.StatusInternalServerError, err)
				return
			}

			ctx := context.WithValue(r.Context(), "role", user.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func RequiredRole(requiredRole Role) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			role := r.Context().Value("role")
			if role == nil {
				err := map[string]string{
					"message": "Role is missing.",
				}
				utils.Reply(w, r, http.StatusForbidden, err)
				return
			}

			key := fmt.Sprintf("%s", role)
			_role := mapRoles[key]

			if requiredRole <= _role {
				next.ServeHTTP(w, r)
			} else {
				err := map[string]string{
					"message": "Insufficient role.",
				}
				utils.Reply(w, r, http.StatusForbidden, err)
				return
			}

		}
		return http.HandlerFunc(fn)
	}
}
