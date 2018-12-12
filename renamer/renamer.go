package renamer

import (
	"fmt"
	"path/filepath"
	"strings"
)

//TvShowDetails holds the information about the show
type TvShowDetails struct {
	Name         string
	Path         string
	Season       int
	episode      int
	extension    string
	ComputedName string
}

type fileDetails struct {
	filename  string
	path      string
	extension string
}

type parsedFileDetails struct {
	season  int
	episode int
	name    string
}

func (p *parsedFileDetails) validSeasonFound() bool {
	return p.season == 0 && p.episode == 0
}

//NewTvShowDetails constuctor
func NewTvShowDetails(parsedFileDetails parsedFileDetails, extension, path string) *TvShowDetails {
	computedName := fmt.Sprintf("%v S%02dE%02d%v", strings.Title(parsedFileDetails.name), parsedFileDetails.season, parsedFileDetails.episode, extension)
	return &TvShowDetails{Name: parsedFileDetails.name,
		Season: parsedFileDetails.season, episode: parsedFileDetails.episode, extension: extension, ComputedName: computedName, Path: path}
}

//GetTvShowDetails renames the file so it's a clean TV Show Name
func GetTvShowDetails(path string) *TvShowDetails {
	return findDetails(getFileDetails(path))
}

func getFileDetails(path string) fileDetails {

	fName := filepath.Base(path)
	extName := filepath.Ext(path)
	bName := fName[:len(fName)-len(extName)]

	return fileDetails{filename: bName, path: path, extension: extName}
}

func findName(filename string) string {

	if strings.Contains(filename, ".") {
		return strings.Split(filename, ".")[0]
	}
	return filename
}

func findDetails(file fileDetails) *TvShowDetails {

	parsedFileDetails := extractViaRegex("S\\d+E\\d+", extraDigitsSE, file.filename)

	if parsedFileDetails.validSeasonFound() {
		parsedFileDetails = extractViaRegex("\\d+X\\d+", extraDigitsSE, file.filename)
	}

	if parsedFileDetails.validSeasonFound() {
		parsedFileDetails = extractViaRegex("\\d\\d\\d\\d", extraDigits4d, file.filename)
	}

	if parsedFileDetails.validSeasonFound() {
		parsedFileDetails = extractViaRegex("\\d\\d\\d", extractThreeDigit, file.filename)
	}
	return NewTvShowDetails(parsedFileDetails, file.extension, file.path)
}

func removeSymbols(filename string) string {
	result := strings.Trim(strings.Replace(filename, ".", " ", -1), " ")
	return strings.Trim(strings.Replace(result, "_", " ", -1), " ")
}
