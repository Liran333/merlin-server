package utils

import (
	"strings"

	"code.gitea.io/gitea/modules/structs"
)

const reId = 2

func GetOrgRepo(r *structs.Repository) (string, string) {
	if r == nil {
		return "", ""
	}

	repoId := strings.Split(r.FullName, "/")
	if len(repoId) != reId {
		return "", ""
	}

	return repoId[0], repoId[1]
}
