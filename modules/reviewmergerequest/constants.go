package reviewmergerequest

var repoToRecommendedReviewerMapping = map[string][]string{
	"backend": {"raymond", "fazli"},
	"wsa":     {"sentosa", "fahmi"},
	"bridge":  {"shannon", "hai"},
	"fee":     {"fahmi"},
	"item":    {"shannon"},
	"channel": {"sentosa"},
	"common":  {"fahmi"},
	"buyer":   {"gabriel"},
	"seller":  {"fahmi"},
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
