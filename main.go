package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/arthurfalcao/jellyfin-live-tv-favorites/infra/jellyfin"
	"github.com/arthurfalcao/jellyfin-live-tv-favorites/usecase"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
}

func main() {
	client := setupJellyfinClient()
	useCase := setupFavoriteChannelsUseCase(*client)

	err := useCase.FavoriteChannels([]string{"HBO Max", "Premiere", "SporTV", "TNT", "NBA League Pass"})
	if err != nil {
		log.Fatalf("error to favorite channels: %v", err)
	}
}

func setupFavoriteChannelsUseCase(jellyfinClient jellyfin.Client) usecase.UseCaseChannel {
	return usecase.NewUseCaseChannel(jellyfinClient)
}

func setupJellyfinClient() *jellyfin.Client {
	config := jellyfin.ClientConfig{
		BaseURL: os.Getenv("JELLYFIN_BASE_URL"),
		ApiKey:  os.Getenv("JELLYFIN_API_KEY"),
		UserID:  os.Getenv("JELLYFIN_USER_ID"),
	}
	return jellyfin.NewClient(config)
}
