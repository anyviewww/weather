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
    "github.com/rs/cors" // For CORS support
)

// WeatherResponse defines the structure for parsing the OpenWeather API response
type WeatherResponse struct {
    Main struct {
        Temp float64 `json:"temp"` // Temperature in Celsius
    } `json:"main"`
    Name string `json:"name"` // City name
}

// getWeather fetches the current temperature for a given city using the OpenWeather API
func getWeather(city string) (float64, error) {
    // Load the API key from the environment variable
    apiKey := os.Getenv("OPENWEATHER_API_KEY")
    if apiKey == "" {
        return 0, fmt.Errorf("OPENWEATHER_API_KEY is not set")
    }

    // Construct the URL for the OpenWeather API request
    url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)

    // Perform an HTTP GET request to the OpenWeather API
    resp, err := http.Get(url)
    if err != nil {
        return 0, fmt.Errorf("failed to connect to weather API: %w", err)
    }
    defer resp.Body.Close()

    // Check if the response status is OK (HTTP 200)
    if resp.StatusCode != http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        return 0, fmt.Errorf("weather API error: %s, response: %s", resp.Status, body)
    }

    // Read the response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return 0, fmt.Errorf("failed to read weather API response: %w", err)
    }

    // Parse the JSON response into the WeatherResponse struct
    var weatherResp WeatherResponse
    err = json.Unmarshal(body, &weatherResp)
    if err != nil {
        return 0, fmt.Errorf("failed to parse weather API response: %w", err)
    }

    return weatherResp.Main.Temp, nil
}

// getTime returns the current time in UTC
func getTime() string {
    loc, err := time.LoadLocation("UTC") // Use UTC as fallback
    if err != nil {
        loc, _ = time.LoadLocation("UTC")
    }
    return time.Now().In(loc).Format(time.RFC1123)
}

// weatherHandler handles incoming HTTP requests and returns weather and time data
func weatherHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r) // Extract route variables (e.g., city name)
    city := vars["city"]

    // Fetch the current temperature for the given city
    temp, err := getWeather(city)
    if err != nil {
        log.Printf("Error fetching weather data for city %s: %v", city, err)
        http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
        return
    }

    // Get the current time
    time := getTime()

    // Prepare the JSON response
    response := map[string]interface{}{
        "city":    city,
        "weather": temp,
        "time":    time,
    }

    // Set the response headers and encode the JSON response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func main() {
    // Check if the OPENWEATHER_API_KEY environment variable is set
    apiKey := os.Getenv("OPENWEATHER_API_KEY")
    if apiKey == "" {
        log.Fatal("OPENWEATHER_API_KEY is not set. Please set it in the .env file or as an environment variable.")
    }

    // Set up routing using Gorilla Mux
    r := mux.NewRouter()
    r.HandleFunc("/weather/{city}", weatherHandler).Methods("GET")

    // Add CORS middleware to allow cross-origin requests
    corsHandler := cors.New(cors.Options{
        AllowedOrigins: []string{"*"}, // Allow all origins
        AllowedMethods: []string{"GET"},
    })
    handler := corsHandler.Handler(r)

    // Get the server port from the environment variable or use a default value
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // Start the HTTP server
    log.Printf("Starting server on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, handler))
}
