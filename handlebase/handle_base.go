package handlebase

import . "http/httphandler"

func HandlerMapConfig()(map[string]interface{}){
	handlerMap := make(map[string]interface{})


	handlerMap["handler1"] = Handler1
	handlerMap["default"] = HandlerDefault

	return handlerMap
}