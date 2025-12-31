package excel

import (
	"fmt"
	"regexp"
	"sort"
	"time"

	"fllnexusreports/matchespkg"

	"github.com/xuri/excelize/v2"
)

type excelInfo struct {
	excelFile            *excelize.File
	formatingCreated     bool
	fontSize             float64
	fontSizeLarge        float64
	fontFamily           string
	styleBold            int
	styleBoldCenterLarge int
	styleBoldBorder      int
	styleYellow          int
	styleBorderFull      int
	styleBorderLeft      int
	styleGreyBar         int
	styleGreyBarBorder   int
	styleNumberCenter    int
	styleDate            int
	// Schedule references
	SCHEDULE_TEAM_NUMREF        string
	SCHEDULE_TEAM_NAMEREF       string
	SCHEDULE_COACHES_MEETINGREF string
	SCHEDULE_JUDGING_STARTREF   string
	SCHEDULE_JUDGING_COLOR_REF  string
	SCHEDULE_ROUND_REF_START    string

	// Judging references
	JUDGING_TEAM_NUMREF     string
	JUDGING_TEAM_NAMEREF    string
	JUDGING_START_TIME_REF  string
	JUDGING_ROOM_REF        string
	JUDGING_COLUMNBREAK_REF string
	// Match references
	MATCH_TEAM_NUMREF     string
	MATCH_TEAM_NAMEREF    string
	MATCH_ROUND_REF       string
	MATCH_START_REF       string
	MATCH_TABLE_REF       string
	MATCH_COLUMNBREAK_REF string
	// Emcee references
	EMCEE_TEAM_NUMREF     string
	EMCEE_TEAM_NAMEREF    string
	EMCEE_ROOM_REF        string
	EMCEE_ROUND_REF       string
	EMCEE_START_REF       string
	EMCEE_TABLE_REF       string
	EMCEE_COLUMNBREAK_REF string
}

var (
	dateFmt                   = "mm/dd/yyyy"
	letterSize                = 1 // Letter
	legalSize                 = 5 // Legal
	orientationLandscape      = "landscape"
	orientationPortrait       = "portrait"
	emceePrintScale      uint = 100 // percent
	judgingPrintScale    uint = 125 // percent
	matchPrintScale      uint = 100 // percent
	schedulePrintScale   uint = 77  // percent
)

func CreateExcelFile() excelInfo {
	e := excelInfo{}
	e.excelFile = excelize.NewFile()
	if e.excelFile == nil {
		panic("failed to create excel file")
	}

	fmt.Println("Excel file created successfully")
	if !e.formatingCreated {
		if err := e.CreateFormatting(); err != nil {
			panic(fmt.Sprintf("failed to create formatting: %v", err))
		}
		fmt.Println("Formatting created successfully")
	} else {
		fmt.Println("Formatting not created, using default styles")
	}
	return e
}

func (e *excelInfo) SaveExcelFile(filePath string) error {
	if e.excelFile.DeleteSheet("Sheet1") != nil { // Remove default sheet
		fmt.Println("failed to delete default sheet")
	}

	if e.excelFile == nil {
		return fmt.Errorf("excel file is not initialized")
	}
	if err := e.excelFile.SaveAs(filePath); err != nil {
		return fmt.Errorf("failed to save excel file: %w", err)
	}
	return nil
}

func (e *excelInfo) CreateFormatting() error {
	if e.excelFile == nil {
		return fmt.Errorf("excel file is not initialized")
	}

	e.fontSize = 11.0
	e.fontSizeLarge = 16.0
	e.fontFamily = "Calibri"
	// Bold
	style, err := e.excelFile.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold:   true,
			Italic: false,
			Family: e.fontFamily,
			Size:   e.fontSize,
			Color:  "#000000",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to bold create style: %w", err)
	}
	e.styleBold = style

	// Bold Center
	style, err = e.excelFile.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "",
		},
		Font: &excelize.Font{
			Bold:   true,
			Italic: false,
			Family: e.fontFamily,
			Size:   e.fontSizeLarge,
			Color:  "#000000",
		},
	})

	if err != nil {
		return fmt.Errorf("failed to bold create style: %w", err)
	}
	e.styleBoldCenterLarge = style

	// Bold Border
	style, err = e.excelFile.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "",
		},
		Font: &excelize.Font{
			Bold:   true,
			Italic: false,
			Family: e.fontFamily,
			Size:   e.fontSize,
			Color:  "#000000",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 2},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to bold create style: %w", err)
	}
	e.styleBoldBorder = style

	// yellow background
	style, err = e.excelFile.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FFFF00"},
			Pattern: 1,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to yellow create style: %w", err)
	}
	e.styleYellow = style

	// Left Border
	style, err = e.excelFile.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "#000000", Style: 2},
			{Type: "left", Color: "#000000", Style: 2},
			{Type: "right", Color: "#000000", Style: 2},
			{Type: "bottom", Color: "#000000", Style: 2},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create border: %w", err)
	}
	e.styleBorderFull = style

	// Left Border
	style, err = e.excelFile.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 2},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create border: %w", err)
	}
	e.styleBorderLeft = style

	// Gray Bar
	style, err = e.excelFile.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#BFBFBF"},
			Pattern: 1,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create border: %w", err)
	}
	e.styleGreyBar = style

	// Gray Bar with Border
	style, err = e.excelFile.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#BFBFBF"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left",
				Color: "#000000",
				Style: 2,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create border: %w", err)
	}
	e.styleGreyBarBorder = style

	// Number Center
	style, err = e.excelFile.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "",
		},
		Font: &excelize.Font{
			Bold:   false,
			Italic: false,
			Family: e.fontFamily,
			Size:   e.fontSize,
			Color:  "#000000",
		},
		Border: []excelize.Border{
			{Type: "top", Color: "#000000", Style: 2},
			{Type: "left", Color: "#000000", Style: 2},
			{Type: "right", Color: "#000000", Style: 2},
			{Type: "bottom", Color: "#000000", Style: 2},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to bold create style: %w", err)
	}
	e.styleNumberCenter = style

	// Number Center
	style, err = e.excelFile.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   false,
			Italic: false,
			Family: e.fontFamily,
			Size:   e.fontSize,
			Color:  "#000000",
		},
		CustomNumFmt: &dateFmt,
	})
	if err != nil {
		return fmt.Errorf("failed to bold create style: %w", err)
	}
	e.styleDate = style

	e.formatingCreated = true
	return nil
}

