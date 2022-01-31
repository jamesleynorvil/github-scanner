package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Item struct {
	HtmlUrl string `json:"html_url"`
}

type ResponseData struct {
	Repos []Item `json:"Items"`
}

func GetMatchedReposList(query string, limit int) ([]string, error) {
	matchedRepos := []string{}
	params := url.Values{}
	params.Add("q", query)
	res, err := http.Get("https://api.github.com/search/repositories?" + params.Encode())
	if err != nil {
		return matchedRepos, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	if res.StatusCode != http.StatusOK {
		errorMsg := fmt.Sprintf("Unexpected status code: %d", res.StatusCode)
		return matchedRepos, errors.New(errorMsg)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return matchedRepos, err
	}
	responseData := ResponseData{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return matchedRepos, err
	}
	for i, repo := range responseData.Repos {
		if i+1 > limit {
			break
		}
		matchedRepos = append(matchedRepos, fmt.Sprintf("%s.git", repo.HtmlUrl))
	}
	return matchedRepos, nil
}
