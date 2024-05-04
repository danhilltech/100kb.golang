//go:build !cuda
// +build !cuda

package ai

import "C"
import (
	"fmt"
	"unsafe"
)

type SentenceEmbeddingModel struct {
	model unsafe.Pointer
}

type KeywordExtractionModel struct {
	model unsafe.Pointer
}

type ZeroShotModel struct {
	model unsafe.Pointer
}

const maxWordCount = 256

var ErrNotImplemented = fmt.Errorf("Not implemented without CUDA")

// wget https://huggingface.co/skeskinen/ggml/resolve/main/bert-base-uncased/ggml-model-q4_0.bin?download=true
// wget https://huggingface.co/mudler/all-MiniLM-L6-v2/resolve/main/ggml-model-q4_0.bin?download=true -O models/bert.bin

func NewSentenceEmbeddingModel() (*SentenceEmbeddingModel, error) {

	return nil, ErrNotImplemented

}

func (ai *SentenceEmbeddingModel) Embeddings(texts []string) ([]*Embedding, error) {
	return nil, ErrNotImplemented
}

func (ai *SentenceEmbeddingModel) Close() {

}

func NewKeywordExtractionModel() (*KeywordExtractionModel, error) {

	return nil, ErrNotImplemented

}

func (ai *KeywordExtractionModel) Extract(texts []string) ([]*Keywords, error) {
	return nil, ErrNotImplemented
}

func (ai *KeywordExtractionModel) Close() {

}

func NewZeroShotModel() (*ZeroShotModel, error) {
	return nil, ErrNotImplemented

}

func (ai *ZeroShotModel) Predict(texts []string, labels []string) ([]*ZeroShotClassifications, error) {
	return nil, ErrNotImplemented
}

func (ai *ZeroShotModel) Close() {

}
