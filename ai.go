package main

/*
#cgo LDFLAGS: -L./lib -lgobert
#include "./lib/gobert.h"
*/
import "C"
import (
	"unsafe"
)

type AI struct {
	model unsafe.Pointer
}

// wget https://huggingface.co/skeskinen/ggml/resolve/main/bert-base-uncased/ggml-model-q4_0.bin?download=true
// wget https://huggingface.co/mudler/all-MiniLM-L6-v2/resolve/main/ggml-model-q4_0.bin?download=true -O models/bert.bin

func loadAi(model string) (*AI, error) {

	a := C.new_sentence_embedding()

	return &AI{model: unsafe.Pointer(a)}, nil

}

func (ai *AI) Embeddings(texts []string) ([]float32, error) {

	out := make([]float32, 384*len(texts))

	cout := unsafe.Pointer(&out[0])

	ptrs := make([]*C.char, len(texts))

	for i := 0; i < len(texts); i++ {
		ptrs[i] = C.CString(texts[i])
		defer C.free(unsafe.Pointer(ptrs[i]))
	}

	ptr := &ptrs[0]

	C.sentence_embedding((*C.Model)(ai.model), ptr, cout)

	return out, nil
}

func (ai *AI) Close() {
	C.drop_sentence_embedding((*C.Model)(ai.model))
}
