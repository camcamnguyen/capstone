package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	OpenWeatherToken string //exported variables - notice the uppercase letters
	BotToken         string
)

type GifSearch struct {
	Data struct {
		Type             string `json:"type"`
		Id               string `json:"id"`
		Url              string `json:"url"`
		Slug             string `json:"slug"`
		BitlyGifUrl      string `json:"bitly_gif_url"`
		BitlyUrl         string `json:"bitly_url"`
		EmbedUrl         string `json:"embed_url"`
		Username         string `json:"username"`
		Source           string `json:"source"`
		Title            string `json:"title"`
		Rating           string `json:"rating"`
		ContentUrl       string `json:"content_url"`
		SourceTld        string `json:"source_tld"`
		SourcePostUrl    string `json:"source_post_url"`
		IsSticker        int    `json:"is_sticker"`
		ImportDatetime   string `json:"import_datetime"`
		TrendingDatetime string `json:"trending_datetime"`
		Images           struct {
			FixedWidthStill struct {
				Height string `json:"height"`
				Size   string `json:"size"`
				Url    string `json:"url"`
				Width  string `json:"width"`
			} `json:"fixed_width_still"`
		} `json:"images"`
	} `json:"data"`
	Meta struct {
		Msg        string `json:"msg"`
		Status     int    `json:"status"`
		ResponseId string `json:"response_id"`
	} `json:"meta"`
}


func Run() {
	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatal(err)
	}

	discord.AddHandler(newMessage)
	discord.AddHandler(messageCreate)

	//open session
	discord.Open()
	defer discord.Close()

	//run until code is terminated
	fmt.Println("Bot running... ")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == discord.State.User.ID {
		return
	}

	switch {
	case strings.Contains(message.Content, "weather"):
		discord.ChannelMessageSend(message.ChannelID, "I can help with that! Use '!zip <zip code>'")
	case strings.Contains(message.Content, "!help"):
		helpMessage := help()
		discord.ChannelMessageSend(message.ChannelID, "Here is a list of commands!")
		discord.ChannelMessageSendComplex(message.ChannelID, helpMessage)
	case strings.Contains(message.Content, "hi gopher"):
		err := godotenv.Load(".env")
		giphyToken := os.Getenv("GIPHY_TOKEN")
		if err != nil {
			log.Fatal(err)
		}
		discord.ChannelMessageSend(message.ChannelID, "Hello there")
		url := "https://api.giphy.com/v1/gifs/random"
		var result GifSearch

		// 6
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error in making a new Request", err)
		}
		query := req.URL.Query()
		query.Add("api_key", giphyToken)
		query.Add("tag", "hello there")
		req.URL.RawQuery = query.Encode()
		client := http.Client{}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Error in getting a response, ", err)
		}
		body, _ := ioutil.ReadAll(res.Body)
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("Can not unmarshall JSON", err)
		}
		// 7
		discord.ChannelMessageSend(message.ChannelID, result.Data.EmbedUrl)
		res.Body.Close()
		
	case strings.Contains(message.Content, "!zip"):
		currentWeather := getCurrentWeather(message.Content)
		discord.ChannelMessageSendComplex(message.ChannelID, currentWeather)
	case strings.Contains(message.Content, "!reminders"):
		reminders := getReminders(message.Content, message.Author)
		discord.ChannelMessageSendComplex(message.ChannelID, reminders)
	case strings.Contains(message.Content, "!addreminder "):
		setReminder(message.Content, message.Author)
		discord.ChannelMessageSend(message.ChannelID, "Reminder added")
		reminders := getReminders(message.Content, message.Author)
		discord.ChannelMessageSendComplex(message.ChannelID, reminders)
	case strings.Contains(message.Content, "!remindme "):
		r, _ := regexp.Compile("[0-9]+")
		timer := r.FindString(message.Content)
		messageString := string(message.Content)
		var units = ""
		if strings.Contains(messageString, "seconds") {
			units = "seconds"
		} else if strings.Contains(messageString, "minutes") {
			units = "minutes"
		} else if strings.Contains(messageString, "hours") {
			units = "hours"
		}
		discord.ChannelMessageSend(message.ChannelID, "Reminding you in "+timer+" "+units)
		setTimedReminder(message.Content, timer, units)
		reminders := getReminders(message.Content, message.Author)
		discord.ChannelMessageSendComplex(message.ChannelID, reminders)
	case strings.Contains(message.Content, "!clear"):
		clearReminder(message.Author)
		discord.ChannelMessageSend(message.ChannelID, "Reminders cleared	")
	}
}

func messageCreate(s *discordgo.Session, message *discordgo.MessageCreate) {
	// 1
	err := godotenv.Load(".env")
	giphyToken := os.Getenv("GIPHY_TOKEN")
	if err != nil {
		log.Fatal(err)
	}
	// 2
	if message.Author.ID == s.State.User.ID {
		return
	}
	// 3
	command := strings.Split(message.Content, " ")
	// 4
	if command[0] == "!search" && len(command) > 1 {
		url := "https://api.giphy.com/v1/gifs/random"
		var result GifSearch
		// 5
		gifKeyword := strings.Join(command[1:], " ")

		// 6
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error in making a new Request", err)
		}
		query := req.URL.Query()
		query.Add("api_key", giphyToken)
		query.Add("tag", gifKeyword)
		req.URL.RawQuery = query.Encode()
		client := http.Client{}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Error in getting a response, ", err)
		}
		body, _ := ioutil.ReadAll(res.Body)
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("Can not unmarshall JSON", err)
		}
		// 7
		s.ChannelMessageSend(message.ChannelID, result.Data.EmbedUrl)
		res.Body.Close()
	}
}