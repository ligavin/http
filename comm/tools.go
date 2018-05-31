package comm

import (
	"math/rand"
	"time"
)


func InitRand(){
	rand.Seed(time.Now().UnixNano())
}

func GetRandUint32()(uint32){
	InitRand()
	return rand.Uint32()
}

func GetValue(handler Handler, key string)string{

	mp := handler.QueryParams

	v,exist := mp[key]
	if !exist{
		return ""
	}
	return v
}

func Result(ret int, msg string)map[string]interface{}{
	resMap := make(map[string]interface{})
	resMap["ret"]=ret
	resMap["msg"]=msg
	return resMap
}