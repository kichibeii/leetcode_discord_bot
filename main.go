package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Stat struct {
	QuestionID         int    `json:"question_id"`
	QuestionTitle      string `json:"question__title"`
	QuestionTitleSlug  string `json:"question__title_slug"`
	TotalACS           int    `json:"total_acs"`
	TotalSubmitted     int    `json:"total_submitted"`
	FrontendQuestionID int    `json:"frontend_question_id"`
	IsNewQuestion      bool   `json:"is_new_question"`
}

type Difficulty struct {
	Level int `json:"level"`
}

type Item struct {
	Stat       Stat       `json:"stat"`
	Difficulty Difficulty `json:"difficulty"`
	PaidOnly   bool       `json:"paid_only"`
	IsFavor    bool       `json:"is_favor"`
	Frequency  int        `json:"frequency"`
	Progress   int        `json:"progress"`
}

type Data struct {
	Data []Item `json:"data"`
}

// difficulty
const (
	DifficultyEasy   = 1
	DifficultyMedium = 2
	DifficultyHard   = 3
)

func main() {
	// reading file
	var data Data
	content, _ := os.ReadFile("question_list.json")

	err := json.Unmarshal(content, &data)
	if err != nil {
		panic(err)
	}

	// selecting Data
	tempData := Item{}
	loop := true
	for loop {
		rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
		key := rand.Intn(len(data.Data))
		tempData = data.Data[key]
		if !tempData.PaidOnly && tempData.Difficulty.Level == DifficultyEasy {
			data.Data = append(data.Data[:key], data.Data[key+1])
			loop = false
		}
	}

	link := "https://leetcode.com/problems/" + tempData.Stat.QuestionTitleSlug
	difficulty := ""

	now := time.Now()
	dateString := now.Format("2006-01-02")

	switch tempData.Difficulty.Level {
	case DifficultyMedium:
		difficulty = "medium"
	case DifficultyEasy:
		difficulty = "easy"
	case DifficultyHard:
		difficulty = "hard"
	}

	// print message to discord
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./main <configFilePath>")
		return
	}

	configFilePath := os.Args[1]

	keyBot, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	keyAuth := fmt.Sprintf("Bot %s", string(keyBot))
	keyAuth = strings.TrimSuffix(keyAuth, "\n")
	discord, err := discordgo.New(keyAuth)
	if err != nil {
		fmt.Println(err)
		return
	}
	discord.Open()

	message := fmt.Sprintf("Date : %s \nLink : %s \nDifficulty  : %s \n \n <@!1165437206329032734>", dateString, link, difficulty)

	_, err = discord.ChannelMessageSend("1165116574290673785", message)
	fmt.Println("error send message", err)
}
