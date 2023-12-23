package coderepofileadapter

import (
	"strings"

	"github.com/gobwas/glob"
	"github.com/sirupsen/logrus"
)

const (
	gitAttributesFile = ".gitattributes"
	fileType          = "file"
)

func ParseGitAttributesFile(gitAttributeContent string) []string {

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
