package watch_position_db

import "testing"

// Generic helper functions that may be used by more than one implementation of the watch position DB during testing

func makeFakePositionData() ([]string, []string, []uint8) {
	// Fake data in use
	uuids := []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000002"}
	patiendIds := []string{"dmd01", "dmd01", "dmd02"}
	limbs := []uint8{1, 2, 1}
	return uuids, patiendIds, limbs
}

func compareResults(scan map[string]WatchPosition, t *testing.T) {
	origUUIDs, origPatientIds, origLimbs := makeFakePositionData()

	// Check each row in the original data still exists
	for i := 0; i < len(origUUIDs); i++ {
		val, ok := scan[origUUIDs[i]]

		// Check that the 3 values all exist and match up to their original values
		if !ok {
			t.Errorf("UUID %s not found when should've been", origUUIDs[i])
		}
		if val.PatientID != origPatientIds[i] {
			t.Errorf("Mismatching patient identifiers. Expected %s . Got %s .", origPatientIds[i], val.PatientID)
		}
		if val.Limb != uint8(origLimbs[i]) {
			t.Errorf("Mismatching limb identifiers. Expected %d . Got %d .", origLimbs[i], val.Limb)
		}

	}
}
