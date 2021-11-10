package compare

import (
	"encoding/json"
	"github.com/agiledragon/gomonkey"
	"github.com/bobguo/mysql-compare/utils"
	"github.com/pingcap/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"os"
	"reflect"
	"testing"
	"time"
)

var log =  zap.L().With(zap.String("compare", "file"))

func TestUtils_HashPrResDetail_Succ(t *testing.T) {

	rs := new(ResFromFile)
	rs.PrResult = make([][]string, 0)
	rowRes1 := make([]string, 0)

	rowRes1 = append(rowRes1, "hello", "word", "program")

	rowRes2 := make([]string, 0)
	rowRes2 = append(rowRes2, "hello1", "word1", "program1")
	rs.PrResult = append(rs.PrResult, rowRes1, rowRes2)

	v, err := rs.HashPrResDetail()

	ast := assert.New(t)

	ast.GreaterOrEqual(len(v), 0)
	ast.Nil(err)

}

func TestUtils_HashPrResDetail_Fail(t *testing.T) {

	rs := new(ResFromFile)
	rs.PrResult = make([][]string, 0)

	rowRes1 := make([]string, 0)
	rowRes1 = append(rowRes1, "hello", "word", "program")

	rowRes2 := make([]string, 0)
	rowRes2 = append(rowRes2, "hello1", "word1", "program1")

	rs.PrResult = append(rs.PrResult, rowRes1, rowRes2)

	err1 := errors.New("hash string fail")

	patch := gomonkey.ApplyFunc(utils.HashString, func(s string) (uint64, error) {
		return 0, err1
	})
	defer patch.Reset()

	v, err := rs.HashPrResDetail()

	ast := assert.New(t)

	ast.Nil(v)
	ast.Equal(err, err1)

}

func TestUtils_HashRrResDetail_Succ(t *testing.T) {

	rs := new(ResFromFile)
	rs.RrResult = make([][]string, 0)
	rowRes1 := make([]string, 0)

	rowRes1 = append(rowRes1, "hello", "word", "program")

	rowRes2 := make([]string, 0)
	rowRes2 = append(rowRes2, "hello1", "word1", "program1")
	rs.RrResult = append(rs.RrResult, rowRes1, rowRes2)

	v, err := rs.HashRrResDetail()

	ast := assert.New(t)

	ast.GreaterOrEqual(len(v), 0)
	ast.Nil(err)

}

func TestUtils_HashRrResDetail_Fail(t *testing.T) {

	rs := new(ResFromFile)
	rs.RrResult = make([][]string, 0)

	rowRes1 := make([]string, 0)
	rowRes1 = append(rowRes1, "hello", "word", "program")

	rowRes2 := make([]string, 0)
	rowRes2 = append(rowRes2, "hello1", "word1", "program1")

	rs.RrResult = append(rs.RrResult, rowRes1, rowRes2)

	err1 := errors.New("hash string fail")

	patch := gomonkey.ApplyFunc(utils.HashString, func(s string) (uint64, error) {
		return 0, err1
	})
	defer patch.Reset()

	v, err := rs.HashRrResDetail()

	ast := assert.New(t)

	ast.Nil(v)
	ast.Equal(err, err1)

}

func TestUtils_CompareData_Succ(t *testing.T) {
	a := make([][]interface{}, 0)
	b := make([][]interface{}, 0)

	row1 := make([]interface{}, 0)
	row1 = append(row1, uint64(10), uint64(1))

	row2 := make([]interface{}, 0)
	row2 = append(row2, uint64(10), uint64(1))

	a = append(a, row1, row2)
	b = append(b, row1, row2)

	rs := new(ResFromFile)
	rs.Logger = log
	err := rs.CompareData(a, b)

	ast := assert.New(t)
	ast.Nil(err)

}

