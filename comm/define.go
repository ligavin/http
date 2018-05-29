package comm

import (
	"log"
	"net/http"
)

type Handler struct {
	Logger 		*log.Logger
	Writer 		http.ResponseWriter
	Request 	*http.Request
	Seq   		uint32
	QueryParams map[string]string
}

