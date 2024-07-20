package soundcloud

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type scData struct {
	Version      float32 `json:"version"`
	Type         string  `json:"type"`
	ProviderName string  `json:"provider_name"`
	ProviderURL  string  `json:"provider_url"`
	Height       int     `json:"height"`
	Width        string  `json:"width"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	ThumbnailURL string  `json:"thumbnail_url"`
	HTML         string  `json:"html"`
	AuthorName   string  `json:"author_name"`
	AuthorURL    string  `json:"author_url"`
}

type spData struct {
	HTML            string `json:"html"`
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	Version         string `json:"version"`
	ProviderName    string `json:"provider_name"`
	ProviderURL     string `json:"provider_url"`
	Type            string `json:"type"`
	Title           string `json:"title"`
	ThumbnailURL    string `json:"thumbnail_url"`
	ThumbnailWidth  int    `json:"thumbnail_width"`
	ThumbnailHeight int    `json:"thumbnail_height"`
}

type EmbededPlayer struct {
	Id   string `json:"id"`
	Html string `json:"html"`
}

func GetEmbed(link string) (*EmbededPlayer, error) {
	client := &http.Client{}
	msgUrl := url.QueryEscape(link)

	if strings.Contains(msgUrl, "soundcloud") {
		return parseSoundCloudLink(client, msgUrl)
	} else {
		return parseSpotifyLink(client, msgUrl)
	}
}

func parseSoundCloudLink(client *http.Client, scUrl string) (*EmbededPlayer, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://soundcloud.com/oembed?format=json&maxheight=166&url=%s", scUrl), nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	var data scData
	err = json.Unmarshal(body, &data)

	if err != nil {
		return nil, err
	}

	return &EmbededPlayer{
		Id:   fmt.Sprintf("%s/%s", data.AuthorName, data.Title),
		Html: data.HTML,
	}, nil
}
func parseSpotifyLink(client *http.Client, spUrl string) (*EmbededPlayer, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://open.spotify.com/oembed?height=166&url=%s", spUrl), nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	var data spData
	err = json.Unmarshal(body, &data)

	if err != nil {
		return nil, err
	}

	return &EmbededPlayer{
		Id:   fmt.Sprintf("%s/%s", data.ThumbnailURL, data.Title),
		Html: data.HTML,
	}, nil
}