func (e *excelInfo) GetMaxColumnWidth(sheetName string, column string, startRow int, rows int) int {
	maxWidth := 0
	for row := startRow; row <= rows+1; row++ { // +1 for header
		cell, err := e.excelFile.GetCellValue(sheetName, fmt.Sprintf("%s%d", column, row))
		if err != nil {
			fmt.Printf("failed to get cell value at %s%d: %v\n", column, row, err)
			continue
		}
		if len(cell) > maxWidth {
			maxWidth = len(cell)
		}
	}
	return maxWidth
}

func (e *excelInfo) CreateScheduleSheet(sheetName string,
	teams []map[string]string, team_keys map[string]int,
	judging []map[string]string, judging_keys map[string]int, judging_headers map[string]string,
	matches []map[string]string, match_keys map[string]int, match_headers map[string]string) (bool, error) {
	fmt.Printf("Creating sheet: %s\n", sheetName)
	if e.excelFile == nil {
		return false, fmt.Errorf("excel file is not initialized")
	}
	viewIndex, err := e.excelFile.NewSheet(sheetName)
	if err != nil {
		return false, fmt.Errorf("failed to create new sheet: %w", err)
	}

	// sheet level formatting
	// Hide grid lines for Sheet1
	e.excelFile.SetSheetView(sheetName, viewIndex, &excelize.ViewOptions{
		ShowGridLines: &[]bool{false}[0]})
	row := 1

	// Write the headers
	e.SCHEDULE_TEAM_NUMREF = "A%d"
	e.SCHEDULE_TEAM_NAMEREF = "B%d"
	e.SCHEDULE_JUDGING_STARTREF = "C%d"
	e.SCHEDULE_JUDGING_COLOR_REF = "D%d"
	e.SCHEDULE_ROUND_REF_START = "E%d"
	colref, _, err := excelize.CellNameToCoordinates(fmt.Sprintf(e.SCHEDULE_ROUND_REF_START, 1))
	if err != nil {
		return false, fmt.Errorf("failed to get column ref: %w", err)
	}
	lastCol := colref

	// Get a map of teams and matches
	team_match_table := matchespkg.GetTeamMatchTable(teams, matches, match_keys)
	practiceMatchesExist := matchespkg.HasPracticeMatches(team_match_table)
	compMatches := matchespkg.GetMaxCompetitionMatches(team_match_table)
	fmt.Printf("Practice matches exist: %v, comp matches: %d\n", practiceMatchesExist, compMatches)

	e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.SCHEDULE_TEAM_NUMREF, row), "Team #")
	e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.SCHEDULE_TEAM_NAMEREF, row), "Team Name")
	e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.SCHEDULE_JUDGING_STARTREF, row), "Judging Start")
	e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.SCHEDULE_JUDGING_COLOR_REF, row), "Judging Color")
	practiceOffset := 0
	if practiceMatchesExist {
		tableRef, err := excelize.CoordinatesToCellName(colref, row)
		if err != nil {
			return false, fmt.Errorf("failed to get practice table ref: %w", err)
		}
		e.excelFile.SetCellValue(sheetName, tableRef, "Practice Time")
		timeRef, err := excelize.CoordinatesToCellName(colref+1, row)
		if err != nil {
			return false, fmt.Errorf("failed to get practice time ref: %w", err)
		}
		e.excelFile.SetCellValue(sheetName, timeRef, "Practice Table")
		practiceOffset = 2
		lastCol = colref + 1
	}
	for i := 1; i <= compMatches; i++ {
		tableRef, err := excelize.CoordinatesToCellName(colref+practiceOffset+(i-1)*2, row)
		if err != nil {
			return false, fmt.Errorf("failed to get table offset ref: %w", err)
		}
		e.excelFile.SetCellValue(sheetName, tableRef, fmt.Sprintf("Round %d Time", i))
		timeRef, err := excelize.CoordinatesToCellName(colref+1+practiceOffset+(i-1)*2, row)
		if err != nil {
			return false, fmt.Errorf("failed to get time offset ref: %w", err)
		}
		e.excelFile.SetCellValue(sheetName, timeRef, fmt.Sprintf("Round %d Table", i))
		// Keep track of the last column used for the formatting later
		lastCol = colref + 1 + practiceOffset + (i-1)*2
	}

	row++

	sort.Slice(judging, func(i, j int) bool {
		// Normalize from "09:30 AM" to "09:30AM" to parse correctly
		re := regexp.MustCompile(`\s+`) // Matches one or more whitespace characters
		time1 := re.ReplaceAllString(judging[i]["time"], "")
		time2 := re.ReplaceAllString(judging[j]["time"], "")
		itime, err := time.Parse(time.Kitchen, time1)
		if err != nil {
			panic("Invalid judging start time 1 from csv")
		}
		jtime, err := time.Parse(time.Kitchen, time2)
		if err != nil {
			panic("Invalid judging start time 2 from csv")
		}
		return itime.Before(jtime)
	})
	room_order := []string{}
	for r := range judging_keys {
		if (r != "time") && (r != "notes") {
			room_order = append(room_order, r)
		}
	}
	sort.Slice(room_order, func(i, j int) bool {
		return judging_keys[room_order[i]] < judging_keys[room_order[j]]
	})

	for _, judgingSlot := range judging {
		for _, room := range room_order {
			match := team_match_table[judgingSlot[room]]
			if judgingSlot[room] != "" {
				e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.SCHEDULE_TEAM_NUMREF, row), judgingSlot[room])
				e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.SCHEDULE_TEAM_NAMEREF, row), getTeamNameByNumber(teams, judgingSlot[room]))
				e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.SCHEDULE_JUDGING_STARTREF, row), judgingSlot["time"])
				e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.SCHEDULE_JUDGING_COLOR_REF, row), judging_headers[room])
				if match.HasPractice {
					tableRef, err := excelize.CoordinatesToCellName(colref, row)
					if err != nil {
						return false, fmt.Errorf("failed to get practice table ref: %w", err)
					}
					e.excelFile.SetCellValue(sheetName, tableRef, match.Practice.Time)
					timeRef, err := excelize.CoordinatesToCellName(colref+1, row)
					if err != nil {
						return false, fmt.Errorf("failed to get practice time ref: %w", err)
					}
					e.excelFile.SetCellValue(sheetName, timeRef, match_headers[match.Practice.Table])
				}
				for i := 1; i <= compMatches; i++ {
					tableRef, err := excelize.CoordinatesToCellName(colref+practiceOffset+(i-1)*2, row)
					if err != nil {
						return false, fmt.Errorf("failed to get table offset ref: %w", err)
					}
					e.excelFile.SetCellValue(sheetName, tableRef, match.CompMatches[i-1].Time)
					timeRef, err := excelize.CoordinatesToCellName(colref+1+practiceOffset+(i-1)*2, row)
					if err != nil {
						return false, fmt.Errorf("failed to get time offset ref: %w", err)
					}
					e.excelFile.SetCellValue(sheetName, timeRef, match_headers[match.CompMatches[i-1].Table])
				}
				row++
			}
		}
	}

	lastColRef, err := excelize.ColumnNumberToName(lastCol)
	if err != nil {
		return false, fmt.Errorf("failed to get last column ref: %w", err)
	}

	// resize the columns appropriately
	for _, r := range "ACDEFGHIJKLM" {
		c := string(r)
		colWidth := e.GetMaxColumnWidth(sheetName, c, 1, len(teams)+1)
		if err := e.excelFile.SetColWidth(sheetName, c, c, float64(colWidth)); err != nil {
			return false, fmt.Errorf("failed to set column width: %w", err)
		}
	}
	if err = e.excelFile.SetColWidth(sheetName, "B", "B", 50); err != nil {
		return false, fmt.Errorf("failed to set column width: %w", err)
	}

	fmt.Printf("LastColRef is %s/%v\n", lastColRef, lastColRef)
	// Put the gray bars on alternating rows and borders
	for r := 1; r <= len(teams)+1; r++ {
		if r == 1 {
			// Header row
			e.excelFile.SetCellStyle(sheetName, fmt.Sprintf("A%d", r), fmt.Sprintf("%s%d", lastColRef, r), e.styleBold)
			for _, ch := range "CEGIKM" {
				c := string(ch)
				e.excelFile.SetCellStyle(sheetName, fmt.Sprintf("%s%d", c, r), fmt.Sprintf("%s%d", c, r), e.styleBoldBorder)
			}
		} else if (r % 2) == 0 {
			e.excelFile.SetCellStyle(sheetName, fmt.Sprintf("A%d", r), fmt.Sprintf("%s%d", lastColRef, r), e.styleGreyBar)
			// Put the borders dividing the types
			for _, ch := range "CEGIKM" {
				c := string(ch)
				e.excelFile.SetCellStyle(sheetName, fmt.Sprintf("%s%d", c, r), fmt.Sprintf("%s%d", c, r), e.styleGreyBarBorder)
			}
		} else {
			for _, ch := range "CEGIKM" {
				c := string(ch)
				e.excelFile.SetCellStyle(sheetName, fmt.Sprintf("%s%d", c, r), fmt.Sprintf("%s%d", c, r), e.styleBorderLeft)
			}
		}
	}
	// Setup sheet wide formatting
	// printing
	err = e.excelFile.SetPageLayout(sheetName, &excelize.PageLayoutOptions{
		Orientation: &orientationLandscape,
		Size:        &legalSize, // paper size
		AdjustTo:    &schedulePrintScale,
	})
	if err != nil {
		return false, fmt.Errorf("failed to set page layout: %w", err)
	}

	return true, nil
}

