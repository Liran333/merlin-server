/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package primitive

// GetId generates a unique identifier as an int64 value.
func GetId() int64 {
	return node.Generate().Int64()
}
