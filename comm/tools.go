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

func GetValue(mp map[string]string, key string)string{
	v,exist := mp[key]
	if !exist{
		return ""
	}
	return v
}