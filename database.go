package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "test"
)

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

func getJsonObjectFromLeetcode() ProblemsJson {
	req, err := http.NewRequest("GET", "https://leetcode.com/api/problems/algorithms/", nil)
	if err != nil {
		fmt.Println("failed to get reply from leetcode", err)
	}
	req.Header.Set("Cookie", "_gat=1; csrftoken=xKpOqsNQvNKCzhXyEwUVMhThfejftmsmUsbjvMTVOO0awGBRnrP0Ogvad2HSmXpj; _ga=GA1.2.1442407930.1481621879")
	req.Header.Set("Dnt", "1")
	//TODO need figure why uncomment this does work
	// req.Header.Set("Accept-Encoding", "gzip, deflate, sdch, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.6,ja;q=0.4,en;q=0.2")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Referer", "https://leetcode.com/problemset/algorithms/")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Connection", "keep-alive")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("failed to Do request", err)
	}
	defer resp.Body.Close()

	lp := ProblemsJson{}

	err = json.NewDecoder(resp.Body).Decode(&lp)
	if err != nil {
		fmt.Printf("failed to decode json object %v", err)
	}
	return lp
}

func dbINFO() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
}

type PostgresDB struct {
	db     *sql.DB
	logger log.Logger
}

func (psql *PostgresDB) Open() {
	db, err := sql.Open("postgres", dbINFO())
	if err != nil {
		fmt.Println("failed to open sql driver")
	}

	psql.db = db
}

func (psql *PostgresDB) Close() {
	psql.db.Close()
}

func (psql *PostgresDB) write() {
	ls := getJsonObjectFromLeetcode()
	length := len(ls.StatStatusPairs)
	for i := length - 1; i >= 0; i-- {
		data := ls.StatStatusPairs[i]
		//TODO; need think about how to add problem statements and also the solution
		queryStr := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(question_id int, question_title text, difficulty int);`, "leetcode")
		if _, err := psql.db.Exec(queryStr); err != nil {
			fmt.Println("failed to exec", err)
			return
		}
		insertStat := fmt.Sprintf(`INSERT INTO %s (question_id, question_title, difficulty) VALUES(%d, '%s', %d);`, "leetcode",
			data.Stat.QuestionID, data.Stat.QuestionTitleSlug, data.Difficulty.Level)
		if _, err := psql.db.Exec(insertStat); err != nil {
			fmt.Println("failed to exec", err)
			return
		}
	}
}

func (psql *PostgresDB) udpate() {
	//TODO(zhexuany) does not have some good methods about how to update database
	psql.write()
}

// queryByQuestionID can query from db according the questionID
func queryByQuestionID(questionID int, db *sql.DB) *sql.Rows {
	var queryStr = `select * from leetcode where question_id= %d`
	queryStr = fmt.Sprintf(queryStr, questionID)

	rows, err := db.Query(queryStr)
	if err != nil {
		fmt.Println("failed to query from database")
	}

	return rows
}

func (psql *PostgresDB) Query(key interface{}) {
	var rows *sql.Rows
	switch value := key.(type) {
	case int:
		rows = queryByQuestionID(value, psql.db)
	case string:
	default:
		fmt.Println("not allowed type")
	}

	for rows.Next() {
		var questionID int
		var questionTitle string
		var questionDifficulty int
		if err := rows.Scan(&questionID, &questionTitle, &questionDifficulty); err != nil {
			fmt.Println("failed to query from postgres", err)
		}
		res := fmt.Sprintf("question_id:%d question_title:%s difficulty:%d", questionID, questionTitle, questionDifficulty)
		fmt.Println(res)
	}

}
