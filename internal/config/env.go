package config

import "os"

func GetGitlabPrivateToken() string {
	return os.Getenv("GITLAB_PRIVATE_TOKEN")
}
