package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/zhexuany/leetcode-ctl/config"
)

const (
	LEETCODE_BASE_URL = "https://leetcode.com"
)

type Client struct {
	config     *config.Config
	logger     *log.Logger
	httpClient *http.Client
	// postgresDB PostgresDB
}

func NewClient(conf *config.Config) (*Client, error) {
	return &Client{
		logger:     log.New(os.Stderr, "[submit] ", log.LstdFlags),
		config:     conf,
		httpClient: &http.Client{},
		//TODO add db_name etc
		// postgresDB: PostgresDB{},
	}, nil
}

func (c *Client) setReqHeader(req *http.Request) {
	cookieStr := fmt.Sprintf("LEETCODE_SESSION=%s; csrftoken=%s;", c.config.LeetcodeSession, c.config.CsrfToken)
	req.Header.Set("Cookie", cookieStr)
	req.Header.Set("Origin", LEETCODE_BASE_URL)
	req.Header.Set("X-CSRFToken", c.config.CsrfToken)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-US,en;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://leetcode.com/problems/two-sum/")
}

type submissionID struct {
	SubmissionID int `json:"submission_id"`
}

type submissionContent struct {
	QuestionID int    `json:"question_id"`
	Lang       string `json:"lang"`
	TypedCode  string `json:"typed_code"`
	DataInput  string `json:"data_input"`
	TestMode   bool   `json:"test_mode"`
	JudgeType  string `json:"judge_type"`
}

type checkStatus struct {
	Lang           string `json:"lang"`
	TotalTestcases int    `json:"total_testcases"`
	UserID         int    `json:"user_id"`
	CodeOutput     string `json:"code_output"`
	StatusCode     int    `json:"status_code"`
	StatusRuntime  string `json:"status_runtime"`
	CompareResult  string `json:"compare_result"`
	DisplayRuntime string `json:"display_runtime"`
	State          string `json:"state"`
	TotalCorrect   int    `json:"total_correct"`
	RunSuccess     bool   `json:"run_success"`
	JudgeType      string `json:"judge_type"`
	StdOutput      string `json:"std_output"`
	QuestionID     int    `json:"question_id"`
}

func decode(resp *http.Response) (io.ReadCloser, error) {
	if resp.Header.Get("Content-Encoding") == "gzip" {
		resp.Header.Del("Content-Length")
		zr, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		return gzreadCloser{zr, resp.Body}, nil
	}
	return resp.Body, nil
}
func parseFileContents(path string) (*bytes.Buffer, error) {
	sol := submissionContent{}
	sol.QuestionID = 1
	sol.Lang = "golang"
	sol.TestMode = false
	sol.JudgeType = "large"
	// sol.DataInput = "[3,2,4]\n6"
	sol.TypedCode = "func twoSum(nums []int, target int) []int {\r\n   //naive approach\r\n   for i := 0; i < len(nums); i++{\r\n       for j := i + 1; j < len(nums); j++ {\r\n           if nums[i] + nums[j] == target {\r\n               return []int{i, j}\r\n           }\r\n       }\r\n   }\r\n   return nil\r\n}"

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(sol); err != nil {
		return nil, err
	}

	return b, nil
}

// Submit will read contents in such file and submit
func (c *Client) Submit(path string) error {
	s := spinner.New(spinner.CharSets[36], 100*time.Millisecond)
	s.Prefix = "Submitting"
	s.Start()
	b, err := parseFileContents(path)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://leetcode.com/problems/two-sum/submit/", b)
	if err != nil {
		return err
	}

	c.setReqHeader(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	resp.Body, err = decode(resp)
	if err != nil {
		return err
	}

	id := submissionID{}
	if err := json.NewDecoder(resp.Body).Decode(&id); err != nil {
		return err
	}

	c.logger.Println("got id", id.SubmissionID)
	s.Stop()
	checkURL := fmt.Sprintf("https://leetcode.com/submissions/detail/%d/check/", id.SubmissionID)

	c.logger.Println("Waiting for 1 second, since leetcode judger need a few monents to determine your answer is correct or not.")

	s.Prefix = "Checking"
	s.Restart()
	time.Sleep(1 * time.Second)
	checkReq, err := http.NewRequest("GET", checkURL, nil)
	c.setReqHeader(checkReq)
	checkResp, err := http.DefaultClient.Do(checkReq)
	if err != nil {
		c.logger.Println(err)
		return err
	}

	defer checkResp.Body.Close()
	if checkResp.Header.Get("Content-Encoding") == "gzip" {
		checkResp.Header.Del("Content-Length")
		zr, err := gzip.NewReader(checkResp.Body)
		if err != nil {
			return err
		}
		checkResp.Body = gzreadCloser{zr, checkResp.Body}
	}

	status := checkStatus{}
	if err := json.NewDecoder(checkResp.Body).Decode(&status); err != nil {
		return err
	}

	s.Stop()
	if status.RunSuccess {
		c.logger.Printf("Congras. You solve this problem with runtime %s", status.StatusRuntime)
		c.logger.Printf("You passed %d test out of total test %d", status.TotalCorrect, status.TotalTestcases)
		return nil
	}

	c.logger.Println("your answer is not correct and the reason is", status.DisplayRuntime)
	c.logger.Printf("You passed %d test out of total test %d", status.TotalCorrect, status.TotalTestcases)
	return nil
}

type gzreadCloser struct {
	*gzip.Reader
	io.Closer
}

func (gz gzreadCloser) Close() error {
	return gz.Closer.Close()
}

func (c *Client) Search(v interface{}) error {
	return nil
}
