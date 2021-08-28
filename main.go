package main

import "github.com/bwmarrin/discordgo"

func main() {
	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		panic(err)
	}
	err = discord.Open()
	if err != nil {
		panic(err)
	}
	discord.Close()
}
