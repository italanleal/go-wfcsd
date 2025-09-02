package helper

import (
	"fmt"
	"sort"

	"github.com/italanleal/go-wfcsd/models"
)

func OptimizeTopPatterns(tileList []models.Pattern, tileMap map[string][]int, attrList []string, itemList []models.Item, output []bool, numTop int, numIter int) {
	// Select top patterns (make copies)
	if numTop > len(tileList) {
		numTop = len(tileList)
	}
	topPatterns := make([]*models.Pattern, numTop)
	for i := 0; i < numTop; i++ {
		// Create a copy
		topPatterns[i] = new(models.Pattern)
		topPatterns[i].Items = append([]int{}, tileList[i].Items...)
		topPatterns[i].Freq = tileList[i].Freq
		topPatterns[i].Wracc = tileList[i].Wracc
		topPatterns[i].Options = make([]int, 0)

		PopulatePatternOptions(topPatterns[i], tileList, &tileMap, attrList, itemList)
		sort.Slice(topPatterns[i].Options, func(a, b int) bool {
			return tileList[topPatterns[i].Options[a]].Wracc > tileList[topPatterns[i].Options[b]].Wracc
		})

	}
	// for g, gdx := range topPatterns[0].Options {
	// 	tile := tileList[gdx]
	// 	for _, item := range tile.Items {
	// 		attrName := attrList[itemList[item].Attr]
	// 		fmt.Printf("Option %d: ItemIdx=%d, Attr=%s, Value=%s\n", g, item, attrName, itemList[item].Value)
	// 	}

	// }

	// Step 2: Iterative WRAcc improvement
	for iter := 0; iter < numIter; iter++ {
		for _, pattern := range topPatterns {
			if len(pattern.Options) == 0 {
				fmt.Printf("empty option error\n")
				continue
			}

			optionIdx, _ := SliceShift(&pattern.Options)

			origOption := tileList[optionIdx]

			optionPattern := &models.Pattern{
				Items: append([]int{}, origOption.Items...),
				Freq:  origOption.Freq,
				Wracc: origOption.Wracc,
			}
			fmt.Printf("%v\n", len(pattern.Options))
			PrintPatternCompact(pattern, itemList, attrList)
			PrintPatternCompact(optionPattern, itemList, attrList)

			newItems := append([]int{}, pattern.Items...)

			overlapCount := 0
			for _, j := range optionPattern.Items {
				for _, i := range pattern.Items {
					if itemList[i].Attr == itemList[j].Attr {
						overlapCount++
						fmt.Printf("%v pattern attr collides with %v option %v  attr\n", attrList[itemList[i].Attr], optionIdx, attrList[itemList[j].Attr])
						break
					} else {
						newItems, _ = AddUnique(newItems, j)
					}
				}
			}
			if overlapCount > 1 {
				fmt.Println("No new items could be merged")
				continue
			}
			optionPattern.Items = newItems
			CalcPatternStats(optionPattern, itemList, output)
			fmt.Printf("Enhanced Pattern: %v\n", optionPattern.Items)
			fmt.Printf("WRAcc: %.4f, Freq: %f\n", optionPattern.Wracc, optionPattern.Freq)

			if optionPattern.Wracc > pattern.Wracc {
				pattern.Wracc = optionPattern.Wracc
				pattern.Freq = optionPattern.Freq
				pattern.Items = optionPattern.Items
				PopulatePatternOptions(pattern, tileList, &tileMap, attrList, itemList)
			} else {
				pattern.Options = pattern.Options[1:]
			}
		}
	}

	// Step 3: Print top patterns sorted by WRAcc
	sort.Slice(topPatterns, func(i, j int) bool {
		return topPatterns[i].Wracc > topPatterns[j].Wracc
	})

	fmt.Println("Top patterns after optimization:")
	PrintTopPatterns(topPatterns, itemList, attrList)
}

func PrintTopPatterns(topPatterns []*models.Pattern, itemList []models.Item, attrList []string) {
	for i, p := range topPatterns {
		fmt.Printf("%d: WRAcc=%.6f, Freq=%.4f\n", i, p.Wracc, p.Freq)
		fmt.Printf("   Items: ")
		for _, idx := range p.Items {
			item := itemList[idx]
			attrName := attrList[item.Attr]
			fmt.Printf("[%s: %v] ", attrName, item.Value) // replace Value with your actual field
		}
		fmt.Println()
	}
}

func PrintPatternCompact(pattern *models.Pattern, itemList []models.Item, attrList []string) {
	fmt.Printf("WRAcc=%.6f, Freq=%.4f | Items: ", pattern.Wracc, pattern.Freq)
	for _, idx := range pattern.Items {
		item := itemList[idx]
		attrName := attrList[item.Attr]
		fmt.Printf("%s=%v ", attrName, item.Value)
	}
	fmt.Println()
}
