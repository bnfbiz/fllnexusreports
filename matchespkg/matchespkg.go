package matchespkg

import (
	"regexp"
	"sort"
)

type Match struct {
	Table string
	Time  string
}

type Matches struct {
	Count       int
	HasPractice bool
	Practice    Match
	CompMatches []Match
}

type WildCardMatch struct {
	WildCardTable string
	ActiveTable   string
	Time          string
}

func GetTeamMatchTable(teams []map[string]string, matches []map[string]string, matchKeys map[string]int) map[string]Matches {
	teamMatchTable := map[string]Matches{}
	tableOrder := []string{}
	for r := range matchKeys {
		if (r != "time") && (r != "notes") {
			tableOrder = append(tableOrder, r)
		}
	}
	sort.Slice(tableOrder, func(i, j int) bool {
		return matchKeys[tableOrder[i]] < matchKeys[tableOrder[j]]
	})
	for _, team := range teams {
		teamnumber := team["teamnumber"]
		teamMatchTable[teamnumber] = Matches{}
		for _, match := range matches {
			found := false
			tableFound := ""
			for _, table := range tableOrder {
				if match[table] == teamnumber {
					found = true
					tableFound = table
					break
				}
			}
			if found {
				re := regexp.MustCompile(`\s+`) // Matches one or more whitespace characters
				time := re.ReplaceAllString(match["time"], "")
				if match["practice"] == "TRUE" {
					teamMatchTable[teamnumber] = Matches{
						HasPractice: true,
						Practice: Match{
							Table: tableFound,
							Time:  time,
						},
					}
				} else {
					round := teamMatchTable[teamnumber].Count
					teamMatchTable[teamnumber] = Matches{
						HasPractice: teamMatchTable[teamnumber].HasPractice,
						Practice:    teamMatchTable[teamnumber].Practice,
						Count:       round + 1,
						CompMatches: append(
							teamMatchTable[teamnumber].CompMatches,
							Match{
								Table: tableFound,
								Time:  time,
							},
						),
					}
				}
			}
		}
	}
	return teamMatchTable
}

func MergeMaps(dst, src map[string]string) map[string]string {
	if dst == nil {
		dst = make(map[string]string, len(src))
	}
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func HasPracticeMatches(teamMatches map[string]Matches) bool {
	for _, m := range teamMatches {
		if m.HasPractice {
			return true
		}
	}
	return false
}

func GetMaxCompetitionMatches(teamMatches map[string]Matches) int {
	max := 0
	for _, m := range teamMatches {
		if m.Count > max {
			max = m.Count
		}
	}
	return max
}

func FindWildCardMatches(matches []map[string]string, headers map[string]int) []WildCardMatch {
	// determine table pairs
	tablePairs := map[int]string{}
	for header, idx := range headers {
		if idx > 1 {
			tablePairs[idx-2] = header
		}
	}
	wildcardMatches := []WildCardMatch{}
	for _, match := range matches {
		for i := 0; i < len(tablePairs); i += 2 {
			wildcardtable := ""
			activetable := ""
			if match[tablePairs[i]] == "" || match[tablePairs[i+1]] == "" {
				if match[tablePairs[i]] != "" {
					// the other table is the empty one
					wildcardtable = tablePairs[i+1]
					activetable = tablePairs[i]
				}
				if match[tablePairs[i+1]] != "" {
					// the other table is the empty one
					wildcardtable = tablePairs[i]
					activetable = tablePairs[i+1]
				}
				if wildcardtable != "" {
					wildcardMatches = append(wildcardMatches, WildCardMatch{
						WildCardTable: wildcardtable,
						ActiveTable:   activetable,
						Time:          match["time"],
					})
				}
			}
		}
	}
	return wildcardMatches
}

func IsWildCardMatchOnOtherTable(wildcardMatches []WildCardMatch, time string, table string) bool {
	for _, wm := range wildcardMatches {
		// Normalize from "09:30 AM" to "09:30AM" to parse correctly
		re := regexp.MustCompile(`\s+`) // Matches one or more whitespace characters
		time1 := re.ReplaceAllString(wm.Time, "")
		time2 := re.ReplaceAllString(time, "")
		if time1 == time2 && wm.ActiveTable == table {
			return true
		}
	}
	return false
}

func IsWildCardMatchOnTable(wildcardMatches []WildCardMatch, time string, table string) bool {
	for _, wm := range wildcardMatches {
		// Normalize from "09:30 AM" to "09:30AM" to parse correctly
		re := regexp.MustCompile(`\s+`) // Matches one or more whitespace characters
		time1 := re.ReplaceAllString(wm.Time, "")
		time2 := re.ReplaceAllString(time, "")
		if time1 == time2 && wm.WildCardTable == table {
			return true
		}
	}
	return false
}

func GetWildCardMatchOnOtherTable(wildcardMatches []WildCardMatch, time string, table string) WildCardMatch {
	for _, wm := range wildcardMatches {
		// Normalize from "09:30 AM" to "09:30AM" to parse correctly
		re := regexp.MustCompile(`\s+`) // Matches one or more whitespace characters
		time1 := re.ReplaceAllString(wm.Time, "")
		time2 := re.ReplaceAllString(time, "")
		if time1 == time2 && wm.ActiveTable == table {
			return wm
		}
	}
	return WildCardMatch{}
}

func GetWildCardMatchOnTable(wildcardMatches []WildCardMatch, time string, table string) WildCardMatch {
	for _, wm := range wildcardMatches {
		// Normalize from "09:30 AM" to "09:30AM" to parse correctly
		re := regexp.MustCompile(`\s+`) // Matches one or more whitespace characters
		time1 := re.ReplaceAllString(wm.Time, "")
		time2 := re.ReplaceAllString(time, "")
		if time1 == time2 && wm.WildCardTable == table {
			return wm
		}
	}
	return WildCardMatch{}
}
