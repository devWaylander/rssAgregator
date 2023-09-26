package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/kenshaw/snaker"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

type repo struct {
	db *sqlx.DB
}

func RunServer() {
	err := godotenv.Load()
	if err != nil {
		logrus.WithError(err).Panic("The env can't be readed")
	}

	port := os.Getenv("PORT")
	if port == "" {
		logrus.Fatal("PORT is not found in the env")
	}

	db, err := newDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		logrus.WithError(err).Panic("Unable to connect to DB")
	}
	repoCfg := repo{
		db: db,
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/user", repoCfg.handlerUser)

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	logrus.Info("Server strating at port: ", port)
	err = server.ListenAndServe()
	if err != nil {
		logrus.WithError(err).Panic("The server can't start at port: ", port)
	}
}

func newDB(connectionString string) (*sqlx.DB, error) {
	d, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	d.SetMaxOpenConns(5)
	d.SetMaxIdleConns(5)
	d.SetConnMaxLifetime(5 * time.Minute)
	d.MapperFunc(snaker.CamelToSnake)

	return d, nil
}
