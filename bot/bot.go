package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
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

	case checkCommand("online", message):
		runOnline(discord, message)

	case checkCommand("help", message):
		runHelp(discord, message)

	case checkCommand("ping", message):
		runPing(discord, message)
	}

}

func checkCommand(command string, message *discordgo.MessageCreate) bool {
	return strings.Contains(message.Content, BotPrefix+command)
}

func runUpmap(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if len(message.Attachments) < 1 {
		discord.ChannelMessageSend(
			message.ChannelID,
			"```To uppload a map, include the map as an attachment```")
		return
	}
	names := downloadMaps(message.Attachments)
	if len(names) < 1 {
		discord.ChannelMessageSend(
			message.ChannelID,
			"```Failed to upload attachment, is your attachment a map?```")

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
		discord.ChannelMessageSend(
			message.ChannelID,
			"```There are no maps in map dir```")
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
		"!upmap [map files]   uploads given maps\n" +
		"!lsmap               lists all maps\n" +
		"!rmmap [names]       removes maps by names\n" +
		"!online              lists all online RRT members\n" +
		"!ping                pongs\n" +
		"```"
	discord.ChannelMessageSend(message.ChannelID, helpString)
}

func runPing(discord *discordgo.Session, message *discordgo.MessageCreate) {
	discord.ChannelMessageSend(message.ChannelID, "```pong!```")
}

func runOnline(discord *discordgo.Session, message *discordgo.MessageCreate) {
	err := getServers()
	if err != nil {
		log.Println("Timeout from master1.ddnet.org")
		discord.ChannelMessageSend(
			message.ChannelID,
			"```Could not reach master server!```")
		return
	}
	clients := getOnlineByClan("ℜℜͲ")
	var names []string
	for _, client := range clients {
		names = append(names, client.Name)
	}
	if len(names) < 1 {
		discord.ChannelMessageSend(
			message.ChannelID,
			"```0 online RRT members```")
		return
	}
	replyString := "```" +
		strconv.Itoa(len(names)) + " online RRT members:\n" +
		strings.Join(names, "     ") +
		"```"
	discord.ChannelMessageSend(message.ChannelID, replyString)
}
