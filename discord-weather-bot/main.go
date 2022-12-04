package main

import (
	"discord-weather-bot/bot"
	"log"
	"os"
)

func main() {
	//load environment variables
	botToken, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		log.Fatal("Must set Discord token as env variable")
	}
	openWeatherToken, ok := os.LookupEnv("OPENWEATHER_TOKEN")
	if !ok {
		log.Fatal("Must set Open weather token as an env variable")
	}

	//start the bot
	bot.BotToken = botToken
	bot.OpenWeatherToken = openWeatherToken
	bot.Run()
}