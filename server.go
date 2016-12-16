package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"strings"
)

const (
	Leetcode_base_url             = "https://leetcode.com"
	Leetcode_normal_login_action  = "accounts/login/"
	Leetcode_github_logtin_action = "accounts/github/login"
	Leetcode_run_url              = "problems/%s/interpret_solution/"
)

type Server struct {
	handler *handler
	psql    *PostgresDB
}

func NewServer() *Server {
	return &Server{
		handler: newHandler(),
		psql:    &PostgresDB{},
	}
}

func (s *Server) Open() {
	s.psql.Open()
	server := http.Server{
		Addr:    "localhost:40000",
		Handler: s.handler,
	}
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("failed to run server", err)
	}
}

//json Object for replying
type InterpreReplyJson struct {
	CodeOutput     []interface{} `json:"code_output"`
	StatusCode     int           `json:"status_code"`
	StatusRuntime  string        `json:"status_runtime"`
	State          string        `json:"state"`
	TotalCorrect   interface{}   `json:"total_correct"`
	CompileError   string        `json:"compile_error"`
	RunSuccess     bool          `json:"run_success"`
	TotalTestcases interface{}   `json:"total_testcases"`
}

type InterpreJson struct {
	InterpretExpectedID string `json:"interpret_expected_id"`
	InterpretID         string `json:"interpret_id"`
	TestCase            string `json:"test_case"`
}

type handler struct {
	logger *log.Logger
	// cookie *http.Cookie
	client http.Client
}

func newHandler() *handler {
	jar, _ := cookiejar.New(nil)
	h := &handler{
		client: http.Client{
			Jar: jar,
		},
	}
	return h
}

func (h *handler) WrapHandler(name string, hf http.HandlerFunc) http.Handler {
	var handler http.Handler
	handler = http.HandlerFunc(hf)
	return handler
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		switch r.URL.Path {
		case "/login":
			h.WrapHandler("login", h.serveLogin).ServeHTTP(w, r)
		case "/logout":
			h.WrapHandler("logout", h.serveLogout).ServeHTTP(w, r)
		case "/run":
			h.WrapHandler("run", h.serveRun).ServeHTTP(w, r)
		case "/submit":
			h.WrapHandler("submit", h.serveSubmit).ServeHTTP(w, r)
		case "/search":
			h.WrapHandler("search", h.serveSearch).ServeHTTP(w, r)
		}
	case "POST":
		switch r.URL.Path {
		case "/login":
			h.WrapHandler("login", h.serveLogin).ServeHTTP(w, r)
		case "/logout":
			h.WrapHandler("logout", h.serveLogout).ServeHTTP(w, r)
		}
	default:
		http.Error(w, "", http.StatusBadRequest)
	}
}

func (h *handler) isLogged() bool {
	//TODO fix this login later, doesn't seem correct
	return h.client.Jar != nil
}

func (h *handler) serveLogin(w http.ResponseWriter, r *http.Request) {
	//extract value from request
	uname := r.URL.Query().Get("login")
	pass := r.URL.Query().Get("password")

	if uname == "" || pass == "" {
		http.Error(w, "user name or password is not set", http.StatusBadRequest)
		return
	}

	loginURL := Leetcode_base_url + "/" + Leetcode_normal_login_action
	bodyStr := fmt.Sprintf("%s=%s&login=%s&password=%s", "csrfmiddlewaretoken", "qlkaLSFNuSFQ91s3lop8GnqA41lOYE7nU8nPAy0vytyD78yQVE22dCsRYSs6GYfs", uname, pass)
	body := strings.NewReader(bodyStr)
	req, err := http.NewRequest("POST", loginURL, body)
	if err != nil {
		fmt.Println("failed to post request", err)
	}

	req.Header = r.Header
	req.Header.Set("Referer", loginURL)
	if dump, err := httputil.DumpRequest(req, true); err == nil {
		fmt.Printf("%q", dump)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		h.logger.Println("failed to postfrom client", err)
		return
	}

	fmt.Printf("statuscode is %d", resp.StatusCode)
	if resp.StatusCode == http.StatusFound {
		_, _ = w.Write([]byte("login sucess"))
	}

	switch resp.Header.Get("Transfer-Encoding") {
	case "chunked":
		reader := httputil.NewChunkedReader(resp.Body)
		if buf, err := ioutil.ReadAll(reader); err != nil {
			fmt.Println(string(buf))
		}
	default:
		if buf, err := ioutil.ReadAll(resp.Body); err != nil {
			fmt.Println(string(buf))
		}
	}

	defer resp.Body.Close()
}

func (h *handler) serveLogout(w http.ResponseWriter, r *http.Request) {
	//TODO do request
	//make sure set cookie in halder to nil
}

func (h *handler) serveRun(w http.ResponseWriter, r *http.Request) {
	//TODO how to get parameter from r
	params := r.URL.Query()
	problem_title := params.Get("pb")
	if problem_title == "" {
		http.Error(w, "pb must be seted", http.StatusExpectationFailed)
	}
	url := fmt.Sprintf(Leetcode_run_url, problem_title)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
	}
	//make sure cookie is not nil
	if h.isLogged() {
		http.Error(w, "you should login first", http.StatusExpectationFailed)
		return
	}

	//Form header in client side
	req.Header = r.Header

	resp, err := h.client.Do(req)
	if err != nil {
		http.Error(w, "failed to do request", http.StatusExpectationFailed)
		return
	}
	defer resp.Body.Close()

	ij := InterpreJson{}
	err = json.NewDecoder(resp.Body).Decode(&ij)
	if err != nil {
		http.Error(w, "failed to decode json object", http.StatusInternalServerError)
	}

	//TODO parse resp as jsob object and form a request again to get
	// {"interpret_expected_id": "interpret_expected_1481731971.8_209932_4", "interpret_id": "interpret_1481731971.8_209932_4", "test_case": "[3,2,4]\n6"}%
	// data-input
	checkURL := LEETCODE_BASE_URL + "/"
	checkURL += fmt.Sprintf(Leetcode_run_url, ij.InterpretExpectedID)
	req, err = http.NewRequest("GET", checkURL, nil)
	if err != nil {
		http.Error(w, "failed to create request for check result", http.StatusInternalServerError)
	}

	req = setReqHeader(req)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "failed to do request for check result", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	//TODO(zhexuany) do we really need to decode&encode in server side
	irj := InterpreReplyJson{}
	err = json.NewDecoder(resp.Body).Decode(&irj)
	if err != nil {
		http.Error(w, "failed to decode json object", http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(irj); err != nil {
		http.Error(w, "failed to decode json object", http.StatusInternalServerError)
	}
	return
}

func (h *handler) serveSearch(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) serveSubmit(w http.ResponseWriter, r *http.Request) {

}
