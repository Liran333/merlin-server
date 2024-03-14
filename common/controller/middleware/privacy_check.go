/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package middleware

import "github.com/gin-gonic/gin"

// PrivacyCheck is an interface that defines the methods for checking user privacy agreement.
type PrivacyCheck interface {
	Check(*gin.Context)
	CheckOwner(ctx *gin.Context)
	CheckName(ctx *gin.Context)
}