func TestUtils_CompareData_Fail(t *testing.T) {
	a := make([][]interface{}, 0)
	b := make([][]interface{}, 0)

	row1 := make([]interface{}, 0)
	row1 = append(row1, uint64(10), uint64(1))

	row2 := make([]interface{}, 0)
	row2 = append(row2, uint64(11), uint64(1))

	a = append(a, row1)
	b = append(b, row2)

	rs := new(ResFromFile)
	rs.Logger = log

	err := rs.CompareData(a, b)

	ast := assert.New(t)
	ast.NotNil(err)

}

func TestUtils_CompareResDetail_Succ(t *testing.T) {

	rs := new(ResFromFile)
	rs.RrResult = make([][]string, 0)

	rowRes1 := make([]string, 0)
	rowRes1 = append(rowRes1, "hello", "word", "program")

	rowRes2 := make([]string, 0)
	rowRes2 = append(rowRes2, "hello1", "word1", "program1")

	rs.RrResult = append(rs.RrResult, rowRes1, rowRes2)

	rs.RrResult = make([][]string, 0)
	rowRes11 := make([]string, 0)

	rowRes11 = append(rowRes11, "hello", "word", "program")

	rowRes12 := make([]string, 0)
	rowRes12 = append(rowRes12, "hello1", "word1", "program1")
	rs.RrResult = append(rs.RrResult, rowRes11, rowRes12)

	rs.Logger = log

	ok, err := rs.CompareResDetail()

	ast := assert.New(t)

	ast.Nil(err)
	ast.Equal(ok, true)

}

func TestUtils_CompareResDetail_With_HashPrResDetail_Fail(t *testing.T) {

	rs := new(ResFromFile)
	rs.Logger = log

	err1 := errors.New("hash string fail")
	patches := gomonkey.ApplyMethod(reflect.TypeOf(rs), "HashPrResDetail",
		func(_ *ResFromFile) ([][]interface{}, error) {
			return nil, err1
		})
	defer patches.Reset()

	ok, err := rs.CompareResDetail()

	ast := assert.New(t)

	ast.Equal(err, err1)
	ast.Equal(ok, false)

}

func TestUtils_CompareResDetail_With_HashRrResDetail_Fail(t *testing.T) {

	rs := new(ResFromFile)
	rs.Logger = log

	err1 := errors.New("hash string fail")
	patches := gomonkey.ApplyMethod(reflect.TypeOf(rs), "HashPrResDetail",
		func(_ *ResFromFile) ([][]interface{}, error) {
			return nil, nil
		})
	defer patches.Reset()

	patches1 := gomonkey.ApplyMethod(reflect.TypeOf(rs), "HashRrResDetail",
		func(_ *ResFromFile) ([][]interface{}, error) {
			return nil, err1
		})
	defer patches1.Reset()

	ok, err := rs.CompareResDetail()

	ast := assert.New(t)

	ast.Equal(err, err1)
	ast.Equal(ok, false)

}

func TestUtils_CompareResDetail_With_CompareData_Fail(t *testing.T) {

	rs := new(ResFromFile)
	rs.Logger = log

	rs.RrResult = make([][]string, 0)

	rowRes1 := make([]string, 0)
	rowRes1 = append(rowRes1, "hello", "word", "program")

	rowRes2 := make([]string, 0)
	rowRes2 = append(rowRes2, "hello1", "word1", "program1")

	rs.RrResult = append(rs.RrResult, rowRes1, rowRes2)

	rs.RrResult = make([][]string, 0)
	rowRes11 := make([]string, 0)

	rowRes11 = append(rowRes11, "hello", "word", "program")

	rowRes12 := make([]string, 0)
	rowRes12 = append(rowRes12, "hello1", "word1", "program1")

	rs.RrResult = append(rs.RrResult, rowRes11, rowRes12)

	err1 := errors.New("compare row string hash value not equal")

	patches2 := gomonkey.ApplyMethod(reflect.TypeOf(rs), "CompareData",
		func(_ *ResFromFile, a, b [][]interface{}) error {
			return err1
		})
	defer patches2.Reset()

	ok, err := rs.CompareResDetail()

	ast := assert.New(t)

	ast.Equal(err, err1)
	ast.Equal(ok, false)
}

