package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/daniilsolovey/rss-parser/internal/config"
	_ "github.com/lib/pq"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

type Database struct {
	db     *sql.DB
	config *config.Config
}

type ResultNews struct {
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

func Connect(config *config.Config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s "+
			"password=%s sslmode=disable",
		config.Database.Host, config.Database.Port,
		config.Database.User, config.Database.Password,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to open connection to the database",
		)
	}

	err = db.Ping()
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to ping database",
		)
	}

	_, err = db.Exec("create database " + config.Database.Name)
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return nil, karma.Format(
				err,
				"unable to create database",
			)
		}
	}

	db.Close()
	psqlInfoWithDatabase := fmt.Sprintf(
		"host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
		config.Database.Host, config.Database.Port,
		config.Database.User, config.Database.Password, config.Database.Name,
	)

	db, err = sql.Open("postgres", psqlInfoWithDatabase)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to open connection to the database",
		)
	}

	log.Info("successfully connection")
	return db, nil
}

func (database *Database) CreateTable() {
	_ = database.db.QueryRow(
		"create table articles(id serial primary key, title text, link text, author text, date text);",
	)
}

func (database *Database) InsertNewsIntoTable(tableName, title, link, author, date string) {
	query := "insert into " + tableName +
		"(" + "title" + ", " + "link" + ", " + "author" + ", " + "date" + ")" +
		" " + "values" + " " +
		"(" + "'" + title + "', " + "'" + link + "', " + "'" + author + "', " +
		"'" + date + "'" + ")" +
		";"
	_ = database.db.QueryRow(query)
}

func (database *Database) GetAllRecords() ([]ResultNews, error) {
	query := "select title, link, author, date from " +
		database.config.Database.TableName + ";"
	rows, err := database.db.Query(query)
	if err != nil {
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

	return resultMain, nil
}
