package server

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Data interface{} `json:"data"`
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
}

func ResponseOK(w http.ResponseWriter, data interface{}, msg string) {
	response := Response{
		Data: data,
		Code: "0000",
		Msg:  msg,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func ResponseErr(w http.ResponseWriter, code string, msg string, httpCode int) {
	response := Response{
		Data: "",
		Code: code,
		Msg:  msg,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(response)
}
