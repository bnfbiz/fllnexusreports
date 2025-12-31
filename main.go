package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"fllnexusreports/csvmap"
	"fllnexusreports/excel"
)

func printSample(title string, rows []map[string]string, keys map[string]int, n int) {
	fmt.Printf("--- %s (%d rows) ---\n", title, len(rows))
	end := n
	if end > len(rows) {
		end = len(rows)
	}
	for i := 0; i < end; i++ {
		r := rows[i]
		// print keys
		parts := []string{}
		for i, k := range keys {
			fmt.Printf("i=%v, k=%v\n", i, k)
			parts = append(parts, fmt.Sprintf("%s=%s", i, r[i]))
		}
		fmt.Println(strings.Join(parts, ", "))
		_ = parts
		_ = r
	}
}

func main() {
	var spreadsheetName string
	var teamcsv string
	var judgingschedulecsv string
	var gameschedulecsv string

	flag.StringVar(&teamcsv, "teamCSVFilename", "teams.csv", "Name of the team CSV file")
	flag.StringVar(&judgingschedulecsv, "judgingScheduleCSVFilename", "judging_schedule.csv", "Name of the judging schedule CSV file")
	flag.StringVar(&gameschedulecsv, "gameScheduleCSVFilename", "game_schedule.csv", "Name of the game schedule CSV file")
	flag.StringVar(&spreadsheetName, "spreadsheet", "", "If you want a spreadsheet with reports provide the name of the spreadsheet to create")
	flag.Parse()

	teams, team_keys, _, err := csvmap.ReadCSVToMaps(teamcsv)
	if err != nil {
		log.Fatalf("teams read error: %v", err)
	}
	printSample("Teams", teams, team_keys, 5)

	judging, judging_keys, judging_headers, err := csvmap.ReadCSVToMaps(judgingschedulecsv)
	if err != nil {
		log.Fatalf("judging read error: %v", err)
	}
	printSample("Judging", judging, judging_keys, 5)

	games, game_keys, game_headers, err := csvmap.ReadCSVToMaps(gameschedulecsv)
	if err != nil {
		log.Fatalf("game schedule read error: %v", err)
	}
	printSample("Games", games, game_keys, 5)

	if spreadsheetName != "" {
		e := excel.CreateExcelFile()
		if err != nil {
			log.Fatalf("Failed to create Excel file: %v", err)
		}

		// Create the sheets
		_, err := e.CreateScheduleSheet("Schedule", teams, team_keys, judging, judging_keys, judging_headers, games, game_keys, game_headers)
		if err != nil {
			log.Fatalf("Failed to create schedule sheet: %v", err)
		}

		_, err = e.CreateJudgeQueueSheet("JudgeQueueing", teams, team_keys, judging, judging_keys, judging_headers, games, game_keys, game_headers)
		if err != nil {
			log.Fatalf("Failed to create judge queueing sheet: %v", err)
		}

		_, err = e.CreateMatchQueueSheet("MatchQueueing", teams, team_keys, judging, judging_keys, judging_headers, games, game_keys, game_headers)
		if err != nil {
			log.Fatalf("Failed to create Match queueing sheet: %v", err)
		}

		_, err = e.CreateEmceeSheet("EmceeReport", teams, team_keys, judging, judging_keys, judging_headers, games, game_keys, game_headers)
		if err != nil {
			log.Fatalf("Failed to create Match queueing sheet: %v", err)
		}

		fmt.Printf("Writing the spreadsheet: %s\n", spreadsheetName)
		if err := e.SaveExcelFile(spreadsheetName); err != nil {
			log.Fatalf("Failed to save Excel file: %v", err)
		}
		log.Println("Excel file saved successfully")
	} else {
		fmt.Println("No spreadsheet name provided, skipping Excel file creation")
	}
}
