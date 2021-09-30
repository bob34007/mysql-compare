package utils

import (
	"sort"
)

func Sort2DSlice(a [][]interface{}) {
	var less = func(i, j int) bool {
		m := a[i][0]
		n := a[j][0]
		if m ==nil{
			return true
		}
		if n ==nil {
			return false
		}
		switch m.(type) {
		case int:
			return m.(int) < n.(int)
		case uint:
			return m.(uint) < n.(uint)
		case int32:
			return m.(int32) < n.(int32)
		case uint32:
			return m.(uint32) < n.(uint32)
		case int64:
			return m.(int64) < n.(int64)
		case uint64:
			return m.(uint64) < n.(uint64)
		case string:
			return m.(string) < n.(string)
		default:
			//fmt.Println("unsport type : " + reflect.ValueOf(m).Type().String())
			//Unsupported types so far we return true
			return false
		}
		return true
	}
	sort.Slice(a, less)
}