func (e *excelInfo) CreateJudgeQueueSheet(sheetName string,
	teams []map[string]string, team_keys map[string]int,
	judging []map[string]string, judging_keys map[string]int, judging_headers map[string]string,
	matches []map[string]string, match_keys map[string]int, match_headers map[string]string) (bool, error) {

	fmt.Printf("Creating sheet: %s\n", sheetName)
	if e.excelFile == nil {
		return false, fmt.Errorf("excel file is not initialized")
	}
	viewIndex, err := e.excelFile.NewSheet(sheetName)
	if err != nil {
		return false, fmt.Errorf("failed to create new sheet: %w", err)
	}

	// sheet level formatting
	// Hide grid lines for Sheet1
	e.excelFile.SetSheetView(sheetName, viewIndex, &excelize.ViewOptions{
		ShowGridLines: &[]bool{false}[0]})
	row := 1

	// Write the headers
	e.JUDGING_TEAM_NUMREF = "A%d"
	e.JUDGING_TEAM_NAMEREF = "B%d"
	e.JUDGING_START_TIME_REF = "C%d"
	e.JUDGING_ROOM_REF = "D%d"
	e.JUDGING_COLUMNBREAK_REF = "E%d"

	e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.JUDGING_TEAM_NUMREF, row), "Per Judging Room Schedule")

	e.excelFile.MergeCell(sheetName, fmt.Sprintf(e.JUDGING_TEAM_NUMREF, row), fmt.Sprintf(e.JUDGING_ROOM_REF, row))
	row++
	e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.JUDGING_TEAM_NUMREF, row), "Team #")
	e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.JUDGING_TEAM_NAMEREF, row), "Team Name")
	e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.JUDGING_START_TIME_REF, row), "Judging Start Time")
	e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.JUDGING_ROOM_REF, row), "Judging Room")
	e.excelFile.SetCellStyle(sheetName, fmt.Sprintf(e.JUDGING_TEAM_NUMREF, row-1), fmt.Sprintf(e.JUDGING_ROOM_REF, row), e.styleBold)
	row++

	sort.Slice(judging, func(i, j int) bool {
		// Normalize from "09:30 AM" to "09:30AM" to parse correctly
		re := regexp.MustCompile(`\s+`) // Matches one or more whitespace characters
		time1 := re.ReplaceAllString(judging[i]["time"], "")
		time2 := re.ReplaceAllString(judging[j]["time"], "")
		itime, err := time.Parse(time.Kitchen, time1)
		if err != nil {
			panic("Invalid judging start time 1 from csv")
		}
		jtime, err := time.Parse(time.Kitchen, time2)
		if err != nil {
			panic("Invalid judging start time 2 from csv")
		}
		return itime.Before(jtime)
	})
	room_order := []string{}
	for r := range judging_keys {
		if (r != "time") && (r != "notes") {
			room_order = append(room_order, r)
		}
	}
	sort.Slice(room_order, func(i, j int) bool {
		return judging_keys[room_order[i]] < judging_keys[room_order[j]]
	})

	// Get a map of teams and matches
	for _, room := range room_order {
		for _, judgingSlot := range judging {
			if judgingSlot[room] != "" {
				e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.JUDGING_TEAM_NUMREF, row), judgingSlot[room])
				e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.JUDGING_TEAM_NAMEREF, row), getTeamNameByNumber(teams, judgingSlot[room]))
				e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.JUDGING_START_TIME_REF, row), judgingSlot["time"])
				e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.JUDGING_ROOM_REF, row), judging_headers[room])
				row++
			}
		}
		e.excelFile.InsertPageBreak(sheetName, fmt.Sprintf(e.JUDGING_TEAM_NUMREF, row))
	}

	// Put the borders
	err = e.excelFile.SetCellStyle(sheetName, fmt.Sprintf(e.JUDGING_TEAM_NUMREF, 1), fmt.Sprintf(e.JUDGING_ROOM_REF, 1), e.styleBoldCenterLarge)
	if err != nil {
		return false, fmt.Errorf("failed to set bold style: %w", err)
	}
	err = e.excelFile.SetCellStyle(sheetName, fmt.Sprintf(e.JUDGING_TEAM_NUMREF, 3), fmt.Sprintf(e.JUDGING_TEAM_NUMREF, row-1), e.styleNumberCenter)
	if err != nil {
		return false, fmt.Errorf("failed to set bold style: %w", err)
	}
	err = e.excelFile.SetCellStyle(sheetName, fmt.Sprintf(e.JUDGING_TEAM_NAMEREF, 3), fmt.Sprintf(e.JUDGING_ROOM_REF, row-1), e.styleBorderFull)
	if err != nil {
		return false, fmt.Errorf("failed to set bold style: %w", err)
	}
	// resize the columns appropriately
	for _, r := range "ACDEFGHIJKLM" {
		c := string(r)
		colWidth := e.GetMaxColumnWidth(sheetName, c, 2, len(teams))
		if err := e.excelFile.SetColWidth(sheetName, c, c, float64(colWidth)); err != nil {
			return false, fmt.Errorf("failed to set column width: %w", err)
		}
	}
	if err = e.excelFile.SetColWidth(sheetName, "B", "B", 50); err != nil {
		return false, fmt.Errorf("failed to set column width: %w", err)
	}

	// Setup sheet wide formatting
	// printing
	err = e.excelFile.SetPageLayout(sheetName, &excelize.PageLayoutOptions{
		Orientation: &orientationLandscape,
		Size:        &letterSize, // paper size
		AdjustTo:    &judgingPrintScale,
	})
	if err != nil {
		return false, fmt.Errorf("failed to set page layout: %w", err)
	}
	err = e.excelFile.SetDefinedName(&excelize.DefinedName{
		Name:     "_xlnm.Print_Area",
		RefersTo: fmt.Sprintf("%s!$A$1:$D$%d", sheetName, row-1),
		Scope:    sheetName,
	})
	if err != nil {
		return false, fmt.Errorf("failed to set print area: %w", err)
	}
	err = e.excelFile.SetDefinedName(&excelize.DefinedName{
		Name:     "_xlnm.Print_Titles",
		RefersTo: fmt.Sprintf("%s!$1:$2", sheetName),
		Scope:    sheetName,
	})
	if err != nil {
		return false, fmt.Errorf("failed to set print titles: %w", err)
	}
	IsTrue := true
	e.excelFile.SetSheetView(sheetName, viewIndex, &excelize.ViewOptions{
		ShowGridLines: &IsTrue,
	})

	return true, nil
}

