package url

import (
	"errors"
	"fmt"
	"strings"
)

func ExplodeUrls(value string) ([]string, error) {
	urls := strings.Split(value, ",")
	valid := []string{}
	invalid := []string{}

	for _, url := range urls {
		if !ValidUrl([]byte(url)) {
			invalid = append(invalid, url)
			continue
		}
		valid = append(valid, url)
	}

	if len(invalid) > 0 {
		if len(valid) == 0 {
			return []string{}, errors.New("every url has an invalid http syntax")
		}

		fmt.Println("The following urls have an invalid http syntax and will not proceed:")
		for _, url := range invalid {
			fmt.Println(url)
		}
		fmt.Println()
	}

	return valid, nil
}
