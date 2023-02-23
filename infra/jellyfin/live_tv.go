package jellyfin

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Channel struct {
	ID       string `json:"Id"`
	Name     string `json:"Name"`
	UserData struct {
		IsFavorite bool `json:"IsFavorite"`
	} `json:"UserData"`
}

func (c *Client) GetChannels() ([]Channel, error) {
	url := fmt.Sprintf("%s/LiveTv/Channels", c.config.BaseURL)

	req, err := c.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return []Channel{}, fmt.Errorf("client :: could not create request: %v", err)
	}

	q := req.URL.Query()
	q.Add("UserId", c.config.UserID)
	req.URL.RawQuery = q.Encode()

	res, err := c.client.Do(req)
	if err != nil {
		return []Channel{}, fmt.Errorf("client :: error making http request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return []Channel{}, fmt.Errorf("client :: unexpected status code: %d", res.StatusCode)
	}

	type channelsResponse struct {
		Items            []Channel `json:"Items"`
		TotalRecordCount int       `json:"TotalRecordCount"`
		StartIndex       int       `json:"StartIndex"`
	}

	response := channelsResponse{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return []Channel{}, fmt.Errorf("client :: error decoding response body: %v", err)
	}

	return response.Items, nil
}

func (c *Client) GetChannel(channelID string) (Channel, error) {
	url := fmt.Sprintf("%s/LiveTv/Channels/%s", c.config.BaseURL, channelID)

	req, err := c.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Channel{}, fmt.Errorf("client :: could not create request: %v", err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return Channel{}, fmt.Errorf("client :: error making http request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Channel{}, fmt.Errorf("client :: unexpected status code: %d", res.StatusCode)
	}

	channel := Channel{}
	if err := json.NewDecoder(res.Body).Decode(&channel); err != nil {
		return Channel{}, fmt.Errorf("client :: error decoding response body: %v", err)
	}

	return channel, nil
}

func (c *Client) MarkFavoriteItem(itemID string) error {
	url := fmt.Sprintf("%s/Users/%s/FavoriteItems/%s", c.config.BaseURL, c.config.UserID, itemID)

	req, err := c.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("client :: could not create request: %v", err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("client :: error making http request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("client :: unexpected status code: %d", res.StatusCode)
	}

	return nil
}

func (c *Client) UnMarkFavoriteItem(itemID string) error {
	url := fmt.Sprintf("%s/Users/%s/FavoriteItems/%s", c.config.BaseURL, c.config.UserID, itemID)

	req, err := c.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("client :: could not create request: %v", err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("client :: error making http request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("client :: unexpected status code: %d", res.StatusCode)
	}

	return nil
}
