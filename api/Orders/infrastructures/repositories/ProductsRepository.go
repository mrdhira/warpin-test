package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/mrdhira/warpin-test/api/Orders/entities"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// IProductsrepository interface
type IProductsrepository interface {
	GetProductsByID(ctx context.Context, Payload *entities.GetProductsByIDPayload) (Products *entities.Products, err error)
	ProductsUpdate(ctx context.Context, Payload *entities.ProductsUpdatePayload) (err error)
}

// ProductsRepository struct
type ProductsRepository struct {
}

// GetProductsByID func
func (r *ProductsRepository) GetProductsByID(ctx context.Context, Payload *entities.GetProductsByIDPayload) (Products *entities.Products, err error) {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Second * 60,
		}).Dial,
		TLSHandshakeTimeout: time.Second * 60,
	}

	var Client = &http.Client{
		Timeout:   time.Second * 60,
		Transport: netTransport,
	}

	BaseURL := viper.GetString("services.products.url")
	PathURL := "/products/" + strconv.Itoa(Payload.ProductID)

	RequestHTTP, err := http.NewRequest("GET", BaseURL+PathURL, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "err creating new request get product by id to product service",
		}).Error(err)
		return
	}

	RequestHTTP.Header.Set("Authorization", "INTERNAL-SERVICES")

	ResponseHTTP, err := Client.Do(RequestHTTP)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "err performing request get product by id to product service",
		}).Error(err)
		return
	}
	defer ResponseHTTP.Body.Close()

	ResponseBody, _ := ioutil.ReadAll(ResponseHTTP.Body)
	var GetProductsByIDResponse *entities.GetProductsByIDResponse
	json.Unmarshal(ResponseBody, &GetProductsByIDResponse)

	log.WithFields(log.Fields{
		"event": "response from product service for get products by id",
		"data":  string(ResponseBody),
	})

	if GetProductsByIDResponse.Code != 200 {
		return nil, errors.New(GetProductsByIDResponse.Message)
	}
	return GetProductsByIDResponse.Data, nil
}

// ProductsUpdate func
func (r *ProductsRepository) ProductsUpdate(ctx context.Context, Payload *entities.ProductsUpdatePayload) (err error) {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Second * 60,
		}).Dial,
		TLSHandshakeTimeout: time.Second * 60,
	}

	var Client = &http.Client{
		Timeout:   time.Second * 60,
		Transport: netTransport,
	}

	BaseURL := viper.GetString("services.products.url")
	PathURL := "/products/internal/" + strconv.Itoa(Payload.ProductID)

	RequestBody, err := json.Marshal(Payload)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "err when marshal payload products update to product service",
		}).Error(err)
		return
	}

	RequestHTTP, err := http.NewRequest("PUT", BaseURL+PathURL, bytes.NewBuffer(RequestBody))
	if err != nil {
		log.WithFields(log.Fields{
			"event": "err creating new request products update to product service",
		}).Error(err)
		return
	}

	RequestHTTP.Header.Set("Authorization", "INTERNAL-SERVICES")

	ResponseHTTP, err := Client.Do(RequestHTTP)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "err performing request products update to product service",
		}).Error(err)
		return
	}
	defer ResponseHTTP.Body.Close()

	ResponseBody, _ := ioutil.ReadAll(ResponseHTTP.Body)
	var ProductsUpdateResponse *entities.ProductsUpdateResponse
	json.Unmarshal(ResponseBody, &ProductsUpdateResponse)

	log.WithFields(log.Fields{
		"event": "response from product service for update products",
		"data":  string(ResponseBody),
	})

	if ProductsUpdateResponse.Code != 200 {
		return errors.New(ProductsUpdateResponse.Message)
	}
	return
}
