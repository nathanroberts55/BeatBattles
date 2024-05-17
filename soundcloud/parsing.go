package soundcloud

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

type SoundcloudItem struct {
	Id   string `json:"id"`
	Html string `json:"html"`
}

func GetEmbed(link string) (*SoundcloudItem, error) {
	client := &http.Client{}
	scUrl := url.QueryEscape(link)
	req, err := http.NewRequest("GET", fmt.Sprintf("https://soundcloud.com/oembed?format=json&url=%s", scUrl), nil)
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

	return &SoundcloudItem{
		Id:   fmt.Sprintf("%s/%s", data.AuthorName, data.Title),
		Html: data.HTML,
	}, nil
}
