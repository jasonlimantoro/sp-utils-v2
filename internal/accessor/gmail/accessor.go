package gmail

import (
	"context"
	"encoding/json"
	"net/http"

	"git.garena.com/shopee/marketplace-payments/common/errlib"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/config"
)

type Accessor interface {
	CreateDraft(ctx context.Context, body string) error
}

type accessor struct {
	*gmail.Service
}

func NewGmailService() *gmail.Service {
	gmailCredentials := config.GetGmailCredentials()
	cfg, err := google.ConfigFromJSON([]byte(gmailCredentials), gmail.GmailComposeScope)
	if err != nil {
		panic(err)
	}
	client := getClient(cfg)
	ctx := context.Background()

	gmailClient, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		panic(err)
	}

	return gmailClient
}

func getClient(config *oauth2.Config) *http.Client {
	tok, err := getToken()
	if err != nil {
		panic(err)
	}

	return config.Client(context.Background(), tok)
}

func getToken() (*oauth2.Token, error) {
	gmailToken := config.GetGmailToken()
	tok := &oauth2.Token{}
	err := json.Unmarshal([]byte(gmailToken), tok)
	return tok, err
}

func NewAccessor(service *gmail.Service) *accessor {
	return &accessor{Service: service}
}

func (a accessor) CreateDraft(ctx context.Context, body string) error {
	draft := &gmail.Draft{
		Message: &gmail.Message{
			Raw: body,
		},
	}

	if _, err := a.Service.Users.Drafts.Create("me", draft).Do(); err != nil {
		return errlib.WrapFunc(err)
	}

	return nil
}
