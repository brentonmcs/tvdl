package main

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
	feedprocessor "tvdownload/feedProcessor"
	"tvdownload/mover"
	"tvdownload/renamer"
	"tvdownload/transmission"

	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("tvproccessor")

func main() {

	directoryPtr := flag.String("directory", "/mnt/library/TvShows", "the base directory for show")
	downloadDir := flag.String("downloads", "/mnt/library/downloads", "Directory of the downloads")
	torrentServer := flag.String("server", "192.168.0.230", "the bittorrent remote server")
	username := flag.String("username", "pi", "the bittorrent remote server")
	password := flag.String("password", "pi", "the bittorrent remote server")

	minutesBetween := flag.Int64("mintes", 5, "the minutes between checks")
	flag.Parse()

	transmissionHandle := transmission.NewTransmissonHanlde(torrentServer, username, password, log)
	processHandle := feedprocessor.NewFeedProcessorHandle(transmissionHandle, *directoryPtr, log)
	moveHandler := mover.NewMoveShowHandler(*directoryPtr)

	minutes := time.Duration(*minutesBetween) * time.Minute
	feedTicker := time.NewTicker(minutes)
	moveTicker := time.NewTicker(minutes)

	processHandle.ProcessFeed()
	startMoving(downloadDir, moveHandler)

	for {
		select {
		case <-moveTicker.C:
			log.Info("Starting Move Job")
			go startMoving(downloadDir, moveHandler)

		case t := <-feedTicker.C:
			log.Info("Starting Feed Processing at", t)
			err := processHandle.ProcessFeed()

			if err != nil {
				log.Critical(err.Error())
			}
		}
	}
}

func startMoving(downloadDir *string, moveHandler *mover.MoveShowHandle) {

	downloadedShows, _ := getDownloadedShows(*downloadDir)

	for _, s := range downloadedShows {
		log.Infof("Processing filename : %s", s)
		moveHandler.MoveTvShowHome(renamer.GetTvShowDetails(s))
	}
}

func getDownloadedShows(downloadedDirectory string) ([]string, error) {
	files, err := ioutil.ReadDir(downloadedDirectory)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	var result []string
	for _, f := range files {
		if strings.Contains(f.Name(), ".mkv") || strings.Contains(f.Name(), ".avi") {
			result = append(result, filepath.Join(downloadedDirectory, f.Name()))
		}
	}
	return result, nil
}
