package trello

import "time"

type Card struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Badges  struct {
		AttachmentsByType struct {
			Trello struct {
				Board int `json:"board"`
				Card  int `json:"card"`
			} `json:"trello"`
		} `json:"attachmentsByType"`
		Location           bool   `json:"location"`
		Votes              int    `json:"votes"`
		ViewingMemberVoted bool   `json:"viewingMemberVoted"`
		Subscribed         bool   `json:"subscribed"`
		Fogbugz            string `json:"fogbugz"`
		CheckItems         int    `json:"checkItems"`
		CheckItemsChecked  int    `json:"checkItemsChecked"`
		Comments           int    `json:"comments"`
		Attachments        int    `json:"attachments"`
		Description        bool   `json:"description"`
		Due                string `json:"due"`
		DueComplete        bool   `json:"dueComplete"`
	} `json:"badges"`
	CheckItemStates  []string  `json:"checkItemStates"`
	Closed           bool      `json:"closed"`
	Coordinates      string    `json:"coordinates"`
	CreationMethod   string    `json:"creationMethod"`
	DateLastActivity time.Time `json:"dateLastActivity"`
	Desc             string    `json:"desc"`
	DescData         struct {
		Emoji struct {
		} `json:"emoji"`
	} `json:"descData"`
	Due          string `json:"due"`
	DueReminder  string `json:"dueReminder"`
	Email        string `json:"email"`
	IDBoard      string `json:"idBoard"`
	IDChecklists []struct {
		ID string `json:"id"`
	} `json:"idChecklists"`
	IDLabels []struct {
		ID      string `json:"id"`
		IDBoard string `json:"idBoard"`
		Name    string `json:"name"`
		Color   string `json:"color"`
	} `json:"idLabels"`
	IDList         string   `json:"idList"`
	IDMembers      []string `json:"idMembers"`
	IDMembersVoted []string `json:"idMembersVoted"`
	IDShort        int      `json:"idShort"`
	Labels         []string `json:"labels"`
	Limits         struct {
		Attachments struct {
			PerBoard struct {
				Status    string `json:"status"`
				DisableAt int    `json:"disableAt"`
				WarnAt    int    `json:"warnAt"`
			} `json:"perBoard"`
		} `json:"attachments"`
	} `json:"limits"`
	LocationName          string `json:"locationName"`
	ManualCoverAttachment bool   `json:"manualCoverAttachment"`
	Name                  string `json:"name"`
	Pos                   int    `json:"pos"`
	ShortLink             string `json:"shortLink"`
	ShortURL              string `json:"shortUrl"`
	Subscribed            bool   `json:"subscribed"`
	URL                   string `json:"url"`
	Cover                 struct {
		Color                string `json:"color"`
		IDUploadedBackground bool   `json:"idUploadedBackground"`
		Size                 string `json:"size"`
		Brightness           string `json:"brightness"`
		IsTemplate           bool   `json:"isTemplate"`
	} `json:"cover"`
}

type List struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Closed     bool        `json:"closed"`
	IDBoard    string      `json:"idBoard"`
	Pos        float64     `json:"pos"`
	Subscribed bool        `json:"subscribed"`
	SoftLimit  interface{} `json:"softLimit"`
}

type CreateCardRequest struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type CreateListRequest struct {
	Name string `json:"name"`
	Pos  string `json:"pos"`
}

type GetCardActionsRequest struct {
	CardID string
	Filter string `url:"filter"`
}

type CardAction struct {
	ID              string `json:"id"`
	IDMemberCreator string `json:"idMemberCreator"`
	Data            struct {
		Text     string `json:"text"`
		TextData struct {
			Emoji struct {
			} `json:"emoji"`
		} `json:"textData"`
		Card struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			IDShort   int    `json:"idShort"`
			ShortLink string `json:"shortLink"`
		} `json:"card"`
		Board struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			ShortLink string `json:"shortLink"`
		} `json:"board"`
		List struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"list"`
	} `json:"data"`
	AppCreator interface{} `json:"appCreator"`
	Type       string      `json:"type"`
	Date       time.Time   `json:"date"`
	Limits     struct {
		Reactions struct {
			PerAction struct {
				Status    string `json:"status"`
				DisableAt int    `json:"disableAt"`
				WarnAt    int    `json:"warnAt"`
			} `json:"perAction"`
			UniquePerAction struct {
				Status    string `json:"status"`
				DisableAt int    `json:"disableAt"`
				WarnAt    int    `json:"warnAt"`
			} `json:"uniquePerAction"`
		} `json:"reactions"`
	} `json:"limits"`
	MemberCreator struct {
		ID               string      `json:"id"`
		ActivityBlocked  bool        `json:"activityBlocked"`
		AvatarHash       interface{} `json:"avatarHash"`
		AvatarURL        interface{} `json:"avatarUrl"`
		FullName         string      `json:"fullName"`
		IDMemberReferrer string      `json:"idMemberReferrer"`
		Initials         string      `json:"initials"`
		NonPublic        struct {
		} `json:"nonPublic"`
		NonPublicAvailable bool   `json:"nonPublicAvailable"`
		Username           string `json:"username"`
	} `json:"memberCreator"`
}

const (
	ActionCommentCard = "commentCard"
)
