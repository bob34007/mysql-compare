/*
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
 */

/**
 * @Author: guobob
 * @Description:
 * @File:  config_test.go
 * @Version: 1.0.0
 * @Date: 2021/11/26 15:57
 */

package utils

import (
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/pingcap/errors"
	"github.com/stretchr/testify/assert"
)

func Test_CheckConfig_CheckBackDir_len_zero(t *testing.T) {
	cfg := &Config{
		BackDir: "",
	}
	err := cfg.CheckConfig()

	assert.New(t).NotNil(err)
}

func Test_CheckConfig_CheckBackDir_fail(t *testing.T) {
	cfg := &Config{
		BackDir: "./",
	}

	err := errors.New("check backdir fail")
	patch := gomonkey.ApplyFunc(CheckDirExist, func(path string) (bool, error) {
		return false, err
	})
	defer patch.Reset()

	err1 := cfg.CheckConfig()

	assert.New(t).Equal(err, err1)
}

func Test_CheckConfig_CheckDataDir_fail(t *testing.T) {
	cfg := &Config{
		BackDir: "./",
	}

	err := errors.New("check data dir fail")

	outputs := []gomonkey.OutputCell{
		{Values: gomonkey.Params{true, nil}},
		{Values: gomonkey.Params{false, err}},
	}
	patches := gomonkey.ApplyFuncSeq(CheckDirExist, outputs)
	patches.Reset()

	err1 := cfg.CheckConfig()

	assert.New(t).Equal(err, err1)
}

func Test_CheckConfig_CheckMaxGoroutines_fail(t *testing.T) {
	cfg := &Config{
		BackDir:       "./",
		MaxGoroutines: -1,
	}

	patch := gomonkey.ApplyFunc(CheckDirExist, func(path string) (bool, error) {
		return true, nil
	})
	defer patch.Reset()

	err1 := cfg.CheckConfig()

	assert.New(t).NotNil(err1)
}

func Test_CheckConfig_CheckMaxGoroutines_succ(t *testing.T) {
	cfg := &Config{
		BackDir:       "./",
		MaxGoroutines: 1,
	}

	patch := gomonkey.ApplyFunc(CheckDirExist, func(path string) (bool, error) {
		return true, nil
	})
	defer patch.Reset()

	err1 := cfg.CheckConfig()

	assert.New(t).Nil(err1)
}
