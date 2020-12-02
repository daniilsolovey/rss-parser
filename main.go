package main

import (
	"net/http"
	"time"

	"github.com/daniilsolovey/rss-parser/internal/config"
	"github.com/daniilsolovey/rss-parser/internal/database"
	"github.com/daniilsolovey/rss-parser/internal/handler"
	"github.com/daniilsolovey/rss-parser/internal/operator"
	"github.com/docopt/docopt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

var version = "[manual build]"

var usage = `rss-parser
Receive news and write in database

Usage:
  rss-parser [options] -u <url>

Options:
  -u --url <url>      rss url in format: 'https://www.news_site.com/world/rss'
  -c --config <path>  Read specified config file. [default: config.toml]
  --debug             Enable debug messages.
  -v --version        Print version.
  -h --help           Show this help.
`

const (
	RECEIVING_NEWS_INTERVAL = 10 * time.Minute
)

type Options struct {
	URL    string `docopt:"--url"`
	Config string `docopt:"--config"`
	Debug  bool   `docopt:"--debug"`
}

func main() {
	args, err := docopt.ParseArgs(
		usage,
		nil,
		"rss-parser version: "+version,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof(
		karma.Describe("version", version),
		"rss-parser started",
	)

	if args["--debug"].(bool) {
		log.SetLevel(log.LevelDebug)
	}

	var opts Options
	err = args.Bind(&opts)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof(nil, "loading configuration file: %q", args["--config"].(string))

	config, err := config.Load(args["--config"].(string))
	if err != nil {
		log.Fatal(err)
	}

	log.Info("creating database")
	err = database.CreateDatabase(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof(
		karma.Describe("database", config.Database.Name),
		"database successfully created",
	)

	log.Infof(
		karma.Describe("database", config.Database.Name),
		"connecting to the database",
	)

	db, err := database.Connect(config)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Info("successful connection to the database")

	database := database.NewDatabase(db, config)
	handler := handler.NewHandler(database, config)
	operator := operator.NewOperator(database, config)

	err = operator.CreateTable()
	if err != nil {
		log.Fatal(err)
	}

	url := opts.URL

	go func() {
		for {
			err = operator.AddNewsToDatabase(url)
			if err != nil {
				log.Error(err)
			}
			time.Sleep(RECEIVING_NEWS_INTERVAL)
		}
	}()

	err = router(handler, config)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func router(handler *handler.Handler, config *config.Config) error {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(middleware.Timeout(60 * time.Second))

	router.Route("/", func(r chi.Router) {
		r.Get("/news", handler.GetMainPage)
		r.Post("/news", handler.FindNews)
	})

	err := http.ListenAndServe(config.HTTPPort, router)
	if err != nil {
		return err
	}

	return nil
}