func (e *excelInfo) CreateMatchQueueSheet(sheetName string,
	teams []map[string]string, team_keys map[string]int,
	judging []map[string]string, judging_keys map[string]int, judging_headers map[string]string,
	matches []map[string]string, match_keys map[string]int, match_headers map[string]string) (bool, error) {

	fmt.Printf("Creating sheet: %s\n", sheetName)
	if e.excelFile == nil {
		return false, fmt.Errorf("excel file is not initialized")
	}
	viewIndex, err := e.excelFile.NewSheet(sheetName)
	if err != nil {
		return false, fmt.Errorf("failed to create new sheet: %w", err)
	}

	// sheet level formatting
	// Hide grid lines for Sheet1
	e.excelFile.SetSheetView(sheetName, viewIndex, &excelize.ViewOptions{
		ShowGridLines: &[]bool{false}[0]})
	row := 1

	// Write the headers
	e.MATCH_TEAM_NUMREF = "A%d"
	e.MATCH_TEAM_NAMEREF = "B%d"
	e.MATCH_ROUND_REF = "C%d"
	e.MATCH_START_REF = "D%d"
	e.MATCH_TABLE_REF = "E%d"
	e.MATCH_COLUMNBREAK_REF = "F%d"

	e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.MATCH_TEAM_NUMREF, row), "Per Table Match Schedule")
	e.excelFile.MergeCell(sheetName, fmt.Sprintf(e.MATCH_TEAM_NUMREF, row), fmt.Sprintf(e.MATCH_TABLE_REF, row))

	row++
	e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.MATCH_TEAM_NUMREF, row), "Team #")
	e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.MATCH_TEAM_NAMEREF, row), "Team Name")
	e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.MATCH_ROUND_REF, row), "Round")
	e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.MATCH_START_REF, row), "Match Start")
	e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.MATCH_TABLE_REF, row), "Table")
	e.excelFile.SetCellStyle(sheetName, fmt.Sprintf(e.MATCH_TEAM_NUMREF, row-1), fmt.Sprintf(e.MATCH_TABLE_REF, row), e.styleBold)
	row++

	sort.Slice(matches, func(i, j int) bool {
		// Normalize from "09:30 AM" to "09:30AM" to parse correctly
		re := regexp.MustCompile(`\s+`) // Matches one or more whitespace characters
		time1 := re.ReplaceAllString(matches[i]["time"], "")
		time2 := re.ReplaceAllString(matches[j]["time"], "")
		itime, err := time.Parse(time.Kitchen, time1)
		if err != nil {
			panic("Invalid judging start time 1 from csv")
		}
		jtime, err := time.Parse(time.Kitchen, time2)
		if err != nil {
			panic("Invalid judging start time 2 from csv")
		}
		return itime.Before(jtime)
	})
	table_order := []string{}
	for r := range match_keys {
		if (r != "time") && (r != "practice") && (r != "notes") {
			table_order = append(table_order, r)
		}
	}
	sort.Slice(table_order, func(i, j int) bool {
		return match_keys[table_order[i]] < match_keys[table_order[j]]
	})

	rounds := map[string]int{}
	fmt.Printf("Table order: %v\n", table_order)
	for _, table := range table_order {
		for _, match := range matches {
			if match[table] != "" {
				e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.MATCH_TEAM_NUMREF, row), match[table])
				e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.MATCH_TEAM_NAMEREF, row), getTeamNameByNumber(teams, match[table]))
				if match["practice"] == "TRUE" {
					e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.MATCH_ROUND_REF, row), "Practice")
				} else {
					rounds[match[table]]++
					e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.MATCH_ROUND_REF, row), fmt.Sprintf("Comp"))
				}
				e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.MATCH_START_REF, row), match["time"])
				e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.MATCH_TABLE_REF, row), match_headers[table])
				row++
			}
		}
		e.excelFile.InsertPageBreak(sheetName, fmt.Sprintf(e.MATCH_TEAM_NUMREF, row))
	}

	// Put the borders
	err = e.excelFile.SetCellStyle(sheetName, fmt.Sprintf(e.MATCH_TEAM_NUMREF, 1), fmt.Sprintf(e.MATCH_TEAM_NUMREF, 1), e.styleBoldCenterLarge)
	if err != nil {
		return false, fmt.Errorf("failed to set bold style: %w", err)
	}
	err = e.excelFile.SetCellStyle(sheetName, fmt.Sprintf(e.MATCH_TEAM_NUMREF, 3), fmt.Sprintf(e.MATCH_TEAM_NUMREF, row-1), e.styleNumberCenter)
	if err != nil {
		return false, fmt.Errorf("failed to set bold style: %w", err)
	}
	err = e.excelFile.SetCellStyle(sheetName, fmt.Sprintf(e.MATCH_TEAM_NAMEREF, 3), fmt.Sprintf(e.MATCH_TABLE_REF, row-1), e.styleBorderFull)
	if err != nil {
		return false, fmt.Errorf("failed to set bold style: %w", err)
	}
	// resize the columns appropriately
	for _, r := range "ACDE" {
		c := string(r)
		colWidth := e.GetMaxColumnWidth(sheetName, c, 2, len(teams))
		if err := e.excelFile.SetColWidth(sheetName, c, c, float64(colWidth+1.0)); err != nil {
			return false, fmt.Errorf("failed to set column width: %w", err)
		}
	}
	if err = e.excelFile.SetColWidth(sheetName, "B", "B", 50); err != nil {
		return false, fmt.Errorf("failed to set column width: %w", err)
	}

	// Setup sheet wide formatting
	// printing
	err = e.excelFile.SetPageLayout(sheetName, &excelize.PageLayoutOptions{
		Orientation: &orientationPortrait,
		Size:        &letterSize, // paper size
		AdjustTo:    &matchPrintScale,
	})
	if err != nil {
		return false, fmt.Errorf("failed to set page layout: %w", err)
	}
	err = e.excelFile.SetDefinedName(&excelize.DefinedName{
		Name:     "_xlnm.Print_Area",
		RefersTo: fmt.Sprintf("%s!$A$1:$E$%d", sheetName, row-1),
		Scope:    sheetName,
	})
	if err != nil {
		return false, fmt.Errorf("failed to set print area: %w", err)
	}
	err = e.excelFile.SetDefinedName(&excelize.DefinedName{
		Name:     "_xlnm.Print_Titles",
		RefersTo: fmt.Sprintf("%s!$1:$2", sheetName),
		Scope:    sheetName,
	})
	if err != nil {
		return false, fmt.Errorf("failed to set print titles: %w", err)
	}
	IsTrue := true
	e.excelFile.SetSheetView(sheetName, viewIndex, &excelize.ViewOptions{
		ShowGridLines: &IsTrue,
	})

	return true, nil
}

