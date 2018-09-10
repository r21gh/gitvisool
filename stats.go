package main

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"time"
)

const (
	daysInLastSixMonths = 183
	outOfRange = 99999
)

// stats calculates and prints the stats
func stats(email string) {
	commits := processRepositories(email)
	printCommitsStats(commits)
}

// processRepositories given an user email, returns the
// commits made in the last 6 months
func processRepositories(email string) map[int]int {
	filePath := getDotFilePath()
	repos := parseFileLinesToSlice(filePath)
	daysInMap := daysInLastSixMonths

	commits := make(map[int]int, daysInMap)
	for i := daysInMap; i > 0; i-- {
		commits[i] = 0
	}

	for _, path := range repos {
		commits = fillCommits(email, path, commits)
	}

	return commits
}


// fillCommits given a repository found in `path`, gets the commits and
// puts them in the `commits` map, returning it when completed
func fillCommits(email string, path string, commits map[int]int) map[int]int {
	// instantiate a git repo object from path
	repo, err := git.PlainOpen(path)
	if err != nil {
		panic(err)
	}

	// get the HEAD reference
	ref, err := repo.Head()
	if err != nil {
		panic(err)
	}

	// get the commits history starting from HEAD
	iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		panic(err)
	}

	// iterate the commits
	offset := calcOffset()
	err = iterator.ForEach(func(c *object.Commit) error {
		daysAge := countDaysSinceDate(c.Author.When) + offset

		if c.Author.Email != email {
			return nil
		}

		if daysAge != outOfRange {
			commits[daysAge]++
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return commits
}

// getBeginningOfDay given a time.Time calculates the start time of that day
func getBeginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return startOfDay
}

// countDaysSinceDate counts how many days passed since the passed `data`
func countDaysSinceDate(date time.Time) int {
	days := 0
	now := getBeginningOfDay(time.Now())
	for date.Before(now) {
		date = date.Add(time.Hour * 24)
		days++
		if days > daysInLastSixMonths {
			return outOfRange
		}
	}
	return days
}

// calcOffset determines and returns the amount of days missing to fill
// the last row of the stats graph
func calcOffset() int {
	var offset int
	weekday := time.Now().Weekday()
	switch weekday {
	case time.Sunday:
		offset = 7
	case time.Monday:
		offset = 6
	case time.Tuesday:
		offset = 5
	case time.Wednesday:
		offset = 4
	case time.Thursday:
		offset = 3
	case time.Friday:
		offset = 2
	case time.Saturday:
		offset = 1
	}
	return offset
}

