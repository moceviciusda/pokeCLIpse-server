package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/moceviciusda/pokeCLIpse-server/internal/database"
	"github.com/moceviciusda/pokeCLIpse-server/internal/pokeapi"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB            *database.Queries
	jwtSecret     string
	pokeapiClient pokeapi.Client
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		// log.Fatal("Failed to load the environment")
		log.Println("No .env file found, will use environment variables")
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not found in the environment")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not found in the environment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to DB")
	}

	apiCfg := apiConfig{
		DB:            database.New(conn),
		jwtSecret:     jwtSecret,
		pokeapiClient: pokeapi.NewClient(5*time.Minute, 15*time.Second),
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
	v1Router.Get("/error", handlerError)

	v1Router.HandleFunc("/battle/simulate", apiCfg.handlerSimulateBattle)

	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Post("/login", apiCfg.handlerLogin)
	v1Router.Handle("/starter", apiCfg.middlewareAuth(apiCfg.handlerSelectStarterPokemon))

	v1Router.Get("/location", apiCfg.middlewareAuth(apiCfg.hadlerGetUserLocation))
	v1Router.Put("/location/next", apiCfg.middlewareAuth(apiCfg.handlerNextLocation))
	v1Router.Put("/location/previous", apiCfg.middlewareAuth(apiCfg.handlerPreviousLocation))
	v1Router.Handle("/location/search", apiCfg.middlewareAuth(apiCfg.handlerSearchForPokemon))

	v1Router.Get("/pokemon/party", apiCfg.middlewareAuth(apiCfg.handlerGetPokemonParty))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port: %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
