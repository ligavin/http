package comm

import (
	"log"
	"net/http"
	"net/url"
)

type HandlerHead struct {
	Logger 		*log.Logger
	Writer 		http.ResponseWriter
	Request 	*http.Request
	Query 		url.Values
	Seq   		uint32
}

