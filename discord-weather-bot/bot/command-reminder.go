package bot

import (

	//"os/signal"
	//"strings"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)
func help() *discordgo.MessageSend {
	var helpString string

	helpString += "1.'!zip': Returns weather stats for a zip code\n2.'!reminders': Returns a list of all your reminders\n"
	helpString += "3.'!addreminders': Add a reminder to the list\n4.'!remindme': Sets a timed reminder \n5.'!clear' Clears all reminders"
	embed := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{{
			Type:  discordgo.EmbedTypeRich,
			Title: "Help:",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Commands",
					Value:  helpString,
					Inline: true,
				},
			},
		},
		},
	}
	return embed
}
func getReminders(message string, author *discordgo.User) *discordgo.MessageSend {
	var all_reminders string
	path, err := os.Getwd()

	username := author.Username
	path += "/bot/" + username + "reminders.csv"
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	file.Close()
	readFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(readFile)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		for value := range record {
			all_reminders += record[value] + "\n"
		}
	}
	if all_reminders == "" {
		fmt.Println("Empty")
		embed := &discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{{
				Type:  discordgo.EmbedTypeRich,
				Title: "To do:",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "reminders",
						Value:  "No Reminders! Enjoy your day",
						Inline: true,
					},
				},
			},
			},
		}
		readFile.Close()
		return embed
	} else {
		embed := &discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{{
				Type:  discordgo.EmbedTypeRich,
				Title: "To do:",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "reminders",
						Value:  all_reminders,
						Inline: true,
					},
				},
			},
			},
		}
		readFile.Close()
		return embed
	}

}

func setReminder(message string, author *discordgo.User) {
	path, err := os.Getwd()

	//time.Sleep(8 * time.Second)
	username := author.Username
	path += "/bot/" + username + "reminders.csv"
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	readFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	filedata, err := csv.NewReader(readFile).ReadAll()
	if err != nil {
		log.Println(err)
		return
	}
	listIndex := len(filedata)
	listIndex++

	writer := csv.NewWriter(file)
	substring := strconv.Itoa(listIndex) + ". " + message[13:]

	data := []string{substring}
	writer.Write(data)
	writer.Flush()
	file.Close()
}

func setTimedReminder(message string, timer string, units string) {

	intVar, err := strconv.Atoi(timer)
	if err != nil {
		log.Fatal(err)
	}
	if units == "seconds" {
		time.Sleep(time.Duration(intVar) * time.Second)
	}
	if units == "minutes" {
		time.Sleep(time.Duration(intVar) * time.Minute)
	}
	if units == "hours" {
		time.Sleep(time.Duration(intVar) * time.Hour)
	}
}

func clearReminder(author *discordgo.User) {
	path, err := os.Getwd()
	username := author.Username
	path += "/bot/" + username + "reminders.csv"

	writefile, err := os.Create(path)
	writefile.Close()
	if err != nil {
		log.Fatal(err)
	}

}
