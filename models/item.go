package models

type Item struct {
	Attr  int
	Value string // melhor string, pq CSV pode ter textos
	Index []int  // índices onde ocorre
	SuppP float64
}
