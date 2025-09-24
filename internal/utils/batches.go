package utils

func CreateBatches(recordsLength, batchSize int) int {

	expectedBatches := recordsLength / batchSize

	if recordsLength%batchSize != 0 {
		expectedBatches++
	}

	return expectedBatches
}
