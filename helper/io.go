package helper

import (
	"encoding/gob"
	"os"

	"github.com/italanleal/go-wfcsd/models"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// SaveItemList salva itemList em arquivo
func SaveItemList(filename string, itemList []models.Item) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	return enc.Encode(itemList)
}

// LoadItemList carrega itemList de arquivo
func LoadItemList(filename string) ([]models.Item, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var itemList []models.Item
	dec := gob.NewDecoder(file)
	err = dec.Decode(&itemList)
	return itemList, err
}

// SaveAttrList salva lista de atributos
func SaveAttrList(filename string, attrList []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	return enc.Encode(attrList)
}

// LoadAttrList carrega lista de atributos
func LoadAttrList(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var attrList []string
	dec := gob.NewDecoder(file)
	err = dec.Decode(&attrList)
	return attrList, err
}

// SaveTileList salva tileList
func SaveTileList(filename string, tileList []models.Pattern) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	return enc.Encode(tileList)
}

// LoadTileList carrega tileList
func LoadTileList(filename string) ([]models.Pattern, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tileList []models.Pattern
	dec := gob.NewDecoder(file)
	err = dec.Decode(&tileList)
	return tileList, err
}

// SaveTileMap salva tileMap
func SaveTileMap(filename string, tileMap map[string][]int) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	return enc.Encode(tileMap)
}

// LoadTileMap carrega tileMap
func LoadTileMap(filename string) (map[string][]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tileMap map[string][]int
	dec := gob.NewDecoder(file)
	err = dec.Decode(&tileMap)
	return tileMap, err
}
