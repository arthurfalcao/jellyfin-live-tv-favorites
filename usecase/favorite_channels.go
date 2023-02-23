package usecase

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/arthurfalcao/jellyfin-live-tv-favorites/infra/jellyfin"
)

type UseCaseChannel struct {
	jellyfinClient jellyfin.Client
}

func NewUseCaseChannel(jellyfinClient jellyfin.Client) UseCaseChannel {
	return UseCaseChannel{
		jellyfinClient: jellyfinClient,
	}
}

func (u UseCaseChannel) FavoriteChannels(favoriteChannels []string) error {
	channels, err := u.jellyfinClient.GetChannels()
	if err != nil {
		return fmt.Errorf("error getting channel: %v", err)
	}

	wg := sync.WaitGroup{}

	for _, channel := range channels {
		if !u.shouldFavorite(channel.Name, favoriteChannels) {
			continue
		}

		if channel.UserData.IsFavorite {
			log.Printf("channel %s is already favorite", channel.Name)
			continue
		}

		log.Printf("channel %s is not favorite", channel.Name)
		wg.Add(1)

		go func(channel jellyfin.Channel) {
			defer wg.Done()

			err := u.jellyfinClient.MarkFavoriteItem(channel.ID)
			if err != nil {
				log.Printf("error marking favorite item: %v", err)
			}
		}(channel)
	}

	wg.Wait()

	return nil
}

func (u UseCaseChannel) shouldFavorite(channelName string, favoriteChannels []string) bool {
	for _, favoriteChannel := range favoriteChannels {
		if !strings.Contains(channelName, favoriteChannel) {
			continue
		}

		return true
	}

	return false
}
