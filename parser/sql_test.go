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
 * @File:  sql_test.go
 * @Version: 1.0.0
 * @Date: 2021/11/8 10:47
 */

package parser

import (
	_ "github.com/pingcap/tidb/parser/test_driver"
	//"github.com/pingcap/tidb/parser"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

var log =  zap.L().With(zap.String("util", "file"))





func TestParse_With_Succ(t *testing.T){
	sql :="select * from t1 where id =1;"
	_,err := parse(sql)
	assert.New(t).Nil(err)
}

func TestParse_With_fail(t *testing.T){

	sql:="select id from "
	_,err := parse(sql)

	assert.New(t).NotNil(err)
}

func TestParse_With_Result_len_zero(t *testing.T){

	sql:=" "
	_,err := parse(sql)
	//fmt.Println(err.Error())
	assert.New(t).NotNil(err)
}

func TestParse_With_Prepare_Succ(t *testing.T){
	sql:="select * from t where id =?;"

	_,err := parse(sql)

	assert.New(t).Nil(err)
}

func TestCheckIsSelectStmt(t *testing.T) {
	type args struct {
		sql string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name : "sql is select ",
			args : args{
				sql :"select * from t1 where id =10;",
			},
			want: true ,
			wantErr: false,
		},
		{
			name : "sql is  select  for update ",
			args : args{
				sql :"select * from t1 where id =10 for update;",
			},
			want: true ,
			wantErr: false,
		},
		{
			name : "sql is not select ",
			args : args{
				sql :"insert into t (id,name) values (1,'aaa');",
			},
			want: false ,
			wantErr: false,
		},
		{
			name : "insert select  ",
			args : args{
				sql :"insert into t (id,name) select id ,name from t;",
			},
			want: false ,
			wantErr: false,
		},
		{
			name : "sql is not select  ",
			args : args{
				sql :"update t set id =100 where id =10;",
			},
			want: false ,
			wantErr: false,
		},
		{
			name : "select into outfile ",
			args : args{
				sql :"select id,name from test where id >100 into outfile 'test.txt';",
			},
			want: false  ,
			wantErr: false,
		},
		{
			name : "parse sql fail ",
			args : args{
				sql :"select * from ;",
			},
			want: false ,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckIsSelectStmt(tt.args.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckIsSelectStmt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckIsSelectStmt() got = %v, want %v", got, tt.want)
			}
		})
	}
}