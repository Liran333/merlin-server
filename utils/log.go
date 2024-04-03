/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package utils provides utility functions for various purposes.
package utils

import "github.com/sirupsen/logrus"

// DoLog logs the user's action, access, and result using logrus.
func DoLog(userid, username, action, access, result string) {
	logrus.Infof("| userid: %s | username: %s | action: %s | access: %s | result: %s |",
		userid, username, action, access, result)
}