/*
func TestUtils_GetResFromFile_With_ReadBytes_Fail (t *testing.T){

	rs := new(ResFromFile)
	rs.Logger = logger

	file := new(os.File)
	rs.File = file

	err1:= errors.New("io.EOF")

	s,err := rs.GetResFromFile()

	ast :=assert.New(t)
	ast.Equal(err,err1)
	ast.Nil(s)
}
*/
/*
func TestUtils_GetResFromFile_Succ (t *testing.T){

	rs := new(ResFromFile)
	rs.Logger = logger

	file := new(os.File)
	rs.File = file

	ss := make([]byte,0)
	var a =uint64(8)
	var b ="abcdefgh"

	l := make([]byte, 8)
	binary.BigEndian.PutUint64(l, uint64(a))
	ss=append(ss,l...)
	ss=append(ss,[]byte(b)...)
	ss=append(ss,'\n')


	s,err := rs.GetResFromFile()

	ast :=assert.New(t)
	ast.Equal(len(s),8)
	ast.Nil(err)
}
*/

func TestUtils_NewResForWriteFile(t *testing.T) {
	file := new(os.File)
	rs := NewResForWriteFile(file, log)
	ast := assert.New(t)
	ast.NotNil(rs)
}

func TestUtils_UnMarshalToStruct_fail(t *testing.T) {
	rs := new(ResFromFile)
	rs.Logger = log

	file := new(os.File)
	rs.File = file

	var s []byte = nil
	err := rs.UnMarshalToStruct(s)

	ast := assert.New(t)
	ast.NotNil(err)

}

func InitResFromFile() *ResFromFile {
	rs := new(ResFromFile)
	rs.RrResult = make([][]string, 0)

	rowRes1 := make([]string, 0)
	rowRes1 = append(rowRes1, "hello", "word", "program")

	rowRes2 := make([]string, 0)
	rowRes2 = append(rowRes2, "hello1", "word1", "program1")

	rs.RrResult = append(rs.RrResult, rowRes1, rowRes2)

	rs.PrResult = make([][]string, 0)
	rowRes11 := make([]string, 0)

	rowRes11 = append(rowRes11, "hello", "word", "program")

	rowRes12 := make([]string, 0)
	rowRes12 = append(rowRes12, "hello1", "word1", "program1")
	rs.PrResult = append(rs.PrResult, rowRes11, rowRes12)

	rs.Logger = log

	rs.Type = 1
	rs.StmtID = 10
	rs.Params = make([]interface{}, 0)
	rs.Params = append(rs.Params, "abc")
	rs.DB = "test"
	rs.Query = "select * from test.test where name=?"
	var timeNow  = uint64(time.Now().UnixNano())
	rs.PrBeginTime = timeNow
	rs.PrEndTime = timeNow
	rs.PrErrorNo = 1205
	rs.PrErrorDesc = "lock wait timeout"

	rs.RrEndTime = timeNow
	rs.RrBeginTime = timeNow
	rs.RrErrorNo = 1205
	rs.RrErrorDesc = "lock wait timeout"
	return rs
}

func TestUtils_UnMarshalToStruct_Succ(t *testing.T) {
	rs := InitResFromFile()
	s, _ := json.Marshal(rs)

	file := new(os.File)
	rs1 := NewResForWriteFile(file, log)
	err := rs1.UnMarshalToStruct(s)

	ast := assert.New(t)
	ast.Nil(err)
}

