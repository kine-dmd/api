package apple_watch_3

import "github.com/kine-dmd/api/watch_position_db"

/*****************************************
Common testing functionality only
******************************************/

func makeFakeUnparsedDataStruct(patientId string, limb uint8, rawData []byte) UnparsedAppleWatch3Data {
	watchData := UnparsedAppleWatch3Data{
		watch_position_db.WatchPosition{
			patientId, limb,
		},
		rawData,
	}
	return watchData
}
