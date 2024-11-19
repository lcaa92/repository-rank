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

type RankService struct {
	Repositories map[string]*RepositoryActivity
}

func NewRank() *RankService {
	var rank RankService
	rank.Repositories = make(map[string]*RepositoryActivity)
	return &rank
}

func (rank *RankService) GetOrNewRespositoryActivity(repository string) *RepositoryActivity {
	activity, ok := rank.Repositories[repository]
	if ok {
		return activity
	}
	rank.Repositories[repository] = &RepositoryActivity{
		Repository: repository,
		Score:      0.00,
	}
	return rank.Repositories[repository]
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

		activity := rank.GetOrNewRespositoryActivity(repository)
		activity.Commits += 1
		activity.appendUser(username)
		activity.Files += files
		activity.Additions += additions
		activity.Deletions += deletions

		if activity.MinTimestamp == 0 || activity.MinTimestamp > timestamp {
			activity.MinTimestamp = timestamp
		}

		if activity.MaxTimestamp == 0 || activity.MaxTimestamp < timestamp {
			activity.MaxTimestamp = timestamp
		}
	}
	fmt.Printf("Total repositorios: %d\n", len(rank.Repositories))
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
		data.Score = float64(data.ActivityPeriod()) / float64(maxPeriod) * float64(data.Commits)
		fmt.Printf("%v - %d - %d - %d\n", data, data.ActivityPeriod(), maxPeriod, data.Commits)

	}
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
