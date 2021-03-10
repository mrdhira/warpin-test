package pkg

import (
	"encoding/json"
	"net/http"
)

// Response func
func Response(res http.ResponseWriter, HttpStatus int, Data interface{}) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(HttpStatus)
	json.NewEncoder(res).Encode(Data)
}

// JSONResponse struct
type JSONResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
}
