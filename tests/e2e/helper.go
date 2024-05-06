/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package e2e

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getInt64(t *testing.T, data interface{}) int64 {
	if d, ok := data.(float64); ok {
		assert.Truef(t, ok, "data can't convert to int64")
		return int64(d)
	}
	d, ok := data.(int32)
	assert.Truef(t, ok, "data can't convert to int32")
	return int64(d)
}

func getString(t *testing.T, data interface{}) string {
	d, ok := data.(string)
	assert.Truef(t, ok, "data can't convert to string")

	return d
}

func getData(t *testing.T, data interface{}) map[string]interface{} {
	assert.NotNil(t, data)

	if d, ok := data.(map[string]interface{}); ok {
		assert.Truef(t, ok, "data is not a map[string]interface{} actual: %T", data)

		return d
	}

	d := make(map[string]interface{})
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	typ := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := value.Field(i)
		fieldName := field.Tag.Get("json")

		// Check if the field is exported
		if field.PkgPath != "" {
			continue
		}

		// Remove the omitempty option from the JSON tag
		fieldName = strings.Split(fieldName, ",")[0]

		if fieldName != "" {
			d[fieldName] = fieldValue.Interface()
		}
	}

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

func getArrarys(t *testing.T, data interface{}) []map[string]interface{} {
	assert.NotNil(t, data)

	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice {
		assert.Errorf(t, nil, "data is not a slice actual: %T", data)
	}

	var out []map[string]interface{}

	for i := 0; i < value.Len(); i++ {
		item := value.Index(i)

		fields := make(map[string]interface{})
		structType := item.Type()
		for j := 0; j < item.NumField(); j++ {
			field := structType.Field(j)
			fieldValue := item.Field(j).Interface()

			tag := field.Tag.Get("json")
			tagParts := strings.Split(tag, ",")
			fieldName := tagParts[0]

			fields[fieldName] = fieldValue
		}

		out = append(out, fields)
	}

	return out
}
