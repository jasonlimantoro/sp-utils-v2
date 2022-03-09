package reviewmergerequest

import "strings"

var repoToRecommendedReviewerMapping = map[string][]string{
	"core-services/shopee_backend":        {"raymond", "fazli"},
	"core-services/beeshop_web":           {"sentosa", "fahmi"},
	"shopee/pl/marketplace-payment":       {"shannon", "hai"},
	"shopee/pl/mpp-fee-service":           {"fahmi"},
	"shopee/marketplace-payments/item":    {"shannon"},
	"shopee/marketplace-payments/channel": {"sentosa"},
	"shopee/marketplace-payments/common":  {"fahmi"},
	"shopee/marketplace-payments/buyer":   {"gabriel"},
	"shopee/marketplace-payments/seller":  {"fahmi"},
}

var repoPathToRepoAlias = map[string]string{
	"shopee/pl/marketplace-payment": "bridge",
}

func getRepoName(repo string) string {
	if v, ok := repoPathToRepoAlias[repo]; ok {
		return v
	}

	segments := strings.Split(repo, "/")

	return segments[len(segments)-1]
}

var reviewerToMattermostUsernameMapping = map[string]string{
	"fahmi":   "fahmi.fahmi",
	"sentosa": "sentosa.adjikusuma",
	"shannon": "shannon.wong",
	"fazli":   "roslim",
	"raymond": "laiw",
	"hai":     "buith",
	"gabriel": "gabriel.onglx",
}
