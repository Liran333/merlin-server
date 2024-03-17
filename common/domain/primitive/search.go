/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package primitive

import "fmt"

type SearchKey interface {
	SearchKey() string
}

const (
	SearchKeyMaxLength = 100
	SearchKeyMinLength = 0

	SearchTypeMaxLength = 5
	SearchTypeMinLength = 0

	SearchTypeItemMaxLength = 20
	SearchTypeItemMinLength = 0

	SizeMaxLength = 100
	SizeMinLength = 0
)

func NewSearchKey(v string) (SearchKey, error) {
	n := len(v)
	if n >= SearchKeyMaxLength || n <= SearchKeyMinLength {
		return nil, fmt.Errorf("invalid searchKey length, should between %d and %d", 0, 100)
	}

	return searchKey(v), nil
}

func CreateSearchKey(v string) SearchKey {
	return searchKey(v)
}

func (sk searchKey) SearchKey() string {
	return string(sk)
}

type searchKey string

type SearchType interface {
	SearchType() []string
}

func NewSearchType(v []string) (SearchType, error) {
	n := len(v)
	if n >= SearchTypeMaxLength || n <= SearchTypeMinLength {
		return nil, fmt.Errorf("invalid searchType length, should between %d and %d", 0, 5)
	}

	for _, s := range v {
		if len(s) >= SearchTypeItemMaxLength || len(s) <= SearchTypeItemMinLength {
			return nil, fmt.Errorf("invalid searchType item length, should between %d and %d", 0, 20)
		}
	}

	return searchType(v), nil
}

func CreateSearchType(v []string) SearchType {
	return searchType(v)
}

func (st searchType) SearchType() []string {
	return []string(st)
}

type searchType []string

type Size interface {
	Size() int
}

func NewSize(v int) (Size, error) {
	if v >= SizeMaxLength || v <= SizeMinLength {
		return nil, fmt.Errorf("invalid size, should between %d and %d", 1, 100)
	}

	return size(v), nil
}

func CreateSize(v int) Size {
	return size(v)
}

func (s size) Size() int {
	return int(s)
}

type size int
