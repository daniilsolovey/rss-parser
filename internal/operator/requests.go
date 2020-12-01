package operator

import (
	"github.com/daniilsolovey/rss-parser/internal/database"
	"github.com/mmcdole/gofeed"
	"github.com/reconquest/karma-go"
)

func (operator *Operator) getNews(url string) ([]*database.ResultNews, error) {
	feedParser := gofeed.NewParser()
	feed, err := feedParser.ParseURL(url)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to parse url, url: %s",
			url,
		)
	}

	var result []*database.ResultNews
	for _, item := range feed.Items {
		var date string
		if item.UpdatedParsed == nil {
			date = "null"
		} else {
			date = item.UpdatedParsed.String()
		}

		item := database.ResultNews{
			Title:  item.Title,
			Link:   item.Link,
			Author: item.Author.Name,
			Date:   date,
		}

		result = append(result, &item)
	}

	return result, nil
}
