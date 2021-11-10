/*******************************************************************************
 * Copyright (c)  2021 PingCAP, Inc.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 ******************************************************************************/

/**
 * @Author: guobob
 * @Description:
 * @File:  sqlstat_test.go
 * @Version: 1.0.0
 * @Date: 2021/11/8 11:14
 */

package stat

import (
	"encoding/json"
	"github.com/agiledragon/gomonkey"
	"github.com/pingcap/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	//"strings"
	"testing"
)

var log =  zap.L().With(zap.String("stat", "file"))

func TestSQLResult_String(t *testing.T) {
	type fields struct {
		SQLTemplate                  string
		SQLExecCount                 uint64
		SQLExecSuccCount             uint64
		SQLExecFailCount             uint64
		SQLCompareSuccCount          uint64
		SQLCompareFailCount          uint64
		SQLCompareErrNoFailCount     uint64
		SQLCompareRowCountFailCount  uint64
		SQLCompareRowDetailFailCount uint64
		SQLExecTimeCountPr           uint64
		SQLExecTimeCountRr           uint64
		SQLExecTimeCompareFail       uint64
		SQLExecTimeStandard          uint64
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name : "generate struce to json string ",
			fields :fields{
				SQLTemplate:"select ? from `t` where `a` > ?",
				SQLExecCount: 100,
				SQLExecSuccCount:100,
				SQLExecFailCount: 0,
				SQLCompareFailCount: 10,
				SQLCompareSuccCount: 90,
				SQLCompareErrNoFailCount: 1,
				SQLCompareRowCountFailCount: 3,
				SQLCompareRowDetailFailCount: 2,
				SQLExecTimeCompareFail: 4,
				SQLExecTimeCountPr: 1000000,
				SQLExecTimeCountRr: 1500000,
				SQLExecTimeStandard: 200,
			},
			want:"{\"sql\":\"select ? from `t` where `a` \\u003e ?\",\"sql_exec_count\":100,\"sql_exec_succ_count\":100,\"sql_exec_fail_count\":0,\"sql_compare_succ_count\":90,\"sql_compare_fail_count\":10,\"sql_compare_errno_fail_count\":1,\"sql_compare_rowcount_fail_count\":3,\"sql_compare_rowdetail_fail_count\":2,\"sql_exec_time_count_pr\":1000000,\"sql_exec_time_count_rr\":1500000,\"sql_exec_time_compare_fail\":4,\"sql_exec_time_standard\":200}",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := &SQLResult{
				SQLTemplate:                  tt.fields.SQLTemplate,
				SQLExecCount:                 tt.fields.SQLExecCount,
				SQLExecSuccCount:             tt.fields.SQLExecSuccCount,
				SQLExecFailCount:             tt.fields.SQLExecFailCount,
				SQLCompareSuccCount:          tt.fields.SQLCompareSuccCount,
				SQLCompareFailCount:          tt.fields.SQLCompareFailCount,
				SQLCompareErrNoFailCount:     tt.fields.SQLCompareErrNoFailCount,
				SQLCompareRowCountFailCount:  tt.fields.SQLCompareRowCountFailCount,
				SQLCompareRowDetailFailCount: tt.fields.SQLCompareRowDetailFailCount,
				SQLExecTimeCountPr:           tt.fields.SQLExecTimeCountPr,
				SQLExecTimeCountRr:           tt.fields.SQLExecTimeCountRr,
				SQLExecTimeCompareFail:       tt.fields.SQLExecTimeCompareFail,
				SQLExecTimeStandard:          tt.fields.SQLExecTimeStandard,
			}
			got, err := rs.String()
			if (err != nil) != tt.wantErr {
				t.Errorf("String() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("String() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func InitMap()  {
	var SQL1 = "select * from t1 where id =?"
	AddKey(1000,SQL1, 1, 1, 0, 1, 0, 0, 0, 0, 10, 20,log)
	return
}

func TestPrintMap(t *testing.T) {
	InitMap()
	type args struct {
		log *zap.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:"print struct to log file",
			args:args{
				log:log,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PrintMap(tt.args.log); (err != nil) != tt.wantErr {
				t.Errorf("PrintMap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPrintMap_With_Json_Marshal_fail (t *testing.T){
	InitMap()
	err1:=errors.New("Marshal struct to slice fail")
	patch := gomonkey.ApplyFunc(json.Marshal, func(v interface{}) ([]byte, error) {
		return nil,err1
	})
	defer patch.Reset()
	err:=PrintMap(log)

	ast := assert.New(t)
	ast.Equal(err.Error(),err1.Error())
}

func TestAddKey(t *testing.T) {
	InitMap()
	type args struct {
		key                     uint64
		SQL                     string
		ExecSQL                 uint64
		ExecSQLSucc             uint64
		ExecSQLFail             uint64
		SQLCompareSucc          uint64
		SQLCompareFail          uint64
		SQLCompareErrNoFail     uint64
		SQLCompareRowCountFail  uint64
		SQLCompareRowDetailFail uint64
		SQLExecTimePr           uint64
		SQLExecTimeRr           uint64
		log                     *zap.Logger
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "and new key ",
			args: args{
				key: 2000,
				SQL: "select * from t2 where id =?",
				ExecSQL: 1,
				ExecSQLFail: 0,
				ExecSQLSucc: 1,
				SQLCompareSucc:1,
			},
		},
		{
			name: "and exist key with  no hash collisions",
			args: args{
				key: 1000,
				SQL: "select * from t1 where id >?",
				ExecSQL: 1,
				ExecSQLFail: 0,
				ExecSQLSucc: 1,
				SQLCompareSucc:1,
			},
		},
		{
			name: "and exist key with  hash collisions",
			args: args{
				key: 1000,
				SQL: "select * from t1 where id =?",
				ExecSQL: 1,
				ExecSQLFail: 0,
				ExecSQLSucc: 1,
				SQLCompareSucc:1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}