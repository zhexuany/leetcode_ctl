package html

import (
	"strings"

	"bytes"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
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