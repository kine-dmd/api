package watch_position_db

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

type CSVWatchPositionDB struct {
	filename string
}

func MakeCSVWatchPositionDB() *CSVWatchPositionDB {
	const filename string = "watch_positions.csv"

	// Try and open the file to make sure it exists and we have access to it
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open test file: %s", err)
	}
	_ = file.Close()

	// Create a new database struct object
	csvDB := new(CSVWatchPositionDB)
	csvDB.filename = filename
	return csvDB
}

func (csvDB *CSVWatchPositionDB) GetTableScan() map[string]WatchPosition {
	// Open the CSV file
	csvFile, err := os.Open(csvDB.filename)
	if err != nil {
		log.Printf("Error opening CSV file: %s", err)
		return nil
	}

	// Create a CSV reader and read the lines
	reader := csv.NewReader(csvFile)

	// Parse the rows from the
	parsedRows := make(map[string]WatchPosition)
	for {
		// Try and read a line from the CSV file
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error reading row in CSV watch file: %s", err)
			continue
		}

		// Add the data into a watch position with key equal to the uuid
		parsedRows[line[0]] = WatchPosition{
			line[1], stringToUint8(line[2]),
		}
	}

	return parsedRows
}

func (csvDB *CSVWatchPositionDB) GetWatchPosition(uuid string) (WatchPosition, bool) {
	val, ok := csvDB.GetTableScan()[uuid]
	return val, ok
}

func stringToUint8(limb string) uint8 {
	num, err := strconv.Atoi(limb)
	if err != nil {
		log.Printf("Error converting the CSV limb position to uint: %s", err)
	}
	return uint8(num)
}
