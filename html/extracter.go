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

const html = `
<!DOCTYPE html>
<html>
<body>
  <div class="container"
        ng-app="app"
        ng-controller="AceCtrl as aceCtrl"
        ng-init="aceCtrl.init(
        [{'value': 'cpp', 'text': 'C++', 'defaultCode': 'class Solution {\u000D\u000Apublic:\u000D\u000A    vector\u003Cint\u003E twoSum(vector\u003Cint\u003E\u0026 nums, int target) {\u000D\u000A        \u000D\u000A    }\u000D\u000A}\u003B' },{'value': 'java', 'text': 'Java', 'defaultCode': 'public class Solution {\u000D\u000A    public int[] twoSum(int[] nums, int target) {\u000D\u000A        \u000D\u000A    }\u000D\u000A}' },{'value': 'python', 'text': 'Python', 'defaultCode': 'class Solution(object):\u000D\u000A    def twoSum(self, nums, target):\u000D\u000A        \u0022\u0022\u0022\u000D\u000A        :type nums: List[int]\u000D\u000A        :type target: int\u000D\u000A        :rtype: List[int]\u000D\u000A        \u0022\u0022\u0022\u000D\u000A        ' },{'value': 'c', 'text': 'C', 'defaultCode': '/**\u000D\u000A * Note: The returned array must be malloced, assume caller calls free().\u000D\u000A */\u000D\u000Aint* twoSum(int* nums, int numsSize, int target) {\u000D\u000A    \u000D\u000A}' },{'value': 'csharp', 'text': 'C#', 'defaultCode': 'public class Solution {\u000D\u000A    public int[] TwoSum(int[] nums, int target) {\u000D\u000A        \u000D\u000A    }\u000D\u000A}' },{'value': 'javascript', 'text': 'JavaScript', 'defaultCode': '/**\u000D\u000A * @param {number[]} nums\u000D\u000A * @param {number} target\u000D\u000A * @return {number[]}\u000D\u000A */\u000D\u000Avar twoSum \u003D function(nums, target) {\u000D\u000A    \u000D\u000A}\u003B' },{'value': 'ruby', 'text': 'Ruby', 'defaultCode': '# @param {Integer[]} nums\u000D\u000A# @param {Integer} target\u000D\u000A# @return {Integer[]}\u000D\u000Adef two_sum(nums, target)\u000D\u000A    \u000D\u000Aend' },{'value': 'swift', 'text': 'Swift', 'defaultCode': 'class Solution {\u000D\u000A    func twoSum(_ nums: [Int], _ target: Int) \u002D\u003E [Int] {\u000D\u000A        \u000D\u000A    }\u000D\u000A}' },{'value': 'golang', 'text': 'Go', 'defaultCode': 'func twoSum(nums []int, target int) []int {\u000D\u000A    \u000D\u000A}' },],
        '1_210021',
        1,
        '/problems/two-sum/interpret_solution/',
        '/problems/two-sum/submit/',
        '/submissions/detail/0/',
        '/problems/0/',
        'Two Sum',
        '[3,2,4]\u000A6',
        '{\u000D\u000A  \u0022name\u0022: \u0022twoSum\u0022,\u000D\u000A  \u0022params\u0022: [\u000D\u000A    {\u000D\u000A      \u0022name\u0022: \u0022nums\u0022,\u000D\u000A      \u0022type\u0022: \u0022integer[]\u0022\u000D\u000A    },\u000D\u000A    {\u000D\u000A      \u0022name\u0022: \u0022target\u0022,\u000D\u000A      \u0022type\u0022: \u0022integer\u0022\u000D\u000A    }\u000D\u000A  ],\u000D\u000A  \u0022return\u0022: {\u000D\u000A    \u0022type\u0022: \u0022integer[]\u0022,\u000D\u000A    \u0022size\u0022: 2\u000D\u000A  }\u000D\u000A}',
        true,
        [{&quot;question_title_slug&quot;: &quot;two-sum-ii-input-array-is-sorted&quot;, &quot;question_title&quot;: &quot;Two Sum II - Input array is sorted&quot;, &quot;difficulty&quot;: &quot;E&quot;}, {&quot;question_title_slug&quot;: &quot;two-sum-iii-data-structure-design&quot;, &quot;question_title&quot;: &quot;Two Sum III - Data structure design&quot;, &quot;difficulty&quot;: &quot;E&quot;}],
        '',
        false,
        'large'
        );" ng-cloak>
    <input type='hidden' name='csrfmiddlewaretoken' value='zBpM53wZmqZDVbPBSFdZs1UekDo61il1mUSOF6b1IwN5QAEZaJFza3QfLZDiL1hC' />
    <hr>
    <div class="row" style="margin-bottom:12px;">
      <div class="col-md-12">
        <select class="form-control select mbm" name="lang" ng-model="aceCtrl.selectedLang" ng-options="lang.text for lang in aceCtrl.langs track by lang.value"></select>
        <refresh-button></refresh-button>
        <code-button></code-button>
      </div>
    </div>
  </body>
</html>
`

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

func (e *Extracter) Find() *Extracter {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
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

func QueryByID() {
	// lp := ProblemsJson{}
	// err = json.NewDecoder(resp.Body).Decode(&lp)
	// if err != nil {
	// 	fmt.Printf("failed to decode json object %v", err)
	// }
}
