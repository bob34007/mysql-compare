package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)





func TestUtil_DataCompare_Int_Equal(t *testing.T) {

	var a interface{} = int(1)
	var b interface{} = int(1)

	ok:=CompareInterface(a,b)

	ast := assert.New(t)
	ast.Equal(ok,true )
}

func TestUtil_DataCompare_Int_Not_Equal(t *testing.T) {

	var a interface{} = int(1)
	var b interface{} = int(2)

	ok:=CompareInterface(a,b)

	ast := assert.New(t)
	ast.Equal(ok,false )
}

func TestUtil_DataCompare_Uint_Equal(t *testing.T) {

	var a interface{} = uint(1)
	var b interface{} = uint(1)

	ok:=CompareInterface(a,b)

	ast := assert.New(t)
	ast.Equal(ok,true )
}

func TestUtil_DataCompare_Uint_Not_Equal(t *testing.T) {

	var a interface{} = uint(1)
	var b interface{} = uint(2)

	ok:=CompareInterface(a,b)

	ast := assert.New(t)
	ast.Equal(ok,false )
}

func TestUtil_DataCompare_Int32_Equal(t *testing.T) {

	var a interface{} = int32(1)
	var b interface{} = int32(1)

	ok:=CompareInterface(a,b)

	ast := assert.New(t)
	ast.Equal(ok,true )
}

func TestUtil_DataCompare_Int32_Not_Equal(t *testing.T) {

	var a interface{} = int32(1)
	var b interface{} = int32(2)

	ok:=CompareInterface(a,b)

	ast := assert.New(t)
	ast.Equal(ok,false )
}


func TestUtil_DataCompare_Uint32_Equal(t *testing.T) {

	var a interface{} = uint32(1)
	var b interface{} = uint32(1)

	ok:=CompareInterface(a,b)

	ast := assert.New(t)
	ast.Equal(ok,true )
}

func TestUtil_DataCompare_Uint32_Not_Equal(t *testing.T) {

	var a interface{} = uint32(1)
	var b interface{} = uint32(2)

	ok:=CompareInterface(a,b)

	ast := assert.New(t)
	ast.Equal(ok,false )
}


func TestUtil_DataCompare_Int64_Equal(t *testing.T) {

	var a interface{} = int64(1)
	var b interface{} = int64(1)

	ok:=CompareInterface(a,b)

	ast := assert.New(t)
	ast.Equal(ok,true )
}

func TestUtil_DataCompare_Int64_Not_Equal(t *testing.T) {

	var a interface{} = int64(1)
	var b interface{} = int64(2)

	ok:=CompareInterface(a,b)

	ast := assert.New(t)
	ast.Equal(ok,false )
}


func TestUtil_DataCompare_Uint64_Equal(t *testing.T) {

	var a interface{} = uint64(1)
	var b interface{} = uint64(1)

	ok:=CompareInterface(a,b)

	ast := assert.New(t)
	ast.Equal(ok,true )
}

func TestUtil_DataCompare_Uint64_Not_Equal(t *testing.T) {

	var a interface{} = uint64(1)
	var b interface{} = uint64(2)

	ok:=CompareInterface(a,b)

	ast := assert.New(t)
	ast.Equal(ok,false )
}


func TestUtil_DataCompare_String_Equal(t *testing.T) {

	var a interface{} = "abc"
	var b interface{} = "abc"

	ok:=CompareInterface(a,b)

	ast := assert.New(t)
	ast.Equal(ok,true )
}

func TestUtil_DataCompare_String_Not_Equal(t *testing.T) {

	var a interface{} = "abc"
	var b interface{} = "ABC"

	ok:=CompareInterface(a,b)

	ast := assert.New(t)
	ast.Equal(ok,false )
}

func TestUtil_DataCompare_Type_Not_Equal(t *testing.T) {

	var a interface{} = "abc"
	var b interface{} = uint16(1)

	ok:=CompareInterface(a,b)

	ast := assert.New(t)
	ast.Equal(ok,false )
}

//Unsupported type program functions return interface values that are not equal
func TestUtil_DataCompare_Uint16_Equal(t *testing.T) {

	var a interface{} = uint16(1)
	var b interface{} = uint16(1)

	ok:=CompareInterface(a,b)

	ast := assert.New(t)
	ast.Equal(ok,false )
}

func TestUtil_DataCompare_Uint16_Not_Equal(t *testing.T) {

	var a interface{} = uint16(1)
	var b interface{} = uint16(2)

	ok:=CompareInterface(a,b)

	ast := assert.New(t)
	ast.Equal(ok,false )
}



