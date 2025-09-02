package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/italanleal/go-wfcsd/helper"
	"github.com/italanleal/go-wfcsd/models"
)

var (
	attrList []string
	output   []bool
	itemList []models.Item
	tileList []models.Pattern
	tileMap  map[string][]int
)

const (
	itemFile     = ".\\bin\\itemList.gob"
	attrFile     = ".\\bin\\attrList.gob"
	tileListFile = ".\\bin\\tileList.gob"
	tileMapFile  = ".\\bin\\tileMap.gob"
)

func main() {
	// Check if saved state exists
	if helper.FileExists(itemFile) {

		fmt.Println("Loading saved state from disk...")

		var err error
		attrList, err = helper.LoadAttrList(attrFile)
		if err != nil {
			log.Fatal(err)
		}

		itemList, err = helper.LoadItemList(itemFile)
		if err != nil {
			log.Fatal(err)
		}

		tileList, err = helper.LoadTileList(tileListFile)
		if err != nil {
			log.Fatal(err)
		}

		tileMap, err = helper.LoadTileMap(tileMapFile)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("State loaded successfully.")

	} else {
		fmt.Println("Saved state not found, processing CSV...")

		// Ler CSV
		err := helper.ReadCSV("data\\christensen-pn-freq-2.CSV", "y", &attrList, &output, &itemList)
		if err != nil {
			log.Fatal(err)
		}

		// Calcular suporte positivo
		helper.CalcPositiveSupport(itemList, output)

		// Ordenar por suporte positivo
		sort.Slice(itemList, func(i, j int) bool {
			return itemList[i].SuppP > itemList[j].SuppP
		})

		// Filtrar items por threshold
		threshold := 0.05
		filtered := make([]models.Item, 0, len(itemList))
		for _, it := range itemList {
			if it.SuppP >= threshold {
				filtered = append(filtered, it)
			}
		}
		itemList = filtered

		// Inicializar tileMap antes de gerar combinações
		tileMap = make(map[string][]int)
		helper.GenerateTwoItemCombinations(&itemList, &output, &tileList, &tileMap)
		helper.CalcPatternsStats(tileList, itemList, output)
		// Sort by WRAcc descending
		sort.Slice(tileList, func(i, j int) bool {
			return tileList[i].Wracc > tileList[j].Wracc
		})

		// Salvar estado em disco
		helper.SaveAttrList(attrFile, attrList)
		helper.SaveItemList(itemFile, itemList)
		helper.SaveTileList(tileListFile, tileList)
		helper.SaveTileMap(tileMapFile, tileMap)

		fmt.Println("Processing finished and state saved.")
	}

	fmt.Printf("Total keys: %d\n", len(tileMap))
	fmt.Printf("TileList length: %d\n", len(tileList))

	// === OPTIMIZATION LOOP ===
	topN := 1     // select top 10 patterns
	numIter := 10 // number of optimization iterations

	helper.OptimizeTopPatterns(tileList, tileMap, attrList, itemList, output, topN, numIter)

}
