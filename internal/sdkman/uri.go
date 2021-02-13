package sdkman

import (
	"fmt"
	"net/url"
)

type URI fmt.Stringer

type uri struct {
	url         *url.URL
	scheme      string
	host        string
	segments    []string
	queryString []string
}

func (u *uri) String() string {
	return u.Stringer()
}

func (u *uri) Append(paths ...string) (uri *uri) {
	newUrl := u.url.String()
	for _, path := range paths {
		newUrl = newUrl + "/" + path
	}
	parsedUrl, _ := url.Parse(newUrl)
	u.url = parsedUrl
	return u
}

func (u *uri) Stringer() string {
	return u.url.String()
}
