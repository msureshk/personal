package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// BingAnswer is the struct for the data returned by Bing.
type BingAnswer struct {
	Type         string `json:"_type"`
	QueryContext struct {
		OriginalQuery string `json:"originalQuery"`
	} `json:"queryContext"`
	WebPages struct {
		WebSearchURL          string `json:"webSearchUrl"`
		TotalEstimatedMatches int    `json:"totalEstimatedMatches"`
		Value                 []struct {
			ID               string    `json:"id"`
			Name             string    `json:"name"`
			URL              string    `json:"url"`
			IsFamilyFriendly bool      `json:"isFamilyFriendly"`
			DisplayURL       string    `json:"displayUrl"`
			Snippet          string    `json:"snippet"`
			DateLastCrawled  time.Time `json:"dateLastCrawled"`
			SearchTags       []struct {
				Name    string `json:"name"`
				Content string `json:"content"`
			} `json:"searchTags,omitempty"`
			About []struct {
				Name string `json:"name"`
			} `json:"about,omitempty"`
		} `json:"value"`
	} `json:"webPages"`
	RelatedSearches struct {
		ID    string `json:"id"`
		Value []struct {
			Text         string `json:"text"`
			DisplayText  string `json:"displayText"`
			WebSearchURL string `json:"webSearchUrl"`
		} `json:"value"`
	} `json:"relatedSearches"`
	RankingResponse struct {
		Mainline struct {
			Items []struct {
				AnswerType  string `json:"answerType"`
				ResultIndex int    `json:"resultIndex"`
				Value       struct {
					ID string `json:"id"`
				} `json:"value"`
			} `json:"items"`
		} `json:"mainline"`
		Sidebar struct {
			Items []struct {
				AnswerType string `json:"answerType"`
				Value      struct {
					ID string `json:"id"`
				} `json:"value"`
			} `json:"items"`
		} `json:"sidebar"`
	} `json:"rankingResponse"`
}
//search function searches for the given query term and returns the list of first n urls to be crawled
func search(qt string, cfg *Config) ([]string, error) {

	resultlimit := cfg.Bing.Resultlimit
	endpoint := cfg.Bing.Endpoint
	token := cfg.Bing.Token

	searchTerm := qt

	// Declare a new GET request.
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		panic(err)
	}

	// Add the payload to the request.
	param := req.URL.Query()
	param.Add("q", searchTerm)
	param.Add("count", strconv.Itoa(resultlimit))
	req.URL.RawQuery = param.Encode()

	// Insert the request header.
	req.Header.Add("Ocp-Apim-Subscription-Key", token)

	// Create a new client.
	client := new(http.Client)

	// Send the request to Bing.
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// Close the response.
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Create a new answer.
	ans := new(BingAnswer)
	err = json.Unmarshal(body, &ans)
	if err != nil {
		fmt.Println(err)
	}

	var urllist []string
	for _, result := range ans.WebPages.Value {

		urllist = append(urllist, result.URL)
	}
	return urllist, err
}
