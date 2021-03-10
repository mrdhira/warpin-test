package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/mrdhira/warpin-test/api/Orders/entities"
	"github.com/mrdhira/warpin-test/pkg"
	log "github.com/sirupsen/logrus"
)

// AuthMiddleware func
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Validate Token
		Authorization := req.Header.Get("Authorization")
		TokenData := &entities.TokenClaim{}
		Token, err := jwt.ParseWithClaims(Authorization, TokenData, func(token *jwt.Token) (interface{}, error) {
			if jwt.GetSigningMethod("HS256") != token.Method {
				return nil, fmt.Errorf("Unex[ected signing method: %v", token.Header["alg"])
			}

			return []byte("secret"), nil
		})

		if Token != nil && err == nil {
			TokenDataJSON, _ := json.Marshal(TokenData)
			fmt.Println("Token Verified: ", string(TokenDataJSON))
			context.Set(req, "token", string(TokenDataJSON))
		} else {
			log.WithFields(log.Fields{
				"event": "unauthorized token",
				"data":  Token,
			}).Error(err)
			pkg.Response(res, 401, pkg.JSONResponse{
				Code:    401,
				Message: "Unauthorized",
				Error:   err.Error(),
			})
			return
		}

		next.ServeHTTP(res, req)
	})
}

// AuthAdmniMiddleware func
func AuthAdmniMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Validate Token
		Authorization := req.Header.Get("Authorization")
		if Authorization == "INTERNAL-SERVICES" {
			fmt.Println("coming request from internal services")
			TokenData := &entities.TokenClaim{
				UserID:   0,
				UserRole: "ADMIN",
			}
			TokenDataJSON, _ := json.Marshal(TokenData)
			context.Set(req, "token", string(TokenDataJSON))
		} else {
			TokenData := &entities.TokenClaim{}
			Token, err := jwt.ParseWithClaims(Authorization, TokenData, func(token *jwt.Token) (interface{}, error) {
				if jwt.GetSigningMethod("HS256") != token.Method {
					return nil, fmt.Errorf("Unex[ected signing method: %v", token.Header["alg"])
				}

				return []byte("secret"), nil
			})

			if Token != nil && err == nil {
				TokenDataJSON, _ := json.Marshal(TokenData)
				fmt.Println("Token Verified: ", string(TokenDataJSON))
				if TokenData.UserRole != entities.Admin {
					pkg.Response(res, 403, &pkg.JSONResponse{
						Code:    403,
						Message: "Forbidden Access",
					})
					return
				}
				context.Set(req, "token", string(TokenDataJSON))
			} else {
				log.WithFields(log.Fields{
					"event": "unauthorized token",
					"data":  Token,
				}).Error(err)
				pkg.Response(res, 401, &pkg.JSONResponse{
					Code:    401,
					Message: "Unauthorized",
					Error:   err.Error(),
				})
				return
			}
		}

		next.ServeHTTP(res, req)
	})
}
