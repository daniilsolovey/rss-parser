package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/daniilsolovey/rss-parser/internal/config"
	_ "github.com/lib/pq"
	"github.com/reconquest/karma-go"
)

type Database struct {
	db     *sql.DB
	config *config.Config
}

type ResultNews struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Link   string `json:"link"`
	Author string `json:"author"`
	Date   string `json:"date"`
}

type ResultNewsSQL struct {
	Title  sql.NullString
	Link   sql.NullString
	Author sql.NullString
	Date   sql.NullString
}

func NewDatabase(
	db *sql.DB,
	config *config.Config,
) *Database {
	database := &Database{
		db:     db,
		config: config,
	}

	return database
}

func CreateDatabase(config *config.Config) error {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s "+
			"password=%s sslmode=disable",
		config.Database.Host, config.Database.Port,
		config.Database.User, config.Database.Password,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return karma.Format(
			err,
			"unable to open connection to the database",
		)
	}

	err = db.Ping()
	if err != nil {
		return karma.Format(
			err,
			"unable to ping database",
		)
	}

	_, err = db.Exec("create database " + config.Database.Name)
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return karma.Format(
				err,
				"unable to create database",
			)
		}
	}

	err = db.Close()
	if err != nil {
		return karma.Format(
			err,
			"unable to close connection to the database",
		)
	}

	return nil
}

func Connect(config *config.Config) (*sql.DB, error) {
	psqlInfoWithDatabase := fmt.Sprintf(
		"host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
		config.Database.Host, config.Database.Port,
		config.Database.User, config.Database.Password, config.Database.Name,
	)

	db, err := sql.Open("postgres", psqlInfoWithDatabase)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to open connection to the database",
		)
	}

	return db, nil
}

func (database *Database) CreateTable() error {
	result := database.db.QueryRow(
		"CREATE TABLE articles(id serial primary key, title text, link text, " +
			"author text, date text, unique (title))",
	)

	if result.Err() != nil && !strings.Contains(result.Err().Error(), "already exists") {
		return karma.Describe(
			"table_name", database.config.Database.TableName,
		).Format(
			result.Err(),
			"unable to create table",
		)
	}
	return nil
}

func (database *Database) InsertIntoTable(
	tableName string, record *ResultNews,
) error {
	rows, err := database.db.Query(
		"INSERT INTO "+database.config.Database.TableName+
			" (title, link, author, date) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING",
		record.Title, record.Link, record.Author, record.Date,
	)
	if err != nil || rows.Err() != nil {
		return karma.Describe(
			"table_name", tableName,
		).Describe(
			"title", record.Title,
		).Describe(
			"author", record.Author,
		).Describe(
			"link", record.Link,
		).Format(
			err,
			"unable to insert values to table",
		)
	}

	return nil
}

func (database *Database) GetRecordsByFilter(filter string) (*[]ResultNews, error) {
	rows, err := database.db.Query(
		"SELECT title, link, author, date FROM "+database.config.Database.TableName+
			" WHERE title LIKE $1 OR link LIKE $1 OR author LIKE $1",
		"%"+filter+"%",
	)
	if err != nil || rows.Err() != nil {
		return nil, karma.Format(
			err,
			"unable to get rows from database by filter: %s",
			filter,
		)
	}

	var resultMain []ResultNews
	data := ResultNewsSQL{}
	for rows.Next() {
		err = rows.Scan(&data.Title, &data.Link, &data.Author, &data.Date)
		if err != nil {
			return nil, karma.Format(
				err,
				"unable to scan rows from database",
			)
		}

		dataMain := ResultNews{
			Title:  data.Title.String,
			Link:   data.Link.String,
			Author: data.Author.String,
			Date:   data.Date.String,
		}

		resultMain = append(resultMain, dataMain)
	}

	return &resultMain, nil
}
