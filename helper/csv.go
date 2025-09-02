package helper

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/italanleal/go-wfcsd/models"
)

// ReadCSV preenche attrList, output e itemList a partir do CSV
// Calcula o PositiveSupport de cada item

// ReadCSV rápido usando goroutines e map auxiliar
func ReadCSV(filename, target string, attrList *[]string, output *[]bool, itemList *[]models.Item) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))

	// Ler header
	header, err := reader.Read()
	if err != nil {
		return err
	}
	*attrList = header

	targetCol := -1
	for i, name := range header {
		if strings.TrimSpace(name) == target {
			targetCol = i
			break
		}
	}
	if targetCol == -1 {
		return fmt.Errorf("target column %s not found", target)
	}

	// map temporário para acesso rápido de items
	type key struct {
		attr  string
		value string
	}
	itemMap := make(map[key]*models.Item)

	var mu sync.Mutex
	var wg sync.WaitGroup

	rowIndex := 0

	for {
		row, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}

		wg.Add(1)
		go func(row []string, rowIndex int) {
			defer wg.Done()

			val := strings.TrimSpace(row[targetCol])
			isPositive := strings.EqualFold(val, "p")

			mu.Lock()
			*output = append(*output, isPositive)
			mu.Unlock()

			for colIndex, cell := range row {
				if colIndex == targetCol {
					continue
				}

				cell = strings.TrimSpace(cell)
				k := key{attr: header[colIndex], value: cell}

				mu.Lock()
				it, exists := itemMap[k]
				if !exists {
					it = &models.Item{
						Attr:  colIndex,
						Value: cell,
						Index: []int{rowIndex},
					}
					itemMap[k] = it
				} else {
					it.Index = append(it.Index, rowIndex)
				}
				mu.Unlock()
			}
		}(row, rowIndex)

		rowIndex++
	}

	wg.Wait()

	// copiar para itemList
	for _, it := range itemMap {
		*itemList = append(*itemList, *it)
	}

	return nil
}
