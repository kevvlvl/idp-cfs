package git_api

import (
	"fmt"
	"os"
	"strings"
)

func getAuth(tool string) *GitApiAuth {

	auth := &GitApiAuth{}

	user := os.Getenv(fmt.Sprintf("CFS_CODE_%s_USER", strings.ToUpper(tool)))
	token := os.Getenv(fmt.Sprintf("CFS_CODE_%s_PAT", strings.ToUpper(tool)))

	if user != "" && token != "" {
		auth.codeUser = user
		auth.codeToken = token
		auth.codeDefined = true
	}

	user = os.Getenv(fmt.Sprintf("CFS_GP_%s_USER", strings.ToUpper(tool)))
	token = os.Getenv(fmt.Sprintf("CFS_GP_%s_PAT", strings.ToUpper(tool)))

	if user != "" && token != "" {
		auth.gpUser = user
		auth.gpToken = token
		auth.gpDefined = true
	}

	return auth
}