/*
func TestUtils_CompareRes_Succ(t *testing.T) {

	rs := InitResFromFile()

	scr := rs.CompareRes()

	ast := assert.New(t)

	ast.Equal(scr.ErrCode, 0)

}


func InitResFromFileErrorNoNotEq() *ResFromFile {
	rs := new(ResFromFile)
	rs.RrResult = make([][]string, 0)

	rowRes1 := make([]string, 0)
	rowRes1 = append(rowRes1, "hello", "word", "program")

	rowRes2 := make([]string, 0)
	rowRes2 = append(rowRes2, "hello1", "word1", "program1")

	rs.RrResult = append(rs.RrResult, rowRes1, rowRes2)

	rs.PrResult = make([][]string, 0)
	rowRes11 := make([]string, 0)

	rowRes11 = append(rowRes11, "hello", "word", "program")

	rowRes12 := make([]string, 0)
	rowRes12 = append(rowRes12, "hello1", "word1", "program1")
	rs.PrResult = append(rs.PrResult, rowRes11, rowRes12)

	rs.Logger = log

	rs.Type = 1
	rs.StmtID = 10
	rs.Params = make([]interface{}, 0)
	rs.Params = append(rs.Params, "abc")
	rs.DB = "test"
	rs.Query = "select * from test.test where name=?"
	var timeNow  = uint64(time.Now().UnixNano())
	rs.PrBeginTime = timeNow
	rs.PrEndTime = timeNow
	rs.PrErrorNo = 1206
	rs.PrErrorDesc = "lock wait timeout"

	rs.RrEndTime = timeNow
	rs.RrBeginTime = timeNow
	rs.RrErrorNo = 1205
	rs.RrErrorDesc = "lock wait timeout"
	return rs
}

func TestUtils_CompareRes_With_Errcode1(t *testing.T) {

	rs := InitResFromFileErrorNoNotEq()

	scr := rs.CompareRes()

	ast := assert.New(t)

	ast.Equal(scr.ErrCode, 1)

}

func InitResFromFileExecTimeNotEq() *ResFromFile {
	rs := new(ResFromFile)
	rs.RrResult = make([][]string, 0)

	rowRes1 := make([]string, 0)
	rowRes1 = append(rowRes1, "hello", "word", "program")

	rowRes2 := make([]string, 0)
	rowRes2 = append(rowRes2, "hello1", "word1", "program1")

	rs.RrResult = append(rs.RrResult, rowRes1, rowRes2)

	rs.PrResult = make([][]string, 0)
	rowRes11 := make([]string, 0)

	rowRes11 = append(rowRes11, "hello", "word", "program")

	rowRes12 := make([]string, 0)
	rowRes12 = append(rowRes12, "hello1", "word1", "program1")
	rs.PrResult = append(rs.PrResult, rowRes11, rowRes12)

	rs.Logger = log

	rs.Type = 1
	rs.StmtID = 10
	rs.Params = make([]interface{}, 0)
	rs.Params = append(rs.Params, "abc")
	rs.DB = "test"
	rs.Query = "select * from test.test where name=?"
	var timeNow  = uint64(time.Now().UnixNano())
	rs.PrBeginTime = timeNow
	rs.PrEndTime = timeNow + 200
	rs.PrErrorNo = 1205
	rs.PrErrorDesc = "lock wait timeout"

	rs.RrEndTime = timeNow + 200000000
	rs.RrBeginTime = timeNow
	rs.RrErrorNo = 1205
	rs.RrErrorDesc = "lock wait timeout"
	return rs
}


func TestUtils_CompareRes_With_Errcode2(t *testing.T) {

	rs := InitResFromFileExecTimeNotEq()

	scr := rs.CompareRes()

	ast := assert.New(t)

	ast.Equal(scr.ErrCode, 0)

}


func InitResFromFileResRowLenNotEq() *ResFromFile {
	rs := new(ResFromFile)
	rs.RrResult = make([][]string, 0)

	rowRes1 := make([]string, 0)
	rowRes1 = append(rowRes1, "hello", "word", "program")

	rowRes2 := make([]string, 0)
	rowRes2 = append(rowRes2, "hello1", "word1", "program1")

	rs.RrResult = append(rs.RrResult, rowRes1, rowRes2)

	rs.PrResult = make([][]string, 0)
	rowRes11 := make([]string, 0)

	rowRes11 = append(rowRes11, "hello", "word", "program")

	rs.PrResult = append(rs.PrResult, rowRes11)

	rs.Logger = log

	rs.Type = 1
	rs.StmtID = 10
	rs.Params = make([]interface{}, 0)
	rs.Params = append(rs.Params, "abc")
	rs.DB = "test"
	rs.Query = "select * from test.test where name=?"
	var timeNow  = uint64(time.Now().UnixNano())
	rs.PrBeginTime = timeNow
	rs.PrEndTime = timeNow + 200
	rs.PrErrorNo = 1205
	rs.PrErrorDesc = "lock wait timeout"

	rs.RrEndTime = timeNow + 200
	rs.RrBeginTime = timeNow
	rs.RrErrorNo = 1205
	rs.RrErrorDesc = "lock wait timeout"
	return rs
}

func TestUtils_CompareRes_With_Errcode3(t *testing.T) {

	rs := InitResFromFileResRowLenNotEq()

	scr := rs.CompareRes()

	ast := assert.New(t)

	ast.Equal(scr.ErrCode, 3)

}

func InitResFromFileResRowDetailNotEq() *ResFromFile {
	rs := new(ResFromFile)
	rs.RrResult = make([][]string, 0)

	rowRes1 := make([]string, 0)
	rowRes1 = append(rowRes1, "hello", "word", "program")

	rowRes2 := make([]string, 0)
	rowRes2 = append(rowRes2, "hello1", "word1", "program1")

	rs.RrResult = append(rs.RrResult, rowRes1, rowRes2)

	rs.PrResult = make([][]string, 0)
	rowRes11 := make([]string, 0)

	rowRes11 = append(rowRes11, "hello", "word", "program")

	rowRes12 := make([]string, 0)
	rowRes12 = append(rowRes2, "hello1", "word2", "program1")

	rs.PrResult = append(rs.PrResult, rowRes11, rowRes12)

	rs.Logger = log

	rs.Type = 1
	rs.StmtID = 10
	rs.Params = make([]interface{}, 0)
	rs.Params = append(rs.Params, "abc")
	rs.DB = "test"
	rs.Query = "select * from test.test where name=?"
	var timeNow  = uint64(time.Now().UnixNano())
	rs.PrBeginTime = timeNow
	rs.PrEndTime = timeNow + 200
	rs.PrErrorNo = 1205
	rs.PrErrorDesc = "lock wait timeout"

	rs.RrEndTime = timeNow + 200
	rs.RrBeginTime = timeNow
	rs.RrErrorNo = 1205
	rs.RrErrorDesc = "lock wait timeout"
	return rs
}

func TestUtils_CompareRes_With_Errcode4(t *testing.T) {

	rs := InitResFromFileResRowDetailNotEq()

	scr := rs.CompareRes()

	ast := assert.New(t)

	ast.Equal(scr.ErrCode, 4)

}



func TestUtils_DoCompare_With_GetResFromFile_Fail(t *testing.T) {
	file := new(os.File)
	fn := "result-test"
	var wg sync.WaitGroup
	wg.Add(1)

	filePath:="./"
	backDir:="./"

	rs := NewResForWriteFile(file, log)

	patch := gomonkey.ApplyFunc(NewResForWriteFile, func(file *os.File, log *zap.Logger) *ResFromFile {
		return rs
	})
	defer patch.Reset()

	err := errors.New("get data from file fail")

	patches := gomonkey.ApplyMethod(reflect.TypeOf(rs), "GetResFromFile",
		func(_ *ResFromFile) ([]byte, error) {
			return nil, err
		})
	defer patches.Reset()
	var ct int32=1
	DoCompare(fn, &ct, &wg,filePath,backDir)

	ast := assert.New(t)

	ast.Greater(stat.Statis.MStat["ResultFiles"], uint64(0))
	ast.Greater(stat.Statis.MStat["ReadFailFiles"], uint64(0))

}

func TestUtils_DoCompare_With_GetResFromFile_Fail_EOF(t *testing.T) {
	file := new(os.File)
	fn := "result-test"
	var wg sync.WaitGroup
	wg.Add(1)

	filePath:="./"
	backDir:="./"

	rs := NewResForWriteFile(file, log)

	patch := gomonkey.ApplyFunc(NewResForWriteFile, func(file *os.File, log *zap.Logger) *ResFromFile {
		return rs
	})
	defer patch.Reset()

	patches := gomonkey.ApplyMethod(reflect.TypeOf(rs), "GetResFromFile",
		func(_ *ResFromFile) ([]byte, error) {
			return nil, io.EOF
		})
	defer patches.Reset()

	var ct int32=1
	DoCompare(fn, &ct, &wg,filePath,backDir)

	ast := assert.New(t)

	ast.Greater(stat.Statis.MStat["ResultFiles"], uint64(0))
	ast.Greater(stat.Statis.MStat["ReadSuccFiles"], uint64(0))

}

func TestUtils_DoCompare_With_UnMarshalToStruct_Fail(t *testing.T) {
	file := new(os.File)
	fn := "result-test"
	var wg sync.WaitGroup
	wg.Add(1)

	rs := NewResForWriteFile(file, log)

	patch := gomonkey.ApplyFunc(NewResForWriteFile, func(file *os.File, log *zap.Logger) *ResFromFile {
		return rs
	})
	defer patch.Reset()

	patches := gomonkey.ApplyMethod(reflect.TypeOf(rs), "GetResFromFile",
		func(_ *ResFromFile) ([]byte, error) {
			return nil, nil
		})
	defer patches.Reset()

	err := errors.New("unmarshal to struct fail")
	patches1 := gomonkey.ApplyMethod(reflect.TypeOf(rs), "UnMarshalToStruct",
		func(_ *ResFromFile, s []byte) error {
			return err
		})
	defer patches1.Reset()

	var ct int32=1
	DoCompare(fn, &ct, &wg,"./","./")

	ast := assert.New(t)

	ast.Greater(stat.Statis.MStat["ResultFiles"], uint64(0))
	ast.Greater(stat.Statis.MStat["ReadFailFiles"], uint64(0))

}


func TestUtils_DoCompare_With_CompareRes_Fail_Marshal_Fail(t *testing.T) {
	file := new(os.File)
	fn := "result-test"
	var wg sync.WaitGroup
	wg.Add(1)

	rs := NewResForWriteFile(file, log)
	rs.Query = "select * from test.test where id =?"

	a := gomonkey.OutputCell{
		Values: gomonkey.Params{nil, nil},
		Times:  1,
	}

	b := gomonkey.OutputCell{
		Values: gomonkey.Params{nil, io.EOF},
		Times:  2,
	}
	outputs := make([]gomonkey.OutputCell, 0)
	outputs = append(outputs, a, b)

	patch := gomonkey.ApplyFunc(NewResForWriteFile, func(file *os.File, log *zap.Logger) *ResFromFile {
		return rs
	})
	defer patch.Reset()

	patches := gomonkey.ApplyMethodSeq(reflect.TypeOf(rs), "GetResFromFile",
		outputs)
	defer patches.Reset()

	patches1 := gomonkey.ApplyMethod(reflect.TypeOf(rs), "UnMarshalToStruct",
		func(_ *ResFromFile, s []byte) error {
			return nil
		})
	defer patches1.Reset()

	err := errors.New("marshal data fail")
	patch1 := gomonkey.ApplyFunc(json.Marshal, func(v interface{}) ([]byte, error) {
		return nil, err
	})
	defer patch1.Reset()

	res := new(SqlCompareRes)
	res.Sql = "select * from test.test where id =?"
	res.ErrCode = 1
	res.ErrDesc = "1205-1206"
	res.Values = make([]interface{}, 0)
	res.Values = append(res.Values, "abc")


	var ct int32=1
	DoCompare(fn, &ct, &wg,"./","./")

	ast := assert.New(t)

	ast.Greater(stat.Statis.MStat["ResultFiles"], uint64(0))
	ast.Greater(stat.Statis.MStat["ReadSuccFiles"], uint64(0))

}



func TestUtils_DoCompare_With_CompareRes_Fail_Marshal_Succ(t *testing.T) {
	file := new(os.File)
	fn := "result-test"
	var wg sync.WaitGroup
	wg.Add(1)

	rs := NewResForWriteFile(file, log)

	a := gomonkey.OutputCell{
		Values: gomonkey.Params{nil, nil},
		Times:  1,
	}

	b := gomonkey.OutputCell{
		Values: gomonkey.Params{nil, io.EOF},
		Times:  2,
	}
	outputs := make([]gomonkey.OutputCell, 0)
	outputs = append(outputs, a, b)

	patch := gomonkey.ApplyFunc(NewResForWriteFile, func(file *os.File, log *zap.Logger) *ResFromFile {
		return rs
	})
	defer patch.Reset()

	patches := gomonkey.ApplyMethodSeq(reflect.TypeOf(rs), "GetResFromFile",
		outputs)
	defer patches.Reset()

	patches1 := gomonkey.ApplyMethod(reflect.TypeOf(rs), "UnMarshalToStruct",
		func(_ *ResFromFile, s []byte) error {
			return nil
		})
	defer patches1.Reset()

	//err:=errors.New("marshal data fail")
	patch1 := gomonkey.ApplyFunc(json.Marshal, func(v interface{}) ([]byte, error) {
		return nil, nil
	})
	defer patch1.Reset()

	res := new(SqlCompareRes)
	res.Sql = "select * from test.test wher id =?"
	res.ErrCode = 1
	res.ErrDesc = "1205-1206"
	res.Values = make([]interface{}, 0)
	res.Values = append(res.Values, "abc")

	DoCompare(fn, file, &wg)

	ast := assert.New(t)

	ast.Greater(stat.Statis.MStat["ResultFiles"], uint64(0))
	ast.Greater(stat.Statis.MStat["ReadSuccFiles"], uint64(0))

}
*/

