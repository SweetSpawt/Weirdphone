package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type output struct {
	Word         string `json:"word"`
	Score        int    `json:"score"`
	NumSyllables int    `json:"numSyllables"`
}

type searchResults struct {
	Total     int    `json:"total"`
	TotalHits int    `json:"totalHits"`
	Hits      []Hits `json:"hits"`
}

//Hits is struct of searchResults
type Hits struct {
	ID              int    `json:"id"`
	PageURL         string `json:"pageURL"`
	Type            string `json:"type"`
	Tags            string `json:"tags"`
	PreviewURL      string `json:"previewURL"`
	PreviewWidth    int    `json:"preveiwWidth"`
	PreviewHeight   int    `json:"previewHeight"`
	WebformatURL    string `json:"webformatURL"`
	WebformatWidth  int    `json:"webformatWidth"`
	WebformatHeight int    `json:"webformatHeight"`
	LargeImageURL   string `json:"largeImageURL"`
	FullHDURL       string `json:"fullHDURL"`
	ImageURL        string `json:"imageURL"`
	ImageWidth      int    `json:"imageWidth"`
	ImageHeight     int    `json:"imageHeight"`
	ImageSize       int    `json:"imageSize"`
	Views           int    `json:"views"`
	Downloads       int    `json:"downloads"`
	Favorites       int    `json:"favorites"`
	Likes           int    `json:"likes"`
	Comments        int    `json:"comments"`
	UserID          int    `json:"user_id"`
	User            string `json:"user"`
	UserImageURL    string `json:"userImageURL"`
}

func main() {
	fmt.Print("Enter text: ")
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
		return
	}

	input = strings.TrimSuffix(input, "\n")
	input = input[0 : len(input)-1]
	fmt.Println(input)
	output, err := rhyme(input)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	image, err := imageSearch(output)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	Open(image)
}

func rhyme(input string) (string, error) {
	resp, err := http.Get("https://api.datamuse.com/words?rel_rhy=" + input)
	if err != nil {
		fmt.Println("Error:", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
	}

	var out []output
	json.Unmarshal(body, &out)
	if len(out) == 0 {
		return "", nil
	}

	var n int
	rand.Seed(time.Now().UnixNano())
	n = 0 + rand.Intn(len(out))

	return out[n].Word, nil
}

func imageSearch(output string) (string, error) {

	resp, err := http.Get("https://pixabay.com/api/?key=18226674-3cca4777986de04451f60d6cf&q=" + output + "&image_type=photo")
	if err != nil {
		fmt.Println("Error:", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
	}

	var im searchResults
	json.Unmarshal(body, &im)

	return im.Hits[0].LargeImageURL, nil
}

// Commands returns a list of possible commands to use to open a url.
func Commands() [][]string {
	var cmds [][]string
	if exe := os.Getenv("BROWSER"); exe != "" {
		cmds = append(cmds, []string{exe})
	}
	switch runtime.GOOS {
	case "darwin":
		cmds = append(cmds, []string{"/usr/bin/open"})
	case "windows":
		cmds = append(cmds, []string{"cmd", "/c", "start"})
	default:
		if os.Getenv("DISPLAY") != "" {
			// xdg-open is only for use in a desktop environment.
			cmds = append(cmds, []string{"xdg-open"})
		}
	}
	cmds = append(cmds,
		[]string{"chrome"},
		[]string{"google-chrome"},
		[]string{"chromium"},
		[]string{"firefox"},
	)
	return cmds
}

// Open tries to open url in a browser and reports whether it succeeded.
func Open(url string) bool {
	for _, args := range Commands() {
		cmd := exec.Command(args[0], append(args[1:], url)...)
		if cmd.Start() == nil && appearsSuccessful(cmd, 3*time.Second) {
			return true
		}
	}
	return false
}

// appearsSuccessful reports whether the command appears to have run successfully.
// If the command runs longer than the timeout, it's deemed successful.
// If the command runs within the timeout, it's deemed successful if it exited cleanly.
func appearsSuccessful(cmd *exec.Cmd, timeout time.Duration) bool {
	errc := make(chan error, 1)
	go func() {
		errc <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		return true
	case err := <-errc:
		return err == nil
	}
}
