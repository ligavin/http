package comm

import (
	"net/url"
	"math/rand"
	"time"
)

func GetValue(query url.Values, key string)(string){
	value, exist := query[key]

	if (!exist) {
		return ""
	}
	return value[0]
}

func InitRand(){
	rand.Seed(time.Now().UnixNano())
}

func GetRandUint32()(uint32){
	InitRand()
	return rand.Uint32()
}