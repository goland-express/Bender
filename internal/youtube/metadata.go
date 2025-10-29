package youtube

import (
	"encoding/json"
	"io"
	"log"
	"os/exec"
)

type Thumbnail struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Metadata struct {
	Id         string      `json:"id"`
	Title      string      `json:"title"`
	Url        string      `json:"webpage_url"`
	Duration   float64     `json:"duration"`
	Thumbnail  string      `json:"thumbnail"`
	Thumbnails []Thumbnail `json:"thumbnails"`
}

func FetchMetadata(query string) (*Metadata, error) {
	ytQuery, err := prepareQuery(query)

	if err != nil {
		return nil, err
	}

	cmd := exec.Command("yt-dlp",
		"--quiet",
		"--skip-download",
		"--print", "%(.{id,title,webpage_url,thumbnail,thumbnails,duration})#j",
		"--flat-playlist",
		"--no-playlist",
		ytQuery,
	)

	ytStdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Println("Error getting youtube stdout pipe: ", err)
		return nil, err
	}

	err = cmd.Start()

	var rawMetadata []byte
	var metadata *Metadata

	rawMetadata, err = io.ReadAll(ytStdout)

	if err != nil {
		log.Println("Error on reading all youtube metadata: ", rawMetadata)
		return nil, err
	}

	err = json.Unmarshal(rawMetadata, &metadata)

	if err != nil {
		log.Println("Error parsing youtube metadata json: ", err)
		return nil, err
	}

	return metadata, nil
}
