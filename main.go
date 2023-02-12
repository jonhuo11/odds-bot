package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := os.Getenv("ODDSBOTTOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Panic(err)
	}

	storeDao, err = newMockStore()
	if err != nil {
		log.Panic(err)
	}

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	if err = dg.Open(); err != nil {
		log.Panic(err)
	}
	defer dg.Close()

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, "", &v)
		if err != nil {
			log.Panicf("cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	log.Println("odds-bot is now running, press CTRL-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	for _, v := range registeredCommands {
		err := dg.ApplicationCommandDelete(dg.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("cannot delete '%v' command: %v", v.Name, err)
		}
	}

	log.Println("shut down gracefully")
}
