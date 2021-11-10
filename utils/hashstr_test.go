package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUtil_hashString_Succ (t *testing.T){

	var s = "abcdef"

	_,err:= HashString(s)

	ast := assert.New(t)

	ast.Nil(err)

}

/*
func TestUtil_hashString_With_Write_Fail (t *testing.T){

	var s = "abcdef"

	err1:=errors.New("write data fail , no space left")


	ctrl:=gomock.NewController(t)
	defer ctrl.Finish()

	mocker:=NewMockHash64(ctrl)



	mocker.EXPECT().Write(gomock.Any()).Return(0,err1).AnyTimes()


	n,err:= hashString(s)

	fmt.Println(err,err1)

	ast := assert.New(t)

	ast.Equal(err1,err)
	ast.Equal(n,uint64(0))

}
*/