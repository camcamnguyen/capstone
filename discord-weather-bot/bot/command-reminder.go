package bot

import (
	"fmt"
	"log"
	"os"
	//"os/signal"
	"encoding/csv"
	"io"
	//"strings"
	"github.com/bwmarrin/discordgo"
)



func getReminders(message string, author *discordgo.User) *discordgo.MessageSend {
	//var all_reminders [10]string
	var all_reminders string
	path, err := os.Getwd()
	
	username := author.Username
	path += "/bot/" + username + "reminders.csv"
	file, err := os.Open(path)
    if err != nil {
        log.Fatal(err)
    }

    reader := csv.NewReader(file)
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatal(err)
        }
        for value := range record {
            //fmt.Printf("%s\n", record[value])
			//store the values
			//all_reminders[value] = record[value]
			all_reminders += record[value] + "\n"
        }
    }

	embed := &discordgo.MessageSend {
		Embeds: []*discordgo.MessageEmbed {{
			Type: discordgo.EmbedTypeRich,
			Title: "To do:",
			Fields: []*discordgo.MessageEmbedField {
				{
					Name: "reminders",
					Value: all_reminders,
					Inline: true,
				},
			},
		},
	},
	}
	return embed

}

func setReminder(message string, author *discordgo.User) {
	// r, _ := regexp.Compile(`\d{5}`) //regular expression for zipcode
	// zip := r.FindString(message)
	path, err := os.Getwd()

	username := author.Username
	path += "/bot/" + username + "reminders.csv"
	//path += "/bot/reminders.csv"
	//file, err := os.Open(path + "/bot/reminders.csv")
	//file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)

    if err != nil {
        log.Fatal(err)
    }

	writer := csv.NewWriter(file)
	substring := message[10:len(message)]

	data := []string{substring}
	writer.Write(data)
	writer.Flush()

	fmt.Println(username)
}