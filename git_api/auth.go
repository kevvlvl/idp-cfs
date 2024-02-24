package git_api

import (
	"fmt"
	"os"
	"strings"
)

func getAuth(tool string) *GitApiAuth {

	auth := &GitApiAuth{}
	toolUpper := strings.ToUpper(tool)

	user := os.Getenv(fmt.Sprintf("CFS_CODE_%s_USER", toolUpper))
	token := os.Getenv(fmt.Sprintf("CFS_CODE_%s_PAT", toolUpper))

	if user != "" && token != "" {
		auth.codeUser = user
		auth.codeToken = token
		auth.codeDefined = true
	}

	user = os.Getenv(fmt.Sprintf("CFS_GP_%s_USER", toolUpper))
	token = os.Getenv(fmt.Sprintf("CFS_GP_%s_PAT", toolUpper))

	if user != "" && token != "" {
		auth.gpUser = user
		auth.gpToken = token
		auth.gpDefined = true
	}

	return auth
}
