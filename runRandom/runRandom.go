package runRandom

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://www.acmicpc.net"
var problemURL string = "/problem/tags"
var userURL string = "/user/"

// ExtractedProblemLinks is a struct for the list of all problems
type ExtractedProblemLinks struct {
	tagTitle      string
	problemNumber int
	title         string
}

// RunRandom generates a string of a random BaekJoon problem
func RunRandom(problems []ExtractedProblemLinks, userSolved []int) (string, string) {
	startTime := time.Now()

	var newRand int

	for true {
		newRand = GetRandomIndex(len(problems))
		if problems[newRand].problemNumber == 0 {
			continue
		}
		if !Contains(userSolved, problems[newRand].problemNumber) {
			break
		}
	}
	fmt.Println(baseURL + "/problem/" + strconv.Itoa(problems[newRand].problemNumber))

	elapsedTime := time.Now().Sub(startTime)
	fmt.Println("Time elapsed: ", elapsedTime)

	return problems[newRand].title, baseURL + "/problem/" + strconv.Itoa(problems[newRand].problemNumber)
}

// GetPage retrieves the full list of all problems
func GetPage() []ExtractedProblemLinks {
	var resultProblems []ExtractedProblemLinks
	var targetURL string = baseURL + problemURL
	c := make(chan []ExtractedProblemLinks)
	numberOfTags := 0

	res, err := http.Get(targetURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find("tr")

	searchCards.Each(func(i int, card *goquery.Selection) {
		var address string
		var tagTitle string
		var numberOfProblems int
		a := card.Find("a")
		a.Each(func(ia int, acard *goquery.Selection) {
			if ia == 0 {
				href, ok := a.Attr("href")
				if ok {
					address = baseURL + href
				}
			}
		})

		tdList := card.Find("td")
		tdList.Each(func(j int, td *goquery.Selection) {
			if j == 0 {
				tagTitle = td.Text()
			} else if j == 2 {
				numberOfProblems = atoi(td.Text())
			}
		})
		if numberOfProblems < 5 {
			return
		}
		numberOfTags++
		go extraactSeperateProblems(address, tagTitle, c)
	})

	for i := 0; i < numberOfTags; i++ {
		problems := <-c
		resultProblems = append(resultProblems, problems...)
	}

	return resultProblems
}

func extraactSeperateProblems(url, title string, c chan<- []ExtractedProblemLinks) {
	var resultList []ExtractedProblemLinks

	res, err := http.Get(url)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	trList := doc.Find("tr")

	trList.Each(func(i int, tr *goquery.Selection) {
		var problem ExtractedProblemLinks

		problem.tagTitle = title
		problem.problemNumber = atoi(tr.Find(".list_problem_id").Text())

		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			if j != 1 {
				return
			}
			problem.title = td.Find("a").Text()
		})

		resultList = append(resultList, problem)
	})

	c <- resultList
}

// GetUserSolvedProblemInfo returns an array of all solved problems by the user
func GetUserSolvedProblemInfo(userID string) []int {
	targetURL := baseURL + userURL + userID
	var resultList []int

	res, err := http.Get(targetURL)
	checkErr(err)
	if res.StatusCode != 200 {
		return []int{-1}
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".panel-body").First().Find("a").Each(func(i int, a *goquery.Selection) {
		resultList = append(resultList, atoi(a.Text()))
	})

	return resultList
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
}

// CleanString cleans the given string
func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func atoi(str string) int {
	if len(str) == 0 {
		return 0
	}
	n, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return n
}

// Contains checks if the int array contains a given number; true if it does, false otherwise
func Contains(arr []int, number int) bool {
	for _, a := range arr {
		if a == number {
			return true
		}
	}
	return false
}

// GetRandomIndex returns a random integer from a given range of numbers
func GetRandomIndex(number int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(number)))
	checkErr(err)
	return int(n.Int64())
}

// PrintUserID logs the user ID requested
func PrintUserID(userID string) {
	fmt.Println("User Requesting : ", userID)
}