func Test_GetResFromFile_First_Read_EOF(t *testing.T){
	file := new(os.File)

	rs := NewResForWriteFile(file, log)

	patches1 := gomonkey.ApplyMethod(reflect.TypeOf(file), "Read",
		func(_ *os.File,b []byte) (n int, err error) {
			return 0,io.EOF
		})
	defer patches1.Reset()

	_,err :=rs.GetResFromFile()

	ast:= assert.New(t)
	ast.Equal(err,io.EOF)

}

func Test_GetResFromFile_First_Read_err(t *testing.T){
	file := new(os.File)

	rs := NewResForWriteFile(file, log)

	err1:= errors.New("read data fail")
	patches1 := gomonkey.ApplyMethod(reflect.TypeOf(file), "Read",
		func(_ *os.File,b []byte) (n int, err error) {
			return 0,err1
		})
	defer patches1.Reset()

	_,err :=rs.GetResFromFile()

	ast:= assert.New(t)
	ast.Equal(err,err1)

}

func Test_GetResFromFile_Second_Read_err(t *testing.T){
	file := new(os.File)

	rs := NewResForWriteFile(file, log)
	err1:= errors.New("read data fail")

	//s1:=[]uint8{0,0,0,0,0,0,0,243}

	a := gomonkey.OutputCell{
		Values: gomonkey.Params{8, nil},
		Times:  1,
	}

	b := gomonkey.OutputCell{
		Values: gomonkey.Params{0, err1},
		Times:  2,
	}
	outputs := make([]gomonkey.OutputCell, 0)
	outputs = append(outputs, a, b)

	patches := gomonkey.ApplyMethodSeq(reflect.TypeOf(file), "Read",
		outputs)
	defer patches.Reset()


	_,err :=rs.GetResFromFile()

	ast:= assert.New(t)
	ast.Equal(err,err1)

}

