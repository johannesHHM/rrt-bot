package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"unicode"

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

	case checkCommand("gt", message):
		runTranslate(discord, message)

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
	mapInfos := getMapsList()
	if len(mapInfos) < 1 {
		discord.ChannelMessageSend(
			message.ChannelID,
			"```There are no maps in map dir```")
	} else {
		replyString := "```Maps:\n------------------------------------------------\n"
		for _, mapInfo := range mapInfos {
			if !mapInfo.IsDir() {
				lineString := mapInfo.Name()
				lineString = rightPad(lineString, 24)
				lineString += byteCountIEC(mapInfo.Size()) + " "
				lineString = rightPad(lineString, 36)
				lineString += mapInfo.ModTime().Format("Jan 2 15:04") + "\n"
				replyString += lineString
			}
		}
		replyString += "```"

		discord.ChannelMessageSend(message.ChannelID, replyString)
	}
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
		"-------------------------------------------------\n" +
		"!help                lists all commands\n" +
		"!upmap [map files]    uploads given maps\n" +
		"!lsmap               lists all maps\n" +
		"!rmmap [names]       removes maps by names\n" +
		"!online              lists all online RRT members\n" +
		"!gt [msg] :[code]    translates message\n" +
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
		log.Println("Timeout from master server")
		discord.ChannelMessageSend(
			message.ChannelID,
			"```Could not reach master server!```")
		return
	}
	clients, servers := getOnlineByClan("ℜℜͲ")

	if len(clients) < 1 {
		discord.ChannelMessageSend(
			message.ChannelID,
			"```0 online RRT members```")
		return
	}

	replyString := "```" + strconv.Itoa(len(clients)) + " online RRT members:\n" +
		"-----------------------------------------------\n"

	for i, client := range clients {
		lineString := client.Name
		if client.AFK {
			lineString += " (AFK)"
		}
		lineString = rightPad(lineString, 21)
		lineString += " " + getServerShortName(servers[i])
		lineString += " (" +
			strconv.Itoa(len(servers[i].Info.Clients)) +
			"/" +
			strconv.Itoa(servers[i].Info.MaxClients) +
			")"
		replyString += lineString + "\n"
	}
	replyString += "```"

	discord.ChannelMessageSend(message.ChannelID, replyString)
}

func runTranslate(discord *discordgo.Session, message *discordgo.MessageCreate) {
	result := translateMessage(message.Content)
	if checkWhitespaceOnly(result.Text) {
		discord.ChannelMessageSend(message.ChannelID, "```Translation requires a message```")
		return
	}
	replyString := "```" + strings.ToUpper(result.Dest) + ": " + result.Text + "```"
	discord.ChannelMessageSend(message.ChannelID, replyString)
}

func rightPad(string string, length int) (paddedString string) {
	neededPad := length - len(string)
	if neededPad < 1 {
		return string
	}
	return string + strings.Repeat(" ", neededPad)
}

func byteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func checkWhitespaceOnly(content string) bool {
	for _, c := range content {
		if !unicode.IsSpace(c) {
			return false
		}
	}
	return true
}
