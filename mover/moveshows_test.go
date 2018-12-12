package mover

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"tvdownload/renamer"

	uuid "github.com/satori/go.uuid"
)

type testFindNames struct {
	input    string
	expected string
}

func TestFindShowDirectory(t *testing.T) {

	shows := []string{"Bones", "Suits", "Game of Thrones"}

	handler := MoveShowHandle{homeTvDirectory: "", shows: shows}

	testCase := []testFindNames{
		testFindNames{input: "Bone", expected: "Bones"},
		testFindNames{input: "Bones", expected: "Bones"},
		testFindNames{input: "bone", expected: "Bones"},
		testFindNames{input: "bean", expected: "Bean"},
		testFindNames{input: "suit", expected: "Suits"},
		testFindNames{input: "gameofthrones", expected: "Game of Thrones"},
		testFindNames{input: "GOT", expected: "Game of Thrones"},
		testFindNames{input: "random", expected: "Random"},
	}

	for _, tC := range testCase {
		result := handler.findShowDirectory(tC.input)

		if result != tC.expected {
			t.Fatalf("Find Directory is not correct - expected %v, result %v, input %v", tC.expected, result, tC.input)
		}
	}
}

func TestSearchDirectory(t *testing.T) {

	uuid, _ := uuid.NewV4()
	dirName := "./" + uuid.String()
	defer cleanupFolder(dirName)
	os.Mkdir(dirName, 07777)

	os.Mkdir(dirName+"/bones", 0777)
	os.Mkdir(dirName+"/suits", 0777)
	os.Mkdir(dirName+"/game of thrones", 0777)

	result := NewMoveShowHandler(dirName)

	if len(result.shows) != 3 {
		t.Fatalf("Should be only 3 shows not %d", len(result.shows))
	}

	findResult := result.findShowDirectory("Bones")

	if findResult != "bones" {
		t.Fail()
	}
}

func TestMovingFile(t *testing.T) {
	uuid, _ := uuid.NewV4()
	dirName := "./" + uuid.String()
	os.Mkdir(dirName, 07777)
	defer cleanupFolder(dirName)

	handler := NewMoveShowHandler(dirName)

	_, filename, _, _ := runtime.Caller(0)

	path := filepath.Join(filepath.Dir(filename), "test.txt")
	_, err := os.Create(path)
	if err != nil {
		t.Fatal("File was not created")
	}

	tvShowDetails := &renamer.TvShowDetails{ComputedName: "Test S01E01.txt", Name: "Test", Path: path, Season: 1}
	handler.MoveTvShowHome(tvShowDetails)

	_, err = os.Stat(filepath.Join(dirName, "Test/Season 1/Test S01E01.txt"))

	if err != nil {
		t.Fatal("File has not been moved")
	}
}
func cleanupFolder(dirName string) {
	os.RemoveAll(dirName)
}
