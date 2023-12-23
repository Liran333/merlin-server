package utils

import (
	"strings"

	"code.gitea.io/gitea/modules/structs"
)

func GetOrgRepo(r *structs.Repository) (string, string) {
	if r == nil {
		return "", ""
	}

	split := strings.Split(r.FullName, "/")
	if len(split) != 2 {
		return "", ""
	}

	return split[0], split[1]
}
