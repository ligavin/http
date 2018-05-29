package handlebase

import . "http/httphandler"

func HandlerMapConfig()(map[string]interface{}){
	handlerMap := make(map[string]interface{})

	handlerMap["default"] = HandlerDefault
	handlerMap["handler1"] = Handler1

	handlerMap["upload"] = Upload

	return handlerMap
}