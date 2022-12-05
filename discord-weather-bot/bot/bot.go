package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"github.com/bwmarrin/discordgo"
)

var(
	OpenWeatherToken string //exported variables - notice the uppercase letters
	BotToken string
)

func Run() {
	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatal(err)
	}

	discord.AddHandler(newMessage)

	//open session
	discord.Open()
	defer discord.Close()

	//run until code is terminated
	fmt.Println("Bot running... ")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<- c
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == discord.State.User.ID {
		return
	}

	switch{
	case strings.Contains(message.Content, "weather"):
		discord.ChannelMessageSend(message.ChannelID, "I can help with that! Use '!zip <zip code>'")
	case strings.Contains(message.Content, "gopher"):
		discord.ChannelMessageSend(message.ChannelID, "Hello there")
	case strings.Contains(message.Content, "!zip"):
		currentWeather := getCurrentWeather(message.Content)
		discord.ChannelMessageSendComplex(message.ChannelID, currentWeather)
	}
}
