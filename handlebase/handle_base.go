package handlebase

import . "http/httphandler"

func HandlerMapConfig()(map[string]interface{}){
	handlerMap := make(map[string]interface{})


	handlerMap["HandlerTest1"] = HandlerTest1
	handlerMap["default"] = HandlerDefault

	return handlerMap
}