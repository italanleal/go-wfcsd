package helper

import (
	"fmt"

	"github.com/italanleal/go-wfcsd/models"
)

// func GenerateTwoItemCombinations(itemList *[]models.Item, output *[]bool, tileList *[]models.Pattern, tileMap *map[string][]int) {
// 	n := len(*itemList)
// 	for i := 0; i < n-1; i++ {
// 		for j := i + 1; j < n; j++ {
// 			if (*itemList)[i].Attr != (*itemList)[j].Attr {
// 				pattern := models.Pattern{
// 					Items: []int{i, j},
// 				}
// 				idx := len(*tileList)
// 				*tileList = append(*tileList, pattern)

// 				// passar diretamente os atributos para a função de chave
// 				key := MakePatternKeyAttrs((*itemList)[i].Attr, (*itemList)[j].Attr)
// 				(*tileMap)[key] = append((*tileMap)[key], idx)
// 			}
// 		}
// 	}
// }

func GenerateTwoItemCombinations(itemList *[]models.Item, output *[]bool, tileList *[]models.Pattern, tileMap *map[string][]int) {
	n := len(*itemList)
	seen := make(map[string]bool)

	for i := 0; i < n-1; i++ {
		for j := 0; j < n; j++ {
			if (*itemList)[i].Attr != (*itemList)[j].Attr {
				itemKey := MakePatternKeyAttrs(i, j)
				if seen[itemKey] {
					continue
				}
				seen[itemKey] = true

				// create a unique key for the attribute pair
				attrKey := MakePatternKeyAttrs((*itemList)[i].Attr, (*itemList)[j].Attr)

				pattern := models.Pattern{
					Items: []int{i, j},
				}
				idx := len(*tileList)
				*tileList = append(*tileList, pattern)

				(*tileMap)[attrKey] = append((*tileMap)[attrKey], idx)

			}
		}
	}
}

func MakePatternKeyAttrs(attr1, attr2 int) string {
	if attr1 > attr2 {
		attr1, attr2 = attr2, attr1 // garante ordem consistente
	}
	return fmt.Sprintf("%d:%d", attr1, attr2)
}

// func PopulatePatternOptions(tileList *[]models.Pattern, tileMap *map[string][]int, attrList []string, itemList []models.Item) {
// 	attrIndices := make([]int, len(attrList))
// 	for i := range attrList {
// 		attrIndices[i] = i
// 	}

// 	for pIdx := range *tileList {
// 		fmt.Printf("start processing %d tile adjacency/n", pIdx)
// 		pattern := &(*tileList)[pIdx]

// 		// pegar atributos usados no pattern
// 		usedAttrs := []int{}
// 		for _, itemIdx := range pattern.Items {
// 			usedAttrs = append(usedAttrs, itemList[itemIdx].Attr)
// 		}

// 		// criar lista de atributos restantes
// 		remainingAttrs := []int{}
// 		for _, attr := range attrIndices {
// 			found := false
// 			for _, ua := range usedAttrs {
// 				if ua == attr {
// 					found = true
// 					break
// 				}
// 			}
// 			if !found {
// 				remainingAttrs = append(remainingAttrs, attr)
// 			}
// 		}

// 		optionSet := make(map[int]bool)

// 		// combinar cada item do pattern com os atributos restantes
// 		for _, itemIdx := range pattern.Items {
// 			attr1 := itemList[itemIdx].Attr
// 			for _, remAttr := range remainingAttrs {
// 				key := MakePatternKeyAttrs(attr1, remAttr)

// 				if indices, ok := (*tileMap)[key]; ok {
// 					for _, idx := range indices {
// 						optionSet[idx] = true
// 					}
// 				}
// 			}
// 		}

// 		// preencher Options
// 		pattern.Options = pattern.Options[:0]
// 		for idx := range optionSet {
// 			pattern.Options = append(pattern.Options, idx)
// 		}
// 	}
// }

// PopulatePatternOptionsForSingle calculates the Options for a single pattern
func PopulatePatternOptions(pattern *models.Pattern, tileList []models.Pattern, tileMap *map[string][]int, attrList []string, itemList []models.Item) {
	// Reset options
	pattern.Options = pattern.Options[:0]

	// Build a set of attributes already in the pattern
	attrInPattern := []int{}
	for _, it := range pattern.Items {
		attrInPattern = append(attrInPattern, itemList[it].Attr)
	}

	// Iterate over attribute *values*
	for _, attr1 := range attrInPattern {
		cp := []int{}
		for _, att := range attrInPattern {
			if att != attr1 {
				cp = append(cp, att)
			}
		}
		fmt.Printf("attrInPattern: %v\n", attrInPattern)
		fmt.Printf("cp:            %v\n", cp)

		for attr2, _ := range attrList {
			if SliceContains(cp, attr2) {
				continue
			}
			key := MakePatternKeyAttrs(attr1, attr2)
			if indices, ok := (*tileMap)[key]; ok {
				for _, i := range indices {
					overlap := 0
					for _, pItemIdx := range pattern.Items {
						for _, optItemIdx := range tileList[i].Items {
							if pItemIdx == optItemIdx {
								overlap++
							}
						}
					}

					if overlap == 1 {
						pattern.Options, _ = AddUnique(pattern.Options, i)
					}
				}
			}
		}
	}
}

// MergeOptionsSingleAttrOverlap merges two option slices, keeping only
// options whose pattern in tileList shares exactly one attribute with each other
func MergeOptionsSingleAttrOverlap(opt1, opt2 []int, tileList []models.Pattern, itemList []models.Item) []int {
	merged := append(opt1, opt2...)
	result := []int{}

	for _, optIdx := range merged {
		option := tileList[optIdx]
		attrSet := make(map[int]struct{})
		for _, idx := range option.Items {
			attrSet[itemList[idx].Attr] = struct{}{}
		}

		// Only keep patterns with exactly one attribute (length of attrSet == 1)
		if len(attrSet) == 1 {
			result = append(result, optIdx)
		}
	}

	return result
}

func FilterOptionsBySingleAttrOverlap(pattern *models.Pattern, tileList []models.Pattern, itemList []models.Item) []int {
	result := []int{}

	for _, optIdx := range pattern.Options {
		option := tileList[optIdx]

		// Count how many attributes overlap with pattern.Items
		overlap := 0
		for _, pItemIdx := range pattern.Items {
			attr1 := itemList[pItemIdx].Attr
			for _, optItemIdx := range option.Items {
				attr2 := itemList[optItemIdx].Attr
				if attr1 == attr2 {
					overlap++
				}
			}
		}

		if overlap == 1 {
			result = append(result, optIdx)
		}
	}

	return result
}
