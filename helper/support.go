package helper

import (
	"github.com/italanleal/go-wfcsd/models"
)

// CalcPositiveSupport recalcula o suporte positivo para cada item
// itemList: lista de itens
// output: vetor booleano indicando se a linha é positiva (true = "P")
// func CalcPositiveSupport(itemList []models.Item, output []bool) {
// 	for i := range itemList {
// 		posCount := 0
// 		for _, idx := range itemList[i].Index {
// 			if idx < len(output) && output[idx] {
// 				posCount++
// 			}
// 		}
// 		if len(itemList[i].Index) > 0 {
// 			itemList[i].SuppP = float64(posCount) / float64(len(itemList[i].Index))
// 		} else {
// 			itemList[i].SuppP = 0
// 		}
// 	}
// }

func CalcPositiveSupport(itemList []models.Item, output []bool) {
	// contar total de positivos no dataset
	totalPos := 0
	for _, v := range output {
		if v {
			totalPos++
		}
	}
	if totalPos == 0 {
		totalPos = 1 // evita divisão por zero
	}

	for i := range itemList {
		posCount := 0
		for _, idx := range itemList[i].Index {
			if idx < len(output) && output[idx] {
				posCount++
			}
		}
		itemList[i].SuppP = float64(posCount) / float64(totalPos)
	}
}

func CalcPatternStats(pattern *models.Pattern, itemList []models.Item, output []bool) {
	// total positives in dataset
	totalPos := 0
	for _, v := range output {
		if v {
			totalPos++
		}
	}
	if totalPos == 0 {
		totalPos = 1
	}

	if len(pattern.Items) == 0 {
		pattern.Freq = 0
		pattern.Wracc = 0
		return
	}

	// intersection of all item indices
	coveredIndices := make(map[int]struct{})
	for _, idx := range itemList[pattern.Items[0]].Index {
		coveredIndices[idx] = struct{}{}
	}

	for _, itemIdx := range pattern.Items[1:] {
		newCovered := make(map[int]struct{})
		for _, idx := range itemList[itemIdx].Index {
			if _, ok := coveredIndices[idx]; ok {
				newCovered[idx] = struct{}{}
			}
		}
		coveredIndices = newCovered
	}

	patternSize := len(coveredIndices)
	if patternSize == 0 {
		pattern.Freq = 0
		pattern.Wracc = 0
		return
	}

	// count positives
	posInPattern := 0
	for idx := range coveredIndices {
		if idx < len(output) && output[idx] {
			posInPattern++
		}
	}

	totalRecords := len(output)
	if totalRecords == 0 {
		pattern.Freq = 0
		pattern.Wracc = 0
		return
	}

	pattern.Freq = float64(patternSize) / float64(totalRecords)
	pattern.Wracc = pattern.Freq * (float64(posInPattern)/float64(patternSize) - float64(totalPos)/float64(totalRecords))
}

func CalcPatternsStats(patternList []models.Pattern, itemList []models.Item, output []bool) {
	for i := range patternList {
		CalcPatternStats(&patternList[i], itemList, output)
	}
}
