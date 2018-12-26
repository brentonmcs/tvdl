package mover

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"tvdownload/renamer"

	"github.com/renstrom/fuzzysearch/fuzzy"
)

//MoveShowHandle stores the config and details for moving files
type MoveShowHandle struct {
	homeTvDirectory string
	shows           []string
}

//NewMoveShowHandler creates a new Handle to manage moving shows
func NewMoveShowHandler(homeTvDirectory string) *MoveShowHandle {
	shows := FindDirectoryNames(homeTvDirectory)
	return &MoveShowHandle{homeTvDirectory: homeTvDirectory, shows: shows}
}

//MoveTvShowHome - takes the renamed file and moves it home
func (m *MoveShowHandle) MoveTvShowHome(t *renamer.TvShowDetails) {

	log.Printf("Extension is %s", t.Extension)
	if strings.Contains(t.Extension, "part") {
		log.Printf("Not moving as file is not the full version")
		return
	}
	if t.Season == 0 {
		log.Printf("Not moving file as Season details were not found")
		return
	}

	if t.Name == "" {
		log.Printf("Not moving file as Show %s folder not found", t.ComputedName)
		return
	}

	showDirectory := m.findShowDirectory(t.Name)
	seasonDirectory := m.findSeasonDirectory(showDirectory, t.Season)
	newShowPath := filepath.Join(seasonDirectory, t.ComputedName)

	log.Printf("Moving %s to the new path %s", t.Path, newShowPath)

	err := os.Rename(t.Path, newShowPath)

	if err != nil {
		log.Fatal(err)
	}
}

func (m *MoveShowHandle) createDirectory(name string) string {
	os.Mkdir(m.homeTvDirectory+"/"+name, 07777)
	return name
}
func (m *MoveShowHandle) findShowDirectory(showName string) string {

	matches := fuzzy.RankFindFold(showName, m.shows)
	if len(matches) == 0 {
		showName = strings.Title(showName)
		m.createDirectory(showName)
		return showName
	}
	sort.Sort(matches)
	return matches[0].Target
}

func (m *MoveShowHandle) findSeasonDirectory(showDirectory string, season int) string {
	directoryNames := FindDirectoryNames(m.homeTvDirectory + "/" + showDirectory)

	result := fuzzy.FindFold("Season "+strconv.Itoa(season), directoryNames)

	showPath := showDirectory + "/Season " + strconv.Itoa(season)
	if len(result) == 0 {
		m.createDirectory(showPath)
	}

	return m.homeTvDirectory + "/" + showPath
}

//FindDirectoryNames list of all the TV Show Names
func FindDirectoryNames(baseTvDirectory string) []string {
	files, err := ioutil.ReadDir(baseTvDirectory)
	if err != nil {
		log.Fatal(err)
	}
	var result []string
	for _, f := range files {
		if f.IsDir() {
			result = append(result, f.Name())
		}
	}
	return result
}
