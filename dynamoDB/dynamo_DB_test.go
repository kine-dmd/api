package dynamoDB

import (
	"strconv"
	"testing"
)

/***********************************************************************************

These test rely on a test database on AWS. The test database must be as follows:


             testTable

  pKey (string)  |  col1 (number)
-----------------|-----------------
       "1"       |        1
       "2"       |        2
                 |

***********************************************************************************/

// These tests require read permissions on AWS dynamoDB
func TestReadingFromTestDBReturnsSomething(t *testing.T) {
	// Make a client
	client := DynamoDBClient{}

	// Initialise the connection
	err := client.InitConn("testTable")
	if err != nil {
		t.Fatalf("Received unexpected error when initialising connection: %s", err)
	}

	// Scan the table
	scan := client.GetTableScan()
	if len(scan) != 2 {
		t.Fatalf("Expected 2 rows in the test table - got %d. Check test table on AWS.", len(scan))
	}

	// Check that the pKey (string) and col1 match up
	for _, row := range scan {
		numPKey, _ := strconv.Atoi(row["pKey"].(string))
		// Check the first row
		if float64(numPKey) != row["col1"].(float64) {
			t.Fatal("String interpretation of primary key and col1 do not match up. Check test table on AWS.")
		}

	}
}
