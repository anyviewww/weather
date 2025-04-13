package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Structure for parsing response from OpenWeather API
type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Name string `json:"name"`
}

// Function to get weather data
func getWeather(city string) (float64, error) {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		return 0, fmt.Errorf("OPENWEATHER_API_KEY is not set")
	}

	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to weather API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return 0, fmt.Errorf("weather API error: %s, response: %s", resp.Status, body)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read weather API response: %w", err)
	}

	var weatherResp WeatherResponse
	err = json.Unmarshal(body, &weatherResp)
	if err != nil {
		return 0, fmt.Errorf("failed to parse weather API response: %w", err)
	}

	return weatherResp.Main.Temp, nil
}

// Function to get the current time
func getTime() string {
	loc, err := time.LoadLocation("UTC") // Use UTC as fallback
	if err != nil {
		loc, _ = time.LoadLocation("UTC")
	}
	return time.Now().In(loc).Format(time.RFC1123)
}

// Request handler
func weatherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	city := vars["city"]

	temp, err := getWeather(city)
	if err != nil {
		log.Printf("Error fetching weather data for city %s: %v", city, err)
		http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
		return
	}

	time := getTime()

	response := map[string]interface{}{
		"city":    city,
		"weather": temp,
		"time":    time,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Check for the presence of an API key
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENWEATHER_API_KEY is not set")
	}

	// Routing setup
	r := mux.NewRouter()
	r.HandleFunc("/weather/{city}", weatherHandler).Methods("GET")

	// Add middleware for CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow all domains
		AllowedMethods: []string{"GET"},
	})
	handler := corsHandler.Handler(r)

	// Get the port from an environment variable or use the default value
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
