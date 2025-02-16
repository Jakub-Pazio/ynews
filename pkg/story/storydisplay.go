package story

import (
	"fmt"
	"strings"
)

type DisplayStory struct {
	Story
	Icons []string
}

func (ds DisplayStory) String() string {
	return fmt.Sprintf("%s %s", ds.Title, strings.Join(ds.Icons, ", "))
}

func New(s Story) DisplayStory {
	return DisplayStory{Story: s}
}

func ConvertToDisplay(ss []Story) []DisplayStory {
	result := make([]DisplayStory, len(ss))

	for i, s := range ss {
		result[i] = DisplayStory{Story: s}
	}

	return result
}

func Filter(dss []DisplayStory, filter string) []DisplayStory {
	// when no filter is provided return whole input
	if filter == "" {
		return dss
	}
	var result []DisplayStory

	for _, ds := range dss {
		if strings.Contains(strings.ToLower(ds.Title), strings.ToLower(filter)) {
			result = append(result, ds)
		}
	}
	return result
}

func StartMostPopular(dss []DisplayStory) {
	if len(dss) == 0 {
		return
	}

	maxScore, maxId := dss[0].Score, 0

	for i, ds := range dss {
		if ds.Score > maxScore {
			maxScore, maxId = ds.Score, i
		}
	}

	dss[maxId].Icons = append(dss[maxId].Icons, "🔥")
}
func MarkMostCommented(dss []DisplayStory) {
	if len(dss) == 0 {
		return
	}

	maxScore, maxId := len(dss[0].Kids), 0

	for i, ds := range dss {
		if len(ds.Kids) > maxScore {
			maxScore, maxId = len(ds.Kids), i
		}
	}

	dss[maxId].Icons = append(dss[maxId].Icons, "💬")
}
