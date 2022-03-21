package config

import "os"

func GetGitlabPrivateToken() string {
	return os.Getenv("GITLAB_PRIVATE_TOKEN")
}

func GetTrelloAPIKey() string {
	return os.Getenv("TRELLO_API_KEY")
}

func GetTrelloAPIToken() string {
	return os.Getenv("TRELLO_API_TOKEN")
}
