package main

import (
	"fmt"
	"encoding/json"
)

func main(){
	mp := map[interface{}]interface{}{}

	mp["1"]=2
	mp[2] = "1"

	data, err := json.Marshal(mp)

	fmt.Printf("mp:%v\ndata:%s,err:%s,%v", mp,string(data), err, mp["1"])

}
