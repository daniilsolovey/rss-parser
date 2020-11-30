package operator

import (
	"github.com/daniilsolovey/rss-parser/internal/config"
	"github.com/daniilsolovey/rss-parser/internal/database"
	"github.com/reconquest/karma-go"
)

type Operator struct {
	database *database.Database
	config   *config.Config
}

func NewOperator(
	database *database.Database,
	config *config.Config,
) *Operator {
	operator := &Operator{
		database: database,
		config:   config,
	}

	return operator
}

func (operator *Operator) GetAllRecords() {
	operator.database.GetAllRecords()

}

func (operator *Operator) CreateTable() {
	operator.database.CreateTable()
}

func (operator *Operator) AddNewsToDatabase(url string) error {
	news, err := operator.getNews(url)
	if err != nil {
		return karma.Format(
			err,
			"unable to get news by url: %s",
			url,
		)
	}

	records, err := operator.database.GetAllRecords()
	if err != nil {
		return karma.Format(
			err,
			"unable to get all records from the database",
		)
	}

	for _, item := range news {
		if checkDuplicates(records, item) {
			continue
		}

		operator.database.InsertNewsIntoTable(
			operator.config.Database.TableName,
			item.Title,
			item.Link,
			item.Author,
			item.Date,
		)
	}

	return nil
}

func checkDuplicates(data []database.ResultNews, record database.ResultNews) bool {
	for _, item := range data {
		if record.Title == item.Title {
			return true
		}
	}

	return false
}
