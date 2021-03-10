package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/mrdhira/warpin-test/api/Orders/entities"
	"github.com/mrdhira/warpin-test/api/Orders/usecases"
	"github.com/mrdhira/warpin-test/pkg"
	log "github.com/sirupsen/logrus"
)

var validate *validator.Validate

// OrdersControllers struct
type OrdersControllers struct {
	OrdersUsecase usecases.IOrdersUsecases
}

// InitOrdersControllers func
func InitOrdersControllers() *OrdersControllers {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Init Usecase
	ordersUsecases := usecases.InitOrdersUsecases()

	return &OrdersControllers{
		OrdersUsecase: ordersUsecases,
	}
}

// OrdersListUsers func
func (c *OrdersControllers) OrdersListUsers(res http.ResponseWriter, req *http.Request) {
	var requestBody *entities.OrdersListUsersRequest

	if req.URL.Query().Get("limit") == "" && req.URL.Query().Get("offset") == "" {
		pkg.Response(res, http.StatusUnprocessableEntity, &pkg.JSONResponse{
			Code:    422,
			Message: "Limit dan Offset tidak bisa kosong",
		})
		return
	}

	requestBody = &entities.OrdersListUsersRequest{
		Limit:  req.URL.Query().Get("limit"),
		Offset: req.URL.Query().Get("offset"),
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

	Response, err := c.OrdersUsecase.OrdersListUsers(req.Context(), requestBody)
	if err != nil {
		pkg.Response(res, http.StatusBadRequest, &pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	pkg.Response(res, Response.Code, Response)
	return
}

// OrdersCreate func
func (c *OrdersControllers) OrdersCreate(res http.ResponseWriter, req *http.Request) {
	RawPayload, _ := ioutil.ReadAll(req.Body)
	RawPayloadString := string(RawPayload)
	RawPayloadString = strings.ReplaceAll(RawPayloadString, "\n", "")
	RawPayloadString = strings.ReplaceAll(RawPayloadString, "\\", "")
	RawPayloadString = strings.TrimSpace(RawPayloadString)
	// Log Raw Payload
	log.WithFields(log.Fields{
		"data": RawPayloadString,
	}).Info("POST /orders payload body")

	var requestBody *entities.OrdersCreateRequest
	if err := json.Unmarshal(RawPayload, &requestBody); err != nil {
		log.WithFields(log.Fields{
			"event": "error when unmarshal request payload orders create",
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

	Response, err := c.OrdersUsecase.OrdersCreate(req.Context(), requestBody)
	if err != nil {
		pkg.Response(res, http.StatusBadRequest, &pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	pkg.Response(res, Response.Code, Response)
	return
}

// OrdersUpdate func
func (c *OrdersControllers) OrdersUpdate(res http.ResponseWriter, req *http.Request) {
	RawPayload, _ := ioutil.ReadAll(req.Body)
	RawPayloadString := string(RawPayload)
	RawPayloadString = strings.ReplaceAll(RawPayloadString, "\n", "")
	RawPayloadString = strings.ReplaceAll(RawPayloadString, "\\", "")
	RawPayloadString = strings.TrimSpace(RawPayloadString)
	// Log Raw Payload
	log.WithFields(log.Fields{
		"data": RawPayloadString,
	}).Info("PUT /orders/{id} payload body")

	var requestBody *entities.OrdersUpdateRequest
	if err := json.Unmarshal(RawPayload, &requestBody); err != nil {
		log.WithFields(log.Fields{
			"event": "error when unmarshal request payload update products",
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

	OrderID, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		log.WithFields(log.Fields{
			"event": "error when get order id from params",
		}).Error(err)
		pkg.Response(res, http.StatusBadRequest, &pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	requestBody.OrderID = OrderID

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

	Response, err := c.OrdersUsecase.OrdersUpdate(req.Context(), requestBody)
	if err != nil {
		pkg.Response(res, http.StatusBadRequest, &pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	pkg.Response(res, Response.Code, Response)
	return
}

// OrdersCancel func
func (c *OrdersControllers) OrdersCancel(res http.ResponseWriter, req *http.Request) {
	RawPayload, _ := ioutil.ReadAll(req.Body)
	RawPayloadString := string(RawPayload)
	RawPayloadString = strings.ReplaceAll(RawPayloadString, "\n", "")
	RawPayloadString = strings.ReplaceAll(RawPayloadString, "\\", "")
	RawPayloadString = strings.TrimSpace(RawPayloadString)
	// Log Raw Payload
	log.WithFields(log.Fields{
		"data": RawPayloadString,
	}).Info("PUT /orders/{id}/cancel payload body")

	var requestBody *entities.OrdersCancelRequest

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

	requestBody = &entities.OrdersCancelRequest{
		UserID: TokenData.UserID,
	}

	OrderID, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		log.WithFields(log.Fields{
			"event": "error when get order id from params",
		}).Error(err)
		pkg.Response(res, http.StatusBadRequest, &pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	requestBody.OrderID = OrderID

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

	Response, err := c.OrdersUsecase.OrdersCancel(req.Context(), requestBody)
	if err != nil {
		pkg.Response(res, http.StatusBadRequest, &pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	pkg.Response(res, Response.Code, Response)
	return
}

// OrdersListAdmin func
func (c *OrdersControllers) OrdersListAdmin(res http.ResponseWriter, req *http.Request) {
	var requestBody *entities.OrdersListAdminRequest

	if req.URL.Query().Get("limit") == "" && req.URL.Query().Get("offset") == "" {
		pkg.Response(res, http.StatusUnprocessableEntity, &pkg.JSONResponse{
			Code:    422,
			Message: "Limit dan Offset tidak bisa kosong",
		})
		return
	}

	requestBody = &entities.OrdersListAdminRequest{
		Limit:  req.URL.Query().Get("limit"),
		Offset: req.URL.Query().Get("offset"),
	}

	if req.URL.Query().Get("status") != "" {
		Status, err := strconv.Atoi(req.URL.Query().Get("status"))
		if err != nil {
			log.WithFields(log.Fields{
				"event": "error when parse to int for status query params",
			}).Error(err)
			pkg.Response(res, http.StatusBadRequest, &pkg.JSONResponse{
				Code:    400,
				Message: "Terjadi kesalahan sistem",
				Error:   err.Error(),
			})
			return
		}
		requestBody.Status = Status
	}

	if req.URL.Query().Get("product_id") != "" {
		ProductID, err := strconv.Atoi(req.URL.Query().Get("product_id"))
		if err != nil {
			log.WithFields(log.Fields{
				"event": "error when parse to int for product_id query params",
			}).Error(err)
			pkg.Response(res, http.StatusBadRequest, &pkg.JSONResponse{
				Code:    400,
				Message: "Terjadi kesalahan sistem",
				Error:   err.Error(),
			})
			return
		}
		requestBody.ProductID = ProductID
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

	Response, err := c.OrdersUsecase.OrdersListAdmin(req.Context(), requestBody)
	if err != nil {
		pkg.Response(res, http.StatusBadRequest, &pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	pkg.Response(res, Response.Code, Response)
	return
}

// OrdersApprove func
func (c *OrdersControllers) OrdersApprove(res http.ResponseWriter, req *http.Request) {
	RawPayload, _ := ioutil.ReadAll(req.Body)
	RawPayloadString := string(RawPayload)
	RawPayloadString = strings.ReplaceAll(RawPayloadString, "\n", "")
	RawPayloadString = strings.ReplaceAll(RawPayloadString, "\\", "")
	RawPayloadString = strings.TrimSpace(RawPayloadString)
	// Log Raw Payload
	log.WithFields(log.Fields{
		"data": RawPayloadString,
	}).Info("PUT /orders/internal/{id}/approve payload body")

	var requestBody *entities.OrdersApproveRequest

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

	requestBody = &entities.OrdersApproveRequest{
		UserID: TokenData.UserID,
	}

	OrderID, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		log.WithFields(log.Fields{
			"event": "error when get order id from params",
		}).Error(err)
		pkg.Response(res, http.StatusBadRequest, &pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	requestBody.OrderID = OrderID

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

	Response, err := c.OrdersUsecase.OrdersApprove(req.Context(), requestBody)
	if err != nil {
		pkg.Response(res, http.StatusBadRequest, &pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	pkg.Response(res, Response.Code, Response)
	return
}

// OrdersReject func
func (c *OrdersControllers) OrdersReject(res http.ResponseWriter, req *http.Request) {
	RawPayload, _ := ioutil.ReadAll(req.Body)
	RawPayloadString := string(RawPayload)
	RawPayloadString = strings.ReplaceAll(RawPayloadString, "\n", "")
	RawPayloadString = strings.ReplaceAll(RawPayloadString, "\\", "")
	RawPayloadString = strings.TrimSpace(RawPayloadString)
	// Log Raw Payload
	log.WithFields(log.Fields{
		"data": RawPayloadString,
	}).Info("PUT /orders/internal/{id}/reject payload body")

	var requestBody *entities.OrdersRejectRequest

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

	requestBody = &entities.OrdersRejectRequest{
		UserID: TokenData.UserID,
	}

	OrderID, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		log.WithFields(log.Fields{
			"event": "error when get order id from params",
		}).Error(err)
		pkg.Response(res, http.StatusBadRequest, &pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	requestBody.OrderID = OrderID

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

	Response, err := c.OrdersUsecase.OrdersReject(req.Context(), requestBody)
	if err != nil {
		pkg.Response(res, http.StatusBadRequest, &pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	pkg.Response(res, Response.Code, Response)
	return
}
