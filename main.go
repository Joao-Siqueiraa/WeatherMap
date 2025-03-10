package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Estrutura para armazenar a resposta da API do OpenWeatherMap
type WeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

// Definir sua chave de API aqui
const apiKey = "4870356d56b757acb389414da8291dba"

// Função para buscar os dados do clima de uma cidade
func getWeather(city string) (WeatherResponse, error) {
	// Monta a URL para a requisição à API
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)
	fmt.Println("Requisição para URL:", url) // Log da URL

	// Fazendo a requisição GET para a API externa
	resp, err := http.Get(url)
	if err != nil {
		return WeatherResponse{}, fmt.Errorf("erro ao fazer a requisição para a API externa: %v", err)
	}
	defer resp.Body.Close()

	// Verificando o status da resposta
	if resp.StatusCode != 200 {
		return WeatherResponse{}, fmt.Errorf("erro ao obter dados da API, status: %d", resp.StatusCode)
	}

	// Lendo o corpo da resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return WeatherResponse{}, fmt.Errorf("erro ao ler a resposta: %v", err)
	}

	// Exibindo a resposta da API para depuração
	fmt.Println("Resposta da API:", string(body))

	// Decodificando a resposta JSON
	var weatherData WeatherResponse
	if err := json.Unmarshal(body, &weatherData); err != nil {
		return WeatherResponse{}, fmt.Errorf("erro ao decodificar resposta JSON: %v", err)
	}

	return weatherData, nil
}

// Função para configurar e iniciar o servidor HTTP
func handleRequests() {
	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		// Pega o valor da cidade da query string
		city := r.URL.Query().Get("city")
		if city == "" {
			http.Error(w, "Por favor, forneça o nome da cidade", http.StatusBadRequest)
			return
		}

		// Obtém as informações do clima
		weatherData, err := getWeather(city)
		if err != nil {
			http.Error(w, fmt.Sprintf("Erro ao obter dados do clima: %v", err), http.StatusInternalServerError)
			return
		}

		// Formata a resposta como JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(weatherData)
	})

	// Inicia o servidor na porta 8080
	fmt.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	// Verifica se a chave da API foi configurada corretamente
	if apiKey == "" {
		fmt.Println("Erro: A chave de API não está configurada!")
		os.Exit(1)
	}

	// Inicia o servidor HTTP
	handleRequests()
}
