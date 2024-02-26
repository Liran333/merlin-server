/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getInt64(t *testing.T, data interface{}) int64 {
	d, ok := data.(float64)
	assert.Truef(t, ok, "data can't convert to int64")

	return int64(d)
}

func getString(t *testing.T, data interface{}) string {
	d, ok := data.(string)
	assert.Truef(t, ok, "data can't convert to string")

	return d
}

func getData(t *testing.T, data interface{}) map[string]interface{} {
	assert.NotNil(t, data)

	d, ok := data.(map[string]interface{})
	assert.Truef(t, ok, "data is not a map[string]interface{} actual: %T", data)

	return d
}

func getArrary(t *testing.T, data interface{}) []map[string]interface{} {
	assert.NotNil(t, data)

	datalist, ok := data.([]interface{})
	assert.Truef(t, ok, "data is not a []interface{} actual: %T", data)

	var out []map[string]interface{}
	out = make([]map[string]interface{}, len(datalist))

	for item := range datalist {
		assert.NotNil(t, datalist[item])
		d, ok := datalist[item].(map[string]interface{})
		assert.Truef(t, ok, "data is not a map[string]interface{} actual: %T", datalist[item])
		out = append(out, d)
	}

	return out
}
