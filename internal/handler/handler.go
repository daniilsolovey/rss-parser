package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/reconquest/pkg/log"
	"github.com/rss-parser/internal/config"
	"github.com/rss-parser/internal/database"
)

type Handler struct {
	database *database.Database
	config   *config.Config
}

func NewHandler(
	database *database.Database,
	config *config.Config,
) *Handler {
	return &Handler{
		database: database,
		config:   config,
	}
}

func (handler *Handler) FindNewsInDatabase(w http.ResponseWriter, r *http.Request) {
	var requestedText string
	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("site.html")
		if err != nil {
			log.Errorf(
				err,
				"unable to parse html files",
			)
		}

		w.WriteHeader(http.StatusInternalServerError)

		err = t.Execute(w, nil)
		if err != nil {
			log.Errorf(
				err,
				"unable to print html data",
			)
			w.WriteHeader(http.StatusInternalServerError)

		}

	case "POST":
		r.ParseForm()
		for _, values := range r.Form {
			for _, value := range values {
				requestedText = value
			}
		}

		if requestedText == "" {
			_, err := fmt.Fprintln(w, "write text in the input field")
			if err != nil {
				log.Errorf(
					err,
					"unable to print data",
				)
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		records, err := handler.database.GetAllRecords()
		if err != nil {
			if err != nil {
				log.Errorf(
					err,
					"unable to get records from database",
				)

				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		i := 0
		for _, item := range records {
			i = i + 1
			if strings.Contains(item.Author, requestedText) ||
				strings.Contains(item.Title, requestedText) ||
				strings.Contains(item.Link, requestedText) {
				_, err = fmt.Fprintln(
					w,
					"№ ", i,
					"\n Title:", item.Title,
					"\n Author:", item.Author,
					"\n Date:", item.Date,
					"\n Link:", item.Link,
				)
				if err != nil {
					log.Errorf(
						err,
						"unable to print data",
					)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
		}

	default:
		log.Error(
			errors.New("something went wrong"),
		)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
