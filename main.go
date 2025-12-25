package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"manager-bot/config"

	"github.com/bwmarrin/discordgo"
)

func main() {
	cfg := config.LoadConfig()

	if cfg.DiscordToken == "" {
		log.Fatal("DISCORD_TOKEN is not set")
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// Identify intents
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMembers

	// Register the GuildMemberAdd handler
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
		if cfg.AutoRoleID == "" {
			return
		}

		err := s.GuildMemberRoleAdd(m.GuildID, m.User.ID, cfg.AutoRoleID)
		if err != nil {
			log.Printf("Error adding auto role to %s: %v", m.User.Username, err)
		} else {
			log.Printf("Successfully added role %s to user %s", cfg.AutoRoleID, m.User.Username)
		}
	})

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}

	log.Println("Manager Bot is running...")

	// Health check for Render.com
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Manager Bot is running"))
	})

	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		log.Printf("Starting HTTP server on :%s", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
