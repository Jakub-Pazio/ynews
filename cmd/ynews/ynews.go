package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Jakub-Pazio/ynews/pkg/story"
)

var numberFlag = flag.Int("number", 10, "number of stories to display")
var filterFlag = flag.String("filter", "", "filter for article titles")

func main() {
	flag.Parse()

	stories, err := story.FetchDetailed(*numberFlag)
	if err != nil {
		panic(err)
	}

	displayStories := story.ConvertToDisplay(stories)
	displayStories = story.Filter(displayStories, *filterFlag)
	if len(displayStories) == 0 {
		fmt.Println("Could not search for stories you are looking for 😓")
		os.Exit(1)
	}
	story.StartMostPopular(displayStories)
	story.MarkMostCommented(displayStories)

	for i, ds := range displayStories {
		fmt.Printf("%d: %s\n", i+1, ds)
	}

	userChoice, err := getUserNumber()

	if err != nil {
		fmt.Println("Could not get article number 😔")
		os.Exit(1)
	}

	if userChoice > len(displayStories) || userChoice < 1 {
		fmt.Println("No article with such number")
		os.Exit(1)
	}

	openBrowserWithArticle(displayStories[userChoice-1].Url)
}

func getUserNumber() (int, error) {
	fmt.Print(">")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	input = strings.TrimSpace(input)
	userChoice, err := strconv.Atoi(input)

	if err != nil {
		return 0, err
	}
	return userChoice, nil
}

func openBrowserWithArticle(url string) {
	exec.Command("open", url).Run()
}
