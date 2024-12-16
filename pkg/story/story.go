package story

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	TopStoriesUrl = "https://hacker-news.firebaseio.com/v0/topstories.json"
	ItemUrl       = "https://hacker-news.firebaseio.com/v0/item/"
)

type StoryId int

type Story struct {
	Id    StoryId
	Score int
	Title string
	Url   string
	Kids  []int
}

// FetchDetailed sends size number of concurrent request to ycombinator
// we wait at most 1 second for requests if server does not respond
// the Story returned from function has default values
func FetchDetailed(size int) ([]Story, error) {
	result := make([]Story, size)
	var wg sync.WaitGroup
	wg.Add(size)

	ids, err := FetchTop500Ids()
	if err != nil {
		return result, err
	}
	for i, id := range ids[:size] {
		go func(id StoryId, i int) {
			defer wg.Done()
			story, err := FetchById(id)
			if err != nil {
				fmt.Println("error getting story!") //TODO: change to fmt.Debug
				return
			}
			result[i] = story
		}(id, i)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// Wait for either "timeout" or all Stories fetched
	select {
	case <-done:
	case <-time.After(2 * time.Second):
	}

	return result, nil
}

func FetchTop500Ids() ([]StoryId, error) {
	var stories []StoryId
	resp, err := http.Get(TopStoriesUrl)
	if err != nil {
		return nil, fmt.Errorf("could not get top stories: %s", err)
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("could not fetch top stories: %s", err)
	}

	if err := json.Unmarshal(body, &stories); err != nil {
		return nil, fmt.Errorf("top stories have bad id %s", err)
	}

	return stories, nil
}

func FetchById(id StoryId) (Story, error) {
	var s Story
	resp, err := http.Get(urlForId(id))
	if err != nil {
		return s, fmt.Errorf("faild to get story %d: %s", id, err)
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return s, fmt.Errorf("failed reading data for story %d: %s", id, err)
	}

	if err := json.Unmarshal(body, &s); err != nil {
		return s, fmt.Errorf("could not parse story %d: %s", id, err)
	}

	return s, nil
}

func urlForId(id StoryId) string {
	return ItemUrl + strconv.Itoa(int(id)) + ".json"
}
