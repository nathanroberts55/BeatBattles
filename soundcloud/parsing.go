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
	Width        int     `json:"width"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	ThumbnailURL string  `json:"thumbnail_url"`
	HTML         string  `json:"html"`
	AuthorName   string  `json:"author_name"`
	AuthorURL    string  `json:"author_url"`
}

func GetEmbed(link string) (string, error) {
	client := &http.Client{}
	scUrl := url.QueryEscape(link)
	req, err := http.NewRequest("GET", fmt.Sprintf("https://soundcloud.com/oembed?format=json&url=%s&maxheight=256px&maxwidth=256px", scUrl), nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}
	var data scData
	err = json.Unmarshal(body, &data)

	if err != nil {
		return "", err
	}

	return data.HTML, nil
}
