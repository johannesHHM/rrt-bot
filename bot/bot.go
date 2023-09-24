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
	BotToken string
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
	case strings.Contains(message.Content, BotPrefix + "upmap"):
		if len(message.Attachments) == 0 {
			discord.ChannelMessageSend(message.ChannelID, "To uppload a map, include the map as an attachment")
			break
		}
		names := downloadMaps(message.Attachments)
		if len(names) > 0 {
			discord.ChannelMessageSend(message.ChannelID, "Successfully uploaded:  " + strings.Join(names, "     "))
		} else {
			discord.ChannelMessageSend(message.ChannelID, "Failed to upload attachment, is your attachment a map?")
		}

	case strings.Contains(message.Content, BotPrefix + "ping"):
		discord.ChannelMessageSend(message.ChannelID, "pong!")

	case strings.Contains(message.Content, BotPrefix + "lsmap"):
		names := getMapsList(-1)
		if len(names) > 0 {
		discord.ChannelMessageSend(message.ChannelID, "Maps:  " + strings.Join(names, "     "))
		} else {
			discord.ChannelMessageSend(message.ChannelID, "There are no maps in map dir")
		}
	
	case strings.Contains(message.Content, BotPrefix + "rmmap"):
		successNames, failureNames := removeMaps(message.Content)
		if len(successNames) > 0 {
			discord.ChannelMessageSend(message.ChannelID, "Successfully removed:  " + strings.Join(successNames, "     "))
		}
		if len(failureNames) > 0 {
			discord.ChannelMessageSend(message.ChannelID, "Failed to remove:  " + strings.Join(failureNames, "     "))
		}	
	}
}

func downloadMaps(attachments []*discordgo.MessageAttachment) ([]string) {
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

func downloadMap(attachment *discordgo.MessageAttachment) (bool) {
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
	if err != nil  {
		log.Println("Could not copy res to out")
		stat = false
	}
	return stat
}

func getMapsList(amount int) ([]string) {
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
	content = strings.ReplaceAll(content, BotPrefix + "rmmap", "")
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

func removeMap(name string) (bool) {
	err := os.Remove(BotMapDir + name + ".map")
	if err != nil {
		return false
	} else {
		return true
	}
}
