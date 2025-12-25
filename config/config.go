package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken string
	GuildID      string
	AutoRoleID   string
}

func LoadConfig() *Config {
	// Only load .env if it exists (local development)
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}
	}

	return &Config{
		DiscordToken: os.Getenv("DISCORD_TOKEN"),
		GuildID:      os.Getenv("GUILD_ID"),
		AutoRoleID:   os.Getenv("AUTO_ROLE_ID"),
	}
}
