package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"is_admin"`
}

type Problem struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	InputDesc    string `json:"input_desc"`
	OutputDesc   string `json:"output_desc"`
	SampleInput  string `json:"sample_input"`
	SampleOutput string `json:"sample_output"`
	TestInput    string `json:"test_input"`
	TestOutput   string `json:"test_output"`
	TimeLimit    int    `json:"time_limit"`
	MemoryLimit  int    `json:"memory_limit"`
	CreatedAt    string `json:"created_at"`
}

type Submission struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	ProblemID    int    `json:"problem_id"`
	Code         string `json:"code"`
	Language     string `json:"language"`
	Verdict      string `json:"verdict"`
	Runtime      int    `json:"runtime"`
	Memory       int    `json:"memory"`
	SubmittedAt  string `json:"submitted_at"`
	Username     string `json:"username,omitempty"`
	ProblemTitle string `json:"problem_title,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SubmissionRequest struct {
	ProblemID int    `json:"problem_id"`
	Code      string `json:"code"`
	Language  string `json:"language"`
}

var db *sql.DB
var curDir string

func main() {

	var err error
	curDir, err = os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	db, err = sql.Open("sqlite3", "./judge.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	initDB()
	createTempDir()

	r := mux.NewRouter()

	// Auth routes
	r.HandleFunc("/api/register", registerHandler).Methods("POST")
	r.HandleFunc("/api/login", loginHandler).Methods("POST")

	// Problem routes
	r.HandleFunc("/api/problems", getProblemsHandler).Methods("GET")
	r.HandleFunc("/api/problems/{id}", getProblemHandler).Methods("GET")
	r.HandleFunc("/api/problems", createProblemHandler).Methods("POST")

	// Submission routes
	r.HandleFunc("/api/submit", submitSolutionHandler).Methods("POST")
	r.HandleFunc("/api/submissions", getSubmissionsHandler).Methods("GET")
	r.HandleFunc("/api/submissions/user/{user_id}", getUserSubmissionsHandler).Methods("GET")

	// Serve static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(r)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func initDB() {
	// Users table
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		is_admin BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Problems table
	problemTable := `
	CREATE TABLE IF NOT EXISTS problems (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		input_desc TEXT NOT NULL,
		output_desc TEXT NOT NULL,
		sample_input TEXT NOT NULL,
		sample_output TEXT NOT NULL,
		test_input TEXT NOT NULL,
		test_output TEXT NOT NULL,
		time_limit INTEGER DEFAULT 1000,
		memory_limit INTEGER DEFAULT 256,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Submissions table
	submissionTable := `
	CREATE TABLE IF NOT EXISTS submissions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		problem_id INTEGER NOT NULL,
		code TEXT NOT NULL,
		language TEXT NOT NULL,
		verdict TEXT DEFAULT 'Pending',
		runtime INTEGER DEFAULT 0,
		memory INTEGER DEFAULT 0,
		submitted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (problem_id) REFERENCES problems(id)
	);`

	db.Exec(userTable)
	db.Exec(problemTable)
	db.Exec(submissionTable)

	// Create default admin user
	createDefaultAdmin()
	// Create sample problem
	createSampleProblem()
}

func createDefaultAdmin() {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", "admin").Scan(&count)
	if err != nil || count > 0 {
		return
	}

	hasher := md5.New()
	hasher.Write([]byte("admin123"))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	_, err = db.Exec("INSERT INTO users (username, email, password, is_admin) VALUES (?, ?, ?, ?)",
		"admin", "admin@judge.com", hashedPassword, true)
	if err != nil {
		log.Printf("Error creating admin user: %v", err)
	}
}

func createSampleProblem() {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM problems").Scan(&count)
	if err != nil || count > 0 {
		return
	}

	sampleProblem := `
	INSERT INTO problems (title, description, input_desc, output_desc, sample_input, sample_output, test_input, test_output)
	VALUES (
		'A + B Problem',
		'Given two integers A and B, compute A + B.',
		'Two integers A and B (1 ≤ A, B ≤ 1000)',
		'Output A + B',
		'2 3',
		'5',
		'10 20',
		'30'
	);`

	db.Exec(sampleProblem)
}

func createTempDir() {
	os.MkdirAll("./temp", 0755)
}

// Auth handlers
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	hasher := md5.New()
	hasher.Write([]byte(req.Password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	_, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		req.Username, req.Email, hashedPassword)
	if err != nil {
		http.Error(w, "Username or email already exists", http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	hasher := md5.New()
	hasher.Write([]byte(req.Password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	var user User
	err := db.QueryRow("SELECT id, username, email, is_admin FROM users WHERE username = ? AND password = ?",
		req.Username, hashedPassword).Scan(&user.ID, &user.Username, &user.Email, &user.IsAdmin)

	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Problem handlers
func getProblemsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, description, input_desc, output_desc, sample_input, sample_output, time_limit, memory_limit, created_at FROM problems ORDER BY id")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var problems []Problem
	for rows.Next() {
		var p Problem
		err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.InputDesc, &p.OutputDesc,
			&p.SampleInput, &p.SampleOutput, &p.TimeLimit, &p.MemoryLimit, &p.CreatedAt)
		if err != nil {
			continue
		}
		problems = append(problems, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(problems)
}

func getProblemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var p Problem
	err := db.QueryRow("SELECT id, title, description, input_desc, output_desc, sample_input, sample_output, time_limit, memory_limit, created_at FROM problems WHERE id = ?", id).
		Scan(&p.ID, &p.Title, &p.Description, &p.InputDesc, &p.OutputDesc, &p.SampleInput, &p.SampleOutput, &p.TimeLimit, &p.MemoryLimit, &p.CreatedAt)

	if err != nil {
		http.Error(w, "Problem not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func createProblemHandler(w http.ResponseWriter, r *http.Request) {
	var p Problem
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`INSERT INTO problems (title, description, input_desc, output_desc, sample_input, sample_output, test_input, test_output, time_limit, memory_limit) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Title, p.Description, p.InputDesc, p.OutputDesc, p.SampleInput, p.SampleOutput, p.TestInput, p.TestOutput, p.TimeLimit, p.MemoryLimit)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	p.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// Submission handlers
func submitSolutionHandler(w http.ResponseWriter, r *http.Request) {
	var req SubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Get user ID from header (simplified auth)
	userIDStr := r.Header.Get("User-ID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Insert submission
	result, err := db.Exec("INSERT INTO submissions (user_id, problem_id, code, language) VALUES (?, ?, ?, ?)",
		userID, req.ProblemID, req.Code, req.Language)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	submissionID, _ := result.LastInsertId()

	// Judge the submission
	go judgeSubmission(int(submissionID), req.ProblemID, req.Code, req.Language)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"submission_id": submissionID,
		"message":       "Submission received",
	})
}

func getSubmissionsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT s.id, s.user_id, s.problem_id, s.language, s.verdict, s.runtime, s.memory, s.submitted_at, u.username, p.title
		FROM submissions s
		JOIN users u ON s.user_id = u.id
		JOIN problems p ON s.problem_id = p.id
		ORDER BY s.id DESC
		LIMIT 50
	`)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var submissions []Submission
	for rows.Next() {
		var s Submission
		err := rows.Scan(&s.ID, &s.UserID, &s.ProblemID, &s.Language, &s.Verdict,
			&s.Runtime, &s.Memory, &s.SubmittedAt, &s.Username, &s.ProblemTitle)
		if err != nil {
			continue
		}
		submissions = append(submissions, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(submissions)
}

func getUserSubmissionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	rows, err := db.Query(`
		SELECT s.id, s.user_id, s.problem_id, s.code, s.language, s.verdict, s.runtime, s.memory, s.submitted_at, p.title
		FROM submissions s
		JOIN problems p ON s.problem_id = p.id
		WHERE s.user_id = ?
		ORDER BY s.id DESC
	`, userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var submissions []Submission
	for rows.Next() {
		var s Submission
		err := rows.Scan(&s.ID, &s.UserID, &s.ProblemID, &s.Code, &s.Language,
			&s.Verdict, &s.Runtime, &s.Memory, &s.SubmittedAt, &s.ProblemTitle)
		if err != nil {
			continue
		}
		submissions = append(submissions, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(submissions)
}

// Judge system
func judgeSubmission(submissionID, problemID int, code, language string) {
	// Get problem test cases
	var testInput, testOutput string
	var timeLimit int
	err := db.QueryRow("SELECT test_input, test_output, time_limit FROM problems WHERE id = ?", problemID).
		Scan(&testInput, &testOutput, &timeLimit)
	if err != nil {
		updateVerdict(submissionID, "System Error", 0, 0)
		return
	}

	verdict, runtime, memory := executeCode(code, language, testInput, testOutput, timeLimit)
	updateVerdict(submissionID, verdict, runtime, memory)
}

func executeCode(code, language, input, expectedOutput string, timeLimit int) (string, int, int) {
	tempDir := fmt.Sprintf("/temp/sub_%d", time.Now().UnixNano())
	tempDir = filepath.Join(curDir, tempDir)
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	var filename, compileCmd, runCmd string

	switch language {
	case "cpp":
		filename = "solution.cpp"
		compileCmd = fmt.Sprintf("/opt/homebrew/bin/g++-14 -w -std=c++20 -o %s/solution %s/%s", tempDir, tempDir, filename)
		runCmd = fmt.Sprintf("%s/solution", tempDir)
	case "python":
		filename = "solution.py"
		runCmd = fmt.Sprintf("python3 %s/%s", tempDir, filename)
	case "java":
		filename = "Solution.java"
		compileCmd = fmt.Sprintf("javac -d %s %s/%s", tempDir, tempDir, filename)
		runCmd = fmt.Sprintf("java -cp %s Solution", tempDir)
	default:
		return "Compilation Error : Unknown compiler", 0, 0
	}

	fmt.Println()
	fmt.Println()
	fmt.Println(code)
	fmt.Println(tempDir)

	// Write code to file
	codePath := filepath.Join(tempDir, filename)
	err := ioutil.WriteFile(codePath, []byte(code), 0644)
	// err := os.WriteFile(codePath, []byte(code), 0644)
	// if err != nil {
	// 	return "System Error", 0, 0
	// }

	// Compile if needed
	if compileCmd != "" {
		cmd := exec.Command("bash", "-c", compileCmd)
		cmd.Dir = tempDir
		// if err := cmd.Run(); err != nil {
		// 	fmt.Println(err)
		// 	return "Compilation Error``", 0, 0
		// }

		// capture stderr + stdout
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Compilation failed with error:")
			fmt.Println(string(out)) // this prints compiler error like syntax errors
			return "Compilation Error", 0, 0
		}
	}

	// Run code
	start := time.Now()
	cmd := exec.Command("bash", "-c", runCmd)
	cmd.Dir = tempDir
	cmd.Stdin = strings.NewReader(input)

	output, err := cmd.Output()

	fmt.Println("Running: ", string(output))

	runtime := int(time.Since(start).Milliseconds())

	if err != nil {
		return "Runtime Error", runtime, 0
	}

	// Check output
	actualOutput := strings.TrimSpace(string(output))
	expectedOutput = strings.TrimSpace(expectedOutput)

	if actualOutput == expectedOutput {
		return "Accepted", runtime, 0
	} else if runtime > timeLimit {
		return "Time Limit", runtime, 0
	} else {
		return "Wrong Answer", runtime, 0
	}
}

func updateVerdict(submissionID int, verdict string, runtime, memory int) {
	_, err := db.Exec("UPDATE submissions SET verdict = ?, runtime = ?, memory = ? WHERE id = ?",
		verdict, runtime, memory, submissionID)
	if err != nil {
		log.Printf("Error updating verdict: %v", err)
	}
}
