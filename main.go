package main

import (
	"fmt"
	"net/http"
	"os"

	schnur "github.com/Schnur/cmd"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
    if err != nil {
    	fmt.Println("Error loading .env file")
    }
	
	http.HandleFunc("/strings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			schnur.FilterString(w, r)
		case http.MethodPost:
			schnur.AnalyzeString(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/strings/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			schnur.GetString(w, r)
		case http.MethodDelete:
			schnur.DeleteString(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/strings/filter-by-natural-language", schnur.SearchString)

	// http.HandleFunc("/strings", schnur.FilterString) // GET with query params
	// http.HandleFunc("/strings", schnur.AnalyzeString) // POST
	// http.HandleFunc("/strings/", schnur.GetString) // GET /strings/{value}
	// http.HandleFunc("/strings/filter-by-natural-language", schnur.SearchString) // GET with natural language query
	// http.HandleFunc("/strings/", schnur.DeleteString) // DELETE /strings/{value}
	// r := chi.NewRouter()

	// r.Route("/strings", func(r chi.Router) {
    //     // GET /strings
    //     r.Get("/", schnur.FilterString)

    //     // POST /strings
    //     r.Post("/", schnur.AnalyzeString)

    //     // GET /strings/filter-by-natural-language
    //     r.Get("/filter-by-natural-language", schnur.SearchString)

    //     // GET /strings/{value}
    //     r.Get("/{value}", schnur.GetString)

    //     // DELETE /strings/{value}
    //     r.Delete("/{value}", schnur.DeleteString)
    // })

	port := os.Getenv("PORT")
	fmt.Println("Port from env:", port)

	if port == "" {
		port = "8080" 
	}

	fmt.Println("Server is running on port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}