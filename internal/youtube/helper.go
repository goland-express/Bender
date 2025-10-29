package youtube

import (
	"fmt"
	"regexp"
)

func prepareQuery(query string) (string, error) {
	ytQuery := query

	matched, err :=
		regexp.Match(
			"http(?:s?):\\/\\/(?:www\\.)?youtu(?:be\\.com\\/watch\\?v=|\\.be\\/)([\\w\\-\\_]*)(&(amp;)?‌​[\\w\\?‌​=]*)?",
			[]byte(query),
		)

	if err != nil {
		return query, err
	}

	if !matched {
		ytQuery = fmt.Sprintf("ytsearch1:%s", query)
	}

	return ytQuery, nil
}
