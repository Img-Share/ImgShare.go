package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var commands map[string]interface{}

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

func updatestatus(s *discordgo.Session) {
	s.UpdateGameStatus(0, "All hail the Gopher!")
}

func main() {
	discord, err := discordgo.New("Bot " + Token)

	commands = make(map[string]interface{})
	commands["ping"] = ping
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
