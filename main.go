package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type output struct {
	Word         string `json:"word"`
	Score        int    `json:"score"`
	NumSyllables int    `json:"numSyllables"`
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
	fmt.Println(output)
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
	return out[0].Word, nil
}
