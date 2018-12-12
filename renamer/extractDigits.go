package renamer

import (
	"regexp"
	"strconv"
	"strings"
)

func extractThreeDigit(filename string) (int, int) {
	season, _ := strconv.Atoi(filename[0:1])
	eps, _ := strconv.Atoi(filename[1:])
	return season, eps
}

type extractdigits func(string) (int, int)

func extraDigits4d(fileName string) (int, int) {
	return extractDigitsRegex(fileName, "\\d\\d")
}
func extraDigitsSE(fileName string) (int, int) {
	return extractDigitsRegex(fileName, "\\d+")
}

func extractDigitsRegex(fileName string, regexString string) (int, int) {
	re2 := regexp.MustCompile(regexString)
	seasonNumbers := re2.FindAllString(fileName, -1)
	season, _ := strconv.Atoi(seasonNumbers[0])
	episode, _ := strconv.Atoi(seasonNumbers[1])
	return season, episode
}

func extractViaRegex(regex string, extractD extractdigits, filename string) parsedFileDetails {

	re := regexp.MustCompile(regex)
	seasonResult := re.FindAllString(strings.ToUpper(filename), -1)

	if len(seasonResult) == 0 {
		return parsedFileDetails{name: findName(filename)}
	}

	index := re.FindStringIndex(strings.ToUpper(filename))[0]
	season, episode := extractD(seasonResult[0])
	return parsedFileDetails{name: removeSymbols(filename[0:index]), season: season, episode: episode}
}
