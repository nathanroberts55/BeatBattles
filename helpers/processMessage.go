package helpers

import "regexp"

func ExtractSong(message string) string {
	re := regexp.MustCompile(`https://soundcloud\.com/[^/]+/[^/?]+`)
	match := re.FindString(message)
	return match
}