func TestResFromFile_DetermineNeedCompareResult(t *testing.T) {
	type fields struct {
		Type        uint64
		StmtID      uint64
		Params      []interface{}
		DB          string
		Query       string
		PrBeginTime uint64
		PrEndTime   uint64
		PrErrorNo   uint16
		PrErrorDesc string
		PrResult    [][]string
		RrBeginTime uint64
		RrEndTime   uint64
		RrErrorNo   uint16
		RrErrorDesc string
		RrResult    [][]string
		Logger      *zap.Logger
		File        *os.File
		Pos         int64
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:"sql event is not EventQuery ,EventStmtExecute",
			fields:fields{
				Type:utils.EventQuit,
				Logger:log,
				Query:"",
			},
			want:false,
		},
		{
			name:"sql event is  EventQuery and sql is select  ",
			fields:fields{
				Type:utils.EventQuery,
				Logger:log,
				Query:"select * from t1 where id=1",
			},
			want:true,
		},
		{
			name:"sql event is  EventQuery and sql is select for update ",
			fields:fields{
				Type:utils.EventQuery,
				Logger:log,
				Query:"select * from t1 where id=1 for update ",
			},
			want:true,
		},
		{
			name:"sql event is  EventQuery and sql is insert ",
			fields:fields{
				Type:utils.EventQuery,
				Logger:log,
				Query:"insert into t1 (id,name) values (1,'aa');",
			},
			want:false,
		},
		{
			name:"sql event is  EventQuery and sql parse fail ",
			fields:fields{
				Type:utils.EventQuery,
				Logger:log,
				Query:"insert into t1 ;",
			},
			want:false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := &ResFromFile{
				Type:        tt.fields.Type,
				StmtID:      tt.fields.StmtID,
				Params:      tt.fields.Params,
				DB:          tt.fields.DB,
				Query:       tt.fields.Query,
				PrBeginTime: tt.fields.PrBeginTime,
				PrEndTime:   tt.fields.PrEndTime,
				PrErrorNo:   tt.fields.PrErrorNo,
				PrErrorDesc: tt.fields.PrErrorDesc,
				PrResult:    tt.fields.PrResult,
				RrBeginTime: tt.fields.RrBeginTime,
				RrEndTime:   tt.fields.RrEndTime,
				RrErrorNo:   tt.fields.RrErrorNo,
				RrErrorDesc: tt.fields.RrErrorDesc,
				RrResult:    tt.fields.RrResult,
				Logger:      tt.fields.Logger,
				File:        tt.fields.File,
				Pos:         tt.fields.Pos,
			}
			if got := rs.DetermineNeedCompareResult(); got != tt.want {
				t.Errorf("DetermineNeedCompareResult() = %v, want %v", got, tt.want)
			}
		})
	}
}


