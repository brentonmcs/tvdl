package transmission

import (
	"fmt"
	"os"

	logging "github.com/op/go-logging"

	"github.com/hekmon/transmissionrpc"
)

//TransmissonHanlde holds the client details for Transmission
type TransmissonHanlde struct {
	client *transmissionrpc.Client
	log    *logging.Logger
}

//NewTransmissonHanlde create a new handle
func NewTransmissonHanlde(torrentServer *string, username *string, password *string, log *logging.Logger) *TransmissonHanlde {

	transmissionbt, err := transmissionrpc.New(*torrentServer, *username, *password, nil)
	if err != nil {
		panic(err)
	}
	return &TransmissonHanlde{client: transmissionbt, log: log}
}

//GetActiveTorrentNames returns a list of Torrent Names
func (h *TransmissonHanlde) GetActiveTorrentNames() []string {

	list, err := h.client.TorrentGet([]string{"name"}, nil)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}

	h.log.Infof("Found %d Active torrents", len(list))

	var names []string
	for _, l := range list {
		names = append(names, *l.Name)
	}

	return names
}

//AddTorrent to Transmisson Client
func (h *TransmissonHanlde) AddTorrent(magnet string) {
	torrent, err := h.client.TorrentAdd(&transmissionrpc.TorrentAddPayload{
		Filename: &magnet,
	})
	if err != nil {
		h.log.Error("Error Adding Torrent %s", err.Error())
	} else {
		h.log.Info("Added Torrent %s", *torrent.Name)
	}
}
