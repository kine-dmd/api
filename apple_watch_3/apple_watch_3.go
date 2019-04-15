package apple_watch_3

import "github.com/kine-dmd/api/watch_position_db"

type UnparsedAppleWatch3Data struct {
	WatchPosition watch_position_db.WatchPosition `json:"WatchPosition"`
	RawData       []byte                          `json:"RawData"`
}

type Aw3DataWriter interface {
	writeData(data UnparsedAppleWatch3Data) error
}

const ROW_SIZE_BYTES = 88
