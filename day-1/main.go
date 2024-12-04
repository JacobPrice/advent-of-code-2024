package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	cookie, url := getConfig()

	numbers, err := fetchAndParseInput(cookie, url)
	if err != nil {
		log.Fatalln(err)
	}

	diff := computeDifference(numbers)
	println(diff)
}

func getConfig() (string, string) {
	godotenv.Load()

	cookie := os.Getenv("AOC_COOKIE")
	if cookie == "" {
		log.Fatalln("AOC_COOKIE environment variable is required")
	}

	url := os.Getenv("AOC_URL")
	if url == "" {
		log.Fatalln("AOC_URL environment variable is required")
	}

	return cookie, url
}

func fetchAndParseInput(cookie string, url string) ([][2]int, error) {
	body, err := fetchInput(cookie, url)
	if err != nil {
		return nil, err
	}

	return parseNumbers(body)
}

func fetchInput(cookie string, url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: cookie,
	})

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func parseNumbers(input []byte) ([][2]int, error) {
	lines := strings.Split(string(input), "\n")
	var numbers [][2]int

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}

		n1, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		n2, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}

		numbers = append(numbers, [2]int{n1, n2})
	}

	return numbers, nil
}

func computeDifference(numbers [][2]int) int {
	firstColumn := make([]int, len(numbers))
	secondColumn := make([]int, len(numbers))

	for i, pair := range numbers {
		firstColumn[i] = pair[0]
		secondColumn[i] = pair[1]
	}

	sort.Ints(firstColumn)
	sort.Ints(secondColumn)

	sum := 0
	for i := 0; i < len(numbers); i++ {
		diff := secondColumn[i] - firstColumn[i]
		if diff < 0 {
			diff = -diff
		}
		sum += diff
	}

	return sum
}
