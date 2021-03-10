package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/mrdhira/warpin-test/api/Products/entities"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// IOrdersRepository interface
type IOrdersRepository interface {
	GetOrdersByProductID(ctx context.Context, Payload *entities.GetOrdersByPrductIDPayload) (Orders []*entities.Orders, err error)
}

// OrdersRepository struct
type OrdersRepository struct {
}

// GetOrdersByProductID func
func (r *OrdersRepository) GetOrdersByProductID(ctx context.Context, Payload *entities.GetOrdersByPrductIDPayload) (Orders []*entities.Orders, err error) {
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

	BaseURL := viper.GetString("services.orders.url")
	PathURL := "/orders/internal/"

	RequestHTTP, err := http.NewRequest("GET", BaseURL+PathURL, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "err creating new request get order by product id to order service",
		}).Error(err)
		return
	}

	RequestHTTP.Header.Set("Authorization", "INTERNAL-SERVICES")

	QueryParams := RequestHTTP.URL.Query()
	QueryParams.Add("limit", strconv.Itoa(Payload.Limit))
	QueryParams.Add("offset", strconv.Itoa(Payload.Offset))

	if Payload.Status != 0 {
		QueryParams.Add("status", strconv.Itoa(Payload.Status))
	}

	if Payload.ProductID != 0 {
		QueryParams.Add("product_id", strconv.Itoa(Payload.ProductID))
	}

	RequestHTTP.URL.RawQuery = QueryParams.Encode()

	ResponseHTTP, err := Client.Do(RequestHTTP)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "err performing request get order by product id to order service",
		}).Error(err)
		return
	}
	defer ResponseHTTP.Body.Close()

	ResponseBody, _ := ioutil.ReadAll(ResponseHTTP.Body)
	var GetOrdersByPrductIDResponse *entities.GetOrdersByPrductIDResponse
	json.Unmarshal(ResponseBody, &GetOrdersByPrductIDResponse)

	log.WithFields(log.Fields{
		"event": "response from orders service for get orders list",
		"data":  string(ResponseBody),
	})

	if GetOrdersByPrductIDResponse.Code != 200 {
		return nil, errors.New(GetOrdersByPrductIDResponse.Message)
	}
	return GetOrdersByPrductIDResponse.Data, nil
}
