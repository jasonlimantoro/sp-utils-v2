package repository

var RepoToPathMapping = map[string]string{
	"bridge":         "shopee/pl/marketplace-payment",
	"wsa":            "core-services/beeshop_web",
	"backend":        "core-services/shopee_backend",
	"channel":        "shopee/marketplace-payments/channel",
	"item":           "shopee/marketplace-payments/item",
	"buyer":          "shopee/marketplace-payments/buyer",
	"seller":         "shopee/marketplace-payments/seller",
	"common":         "shopee/marketplace-payments/common",
	"tools":          "shopee/marketplace-payments/tools",
	"admin-bff":      "shopee/marketplace-payments/admin-bff",
	"mall-bff":       "shopee/marketplace-payments/mall-bff",
	"beeshop-common": "beetalk-server-deprecated/beeshop_common",
}

func GetRepoPath(repoName string) string {
	val, ok := RepoToPathMapping[repoName]
	if !ok {
		return repoName
	}

	return val
}
