package riot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type localizedNamesDto struct {
}

type contentItemDto struct {
	Name           string            `json:"name"`
	Id             string            `json:"id"`
	LocalizedNames localizedNamesDto `json:"localizedNames,omitempty"`
	AssetName      string            `json:"assetName"`
	AssetPath      string            `json:"assetPath"`
}

type actDto struct {
	Name           string            `json:"name"`
	LocalizedNames localizedNamesDto `json:"localizedNames,omitempty"`
	Id             string            `json:"id"`
	IsActive       bool              `json:"isActive"`
	Type           string            `json:"type"`
}

type contentDto struct {
	Version      string           `json:"version"`
	Characters   []contentItemDto `json:"characters"`
	Maps         []contentItemDto `json:"maps"`
	Chromas      []contentItemDto `json:"chromas"`
	Skins        []contentItemDto `json:"skins"`
	Equips       []contentItemDto `json:"equips"`
	GameModes    []contentItemDto `json:"gameModes"`
	Sprays       []contentItemDto `json:"sprays"`
	SprayLevels  []contentItemDto `json:"sprayLevels"`
	Charms       []contentItemDto `json:"charms"`
	CharmLevels  []contentItemDto `json:"charmLevels"`
	PlayerCards  []contentItemDto `json:"playerCards"`
	PlayerTitles []contentItemDto `json:"playerTitles"`
	Acts         []actDto         `json:"acts"`
}

type PlayerDto struct {
	Puuid           string `json:"puuid"`
	GameName        string `json:"gameName"`
	TagLine         string `json:"tagLine"`
	LeaderboardRank int    `json:"leaderboardRank"`
	RankedRating    int    `json:"rankedRating"`
	NumberOfWins    int    `json:"numberOfWins"`
}

type LeaderboardDto struct {
	Shard        string      `json:"shard"`
	ActId        string      `json:"actId"`
	TotalPlayers int         `json:"totalPlayers"`
	Players      []PlayerDto `json:"players"`
}

var apiKey string = "RGAPI-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

func GetActiveActId() string {
	var serviceUrl string = fmt.Sprintf("https://latam.api.riotgames.com/val/content/v1/contents?locale=en-US&api_key=%s", apiKey)
	response, err := http.Get(serviceUrl)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var contentDtoJSON contentDto

	errJSON := json.Unmarshal([]byte(responseData), &contentDtoJSON)

	if errJSON != nil {
		panic(errJSON)
	}

	// fmt.Println(string(responseData))
	// fmt.Printf("\n\n json object:::: %+v", contentDtoJSON)

	var activeActId string

	for _, act := range contentDtoJSON.Acts {
		if act.IsActive && act.Type == "act" {
			activeActId = act.Id
		}
	}

	return activeActId
}

func GetLeaderboard(actId string) LeaderboardDto {
	var serviceUrl = fmt.Sprintf("https://latam.api.riotgames.com/val/ranked/v1/leaderboards/by-act/%s?size=200&startIndex=0&api_key=%s", actId, apiKey)
	response, err := http.Get(serviceUrl)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var leaderboardDtoJSON LeaderboardDto

	errJSON := json.Unmarshal([]byte(responseData), &leaderboardDtoJSON)

	if errJSON != nil {
		panic(errJSON)
	}

	/*
		for _, player := range leaderboardDtoJSON.Players {
			fmt.Println("Returning: " + player.Puuid)
		}
	*/

	return leaderboardDtoJSON
}
