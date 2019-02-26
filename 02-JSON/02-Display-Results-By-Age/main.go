package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const IssuesURL = "https://api.github.com/search/issues"

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items []*Issue
	CreatedAt time.Time
}

type Issue struct {
	Number int
	HTMLURL string `json:"html_url"`
	Title string
	State string
	User *User
	CreatedAt time.Time `json:"created_at"`
	Body string
}

type User struct {
	Login string
	HTMLURL string `json:"html_url"`
}

func main() {
	result, err := searchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)

	//figure out a way to not to do this without the use of 3 separate loops as this will not be the most efficient way.

	for _, item := range result.Items {
		then := time.Date(item.CreatedAt.Year(), item.CreatedAt.Month(), item.CreatedAt.Day(), item.CreatedAt.Hour(),
			item.CreatedAt.Minute(), item.CreatedAt.Second(), item.CreatedAt.Nanosecond(), time.UTC)
		y1, m1, d1, _, _, _ := diff(then, time.Now())

		if y1 > 0 {
			fmt.Printf("#%-5d %9.9s %.55s - Created: Years: %v Months: %v, Days: %v agp\n",
				item.Number, item.User.Login, item.Title, y1, m1, d1)
		}
	}
	fmt.Printf("\n")

	for _, item := range result.Items {
		then := time.Date(item.CreatedAt.Year(), item.CreatedAt.Month(), item.CreatedAt.Day(), item.CreatedAt.Hour(),
			item.CreatedAt.Minute(), item.CreatedAt.Second(), item.CreatedAt.Nanosecond(), time.UTC)
		y1, m1, d1, _, _, _ := diff(then, time.Now())

		if m1 < 1 {
			fmt.Printf("#%-5d %9.9s %.55s - Created: Years: %v Months: %v, Days: %v agp\n",
				item.Number, item.User.Login, item.Title, y1, m1, d1)
		}
	}
	fmt.Printf("\n")

	for _, item := range result.Items {
		then := time.Date(item.CreatedAt.Year(), item.CreatedAt.Month(), item.CreatedAt.Day(), item.CreatedAt.Hour(),
			item.CreatedAt.Minute(), item.CreatedAt.Second(), item.CreatedAt.Nanosecond(), time.UTC)
		y1, m1, d1, _, _, _ := diff(then, time.Now())

		if y1 < 1 {
			fmt.Printf("#%-5d %9.9s %.55s - Created: Years: %v Months: %v, Days: %v agp\n",
				item.Number, item.User.Login, item.Title, y1, m1, d1)
		}
	}
}

func searchIssues (terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("Search Query Failed: %s", resp.Status)
	}

	var result IssuesSearchResult
	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, err
}

func diff (a, b time.Time) (year, month, day, hour, min, sec int) {
	//time difference only works if the time zones are the same. This will make sure if they are different, they get converted to the correct time zone
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}

	if a.After(b) {
		a, b = b, a //this will make sure that the original date is after the time now
	}

	y1, M1, d1 := a.Date() //get the year, month and date of the first date entered
	y2, M2, d2 := b.Date() //get the year, month and date of the second date entered

	h1, m1, s1 := a.Clock() //get the hour, min, second of the first time
	h2, m2, s2 := b.Clock() //get the hour, min, second of the second time

	year = int(y2 - y1) //this will get the difference in years
	month = int(M2 - M1) //this will get the difference in months
	day = int(d2 - d1) //this will get the difference in days.
	hour = int(h2 - h1) //this will get the difference in hours
	min = int(m2 - m1) //this will get the difference in minutes
	sec = int(s2 - s1) //this will get the difference in seconds

	//normalize negative values (can never have -10 days, instead)
	if sec < 0 { //if sec is lower than 0 (minus value), add 60 and -1
		sec += 60
		sec--
	}
	if min < 0 {
		min += 60
		min--
	}
	if hour < 0 {
		hour += 60
		hour--
	}
	if day < 0 {
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC) //this gets the amount of days in that month
		day += 32 - t.Day() //this will add 32 to that date and subtract t.Day
		month-- //we've now gone into the next month, therefore subtract 1
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}