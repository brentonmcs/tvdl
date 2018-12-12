package feedprocessor

import (
	"strings"
	"tvdownload/mover"
	"tvdownload/renamer"
	"tvdownload/transmission"

	"github.com/mmcdole/gofeed"
	logging "github.com/op/go-logging"
	"github.com/renstrom/fuzzysearch/fuzzy"
)

//ProcessorHandle Creates a new handle for the processor
type ProcessorHandle struct {
	feedParser  *gofeed.Parser
	transmisson *transmission.TransmissonHanlde
	shows       []string
	log         *logging.Logger
}

//NewFeedProcessorHandle Creates a new handle for the processor
func NewFeedProcessorHandle(transmissionHandle *transmission.TransmissonHanlde, tvshowDirectory string, log *logging.Logger) *ProcessorHandle {
	fp := gofeed.NewParser()
	shows := mover.FindDirectoryNames(tvshowDirectory)
	return &ProcessorHandle{feedParser: fp, transmisson: transmissionHandle, shows: shows, log: log}
}

//ProcessFeed Check the Feed list for any matching shows
func (fp *ProcessorHandle) ProcessFeed() error {

	feed, err := fp.feedParser.ParseURL("https://eztv1.unblocked.ms/ezrss.xml")

	if err != nil {
		return err
	}

	activeTorrents := fp.transmisson.GetActiveTorrentNames()

	fp.log.Infof("Feed contains %d items", len(feed.Items))
	for _, item := range feed.Items {
		fileDetails := renamer.GetTvShowDetails(item.Title)

		if fp.validTvTorrent(fileDetails, activeTorrents, item.Title) {
			fp.transmisson.AddTorrent(item.Extensions["torrent"]["magnetURI"][0].Value)
		}
	}

	return nil
}

func (fp *ProcessorHandle) validTvTorrent(fileDetails *renamer.TvShowDetails, activeTorrents []string, fileName string) bool {
	return matchShowsNameMatchesDirectory(fileDetails.Name, fp.shows) &&
		!checkAlreadyDownloading(activeTorrents, fileDetails) &&
		matchFileType(fileName) &&
		checkAlreadyDownloaded(fileDetails)
}

func checkAlreadyDownloading(activeTorrents []string, fileDetails *renamer.TvShowDetails) bool {
	for _, torrent := range activeTorrents {
		torrentShowDetails := renamer.GetTvShowDetails(torrent)

		if torrentShowDetails.ComputedName == fileDetails.ComputedName {
			return true
		}
	}
	return false
}

func matchFileType(fileName string) bool {
	return strings.Contains(fileName, "265") && strings.Contains(fileName, "720")
}

func matchShowsNameMatchesDirectory(showName string, shows []string) bool {
	return len(fuzzy.RankFindFold(showName, shows)) > 0
}

func checkAlreadyDownloaded(fileDetails *renamer.TvShowDetails) bool {
	return true
}
