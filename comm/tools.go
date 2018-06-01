package comm

import (
	"math/rand"
	"time"
	"os"
)

var randNumber int32 = 0

func InitRand(){
	pid := (int64)(os.Getpid()) * 1e10
	rand.Seed(time.Now().UnixNano() + ((int64)(randNumber)*1e9) + pid)
	randNumber++
}

func GetRandUint32()(uint32){
	InitRand()
	return rand.Uint32()
}

func GetRandUint64()(uint64){
	InitRand()
	return rand.Uint64()
}

func GetMapValue(mp	map[string]string, key string)string{
	v, exist := mp[key]

	if !exist{
		return ""
	}

	return v
}

func GetUrlValue(handler Handler, key string)string{

	mp := handler.QueryParams

	return GetMapValue(mp, key)
}

func Result(ret int, msg string)map[string]interface{}{
	resMap := make(map[string]interface{})
	resMap["ret"]=ret
	resMap["msg"]=msg
	return resMap
}

