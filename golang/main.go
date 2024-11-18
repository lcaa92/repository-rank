package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
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
	return float64(r.Files) / float64(r.Commits)
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
}

func (rank *RankService) CalcRankScore() {
	maxPeriod := 0
	for _, data := range rank.Repositories {
		if maxPeriod < data.ActivityPeriod() {
			maxPeriod = data.ActivityPeriod()
		}
	}
	fmt.Printf("MaxPeriod: %d\n", maxPeriod)
	for _, data := range rank.Repositories {
		data.Score = float64(data.ActivityPeriod()) / float64(maxPeriod) * float64(data.Commits)
	}
}

func (rank *RankService) GetTopActiveRepositories() {
	items := make([]string, 0, len(rank.Repositories))
	for repository := range rank.Repositories {
		items = append(items, repository)
	}

	sort.SliceStable(items, func(i, j int) bool {
		return rank.Repositories[items[i]].Score > rank.Repositories[items[j]].Score
	})

	for _, repository := range items[0:9] {
		fmt.Printf(
			"%s - %f - %f - %d\n",
			rank.Repositories[repository].Repository,
			rank.Repositories[repository].Score,
			rank.Repositories[repository].AverageFilesByCommits(),
			rank.Repositories[repository].Commits,
		)
	}
}

func main() {
	rankSrv := NewRank()
	rankSrv.LoadCsvFile()
	rankSrv.CalcRankScore()
	rankSrv.GetTopActiveRepositories()
}
