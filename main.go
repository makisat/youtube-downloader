package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/kkdai/youtube/v2"
)

var youtubeRegex = regexp.MustCompile(`(?:youtu\.be/|v/|u/\w/|embed/|watch\?v=)([^#&?]*).*`)
var fileNameRegex = regexp.MustCompile(`^[a-zA-Z0-9 ]+$`)

func main() {
	videoURL := flag.String("u", "", "YouTube video URL")
	customVideoTitle := flag.String("t", "", "Custom video title")
	flag.Parse()

	fmt.Println("videoURL:", *videoURL)
	fmt.Println("customVideoTitle:", *customVideoTitle)
	if *videoURL == "" {
		flag.Usage()
		return
	}

	videoID, err := extractYoutubeID(*videoURL)
	if err != nil {
		fmt.Println("err extracting video ID:", err)
	}

	err = downloadYoutube(videoID, *customVideoTitle)
	if err != nil {
		fmt.Println("err downloading video:", err)
	}
}

func downloadYoutube(videoID string, filename string) error {
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		fmt.Println("err getting video:", err)
		return err
	}

	formats := video.Formats.WithAudioChannels()
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		fmt.Println("err getting stream:", err)
		return err
	}

	fmt.Println("Video title:", video.Title)
	if filename == "" {
		fmt.Println("filename not set, using video title")
		if fileNameRegex.MatchString(video.Title) {
			fmt.Println("filename set:", video.Title)
			filename = fmt.Sprintf("%s.mp4", video.Title)
		} else {
			filename = "video.mp4"
		}
	} else {
		fmt.Println("filename set:", filename)
		filename = fmt.Sprintf("%s.mp4", filename)
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("err creating file:", err)
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		fmt.Println("err copying file:", err)
		return err
	}
	return err
}

func extractYoutubeID(url string) (string, error) {
	match := youtubeRegex.FindStringSubmatch(url)
	if len(match) > 1 {
		return match[1], nil
	}
	return "", fmt.Errorf("no match")
}
