package trello

const (
	TrelloHost = "api.trello.com/1"
	BoardID    = "6124940fc33001424232fa48"

	RouteGetListOnBoard    = "boards/%s/lists"
	RouteCreateCardOnList  = "lists/%s/cards"
	RouteCreateListOnBoard = "boards/%s/lists?%s"
)
