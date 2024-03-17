/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/
package resourceadapter

import "github.com/openmerlin/merlin-server/search/domain"

type ResourceAdapter interface {
	Search(opt *domain.SearchOption) (domain.SearchResult, error)
}