func getTeamNameByNumber(teams []map[string]string, teamnumber string) string {
	for _, team := range teams {
		if team["teamnumber"] == teamnumber {
			return team["teamnameoptional"]
		}
	}
	return ""
}

func (e *excelInfo) CreateEmceeSheet(sheetName string,
	teams []map[string]string, team_keys map[string]int,
	judging []map[string]string, judging_keys map[string]int, judging_headers map[string]string,
	matches []map[string]string, match_keys map[string]int, match_headers map[string]string) (bool, error) {
	fmt.Printf("Creating sheet: %s\n", sheetName)
	if e.excelFile == nil {
		return false, fmt.Errorf("excel file is not initialized")
	}
	viewIndex, err := e.excelFile.NewSheet(sheetName)
	if err != nil {
		return false, fmt.Errorf("failed to create new sheet: %w", err)
	}

	// sheet level formatting
	// Hide grid lines for Sheet1
	e.excelFile.SetSheetView(sheetName, viewIndex, &excelize.ViewOptions{
		ShowGridLines: &[]bool{false}[0]})
	row := 1

	// Write the headers
	e.EMCEE_TEAM_NUMREF = "A%d"
	e.EMCEE_TEAM_NAMEREF = "B%d"
	e.EMCEE_ROUND_REF = "C%d"
	e.EMCEE_START_REF = "D%d"
	e.EMCEE_TABLE_REF = "E%d"
	e.EMCEE_COLUMNBREAK_REF = "F%d"

	e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.EMCEE_TEAM_NUMREF, row), "Emcee Sheet")
	e.excelFile.MergeCell(sheetName, fmt.Sprintf(e.EMCEE_TEAM_NUMREF, row), fmt.Sprintf(e.MATCH_TABLE_REF, row))

	row++
	e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.EMCEE_TEAM_NUMREF, row), "Team #")
	e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.EMCEE_TEAM_NAMEREF, row), "Team Name")
	e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.EMCEE_ROUND_REF, row), "Round")
	e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.EMCEE_START_REF, row), "Match Start")
	e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.EMCEE_TABLE_REF, row), "Table")
	e.excelFile.SetCellStyle(sheetName, fmt.Sprintf(e.EMCEE_TEAM_NUMREF, row-1), fmt.Sprintf(e.EMCEE_TABLE_REF, row), e.styleBold)
	row++

	sort.Slice(matches, func(i, j int) bool {
		// Normalize from "09:30 AM" to "09:30AM" to parse correctly
		re := regexp.MustCompile(`\s+`) // Matches one or more whitespace characters
		time1 := re.ReplaceAllString(matches[i]["time"], "")
		time2 := re.ReplaceAllString(matches[j]["time"], "")
		itime, err := time.Parse(time.Kitchen, time1)
		if err != nil {
			panic("Invalid judging start time 1 from csv")
		}
		jtime, err := time.Parse(time.Kitchen, time2)
		if err != nil {
			panic("Invalid judging start time 2 from csv")
		}
		return itime.Before(jtime)
	})
	table_order := []string{}
	for r := range match_keys {
		if (r != "time") && (r != "practice") && (r != "notes") {
			table_order = append(table_order, r)
		}
	}
	sort.Slice(table_order, func(i, j int) bool {
		return match_keys[table_order[i]] < match_keys[table_order[j]]
	})

	lastTime := matches[0]["time"]
	grayLine := false
	count := 0
	for _, match := range matches {
		for _, table := range table_order {
			if match[table] != "" {
				if lastTime != match["time"] && row > 3 {
					lastTime = match["time"]
					row++
					grayLine = !grayLine
					if (count % 15) == 0 {
						fmt.Println("Inserting page break at row:", row)
						e.excelFile.InsertPageBreak(sheetName, fmt.Sprintf(e.EMCEE_COLUMNBREAK_REF, row))
						grayLine = false
						count = 0
					}
				}
				e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.EMCEE_TEAM_NUMREF, row), match[table])
				e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.EMCEE_TEAM_NAMEREF, row), getTeamNameByNumber(teams, match[table]))
				if match["practice"] == "TRUE" {
					e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.EMCEE_ROUND_REF, row), "Practice")
				} else {
					e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.EMCEE_ROUND_REF, row), "Comp")
				}
				e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.MATCH_START_REF, row), match["time"])
				e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.MATCH_TABLE_REF, row), match_headers[table])
				if grayLine {
					e.excelFile.SetCellStyle(sheetName, fmt.Sprintf(e.EMCEE_TEAM_NUMREF, row), fmt.Sprintf(e.EMCEE_TABLE_REF, row), e.styleGreyBar)
				}
				row++
				count++
			}
		}
	}
	// for _, match := range matches {
	// 	if schedule.IsMatch(slotInfo.Slot) {

	// 		if teamList.IsTeamNumber(slotInfo.TeamNumber) {
	// 			row++
	// 			count++
	// 			e.excelFile.SetCellFormula(sheetName, fmt.Sprintf(e.EMCEE_TEAM_NUMREF, row), "='"+SHEETNAME_TEAM_LIST+"'!"+fmt.Sprintf(e.TEAMLIST_TEAMNUMBER_REF, slotInfo.TeamNumber+1))
	// 			e.excelFile.SetCellFormula(sheetName, fmt.Sprintf(e.EMCEE_TEAM_NAMEREF, row), "='"+SHEETNAME_TEAM_LIST+"'!"+fmt.Sprintf(e.TEAMLIST_TEAMNAME_REF, slotInfo.TeamNumber+1))
	// 			e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.EMCEE_ROUND_REF, row), schedule.GetRoundName(slotInfo.Slot))
	// 			e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.EMCEE_START_REF, row), slotInfo.Start.Format("03:04 PM"))
	// 			e.excelFile.SetCellFormula(sheetName, fmt.Sprintf(e.EMCEE_TABLE_REF, row), "='"+SHEETNAME_TOURNAMENT_SETUP+"'!"+fmt.Sprintf(GAMETABLESCOLUMNREF, schedule.GetTableOffset(slotInfo.Location)+GAMETABLESROWREF))
	// 			if grayLine {
	// 				e.excelFile.SetCellStyle(sheetName, fmt.Sprintf(e.EMCEE_TEAM_NUMREF, row), fmt.Sprintf(e.EMCEE_TABLE_REF, row), e.styleGreyBar)
	// 			}
	// 		} else if teamList.IsWildCardTeam(slotInfo.TeamNumber) {
	// 			row++
	// 			count++
	// 			// Wild Card Team so don't print the team number
	// 			e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.EMCEE_TEAM_NAMEREF, row), "Wild Card Team")
	// 			e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.EMCEE_ROUND_REF, row), schedule.GetRoundName(slotInfo.Slot))
	// 			e.excelFile.SetCellValue(sheetName, fmt.Sprintf(e.EMCEE_START_REF, row), slotInfo.Start.Format("03:04 PM"))
	// 			e.excelFile.SetCellFormula(sheetName, fmt.Sprintf(e.EMCEE_TABLE_REF, row), "='"+SHEETNAME_TOURNAMENT_SETUP+"'!"+fmt.Sprintf(GAMETABLESCOLUMNREF, schedule.GetTableOffset(slotInfo.Location)+GAMETABLESROWREF))
	// 			if grayLine {
	// 				e.excelFile.SetCellStyle(sheetName, fmt.Sprintf(e.EMCEE_TEAM_NUMREF, row), fmt.Sprintf(e.EMCEE_TABLE_REF, row), e.styleGreyBar)
	// 			}
	// 		}
	// 	}
	// }
	e.excelFile.InsertPageBreak(sheetName, fmt.Sprintf(e.EMCEE_TEAM_NUMREF, row+1))

	// Put the borders
	err = e.excelFile.SetCellStyle(sheetName, fmt.Sprintf(e.EMCEE_TEAM_NUMREF, 1), fmt.Sprintf(e.EMCEE_TEAM_NUMREF, 1), e.styleBoldCenterLarge)
	if err != nil {
		return false, fmt.Errorf("failed to set bold style: %w", err)
	}
	// resize the columns appropriately
	for _, r := range "ACDE" {
		c := string(r)
		colWidth := e.GetMaxColumnWidth(sheetName, c, 2, len(teams))
		if err := e.excelFile.SetColWidth(sheetName, c, c, float64(colWidth+1.0)); err != nil {
			return false, fmt.Errorf("failed to set column width: %w", err)
		}
	}
	if err = e.excelFile.SetColWidth(sheetName, "B", "B", 50); err != nil {
		return false, fmt.Errorf("failed to set column width: %w", err)
	}

	// Setup sheet wide formatting
	// printing
	err = e.excelFile.SetPageLayout(sheetName, &excelize.PageLayoutOptions{
		Orientation: &orientationPortrait,
		Size:        &letterSize, // paper size
		AdjustTo:    &emceePrintScale,
	})
	if err != nil {
		return false, fmt.Errorf("failed to set page layout: %w", err)
	}
	err = e.excelFile.SetDefinedName(&excelize.DefinedName{
		Name:     "_xlnm.Print_Area",
		RefersTo: fmt.Sprintf("%s!$A$1:$E$%d", sheetName, row),
		Scope:    sheetName,
	})
	if err != nil {
		return false, fmt.Errorf("failed to set print area: %w", err)
	}
	err = e.excelFile.SetDefinedName(&excelize.DefinedName{
		Name:     "_xlnm.Print_Titles",
		RefersTo: fmt.Sprintf("%s!$1:$2", sheetName),
		Scope:    sheetName,
	})
	if err != nil {
		return false, fmt.Errorf("failed to set print titles: %w", err)
	}
	IsTrue := true
	e.excelFile.SetSheetView(sheetName, viewIndex, &excelize.ViewOptions{
		ShowGridLines: &IsTrue,
	})

	return true, nil
}
