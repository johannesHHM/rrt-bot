package bot

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	BotToken  string
	BotMapDir string
	BotPrefix string
)

func Run() {
	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatal(err)
	}

	discord.AddHandler(newMessage)
	discord.Open()

	defer discord.Close()

	fmt.Println("Bot running ...")

	// Wait for os interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == discord.State.User.ID {
		return
	}
	switch {
	case checkCommand("upmap", message):
		runUpmap(discord, message)

	case checkCommand("lsmap", message):
		runLsmap(discord, message)

	case checkCommand("rmmap", message):
		runRmmap(discord, message)

	case checkCommand("help", message):
		runHelp(discord, message)

	case checkCommand("ping", message):
		runPing(discord, message)
	}

}

func downloadMaps(attachments []*discordgo.MessageAttachment) []string {
	var successNames []string
	for _, attachment := range attachments {
		if strings.HasSuffix(attachment.Filename, ".map") {
			success := downloadMap(attachment)
			if success {
				successNames = append(
					successNames,
					strings.TrimSuffix(
						attachment.Filename,
						filepath.Ext(attachment.Filename)))
			}
		}
	}
	return successNames
}

func downloadMap(attachment *discordgo.MessageAttachment) bool {
	stat := true
	res, err := http.DefaultClient.Get(attachment.URL)
	if err != nil {
		log.Println("Could not get attachment from URL")
		stat = false
	}
	defer res.Body.Close()

	out, err := os.Create(BotMapDir + attachment.Filename)
	if err != nil {
		log.Println("Could not create file")
		stat = false
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		log.Println("Could not copy res to out")
		stat = false
	}
	return stat
}

func getMapsList(amount int) []string {
	folder, err := os.Open(BotMapDir)
	if err != nil {
		log.Println("Could not get maps folder")
	}
	defer folder.Close()

	names, err := folder.Readdirnames(amount)
	if err != nil {
		log.Println("Could not get map names")
	}
	for i, name := range names {
		names[i] = strings.TrimSuffix(name, filepath.Ext(name))
	}

	return names
}

func removeMaps(content string) ([]string, []string) {
	var successNames, failureNames []string
	content = strings.ReplaceAll(content, BotPrefix+"rmmap", "")
	names := strings.Fields(content)
	for _, name := range names {
		if removeMap(name) {
			successNames = append(successNames, name)
		} else {
			failureNames = append(failureNames, name)
		}
	}
	return successNames, failureNames
}

func removeMap(name string) bool {
	err := os.Remove(BotMapDir + name + ".map")
	if err != nil {
		return false
	} else {
		return true
	}
}

func checkCommand(command string, message *discordgo.MessageCreate) bool {
	return strings.Contains(message.Content, BotPrefix+command)
}

func runUpmap(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if len(message.Attachments) < 1 {
		discord.ChannelMessageSend(message.ChannelID, "```To uppload a map, include the map as an attachment```")
		return
	}
	names := downloadMaps(message.Attachments)
	if len(names) < 1 {
		discord.ChannelMessageSend(message.ChannelID, "```Failed to upload attachment, is your attachment a map?```")

	}
	replyString := "```" +
		"Successfully uploaded:\n" +
		strings.Join(names, "     ") +
		"```"
	discord.ChannelMessageSend(message.ChannelID, replyString)
}

func runLsmap(discord *discordgo.Session, message *discordgo.MessageCreate) {
	names := getMapsList(-1)
	if len(names) < 1 {
		discord.ChannelMessageSend(message.ChannelID, "```There are no maps in map dir```")
	}
	replyString := "```" +
		"Maps:\n" +
		strings.Join(names, "     ") +
		"```"
	discord.ChannelMessageSend(message.ChannelID, replyString)
}

func runRmmap(discord *discordgo.Session, message *discordgo.MessageCreate) {
	successNames, failureNames := removeMaps(message.Content)
	if len(successNames) > 0 {
		replyString := "```" +
			"Successfully removed:\n" +
			strings.Join(successNames, "     ") +
			"```"
		discord.ChannelMessageSend(message.ChannelID, replyString)
	}
	if len(failureNames) > 0 {
		replyString := "```" +
			"Failed to remove:\n" +
			strings.Join(failureNames, "     ") +
			"```"
		discord.ChannelMessageSend(message.ChannelID, replyString)
	}
}

func runHelp(discord *discordgo.Session, message *discordgo.MessageCreate) {
	helpString := "```" +
		"Commands:\n" +
		"!upmap [map files]   uploads all given maps\n" +
		"!lsmap               lists all maps\n" +
		"!rmmap [map names]   removes given maps\n" +
		"!ping                pongs\n" +
		"```"
	discord.ChannelMessageSend(message.ChannelID, helpString)
}

func runPing(discord *discordgo.Session, message *discordgo.MessageCreate) {
	discord.ChannelMessageSend(message.ChannelID, "```pong!```")
}
