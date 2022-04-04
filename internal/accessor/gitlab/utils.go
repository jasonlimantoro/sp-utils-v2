package gitlab

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"git.garena.com/shopee/marketplace-payments/common/errlib"
)

func paginate(
	endpoint string,
	fn func(nextEndpoint string) (interface{}, error),
	getLinkHeader func(resp interface{}) string,
) error {
	resp, err := fn(endpoint)
	if err != nil {
		return errlib.WrapFunc(err)
	}

	linkHeader := buildLinkHeader(getLinkHeader(resp))

	for linkHeader.Next != "" {
		resp, err := fn(linkHeader.Next)
		if err != nil {
			return errlib.WrapFunc(err)
		}
		linkHeader = buildLinkHeader(getLinkHeader(resp))
	}

	return nil
}

type LinkHeader struct {
	First string
	Next  string
	Last  string
}

func buildLinkHeader(linkHeader string) LinkHeader {
	lh := LinkHeader{}

	relRegex := regexp.MustCompile(`rel="(.*)"`)
	linkRegex := regexp.MustCompile(`<(.*)>;`)

	linkSegments := strings.Split(linkHeader, ",")
	for _, segment := range linkSegments {
		relSubMatch := relRegex.FindStringSubmatch(segment)
		if len(relSubMatch) < 1 {
			continue
		}
		rel := relSubMatch[1]

		linkSubmatch := linkRegex.FindStringSubmatch(segment)
		if len(linkSubmatch) < 1 {
			continue
		}
		link := linkSubmatch[1]

		switch rel {
		case "next":
			lh.Next = link
		case "last":
			lh.Last = link
		case "first":
			lh.First = link
		}
	}

	return lh
}

func defaultGetLinkHeader(httpRes *http.Response) string {
	return httpRes.Header.Get("link")
}

// isValidURL tests a string to determine if it is a well-structured url or not.
func isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}
