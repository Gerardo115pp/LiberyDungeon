package workflows

import (
	"fmt"
	dungeon_models "libery-dungeon-libs/models"
	"path"
	"path/filepath"
	"strings"

	"github.com/Gerardo115pp/patriots_lib/echo"
	"github.com/xrash/smetrics"
)

func MatchWord(word string, posibilites []dungeon_models.Category, similarity_threshold float64) []dungeon_models.Category {
	var starts_with_set []dungeon_models.Category = make([]dungeon_models.Category, 0)
	var contains_set []dungeon_models.Category = make([]dungeon_models.Category, 0)
	var similarity_set []dungeon_models.Category = make([]dungeon_models.Category, 0)

	word_lower := strings.ToLower(word)

	var matches []dungeon_models.Category = make([]dungeon_models.Category, 0)

	// Populate sets
	for _, possibility := range posibilites {

		possibility_lower := strings.ToLower(possibility.Name)

		if possibility_lower == word_lower {
			echo.EchoDebug(fmt.Sprintf("%s is equal to %s", possibility_lower, word_lower))
			matches = append(matches, possibility)
		} else if strings.HasPrefix(possibility_lower, word_lower) && possibility_lower != word_lower {
			echo.EchoDebug(fmt.Sprintf("%s starts with %s", possibility_lower, word_lower))
			starts_with_set = append(starts_with_set, possibility)
		} else if strings.Contains(possibility_lower, word_lower) {
			echo.EchoDebug(fmt.Sprintf("%s contains %s", possibility_lower, word_lower))
			contains_set = append(contains_set, possibility)
		} else {
			similarity := smetrics.JaroWinkler(word_lower, possibility_lower, 0.7, 4)
			if similarity >= similarity_threshold {
				echo.EchoDebug(fmt.Sprintf("Similarity between %s and %s: %f", word_lower, possibility_lower, similarity))
				similarity_set = append(similarity_set, possibility)
			}
		}
	}

	matches = append(matches, starts_with_set...)
	matches = append(matches, contains_set...)
	matches = append(matches, similarity_set...)

	return matches
}

func GetUniqueMatches(word string, posibilites []dungeon_models.Category, similarity_threshold float64) []*dungeon_models.Category {
	var unique_matches []*dungeon_models.Category = make([]*dungeon_models.Category, 0)
	var unique_match *dungeon_models.Category

	matches := MatchWord(word, posibilites, similarity_threshold)

	var appearances map[string]*dungeon_models.Category = make(map[string]*dungeon_models.Category)

	for _, match := range matches {
		lower_name := strings.ToLower(match.Name)

		if other_match, exists := appearances[lower_name]; exists {
			echo.EchoDebug(fmt.Sprintf("Match: %s", match.Name))

			echo.EchoDebug(fmt.Sprintf("Other match: %s", other_match.Fullpath))

			other_fullpath := other_match.Fullpath
			other_fullpath = strings.Trim(other_fullpath, "/")

			other_match.Name = fmt.Sprintf("%s(%s)", match.Name, filepath.Base(path.Dir(other_fullpath)))

			echo.EchoDebug(fmt.Sprintf("match: %s", match.Fullpath))
			unique_match = new(dungeon_models.Category)
			unique_match.CopyContent(match) // Modifying the new match

			match_path := match.Fullpath
			match_path = strings.Trim(match_path, "/")

			unique_match.Name = fmt.Sprintf("%s(%s)", match.Name, filepath.Base(path.Dir(match_path)))

			unique_matches = append(unique_matches, unique_match)
		} else {
			unique_match = new(dungeon_models.Category)
			unique_match.CopyContent(match)

			unique_matches = append(unique_matches, unique_match)

			appearances[lower_name] = unique_match
		}
	}

	return unique_matches
}
