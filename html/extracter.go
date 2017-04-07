package html

import (
	"strings"

	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
)

type Extracter struct {
	html    string
	jsonStr string
	codes   Codes
}

func toJsonStr(str string) string {
	ret := ""
	// Have to replace single quote with double quote
	// cause we need marshal string into json object
	ret = strings.Replace(str, `'`, `"`, -1)
	prefixIdx := strings.Index(ret, "[")
	postfixIdx := strings.Index(ret, ",]")
	ret = ret[prefixIdx:postfixIdx]
	ret += "]"
	return ret
}

func (e *Extracter) Find(name string) *Extracter {
	reqURL := fmt.Sprintf("https://leetcode.com/problems/%s", name)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		fmt.Println("failed to get reply from leetcode", err)
	}
	// req.Header.Set("Accept-Encoding", "gzip, deflate, sdch, br")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Connection", "keep-alive")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("failed to Do request", err)
	}
	defer resp.Body.Close()

	bs, _ := ioutil.ReadAll(resp.Body)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bs)))
	if err != nil {
		panic(err)
	}

	ret := &Extracter{}
	doc.Find(".container").Each(func(i int, s *goquery.Selection) {
		if val, ok := s.Attr("ng-init"); ok {
			ret.jsonStr = toJsonStr(val)
		}
	})

	return ret
}

func (e *Extracter) Json() *Extracter {
	if e.jsonStr == "" {
		return nil
	}

	reader := bytes.NewBufferString(e.jsonStr)
	if err := json.NewDecoder(reader).Decode(&e.codes); err != nil {
		return nil
	}
	return e
}

func (e *Extracter) GetDefaultCode(language string) string {
	return e.codes.getDefaultCode(language)
}

type Codes []struct {
	Value       string `json:"value"`
	Text        string `json:"text"`
	DefaultCode string `json:"defaultCode"`
}

// Supported language: cpp, golang, java, python, swift, ruby, csharp, c and javascript
func (c Codes) getDefaultCode(language string) string {
	for _, val := range c {
		if val.Value == language {
			return val.DefaultCode
		}
	}
	return ""
}

type ProblemsJson struct {
	FrequencyMid    int    `json:"frequency_mid"`
	NumSolved       int    `json:"num_solved"`
	CategorySlug    string `json:"category_slug"`
	StatStatusPairs []struct {
		Status interface{} `json:"status"`
		Stat   struct {
			TotalAcs            int         `json:"total_acs"`
			QuestionTitle       string      `json:"question__title"`
			QuestionArticleSlug interface{} `json:"question__article__slug"`
			TotalSubmitted      int         `json:"total_submitted"`
			QuestionTitleSlug   string      `json:"question__title_slug"`
			QuestionArticleLive interface{} `json:"question__article__live"`
			QuestionHide        bool        `json:"question__hide"`
			QuestionID          int         `json:"question_id"`
		} `json:"stat"`
		IsFavor    bool `json:"is_favor"`
		PaidOnly   bool `json:"paid_only"`
		Difficulty struct {
			Level int `json:"level"`
		} `json:"difficulty"`
		Frequency int `json:"frequency"`
		Progress  int `json:"progress"`
	} `json:"stat_status_pairs"`
	IsPaid        bool   `json:"is_paid"`
	FrequencyHigh int    `json:"frequency_high"`
	UserName      string `json:"user_name"`
	NumTotal      int    `json:"num_total"`
	ListNames     struct {
	} `json:"list_names"`
}

func GetJsonObjectFromLeetcode() error {
	req, err := http.NewRequest("GET", "https://leetcode.com/api/problems/algorithms/", nil)
	if err != nil {
		fmt.Println("failed to get reply from leetcode", err)
	}
	// req.Header.Set("Accept-Encoding", "gzip, deflate, sdch, br")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Connection", "keep-alive")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("failed to Do request", err)
	}
	defer resp.Body.Close()

	bs, _ := ioutil.ReadAll(resp.Body)

	ioutil.WriteFile("problems.json", bs, 0644)
	return nil
}

//QueryByID will return
func QueryByID(id int) string {
	bs, err := ioutil.ReadFile("problems.json")
	if err != nil {
		// panic for now
		panic(err)
	}

	reader := bytes.NewReader(bs)
	lp := ProblemsJson{}
	err = json.NewDecoder(reader).Decode(&lp)
	if err != nil {
		fmt.Printf("failed to decode json object %v", err)
	}
	for _, val := range lp.StatStatusPairs {
		if val.Stat.QuestionID == id {
			return strings.ToLower(strings.Replace(val.Stat.QuestionTitle, " ", "-", -1))
		}
	}
	return ""
}
