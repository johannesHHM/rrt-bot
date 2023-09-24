package bot

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bwmarrin/discordgo"
)

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
