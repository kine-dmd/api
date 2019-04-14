package apple_watch_3

import "github.com/kine-dmd/api/binary_file_appender"

type localFileDataWriter struct {
	fileManager binary_file_appender.BinaryFileManager
}

func MakeStandardLocalFileDataWriter() *localFileDataWriter {
	return MakeLocalFileDataWriter(binary_file_appender.MakeStandardBinaryFileManager())
}

func MakeLocalFileDataWriter(manager binary_file_appender.BinaryFileManager) *localFileDataWriter {
	writer := new(localFileDataWriter)
	writer.fileManager = manager
	return writer
}

func (writer localFileDataWriter) writeData(data UnparsedAppleWatch3Data) error {
	// Settings for the storage of data
	const basePath string = "~/data/"
	limbPositions := []string{"rightHand", "leftHand", "rightLeg", "leftLeg"}

	// Generate the file path
	filePath := basePath + data.WatchPosition.PatientID + "/" + limbPositions[data.WatchPosition.Limb] + ".bin"

	// Write the data to file using the concurrent file manager
	return writer.fileManager.AppendToFile(filePath, data.RawData)
}
