import csv
import math


# Model
class RepositoryActivity:
    repository: str
    commits: int
    users: list[str]
    files: int
    additions: int
    deletions: int
    minTimestamp: int | None
    maxTimestamp: int | None
    score: int

    def __init__(self, repository: str):
        self.repository = repository
        self.commits = 0
        self.users = []
        self.files = 0
        self.additions = 0
        self.deletions = 0
        self.minTimestamp = None
        self.maxTimestamp = None

    def __str__(self):
        return f'{self.score}\t{self.commits}\t{self.activityPeriod}\t{self.files}\t{self.repository}\t\t'

    @property
    def activityPeriod(self) -> int:
        return math.ceil(max((
                (
                    self.maxTimestamp - self.minTimestamp
                ) / 60 / 60 / 24
            ), 1))

    @property
    def averageFilesByCommits(self) -> float:
        return self.files / self.commits


# Service
class RankService:
    repositories: list[RepositoryActivity]
    rankScore: dict[str: float]

    def loadCsvFile(self):
        csvContent = open('commits.csv', 'r')
        data = list(csv.reader(csvContent, delimiter=','))
        activities: dict[str, RepositoryActivity] = {}

        for commit in data[1:]:
            timestamp, username, repository, files, additions, deletions = commit
            activity = activities.get(repository, RepositoryActivity(repository=repository))
            activity.commits += 1
            if username not in activity.users:
                activity.users.append(username)
            activity.files += int(files)
            activity.additions += int(additions)
            activity.deletions += int(deletions)

            if not activity.minTimestamp or int(timestamp) < activity.minTimestamp:
                activity.minTimestamp = int(timestamp)

            if not activity.maxTimestamp or int(timestamp) > activity.minTimestamp:
                activity.maxTimestamp = int(timestamp)

            activities[repository] = activity

        self.repositories = list(activities.values())

    def calcRankScore(self):
        maxPeriod = max([x.activityPeriod for x in self.repositories])
        for rep in self.repositories:
            rep.score = round(rep.activityPeriod / maxPeriod * rep.commits, 2)

    def getTopActiveRepositories(self, quantity: int = 10) -> list:
        return sorted(
            self.repositories,
            key=lambda rep: rep.score,
            reverse=True
        )[:quantity]


def main():
    rankSrv = RankService()
    rankSrv.loadCsvFile()
    rankSrv.calcRankScore()
    [print(x) for x in rankSrv.getTopActiveRepositories()]


if __name__ == '__main__':
    main()
