package domain

type Batch struct {
	Cities []*City
}

const maxBatchSize = 50

func GenerateBatches(cities []*City) []*Batch {
	var batches []*Batch
	var currentBatch []*City

	for i, city := range cities {
		currentBatch = append(currentBatch, city)

		if len(currentBatch) == maxBatchSize || i == len(cities) - 1 {
			batches = append(batches, &Batch{Cities: currentBatch})
			currentBatch = nil
		}
	}
	return batches
}