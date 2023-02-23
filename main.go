package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"

	"github.com/arthurfalcao/jellyfin-live-tv-favorites/infra/jellyfin"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	client := setupJellyfinClient()
	channels, err := client.GetChannels()
	if err != nil {
		return fmt.Errorf("error getting channel: %v", err)
	}

	wg := sync.WaitGroup{}

	favoriteChannels := []string{"HBO Max", "Premiere", "SporTV", "TNT", "NBA League Pass"}

	for _, channel := range channels {
		shouldFavorite := false
		for _, favoriteChannel := range favoriteChannels {
			if strings.Contains(channel.Name, favoriteChannel) {
				shouldFavorite = true
				break
			}
		}

		if !shouldFavorite {
			continue
		}

		fmt.Printf("channel: %v\n", channel)

		if channel.UserData.IsFavorite {
			continue
		}

		wg.Add(1)

		go func(channelID string, wg *sync.WaitGroup) {
			defer wg.Done()
			err := client.MarkFavoriteItem(channelID)
			if err != nil {
				fmt.Printf("error marking favorite item: %v", err)
			}
		}(channel.ID, &wg)
	}

	wg.Wait()

	return nil
}

func setupJellyfinClient() *jellyfin.Client {
	config := jellyfin.ClientConfig{
		BaseURL: os.Getenv("JELLYFIN_BASE_URL"),
		ApiKey:  os.Getenv("JELLYFIN_API_KEY"),
		UserID:  os.Getenv("JELLYFIN_USER_ID"),
	}
	return jellyfin.NewClient(config)
}
