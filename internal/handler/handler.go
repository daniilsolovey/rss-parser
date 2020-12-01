package handler

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/daniilsolovey/rss-parser/internal/config"
	"github.com/daniilsolovey/rss-parser/internal/database"
	"github.com/reconquest/pkg/log"
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

func (handler *Handler) GetMainPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/site.html")
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
}

func (handler *Handler) FindNews(w http.ResponseWriter, r *http.Request) {
	var requestedText string
	err := r.ParseForm()
	if err != nil {
		log.Errorf(
			err,
			"unable to parse form",
		)
	}

	for key, values := range r.Form {
		if key == "News" {
			for _, value := range values {
				requestedText = value
			}
		}
	}

	if requestedText == "" {
		_, err = fmt.Fprintln(w, "write text in the input field")
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
		log.Errorf(
			err,
			"unable to get records from database",
		)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resultTemplate := template.Must(template.ParseFiles("static/search_result.html"))

	i := 0
	for _, item := range *records {
		i = i + 1
		if strings.Contains(item.Author, requestedText) ||
			strings.Contains(item.Title, requestedText) ||
			strings.Contains(item.Link, requestedText) {
			item.Number = i
			err = resultTemplate.ExecuteTemplate(w, "search_result.html", item)
			if err != nil {
				log.Errorf(
					err,
					"unable to execute template",
				)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		}
	}
}
