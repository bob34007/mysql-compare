package utils

import "hash/fnv"

func HashString(s string) (uint64,error) {
	h := fnv.New64a()
	ss := Slice(s)
	_,err:=h.Write(ss)
	if err!=nil{
		return 0,err
	}
	return h.Sum64(),nil
}
