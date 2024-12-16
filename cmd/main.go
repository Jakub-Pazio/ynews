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
	fmt.Println(*filterFlag)
	story.StartMostPopular(displayStories)
	story.MarkMostCommented(displayStories)

	for i, ds := range displayStories {
		fmt.Printf("%d: %s\n", i+1, ds)
	}

	userChoise, err := getUserNumber()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if userChoise > len(displayStories) || userChoise < 1 {
		fmt.Println("No article with such number")
		os.Exit(1)
	}

	openBrowserWithArticle(displayStories[userChoise-1].Url)
}

func getUserNumber() (int, error) {
	fmt.Print(">")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	input = strings.TrimSpace(input)
	userChoise, err := strconv.Atoi(input)

	if err != nil {
		return 0, err
	}
	return userChoise, nil
}

func openBrowserWithArticle(url string) {
	exec.Command("open", url).Run()
}
