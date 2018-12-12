package renamer

import "testing"

type TestName struct {
	path     string
	expected string
}

type TestDetails struct {
	path    string
	season  int
	episode int
}

func TestRenameName(t *testing.T) {
	testCase := []TestName{
		TestName{path: "/User/Home/test.local", expected: "test"},
		TestName{path: "/User/Home/test.S01E01.HDTV.local", expected: "test"},
		TestName{path: "/User/Home/the.big.bang.theory.s9e13.hdtv-lol.mp4", expected: "the big bang theory"},
	}

	for _, c := range testCase {
		result := GetTvShowDetails(c.path)

		t.Logf("Rename Info - expected %v, result %v", c.expected, result.Name)
		if c.expected != result.Name {
			t.Fatalf("Rename is not correct - expected %v, result %v", c.expected, result.Name)
		}
	}

}

func TestRenameSeasonDetails(t *testing.T) {
	testCase := []TestDetails{
		TestDetails{path: "/User/Home/test.S01E01.HDTV.local", season: 1, episode: 1},
		TestDetails{path: "/User/Home/test.s01e01.HDTV.local", season: 1, episode: 1},
		TestDetails{path: "/User/Home/test.02X11.HDTV.local", season: 2, episode: 11},
		TestDetails{path: "/User/Home/test.02x11.HDTV.local", season: 2, episode: 11},
		TestDetails{path: "/User/Home/test.0212.HDTV.local", season: 2, episode: 12},
		TestDetails{path: "/User/Home/test.212.HDTV.local", season: 2, episode: 12},
		TestDetails{path: "/User/Home/the.big.bang.theory.s9e13.hdtv-lol.mp4", season: 9, episode: 13},
		TestDetails{path: "/User/Home/The BlackList S04E12.mp4", season: 4, episode: 12},
	}

	for _, c := range testCase {
		result := GetTvShowDetails(c.path)

		if c.season != 0 && c.season != result.Season {
			t.Fatalf("Rename Season is not correct - expected %v, result %v", c.season, result.Season)
		}

		if c.episode != 0 && c.episode != result.episode {
			t.Fatalf("Rename episode is not correct - expected %v, result %v", c.season, result.episode)
		}
	}
}

func TestEndToEnd(t *testing.T) {
	result := GetTvShowDetails("/user/home/test.S01E01.local")

	if result.ComputedName != "Test S01E01.local" {
		t.Fatal("Computed Name is not matcheds")
	}

	result = GetTvShowDetails("/user/home/test.show.S01E01.local")

	if result.ComputedName != "Test Show S01E01.local" {
		t.Fatal("Computed Name is not matcheds")
	}

	result = GetTvShowDetails("/user/home/test_show.S01E01.local")

	if result.ComputedName != "Test Show S01E01.local" {
		t.Fatal("Computed Name is not matcheds")
	}
}

func TestComputedName(t *testing.T) {
	result := NewTvShowDetails(parsedFileDetails{episode: 1, season: 1, name: "lower"}, ".mp3", "")
	expected := "Lower S01E01.mp3"

	t.Logf("%v", result.ComputedName)
	if result.ComputedName != expected {
		t.Fatalf("Computed Name is not correct - expected %v, result %v", result.ComputedName, expected)
	}
	result = NewTvShowDetails(parsedFileDetails{episode: 1, season: 1, name: "the lower name"}, ".mp3", "")
	expected = "The Lower Name S01E01.mp3"

	t.Logf("%v", result.ComputedName)
	if result.ComputedName != expected {
		t.Fatalf("Computed Name is not correct - expected %v, result %v", result.ComputedName, expected)
	}
}
