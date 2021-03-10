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
	"github.com/mrdhira/warpin-test/api/Products/entities"
	"github.com/mrdhira/warpin-test/api/Products/usecases"
	"github.com/mrdhira/warpin-test/pkg"
	log "github.com/sirupsen/logrus"
)

var validate *validator.Validate

// ProductsControllers struct
type ProductsControllers struct {
	ProductsUsecase usecases.IProductsUsecases
}

// InitProductsControllers func
func InitProductsControllers() *ProductsControllers {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Init Usecase
	productsUsecases := usecases.InitProductsUsecases()

	return &ProductsControllers{
		ProductsUsecase: productsUsecases,
	}
}

// GetProducts func
func (c *ProductsControllers) GetProducts(res http.ResponseWriter, req *http.Request) {
	var requestBody *entities.GetProductsRequest

	if req.URL.Query().Get("limit") == "" && req.URL.Query().Get("offset") == "" {
		pkg.Response(res, http.StatusUnprocessableEntity, &pkg.JSONResponse{
			Code:    422,
			Message: "Limit dan Offset tidak bisa kosong",
		})
		return
	}

	requestBody = &entities.GetProductsRequest{
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

	Response, err := c.ProductsUsecase.GetProducts(req.Context(), requestBody)
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

// GetProductsByID func
func (c *ProductsControllers) GetProductsByID(res http.ResponseWriter, req *http.Request) {
	var requestBody *entities.GetProductsByIDRequest

	ProductID, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		log.WithFields(log.Fields{
			"event": "error when get product id from params",
		}).Error(err)
		pkg.Response(res, http.StatusBadRequest, &pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	requestBody = &entities.GetProductsByIDRequest{
		ProductID: ProductID,
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

	Response, err := c.ProductsUsecase.GetProductsByID(req.Context(), requestBody)
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

// AddProducts func
func (c *ProductsControllers) AddProducts(res http.ResponseWriter, req *http.Request) {
	RawPayload, _ := ioutil.ReadAll(req.Body)
	RawPayloadString := string(RawPayload)
	RawPayloadString = strings.ReplaceAll(RawPayloadString, "\n", "")
	RawPayloadString = strings.ReplaceAll(RawPayloadString, "\\", "")
	RawPayloadString = strings.TrimSpace(RawPayloadString)
	// Log Raw Payload
	log.WithFields(log.Fields{
		"data": RawPayloadString,
	}).Info("POST /products/internal/add payload body")

	var requestBody *entities.AddProductsRequest
	if err := json.Unmarshal(RawPayload, &requestBody); err != nil {
		log.WithFields(log.Fields{
			"event": "error when unmarshal request payload add products",
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

	Response, err := c.ProductsUsecase.AddProducts(req.Context(), requestBody)
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

// UpdateProducts func
func (c *ProductsControllers) UpdateProducts(res http.ResponseWriter, req *http.Request) {
	RawPayload, _ := ioutil.ReadAll(req.Body)
	RawPayloadString := string(RawPayload)
	RawPayloadString = strings.ReplaceAll(RawPayloadString, "\n", "")
	RawPayloadString = strings.ReplaceAll(RawPayloadString, "\\", "")
	RawPayloadString = strings.TrimSpace(RawPayloadString)
	// Log Raw Payload
	log.WithFields(log.Fields{
		"data": RawPayloadString,
	}).Info("PUT /products/internal/{id} payload body")

	var requestBody *entities.UpdateProductsRequest
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

	ProductID, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		log.WithFields(log.Fields{
			"event": "error when get product id from params",
		}).Error(err)
		pkg.Response(res, http.StatusBadRequest, &pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		})
		return
	}

	requestBody.ProductID = ProductID

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

	Response, err := c.ProductsUsecase.UpdateProducts(req.Context(), requestBody)
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
