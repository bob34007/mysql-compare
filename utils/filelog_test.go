package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUtil_Logger_With_Name_Len_Large_Zero (t *testing.T){
	var fn  =FileName("abc")
	logger := fn.Logger()
	ast :=assert.New(t)
	ast.NotNil(logger)
}

func TestUtil_Logger_With_Name_Len_Equal_Zero (t *testing.T){
	var fn  =FileName("")
	logger := fn.Logger()
	ast :=assert.New(t)
	ast.NotNil(logger)
}