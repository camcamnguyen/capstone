package bot

import (
	//"fmt"
	"log"
	"os"
	//"os/signal"
	"encoding/csv"
	"io"
	//"strings"
	"github.com/bwmarrin/discordgo"
)



func getReminders(message string) *discordgo.MessageSend {
	//var all_reminders [10]string
	var all_reminders string
	path, err := os.Getwd();
	f, err := os.Open(path + "/bot/reminders.csv")
    if err != nil {
        log.Fatal(err)
    }

    r := csv.NewReader(f)
    for {
        record, err := r.Read()
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
			all_reminders += record[value]
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

func setReminders() {
	// r, _ := regexp.Compile(`\d{5}`) //regular expression for zipcode
	// zip := r.FindString(message)
}