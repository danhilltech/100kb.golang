package ai

import (
	"github.com/viterin/vek/vek32"
)

func Cosine(a []float32, b []float32) (cosine float32) {
	return vek32.CosineSimilarity(a, b)
}
