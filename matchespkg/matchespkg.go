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

func GetTeamMatchTable(teams []map[string]string, matches []map[string]string, match_keys map[string]int) map[string]Matches {
	team_match_table := map[string]Matches{}
	table_order := []string{}
	for r := range match_keys {
		if (r != "time") && (r != "notes") {
			table_order = append(table_order, r)
		}
	}
	sort.Slice(table_order, func(i, j int) bool {
		return match_keys[table_order[i]] < match_keys[table_order[j]]
	})
	for _, team := range teams {
		teamnumber := team["teamnumber"]
		team_match_table[teamnumber] = Matches{}
		for _, match := range matches {
			found := false
			table_found := ""
			for _, table := range table_order {
				if match[table] == teamnumber {
					found = true
					table_found = table
					break
				}
			}
			if found {
				re := regexp.MustCompile(`\s+`) // Matches one or more whitespace characters
				time := re.ReplaceAllString(match["time"], "")
				if match["practice"] == "TRUE" {
					team_match_table[teamnumber] = Matches{
						HasPractice: true,
						Practice: Match{
							Table: table_found,
							Time:  time,
						},
					}
				} else {
					round := team_match_table[teamnumber].Count
					team_match_table[teamnumber] = Matches{
						HasPractice: team_match_table[teamnumber].HasPractice,
						Practice:    team_match_table[teamnumber].Practice,
						Count:       round + 1,
						CompMatches: append(
							team_match_table[teamnumber].CompMatches,
							Match{
								Table: table_found,
								Time:  time,
							},
						),
					}
				}
			}
		}
	}
	return team_match_table
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

func HasPracticeMatches(team_matches map[string]Matches) bool {
	for _, m := range team_matches {
		if m.HasPractice {
			return true
		}
	}
	return false
}

func GetMaxCompetitionMatches(team_matches map[string]Matches) int {
	max := 0
	for _, m := range team_matches {
		if m.Count > max {
			max = m.Count
		}
	}
	return max
}

func FindWildCardMatches(matches []map[string]string, headers map[string]int) []WildCardMatch {
	// determine table pairs
	table_pairs := map[int]string{}
	for header, idx := range headers {
		if idx > 1 {
			table_pairs[idx-2] = header
		}
	}
	wildcard_matches := []WildCardMatch{}
	for _, match := range matches {
		for i := 0; i < len(table_pairs); i += 2 {
			wildcardtable := ""
			activetable := ""
			if match[table_pairs[i]] == "" || match[table_pairs[i+1]] == "" {
				if match[table_pairs[i]] != "" {
					// the other table is the empty one
					wildcardtable = table_pairs[i+1]
					activetable = table_pairs[i]
				}
				if match[table_pairs[i+1]] != "" {
					// the other table is the empty one
					wildcardtable = table_pairs[i]
					activetable = table_pairs[i+1]
				}
				if wildcardtable != "" {
					wildcard_matches = append(wildcard_matches, WildCardMatch{
						WildCardTable: wildcardtable,
						ActiveTable:   activetable,
						Time:          match["time"],
					})
				}
			}
		}
	}
	return wildcard_matches
}

func IsWildCardMatchOnOtherTable(wildcard_matches []WildCardMatch, time string, table string) bool {
	for _, wm := range wildcard_matches {
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

func IsWildCardMatchOnTable(wildcard_matches []WildCardMatch, time string, table string) bool {
	for _, wm := range wildcard_matches {
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

func GetWildCardMatchOnOtherTable(wildcard_matches []WildCardMatch, time string, table string) WildCardMatch {
	for _, wm := range wildcard_matches {
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

func GetWildCardMatchOnTable(wildcard_matches []WildCardMatch, time string, table string) WildCardMatch {
	for _, wm := range wildcard_matches {
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
