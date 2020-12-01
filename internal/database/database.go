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

	db.Close()
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

func (database *Database) CreateTable() {
	_ = database.db.QueryRow(
		"create table articles(" +
			"id serial primary key, title text, link text, author text, date text, " +
			"unique (title));",
	)
}

func (database *Database) InsertNewsIntoTable(
	tableName string, record *ResultNews,
) {
	query := "insert into " + tableName +
		"(" + "title" + ", " + "link" + ", " + "author" + ", " + "date" + ")" +
		" " + "values" + " " +
		"(" + "'" + record.Title + "', " + "'" + record.Link + "', " + "'" + record.Author + "', " +
		"'" + record.Date + "'" + ")" +
		"ON CONFLICT DO NOTHING ;"
	_ = database.db.QueryRow(query)
}

func (database *Database) GetAllRecords() (*[]ResultNews, error) {
	query := "select title, link, author, date from " +
		database.config.Database.TableName + ";"
	rows, err := database.db.Query(query)
	if err != nil || rows.Err() != nil {
		return nil, karma.Format(
			err,
			"unable to get rows from database",
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
