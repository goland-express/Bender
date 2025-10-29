package youtube

import (
	"encoding/json"
	"io"
	"log"
	"os/exec"
	"time"
)

// ffmpeg slowed + reverb filter "-af", "asetrate=44100*0.75,aresample=44100,aecho=0.8:0.9:200:0.3",

func FetchStreamWithMetadata(query string) (io.ReadCloser, *Metadata, error) {
	ytQuery, err := prepareQuery(query)

	if err != nil {
		return nil, nil, err
	}

	ytCmd := exec.Command("yt-dlp",
		"--quiet",
		"--no-playlist",
		"--no-simulate",
		"--print", "%(.{id,title,webpage_url,thumbnail,duration})#j",
		"--output", "-",
		ytQuery,
	)

	ffmpegCmd := exec.Command("ffmpeg", "-y",
		"-i", "pipe:0",
		"-acodec", "libopus",
		"-probesize", "32",
		"-analyzeduration", "0",
		"-fflags", "nobuffer",
		"-b:a", "32k",
		"-f", "opus",
		"pipe:1",
	)

	ytStdout, err := ytCmd.StdoutPipe()

	if err != nil {
		log.Println("Error getting youtube stdout pipe:", err)

		return nil, nil, err
	}

	ytStderr, err := ytCmd.StderrPipe()

	if err != nil {
		log.Println("Error getting youtube stderr pipe:", err)

		return nil, nil, err
	}

	ffmpegCmd.Stdin = ytStdout

	ffmpegStdout, err := ffmpegCmd.StdoutPipe()

	if err != nil {
		log.Println("Error getting ffmpeg stdout pipe:", err)

		return nil, nil, err
	}

	err = ffmpegCmd.Start()

	if err != nil {
		log.Println("Error starting ffmpeg command:", err)

		return nil, nil, err
	}

	err = ytCmd.Start()

	var metadata Metadata
	var rawMetadata []byte

	metadataBuff := make([]byte, 1024)
	startTime := time.Now()

	for {
		n, err := ytStderr.Read(metadataBuff)

		if n > 0 {
			rawMetadata = append(rawMetadata, metadataBuff[:n]...)
		}

		if err = json.Unmarshal(rawMetadata, &metadata); err != nil && time.Since(startTime) > time.Second*3 {
			return nil, nil, err
		} else {
			break
		}
	}

	if err != nil {
		panic(err)
	}

	return ffmpegStdout, &metadata, nil
}

func FetchStream(query string) (io.ReadCloser, error) {
	ytQuery, err := prepareQuery(query)

	if err != nil {
		return nil, err
	}

	ytCmd := exec.Command("yt-dlp",
		"--quiet",
		"--no-playlist",
		"--no-simulate",
		"--output", "-",
		ytQuery,
	)

	ffmpegCmd := exec.Command("ffmpeg", "-y",
		"-i", "pipe:0",
		"-acodec", "libopus",
		"-probesize", "32",
		"-analyzeduration", "0",
		"-fflags", "nobuffer",
		"-b:a", "32k",
		"-f", "opus",
		"pipe:1",
	)

	ytStdout, err := ytCmd.StdoutPipe()

	if err != nil {
		log.Println("Error getting youtube stdout pipe: ", err)

		return nil, err
	}

	ffmpegCmd.Stdin = ytStdout

	ffmpegStdout, err := ffmpegCmd.StdoutPipe()

	if err != nil {
		log.Println("Error getting ffmpeg stdout pipe: ", err)

		return nil, err
	}

	err = ffmpegCmd.Start()

	if err != nil {
		log.Println("Error starting ffmpeg command: ", err)

		return nil, err
	}

	err = ytCmd.Start()

	if err != nil {
		panic(err)
	}

	return ffmpegStdout, nil
}
