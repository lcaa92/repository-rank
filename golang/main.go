package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

type RepositoryActivity struct {
	Repository   string
	Commits      int
	Users        []string
	Files        int
	Additions    int
	Deletions    int
	MinTimestamp int
	MaxTimestamp int
	Score        float64
}

func NewRepositoryActivity(repository string) *RepositoryActivity {
	return &RepositoryActivity{
		Repository: repository,
		Score:      0.00,
	}
}

func (r *RepositoryActivity) appendUser(username string) {
	for _, name := range r.Users {
		if name == username {
			return
		}
	}
	r.Users = append(r.Users, username)
}

func (r *RepositoryActivity) ActivityPeriod() int {
	if r.MinTimestamp == 0 {
		return 1
	}

	duration := float64(r.MaxTimestamp - r.MinTimestamp)
	return int(math.Ceil(math.Max((duration / 60 / 60 / 24), 1)))
}

func (r *RepositoryActivity) AverageFilesByCommits() float64 {
	return float64(r.Files / r.Commits)
}

func (r *RepositoryActivity) SetScore(score float64) {
	r.Score = score
}

type RankService struct {
	Repositories map[string]*RepositoryActivity
}

func NewRank() *RankService {
	var rank RankService
	rank.Repositories = make(map[string]*RepositoryActivity)
	return &rank
}

func (rank *RankService) LoadCsvFile() {
	fmt.Println("Load CSV FILE")
	file, err := os.Open("../commits.csv")
	if err != nil {
		log.Fatal("Error on open file", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	csvData, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	cont := 0
	for line, commit := range csvData {
		if line == 0 {
			continue
		}
		timestamp, _ := strconv.Atoi(commit[0])
		username := commit[1]
		repository := commit[2]
		files, _ := strconv.Atoi(commit[3])
		additions, _ := strconv.Atoi(commit[4])
		deletions, _ := strconv.Atoi(commit[5])

		_, ok := rank.Repositories[repository]
		if !ok {
			rank.Repositories[repository] = NewRepositoryActivity(repository)
		}
		rank.Repositories[repository].Commits += 1
		rank.Repositories[repository].appendUser(username)
		rank.Repositories[repository].Files += files
		rank.Repositories[repository].Additions += additions
		rank.Repositories[repository].Deletions += deletions

		if rank.Repositories[repository].MinTimestamp == 0 || rank.Repositories[repository].MinTimestamp > timestamp {
			rank.Repositories[repository].MinTimestamp = timestamp
		}

		if rank.Repositories[repository].MaxTimestamp == 0 || rank.Repositories[repository].MaxTimestamp < timestamp {
			rank.Repositories[repository].MaxTimestamp = timestamp
		}

		if cont == 10 {
			break
		}
		cont++

		// rank.Repositories[repository] = activity
	}
}

func (rank *RankService) CalcRankScore() {
	fmt.Println("Calc RankScore")
	maxPeriod := 0
	for _, data := range rank.Repositories {
		if maxPeriod < data.ActivityPeriod() {
			maxPeriod = data.ActivityPeriod()
		}
	}

	for _, data := range rank.Repositories {
		data.SetScore(float64(data.ActivityPeriod() / maxPeriod * data.Commits))
		rank.Repositories[data.Repository] = data
	}
	fmt.Println(rank.Repositories)
}
func (rank *RankService) GetTopActiveRepositories() {
	fmt.Println("Get Top Active Repositories")

}

func main() {
	fmt.Println("Starting")

	rankSrv := NewRank()
	rankSrv.LoadCsvFile()
	rankSrv.CalcRankScore()
}
