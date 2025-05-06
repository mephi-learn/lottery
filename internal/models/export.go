package models

type DrawExport struct {
	DrawID             int64
	LotteryType        string
	WinningCombination []int
	WinnerCount        int
}

type DrawExportResults struct {
	Draws []*DrawExportResult `json:"draws"`
}

type DrawExportResult struct {
	DrawId         int              `json:"draw_id"`
	WinCombination []int            `json:"win_combination"`
	Statistic      map[string]int   `json:"statistic"`
	Tickets        map[string][]int `json:"tickets"`
}
