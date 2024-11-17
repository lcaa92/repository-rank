
# Steps Algorithm
1. Read csv file and group commits by repository (`class RepositoryActivity`).
2. Create new `RankService` with list of group repositories.

3. Calculare repositories score (`rankSrv.calcRankScore()`).

    3.1. Get maximium respository period days (`repository.activityPeriod`) between all repositories. This repository's property is calculated by the diff between last and first commit timestamp in days (minimium is 1).

    3.2 - Calculare score for each repository: `repository.activityPeriod / maxPeriod * repository.commits`

    The ideia of use maximium period days of all repositories is to use amount of days as height when calculating the score to avoid false positive activity frequency.

# Get top 10 most active repositories
The service `RankService` has the method `getTopActiveRepositories` that will return a list of repositories with better score


# How to run

It is necessary python3 to run the script: (script use only python standard library)

```python
python3 repositoryRank.py
```

# ToDo Extras:
- Improve tabular result
