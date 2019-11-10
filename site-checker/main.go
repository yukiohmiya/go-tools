package main

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"os"
	"strconv"
	"strings"
	// "time"
)

func main() {
	var urls []string
	var urlsDone []string
	urlList := make(map[int]map[string]string)
	var baseURL string

	//
	fmt.Println("Enter a url:")
	_, err := fmt.Scanln(&baseURL)
	if err != nil {
		fmt.Println("Error:", err)
	}

	rs := []rune(baseURL)
	if string(rs[len(rs)-1]) != "/" {
		baseURL += "/"
	}

	fmt.Println("baseURL:", baseURL)
	urls = append(urls, baseURL)
	urlsDone = append(urlsDone, baseURL)

	for i := 0; ; i++ {
		urlList[i] = make(map[string]string)

		// fmt.Println(urls)
		if len(urls) == 0 {
			break
		}
		u := urls[0]
		urls = urls[1:]

		// urlList = append(urlList, u)

		res, err := http.Get(u)
		if err != nil {
			fmt.Printf("Failed")
		}
		defer res.Body.Close()

		doc, err := goquery.NewDocument(u)
		if err != nil {
			fmt.Printf("Failed")
		}

		title := doc.Find("title").Text()
		fmt.Println(u, title, res.StatusCode)
		urlList[i]["url"] = u
		urlList[i]["title"] = title
		urlList[i]["statusCode"] = strconv.Itoa(res.StatusCode)

		doc.Find("a").Each(func(_ int, s *goquery.Selection) {
			url, _ := s.Attr("href")
			if strings.Index(url, baseURL) == -1 {
				return
			}
			if arrayContains(urlsDone, url) {
				return
			}

			urls = append(urls, url)
			urlsDone = append(urlsDone, url)
		})
	}
	fmt.Println(urlList)
	file, err := os.OpenFile("./data.csv", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer file.Close()
	err = file.Truncate(0)
	if err != nil {
		fmt.Println("Error:", err)
	}
	writer := csv.NewWriter(file)
	err = writer.Write([]string{"ID", "TITLE", "URL", "STATUS CODE"})
	if err != nil {
		fmt.Println("Error:", err)
	}
	for i := 0; i < len(urlList); i++ {
		err = writer.Write([]string{strconv.Itoa(i),
			urlList[i]["title"],
			urlList[i]["url"],
			urlList[i]["statusCode"]})
		if err != nil {
			fmt.Println("Error:", err)
		}

	}
	writer.Flush()
}

func arrayContains(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}
