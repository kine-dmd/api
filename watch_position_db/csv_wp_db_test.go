package watch_position_db

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
)

func makeCSVFile() {
	// Make a new file
	file, _ := os.Create("watch_positions.csv")

	// Write the data to the file
	origUUIDs, origPatientIds, origLimbs := makeFakePositionData()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	for i := range origUUIDs {
		_ = writer.Write([]string{origUUIDs[i],
			origPatientIds[i],
			strconv.Itoa(int(origLimbs[i]))})
	}

	_ = file.Close()
}

func removeCSVFile() {
	_ = os.Remove("watch_positions.csv")
}

func TestScanCSVDatabase(t *testing.T) {
	// Make a CSV file for testing and delete it after test completion
	makeCSVFile()
	defer removeCSVFile()

	// Make the table and take a scan of it
	csvDB := MakeCSVWatchPositionDB()
	scan := csvDB.GetTableScan()

	// Use the standard compare results function
	compareResults(scan, t)
}
