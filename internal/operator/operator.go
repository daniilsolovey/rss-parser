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

func (operator *Operator) CreateTable() error {
	err := operator.database.CreateTable()
	if err != nil {
		return karma.Format(
			err,
			"unable to create table",
		)
	}

	return nil
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

	for _, item := range news {
		err = operator.database.InsertIntoTable(
			operator.config.Database.TableName,
			item,
		)
		if err != nil {
			return karma.Format(
				err,
				"unable to insert data to table",
			)
		}
	}

	return nil
}