func TestDoComparePre_fail(t *testing.T){
	err:=errors.New("do not have privileges ")
	patch := gomonkey.ApplyFunc(utils.OpenFile, func(fn string) (*os.File,error) {
		return nil, err
	})
	defer patch.Reset()

	f,err1:= DoComparePre("./test",log)

	ast:=assert.New(t)
	ast.Nil(f)
	ast.Equal(err,err1)
}

func TestDoComparePre_succ(t *testing.T){

	patch := gomonkey.ApplyFunc(utils.OpenFile, func(fn string) (*os.File,error) {
		return new(os.File), nil
	})
	defer patch.Reset()

	f,err1:= DoComparePre("./test",log)

	ast:=assert.New(t)
	ast.Nil(err1)
	ast.NotNil(f)
}

func TestDoCompareFinish_with_backdir_len_zero(t *testing.T) {
	file := new(os.File)
	DoCompareFinish(file,log,"./","","test")
}

func TestDoCompareFinish_with_move_fail(t *testing.T) {
	file := new(os.File)
	err := errors.New("do not have privileges ")
	patch := gomonkey.ApplyFunc( utils.MoveFileToBackupDir, func(dataDir  ,fileName ,backupDir string ) error {
		return err
	})
	defer patch.Reset()

	DoCompareFinish(file,log,"./","","test")
}

func TestDoCompareFinish_with_move_succ(t *testing.T) {
	file := new(os.File)

	patch := gomonkey.ApplyFunc( utils.MoveFileToBackupDir, func(dataDir  ,fileName ,backupDir string ) error {
		return nil
	})
	defer patch.Reset()

	DoCompareFinish(file,log,"./","","test")
}
