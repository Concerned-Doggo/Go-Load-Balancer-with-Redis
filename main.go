package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	servers := []Server{
		NewSimpleServer(os.Getenv("SERVER_1_API_URL")),
		NewSimpleServer(os.Getenv("SERVER_2_API_URL")),
		NewSimpleServer(os.Getenv("SERVER_3_API_URL")),
		NewSimpleServer(os.Getenv("SERVER_4_API_URL")),
		NewSimpleServer(os.Getenv("SERVER_5_API_URL")),
	}


	mux := http.NewServeMux()

    cors := cors.New(cors.Options{
		AllowedOrigins: []string{ "https://crypto-insight1.netlify.app/", "http://localhost:5173", "http://localhost:4173"},
        AllowedMethods: []string{
            http.MethodGet,
        },
        AllowedHeaders:   []string{"Content-Type"},
    })

	ConnectRedis()
	lb := NewLoadBalancer("8080", servers)
	

	handleCoin := func(rw http.ResponseWriter, r *http.Request) {
		lb.ServeProxy(rw, r)
	}
	handleChart := func (rw http.ResponseWriter, r *http.Request) {
		lb.ChartServerProxy(rw, r)
	}

	mux.HandleFunc("/{name}", handleCoin)
    mux.HandleFunc("/{name}/market_chart", handleChart)

	handler := cors.Handler(mux)

	fmt.Printf("serving request at 'localhost:%s'", lb.Port)
	http.ListenAndServe(":"+lb.Port, handler)
}

