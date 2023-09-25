package cmd

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gitub.com/devWaylander/rssagg/internal/handlers"
)

func RunServer() {
	err := godotenv.Load()
	if err != nil {
		logrus.WithError(err).Panic("The env can't be readed")
	}

	port := os.Getenv("PORT")
	if port == "" {
		logrus.Fatal("PORT is not found in the env")
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
	v1Router.Get("/healthz", handlers.Readiness)
	v1Router.Get("/err", handlers.Err)

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

// func newDB(connectionString string) (*sqlx.DB, error) {
// 	d, err := sqlx.Connect("postgres", connectionString)
// 	if err != nil {
// 		return nil, err
// 	}

// 	d.SetMaxOpenConns(5)
// 	d.SetMaxIdleConns(5)
// 	d.SetConnMaxLifetime(5 * time.Minute)
// 	d.MapperFunc(snaker.CamelToSnake)

// 	return d, nil
// }
