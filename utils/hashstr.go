package utils

import "hash/fnv"

func hashString(s string) (uint64,error) {
	h := fnv.New64a()
	ss := Slice(s)
	_,err:=h.Write(ss)
	//fmt.Println(err)
	if err!=nil{
		return 0,err
	}
	return h.Sum64(),nil
}
