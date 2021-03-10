package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gorilla/context"
	"github.com/mrdhira/warpin-test/api/Users/entities"
	"github.com/mrdhira/warpin-test/api/Users/usecases"
	"github.com/mrdhira/warpin-test/pkg"
	log "github.com/sirupsen/logrus"
)

var validate *validator.Validate

// UsersControllers struct
type UsersControllers struct {
	UsersUsecase usecases.IUsersUsecases
}

// InitUsersControllers func
func InitUsersControllers() *UsersControllers {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Init Usecase
	usersUsecases := usecases.InitUsersUsecases()

	return &UsersControllers{
		UsersUsecase: usersUsecases,
	}
}

// Register func
func (c *UsersControllers) Register(res http.ResponseWriter, req *http.Request) {
	RawPayload, _ := ioutil.ReadAll(req.Body)
	RawPayloadString := string(RawPayload)
	RawPayloadString = strings.ReplaceAll(RawPayloadString, "\n", "")
	RawPayloadString = strings.ReplaceAll(RawPayloadString, "\\", "")
	RawPayloadString = strings.TrimSpace(RawPayloadString)
	// Log Raw Payload
	log.WithFields(log.Fields{
		"data": RawPayloadString,
	}).Info("POST /users/register payload body")

	var requestBody *entities.RegisterRequest
	if err := json.Unmarshal(RawPayload, &requestBody); err != nil {
		log.WithFields(log.Fields{
			"event": "error when unmarshal request payload register",
		}).Error(err)
		pkg.Response(res, http.StatusBadRequest, &pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	// Payload Validation
	if err := validate.Struct(requestBody); err != nil {
		errField := map[string]string{}
		errFields := []map[string]string{}

		for _, e := range err.(validator.ValidationErrors) {
			errField[e.Field()] = fmt.Sprintf("%s failed on the %s tag", e.Field(), e.Tag())
		}
		errFields = append(errFields, errField)

		log.WithFields(log.Fields{
			"event":            "payload validation error",
			"validation_error": errFields,
		})

		pkg.Response(res, http.StatusUnprocessableEntity, &pkg.JSONResponse{
			Code:    422,
			Message: "payload validation error",
			Error:   err.Error(),
			Data:    errFields,
		})
		return
	}

	Response, err := c.UsersUsecase.Register(req.Context(), requestBody)
	if err != nil {
		pkg.Response(res, http.StatusBadRequest, pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	pkg.Response(res, Response.Code, Response)
	return
}

// Login func
func (c *UsersControllers) Login(res http.ResponseWriter, req *http.Request) {
	RawPayload, _ := ioutil.ReadAll(req.Body)
	RawPayloadString := string(RawPayload)
	RawPayloadString = strings.ReplaceAll(RawPayloadString, "\n", "")
	RawPayloadString = strings.ReplaceAll(RawPayloadString, "\\", "")
	RawPayloadString = strings.TrimSpace(RawPayloadString)
	// Log Raw Payload
	log.WithFields(log.Fields{
		"data": RawPayloadString,
	}).Info("POST /users/register payload body")

	var requestBody *entities.LoginRequest
	if err := json.Unmarshal(RawPayload, &requestBody); err != nil {
		log.WithFields(log.Fields{
			"event": "error when unmarshal request payload login",
		}).Error(err)
		pkg.Response(res, http.StatusBadRequest, &pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	// Payload Validation
	if err := validate.Struct(requestBody); err != nil {
		errField := map[string]string{}
		errFields := []map[string]string{}

		for _, e := range err.(validator.ValidationErrors) {
			errField[e.Field()] = fmt.Sprintf("%s failed on the %s tag", e.Field(), e.Tag())
		}
		errFields = append(errFields, errField)

		log.WithFields(log.Fields{
			"event":            "payload validation error",
			"validation_error": errFields,
		})

		pkg.Response(res, http.StatusUnprocessableEntity, &pkg.JSONResponse{
			Code:    422,
			Message: "payload validation error",
			Error:   err.Error(),
			Data:    errFields,
		})
		return
	}

	Response, err := c.UsersUsecase.Login(req.Context(), requestBody)
	if err != nil {
		pkg.Response(res, http.StatusBadRequest, pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	pkg.Response(res, Response.Code, Response)
	return
}

// Profile func
func (c *UsersControllers) Profile(res http.ResponseWriter, req *http.Request) {
	TokenJSON := context.Get(req, "token").(string)

	var requestBody *entities.ProfileRequest
	if err := json.Unmarshal([]byte(TokenJSON), &requestBody); err != nil {
		log.WithFields(log.Fields{
			"event": "error when unmarshal token data",
		}).Error(err)
		pkg.Response(res, http.StatusBadRequest, &pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	Response, err := c.UsersUsecase.Profile(req.Context(), requestBody)
	if err != nil {
		pkg.Response(res, http.StatusBadRequest, pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	pkg.Response(res, Response.Code, Response)
	return
}

// UpdateProfile func
func (c *UsersControllers) UpdateProfile(res http.ResponseWriter, req *http.Request) {
	RawPayload, _ := ioutil.ReadAll(req.Body)
	RawPayloadString := string(RawPayload)
	RawPayloadString = strings.ReplaceAll(RawPayloadString, "\n", "")
	RawPayloadString = strings.ReplaceAll(RawPayloadString, "\\", "")
	RawPayloadString = strings.TrimSpace(RawPayloadString)
	// Log Raw Payload
	log.WithFields(log.Fields{
		"data": RawPayloadString,
	}).Info("POST /users/register payload body")

	var requestBody *entities.UpdateProfileRequest
	if err := json.Unmarshal(RawPayload, &requestBody); err != nil {
		log.WithFields(log.Fields{
			"event": "error when unmarshal request payload update profile",
		}).Error(err)
		pkg.Response(res, http.StatusBadRequest, &pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	TokenJSON := context.Get(req, "token").(string)
	var TokenData *entities.TokenClaim
	if err := json.Unmarshal([]byte(TokenJSON), &TokenData); err != nil {
		log.WithFields(log.Fields{
			"event": "error when unmarshal token",
		}).Error(err)
		pkg.Response(res, http.StatusBadRequest, &pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	requestBody.UserID = TokenData.UserID

	// Payload Validation
	if err := validate.Struct(requestBody); err != nil {
		errField := map[string]string{}
		errFields := []map[string]string{}

		for _, e := range err.(validator.ValidationErrors) {
			errField[e.Field()] = fmt.Sprintf("%s failed on the %s tag", e.Field(), e.Tag())
		}
		errFields = append(errFields, errField)

		log.WithFields(log.Fields{
			"event":            "payload validation error",
			"validation_error": errFields,
		})

		pkg.Response(res, http.StatusUnprocessableEntity, &pkg.JSONResponse{
			Code:    422,
			Message: "payload validation error",
			Error:   err.Error(),
			Data:    errFields,
		})
		return
	}

	Response, err := c.UsersUsecase.UpdateProfile(req.Context(), requestBody)
	if err != nil {
		pkg.Response(res, http.StatusBadRequest, pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	pkg.Response(res, Response.Code, Response)
	return
}

// UpdatePassword func
func (c *UsersControllers) UpdatePassword(res http.ResponseWriter, req *http.Request) {
	RawPayload, _ := ioutil.ReadAll(req.Body)
	RawPayloadString := string(RawPayload)
	RawPayloadString = strings.ReplaceAll(RawPayloadString, "\n", "")
	RawPayloadString = strings.ReplaceAll(RawPayloadString, "\\", "")
	RawPayloadString = strings.TrimSpace(RawPayloadString)
	// Log Raw Payload
	log.WithFields(log.Fields{
		"data": RawPayloadString,
	}).Info("POST /users/register payload body")

	var requestBody *entities.UpdatePasswordRequest
	if err := json.Unmarshal(RawPayload, &requestBody); err != nil {
		log.WithFields(log.Fields{
			"event": "error when unmarshal request payload update password",
		}).Error(err)
		pkg.Response(res, http.StatusBadRequest, &pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	TokenJSON := context.Get(req, "token").(string)
	var TokenData *entities.TokenClaim
	if err := json.Unmarshal([]byte(TokenJSON), &TokenData); err != nil {
		log.WithFields(log.Fields{
			"event": "error when unmarshal token",
		}).Error(err)
		pkg.Response(res, http.StatusBadRequest, &pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	requestBody.UserID = TokenData.UserID

	// Payload Validation
	if err := validate.Struct(requestBody); err != nil {
		errField := map[string]string{}
		errFields := []map[string]string{}

		for _, e := range err.(validator.ValidationErrors) {
			errField[e.Field()] = fmt.Sprintf("%s failed on the %s tag", e.Field(), e.Tag())
		}
		errFields = append(errFields, errField)

		log.WithFields(log.Fields{
			"event":            "payload validation error",
			"validation_error": errFields,
		})

		pkg.Response(res, http.StatusUnprocessableEntity, &pkg.JSONResponse{
			Code:    422,
			Message: "payload validation error",
			Error:   err.Error(),
			Data:    errFields,
		})
		return
	}

	Response, err := c.UsersUsecase.UpdatePassword(req.Context(), requestBody)
	if err != nil {
		pkg.Response(res, http.StatusBadRequest, pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	pkg.Response(res, Response.Code, Response)
	return
}
