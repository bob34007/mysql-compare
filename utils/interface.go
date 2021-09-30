package utils

import (
	"fmt"
	"reflect"
)


//Compares whether the interface values are equal
func CompareInterface(a,b interface{}) bool{
	t1 := reflect.ValueOf(a).Type().String()
	t2 := reflect.ValueOf(b).Type().String()
	if t1!=t2{
		logStr := fmt.Sprintf("import param type is not same , %v-%v ,%v-%v" , t1,t2,a,b)
		log.Warn(logStr)
		return false
	}
	switch a.(type){
	case int :
		return a.(int)==b.(int)
	case uint :
		return a.(uint)==b.(uint)
	case int32 :
		return a.(int32)==b.(int32)
	case uint32 :
		return a.(uint32)==b.(uint32)
	case int64:
		return a.(int64) == b.(int64)
	case uint64:
		return a.(uint64)==b.(uint64)
	case string :
		return a.(string)==b.(string)
	default :
		return false
	}
}


