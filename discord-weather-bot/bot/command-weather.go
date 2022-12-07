package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"
	"github.com/bwmarrin/discordgo"
)

const URL string = "https://api.openweathermap.org/data/2.5/weather?"

type WeatherData struct {
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Name string `json:"name"`
}

func getCurrentWeather(message string) *discordgo.MessageSend { 
	r, _ := regexp.Compile(`\d{5}`) //regular expression for zipcode
	zip := r.FindString(message)

	if zip == "" {
		return &discordgo.MessageSend {
			Content:"Sorry that ZIP code doesn't look right",
		}
	}

	weatherURL := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?zip=%s&units=imperial&appid=a5c2fd3b3e84d3caf4fa23b8cfded36e", zip)

	// Create new HTTP client & set timeout
	client := &http.Client{Timeout: 5 * time.Second}

	// Query OpenWeather API
	//"https://api.openweathermap.org/data/2.5/weather?zip=91016&units=imperial&appid=a5c2fd3b3e84d3caf4fa23b8cfded36e"
	response, err := client.Get(weatherURL)
	if err != nil {
		fmt.Println(URL)
		fmt.Println(zip)
		fmt.Println(OpenWeatherToken)
		fmt.Println(weatherURL)
		fmt.Println(response)
		return &discordgo.MessageSend{
			Content: "Sorry, there was an error trying to get the weather",
		}
	}

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var data WeatherData
	json.Unmarshal([]byte(body), &data)


	city := data.Name
	conditions := data.Weather[0].Description
	temperature := strconv.FormatFloat(data.Main.Temp, 'f', 2, 64)
	humidity := strconv.Itoa(data.Main.Humidity)
	wind := strconv.FormatFloat(data.Wind.Speed, 'f', 2, 64)

	embed := &discordgo.MessageSend {
		Embeds: []*discordgo.MessageEmbed {{
			Type: discordgo.EmbedTypeRich,
			Title: "Current Weather",
			Description: "Weather for " + city,
			Fields: []*discordgo.MessageEmbedField {
				{
					Name: "Conditions",
					Value: conditions,
					Inline: true,
				},
				{
					Name: "Temperature",
					Value: temperature + "F",
					Inline: true,
				},
				{
					Name: "Humidity",
					Value: humidity + "%",
					Inline: true,
				},
				{
					Name: "Wind",
					Value: wind + " mph",
					Inline: true,
				},
			},
		},
	},
	}

	return embed

	// fmt.Println(response)
	// return &discordgo.MessageSend{
	// 	Content: "executed all functions",
	// }
}