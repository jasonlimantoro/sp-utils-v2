package gmail

type CreateDraftRequest struct {
	Subject     string
	To          string
	CC          string
	BCC         string
	ContentType string
	Body        string
}
