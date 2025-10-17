package url

import (
	"net/http"
	"regexp"
	"time"
)

var (
	expression = `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`
	httpClient = &http.Client{Timeout: 5 * time.Second}
)

func ValidUrl(i []byte) bool {
	r, err := regexp.Match(expression, i)
	if err != nil {
		return false
	}
	return r
}
