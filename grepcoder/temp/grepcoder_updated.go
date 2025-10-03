package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
)

// const (
// 	SET_USER = 1
// 	TG       = 100
// 	SUB_STAT = 101
// 	PB_SYNC  = 200
// 	SUB_SYNC = 201
// )

const (
	problemSet           = "gstore.json"
	userSubmissionStatus = "userstatus.json"
	userId               = "user.json"
	kenkooApiUrl         = "https://kenkoooo.com/atcoder/atcoder-api/v3/user/submissions?user=%s&from_second=%d"
)

const (
	emptyJsonErr = "unexpected end of JSON input"
)

type UserInfo struct {
	User      string `json:"user"`
	Timestamp int    `json"ts"`
	Updated   bool
}

func NewUserInfo() (*UserInfo, error) {
	var user UserInfo

	buf, err := os.ReadFile(userId)
	if err != nil {
		return &user, err
	}

	if err := json.Unmarshal(buf, &user); err != nil {
		return &user, err
	}

	return &user, nil
}

func (u *UserInfo) UpdateUser(user string) {
	u.User = user
}

func (u *UserInfo) UpdateTimestamp(ts int) {
	u.Timestamp = ts
}

func (u *UserInfo) SetUpdate() {
	u.Updated = true
}

func (u *UserInfo) UpdateDB() (err error) {
	if !u.Updated {
		return
	}

	if err = os.Truncate(userId, 0); err != nil {
		return
	}

	var buf []byte
	buf, err = json.Marshal(buf)
	if err != nil {
		return
	}

	if err = os.WriteFile(userId, buf, 0644); err != nil {
		return
	}
	return
}

type problemModel struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	ContestId    string   `json:"contest_id"`
	EditorialUrl string   `json:"edi_url"`
	Tags         []string `json:"tag"`
	Status       string
}

func (p *problemModel) SetStatus(stat string) {
	p.Status = stat
}

type subStatus struct {
	Id     string `json:"problem_id"`
	Result string `json:"result"`
	TS     int    `json:"epoch_second"`
}

var pl []problemModel
var us []subStatus
var subStat map[string]string

func init() {
	subStat = make(map[string]string)
}

func Setup() {
	f, err := os.OpenFile(problemSet, os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(buf, &pl); err != nil {
		panic(err)
	}

	statusF, err := os.OpenFile(userSubmissionStatus, os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer statusF.Close()

	buf, err = io.ReadAll(statusF)
	if err != nil {
		panic(err)
	}

	if len(buf) == 0 {
		return
	}

	if err := json.Unmarshal(buf, &us); err != nil {
		panic(err)
	}

	for _, stat := range us {
		if v, ok := subStat[stat.Id]; !ok {
			subStat[stat.Id] = stat.Result
		} else if ok && v != "AC" && stat.Result == "AC" {
			subStat[stat.Id] = stat.Result
		}
	}

	for i := 0; i < len(pl); i++ {
		if v, ok := subStat[pl[i].Id]; ok {
			pl[i].SetStatus(v)
		}
	}
}

func SyncSubmission(user string, prev int) {
	for {
		url := fmt.Sprintf(kenkooApiUrl, strings.TrimSpace(string(user)), prev)
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()

		if err := json.Unmarshal(buf, &us); err != nil {
			panic(err)
		}

		for _, stat := range us {
			if v, ok := subStat[stat.Id]; !ok {
				subStat[stat.Id] = stat.Result
			} else if ok && v != "AC" && stat.Result == "AC" {
				subStat[stat.Id] = stat.Result
			}
		}

		for i := 0; i < len(pl); i++ {
			if v, ok := subStat[pl[i].Id]; ok {
				pl[i].SetStatus(v)
			}
		}

		if prev == us[len(us)-1].TS {
			break
		}

		prev = us[len(us)-1].TS
	}

	us = make([]subStatus, len(subStat))
	for key, value := range subStat {
		us = append(us, subStatus{Id: key, Result: value})
	}
}

func WriteSynacedSubmission() {
	temp, err := os.OpenFile("temp_"+userSubmissionStatus, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer temp.Close()

	buf, err := json.Marshal(us)
	if err != nil {
		panic(err)
	}

	n, err := temp.Write(buf)
	if n != len(buf) || err != nil {
		panic(err)
	}

	if err := os.Remove(userSubmissionStatus); err != nil {
		panic(err)
	}

	if err := os.Rename(temp.Name(), userSubmissionStatus); err != nil {
		panic(err)
	}
}

func Seperate(s string) (key, value string) {
	if len(s) == 0 {
		return
	}
	temp := strings.Split(s, "=")
	key = strings.TrimSpace(temp[0])

	if len(temp) >= 2 {
		value = strings.TrimSpace(temp[1])
	}
	return
}

func FindProblemByTags(tags, stat string) {
	tgSplit := strings.Split(tags, ",")
	mp := make([]int, 0, len(pl))
	for _, tag := range tgSplit {
		tag = strings.TrimSpace(tag)
		for i := 0; i < len(pl); i++ {
			for j := 0; j < len(pl[i].Tags); j++ {
				if pl[i].Tags[j] == tag && ((stat == "ALL") || (stat == pl[i].Status)) {
					mp = append(mp, i)
					break
				}
			}
		}
	}

	sort.Slice(mp, func(i, j int) bool {
		return pl[mp[i]].Id < pl[mp[j]].Id
	})

	for _, key := range mp {
		status := pl[key].Status
		seq := "\033[0m"
		switch status {
		case "AC":
			seq = "\033[32m"
		case "WA":
			fallthrough
		case "RE":
			fallthrough
		case "TLE":
			seq = "\033[31m"
		default:
			status = "NA"
		}

		fmt.Fprintln(os.Stdout, seq, "[", status, "]", pl[key].Id, pl[key].Tags, "\033[0m")

	}
}

func main() {

	userInfo, err := NewUserInfo()
	if err != nil && err.Error() != emptyJsonErr {
		panic(err)
	}
	_ = userInfo

	Setup()

	reader := bufio.NewReader(os.Stdin)

loop:
	for {

		fmt.Fprintf(os.Stdout, ">>> ")
		tg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error Reading..")
			continue
		}

		// cmd
		cmds := strings.Split(tg, ";")

		key, value := Seperate(cmds[0])

		switch key {
		case "TAG":
			var stat, ss string
			if len(cmds) >= 2 {
				ss, stat = Seperate(cmds[1])
				if ss != "SUB_STAT" {
					stat = "ALL"
				}
			}
			FindProblemByTags(value, stat)
		case "SET_USER":
		case "SYNC":
		case "SYNC_PL":
		case "SUB_STAT":
		case "QUIT":
			break loop
		default:
		}
	}

	if err := userInfo.UpdateDB(); err != nil {
		panic(err)
	}

	WriteSynacedSubmission()

}
