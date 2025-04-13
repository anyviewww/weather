# Weather App
A service for displaying the current temperature and time in a selected city.

## Features
* Enter a city name to get the current weather and local time.
* Backend powered by Go and OpenWeather API.
* Frontend built with React.

## Prerequisites
* Docker and Docker Compose must be installed on your system.
* Obtain an API key from OpenWeatherMap .

## Setup
### 1) Create a .env file
In the root directory of the project, create a file named .env and add your OpenWeather API key:
```bash
OPENWEATHER_API_KEY=your_api_key_here
```
Replace your_api_key_here with your actual API key.

### 2) Verify .gitignore
Ensure that the .env file is ignored by Git to prevent accidental exposure of your API key. The .gitignore file should include:
```bash
.env
```

### 3) Build and Run the Application
Start the application using Docker Compose:
```bash
docker-compose up --build
```

### 4) Access the Application
Once the containers are running, open your browser and navigate to:
```bash
http://localhost:3000
```
## Deployment
To deploy the application on a remote server:

1) Copy the project files to the server (e.g., using scp):
```bash
scp -r ./weather-app root@your-server-ip:/root/
```

2) Copy the .env file to the server:
```bash
scp .env root@your-server-ip:/root/weather-app/
```
3) SSH into the server and start the application:
```bash
ssh root@your-server-ip
cd /root/weather-app/
docker-compose up --build -d
```
4) Open the app in your browser:
```bash
http://your-server-ip:3000
```

## Notes
* If you change the .env file, restart the containers to apply the changes:
```bash

docker-compose down
docker-compose up --build -d
```
* Ensure that the server's firewall allows traffic on ports 3000 (frontend) and 8080 (backend).
