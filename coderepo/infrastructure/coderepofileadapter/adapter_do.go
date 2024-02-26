/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package coderepofileadapter

import (
	"fmt"
	"strings"

	"github.com/gobwas/glob"
	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/utils"
)

const (
	gitAttributesFile = ".gitattributes"
	fileType          = "file"
	dirType           = "dir"
)

func parseGitAttributesFile(gitAttributeContent string) []string {

	rows := strings.Split(gitAttributeContent, "\n")

	matchStr := make([]string, 0, len(rows))

	for _, v := range rows {
		split := strings.Split(v, " ")
		if len(split) == 0 {
			continue
		}
		matchStr = append(matchStr, split[0])
	}

	return matchStr
}

func checkLfs(mathStr []string, gitAttribute, name string) bool {
	if gitAttribute == "" || name == gitAttributesFile {
		return false
	}
	for _, v := range mathStr {
		g, err := glob.Compile(v)
		if err != nil {
			logrus.Errorf("compile wildcard character of %s err: %s", v, err.Error())
			continue
		}

		if g.Match(name) {
			return true
		}
	}

	return false
}

func getLfsUrl(codeRepoFile *domain.CodeRepoFile, downloadURL string) (lfsURL string, err error) {
	hostName, err := utils.ExtractDomain(downloadURL)

	if err != nil {
		return "", err
	}

	lfsURL = fmt.Sprintf("https://%s/%s/%s/media/branch/%s/%s",
		hostName,
		codeRepoFile.Owner.Account(),
		codeRepoFile.Name.MSDName(),
		codeRepoFile.Ref.FileRef(),
		codeRepoFile.FilePath.FilePath(),
	)

	return lfsURL, nil
}
