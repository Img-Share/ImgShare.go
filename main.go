package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var commands map[string]interface{}

func saveFile(fp string, url string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	handle, err := os.Create(fp)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	defer handle.Close()
	_, err = io.Copy(handle, resp.Body)

	if err != nil {
		panic(err)
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	iscmd := strings.HasPrefix(m.Content, "#")
	if iscmd {
		// It's a command
		if val, ok := commands[m.Content[1:]]; ok {
			val.(func(*discordgo.Session, *discordgo.MessageCreate))(s, m)
		} else {
			s.ChannelMessageSend(m.ChannelID, "Unknown command.")
		}
	}
}

func ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "pong")
	println(m.Attachments[0].ProxyURL)
}

func postmeme(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Attachments) < 1 {
		s.ChannelMessageSend(m.ChannelID, "No image found.")
		return
	}
	go saveFile("memes/"+m.Attachments[0].Filename, m.Attachments[0].ProxyURL)
	s.ChannelMessageSend(m.ChannelID, "Meme posted!")
}

func updatestatus(s *discordgo.Session) {
	s.UpdateGameStatus(0, "بيننا مجانا 2021-2020 البنجابية الحرة")
}

func main() {
	discord, err := discordgo.New("Bot " + Token)

	commands = make(map[string]interface{})
	commands["ping"] = ping
	commands["postmeme"] = postmeme
	if err != nil {
		panic(err)
	}

	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		panic(err)
	}
	go updatestatus(discord)
	fmt.Println("Now running bot. Press CTRL+C to quit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	discord.Close()
}
