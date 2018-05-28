package httphandler

import (
	"fmt"
)
import . "http/comm"


func HandlerTest1(handler HandlerHead){
	w := handler.Writer
	r := handler.Request
	querys := r.URL.Query()

	fmt.Fprintf(w,"<h1>Hello%s!</h1>",r.URL.Path[1:])
	Debug(handler,"here,querys:%v\n", querys)

	Debug(handler,"test:%s", GetValue(querys, "test"))

}
