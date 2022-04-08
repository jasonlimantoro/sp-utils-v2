package gitlab

const (
	GitlabHost = "git.garena.com"
)

const (
	RouteGetProjectsByName  = "api/v4/projects/%s"
	RouteCreateMergeRequest = "api/v4/projects/%d/merge_requests"
	RouteListMergeRequests  = "api/v4/projects/%d/merge_requests?%s"
	RouteGetBranches        = "api/v4/projects/%d/repository/branches?%s"
)
