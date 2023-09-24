package bot

import (
	"encoding/json"
	"net/http"
	"time"
)

var (
	URLHttpMaster string
	servers       Servers
)

type Servers struct {
	Servers []Server `json:"servers"`
}

type Server struct {
	Addresses []string `json:"addresses"`
	Location  string   `json:"location"`
	Info      Info     `json:"info"`
}

type Info struct {
	MaxClients      int      `json:"max_clients"`
	MaxPlayers      int      `json:"max_players"`
	Passworded      bool     `json:"passworded"`
	GameType        string   `json:"game_type"`
	Name            string   `json:"name"`
	Map             Map      `json:"map"`
	Version         string   `json:"version"`
	ClientScoreKind string   `json:"client_score_kind"`
	Clients         []Client `json:"clients"`
}

type Map struct {
	Name   string `json:"name"`
	SHA256 string `json:"sha256"`
	Size   int    `json:"size"`
}

type Client struct {
	Name     string `json:"name"`
	Clan     string `json:"clan"`
	Country  int    `json:"country"`
	Score    int    `json:"score"`
	IsPlayer bool   `json:"is_player"`
	Skin     Skin   `json:"skin"`
	AFK      bool   `json:"afk"`
	Team     int    `json:"team"`
}

type Skin struct {
	Name      string `json:"name"`
	ColorBody int    `json:"color_body"`
	ColorFeet int    `json:"color_feet"`
}

func getServers() error {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(URLHttpMaster)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(&servers)
}

func getOnlineByClan(clan string) (clients []Client) {
	for _, server := range servers.Servers {
		for _, client := range server.Info.Clients {
			if client.Clan == clan {
				clients = append(clients, client)
			}
		}
	}
	return clients
}
