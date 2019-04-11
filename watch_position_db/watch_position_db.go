package watch_position_db

type WatchPosition struct {
	PatientID string `json:"PatientID"`
	Limb      uint8  `json:"Limb"`
}

// The representation we display to others
type WatchPositionDatabase interface {
	GetWatchPosition(uuid string) (WatchPosition, bool)
	GetTableScan() map[string]WatchPosition
}
