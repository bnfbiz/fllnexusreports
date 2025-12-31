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
