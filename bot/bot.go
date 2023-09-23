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
	case strings.Contains(message.Content, BotPrefix + "upmap") && len(message.Attachments) > 0:	
		names := downloadMaps(message.Attachments)
		discord.ChannelMessageSend(message.ChannelID, "Successfully uploaded " + strings.Join(names, ", "))
	case strings.Contains(message.Content, BotPrefix + "ping"):
		discord.ChannelMessageSend(message.ChannelID, "pong!")	
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
