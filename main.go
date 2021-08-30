package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime"
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
		if val, ok := commands[strings.Split(m.Content, " ")[0][1:]]; ok {
			val.(func(*discordgo.Session, *discordgo.MessageCreate, []string))(s, m, strings.Split(m.Content, " ")[1:])
		}
	}
}

func ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "pong")
	println(m.Attachments[0].ProxyURL)
}

func postmeme(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(m.Attachments) < 1 {
		s.ChannelMessageSend(m.ChannelID, "No image found.")
		return
	}
	ext := m.Attachments[0].Filename[strings.LastIndex(m.Attachments[0].Filename, "."):]
	fn := m.Attachments[0].Filename
	println(args)
	if len(args) > 0 {
		fn = args[0] + ext
		println(fn)
	}
	go saveFile("memes/"+fn, m.Attachments[0].ProxyURL)
	s.ChannelMessageSend(m.ChannelID, "Meme posted!")
}

func updatestatus(s *discordgo.Session) {
	s.UpdateGameStatus(0, fmt.Sprintf("%d goroutines", runtime.NumGoroutine()))
}

func runtimestats(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```ngoroutines: %d\nalloc: %v\ntotalalloc: %v\nsys: %v\nnumgc: %v\n```", runtime.NumGoroutine(), ms.Alloc, ms.TotalAlloc, ms.Sys, ms.NumGC))
}

func main() {
	discord, err := discordgo.New("Bot " + Token)

	commands = make(map[string]interface{})
	commands["ping"] = ping
	commands["postmeme"] = postmeme
	commands["runtimestats"] = runtimestats

	if err != nil {
		panic(err)
	}

	discord.AddHandler(messageCreate)

	err = discord.Open()
	defer discord.Close()
	if err != nil {
		panic(err)
	}
	go updatestatus(discord)
	fmt.Println("Now running bot. Press CTRL+C to quit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}
