package ai

/*
#cgo LDFLAGS: -L../../lib -lgobert
#include "../../lib/gobert.h"
*/
import "C"
import (
	"strings"
	"unsafe"

	"google.golang.org/protobuf/proto"
)

type SentenceEmbeddingModel struct {
	model unsafe.Pointer
}

type KeywordExtractionModel struct {
	model unsafe.Pointer
}

const maxWordCount = 256

// wget https://huggingface.co/skeskinen/ggml/resolve/main/bert-base-uncased/ggml-model-q4_0.bin?download=true
// wget https://huggingface.co/mudler/all-MiniLM-L6-v2/resolve/main/ggml-model-q4_0.bin?download=true -O models/bert.bin

func NewSentenceEmbeddingModel() (*SentenceEmbeddingModel, error) {

	a := C.new_sentence_embedding()

	return &SentenceEmbeddingModel{model: unsafe.Pointer(a)}, nil

}

func (ai *SentenceEmbeddingModel) Embeddings(texts []string) ([]*Embedding, error) {

	req := SentenceEmbeddingRequest{}
	textsTrimmed := make([]string, len(texts))

	for i := 0; i < len(texts); i++ {
		in := texts[i]
		cut := strings.Split(in, " ")
		l := min(len(cut), maxWordCount)

		textsTrimmed[i] = strings.Join(cut[0:l], " ")
	}

	req.Texts = textsTrimmed

	reqBytes, err := proto.Marshal(&req)
	if err != nil {
		return nil, err
	}

	var outSize uintptr

	reqSize := uintptr(len(reqBytes))

	coutSize := unsafe.Pointer(&outSize)
	creqSize := unsafe.Pointer(&reqSize)
	reqPtr := unsafe.Pointer(&reqBytes[0])

	cout := C.sentence_embedding((*C.SharedSentenceEmbeddingModel)(ai.model), (*C.uchar)(reqPtr), (*C.size_t)(creqSize), (*C.size_t)(coutSize))
	if outSize > 0 {
		defer C.drop_bytesarray(cout)
	}

	var chunks SentenceEmbeddingResponse

	protoBuf := unsafe.Slice((*byte)(cout), outSize)

	err = proto.Unmarshal(protoBuf, &chunks)
	if err != nil {
		return nil, err
	}

	return chunks.Texts, nil
}

func (ai *SentenceEmbeddingModel) Close() {
	C.drop_sentence_embedding((*C.SharedSentenceEmbeddingModel)(ai.model))
}

func NewKeywordExtractionModel() (*KeywordExtractionModel, error) {

	a := C.new_keyword_extraction()

	return &KeywordExtractionModel{model: unsafe.Pointer(a)}, nil

}

func (ai *KeywordExtractionModel) Extract(texts []string) ([]*Keywords, error) {

	req := KeywordRequest{}

	textsTrimmed := make([]string, len(texts))

	for i := 0; i < len(texts); i++ {
		in := texts[i]
		cut := strings.Split(in, " ")
		l := min(len(cut), maxWordCount)

		textsTrimmed[i] = strings.Join(cut[0:l], " ")
	}

	req.Texts = textsTrimmed

	reqBytes, err := proto.Marshal(&req)
	if err != nil {
		return nil, err
	}

	var outSize uintptr

	reqSize := uintptr(len(reqBytes))

	coutSize := unsafe.Pointer(&outSize)
	creqSize := unsafe.Pointer(&reqSize)
	reqPtr := unsafe.Pointer(&reqBytes[0])

	cout := C.keyword_extraction((*C.SharedKeywordExtractionModel)(ai.model), (*C.uchar)(reqPtr), (*C.size_t)(creqSize), (*C.size_t)(coutSize))
	if outSize > 0 {
		defer C.drop_bytesarray(cout)
	}

	var chunks KeywordResponse

	protoBuf := unsafe.Slice((*byte)(cout), outSize)

	err = proto.Unmarshal(protoBuf, &chunks)
	if err != nil {
		return nil, err
	}

	return chunks.Texts, nil
}

func (ai *KeywordExtractionModel) Close() {
	C.drop_keyword_extraction((*C.SharedKeywordExtractionModel)(ai.model))
}
