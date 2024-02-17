package svm

// #cgo LDFLAGS: -L. -lsvm
// #include "svm.h"
// #include <stdlib.h>
import "C"
import (
	"unsafe"
)

/**
* Problem
* l: number of observations
* y: floats of target values
* x: nodes of input features
 */

const C_SVC = 0
const NU_SVC = 1
const ONE_CLASS = 2
const EPSILON_SVR = 3
const NU_SVR = 4

type Model struct {
	model         *C.struct_svm_model
	modelType     int
	modelNClasses int
}

type Observation struct {
	Ref      string
	Value    float32
	Features map[string]float32
}

func Train(obs []*Observation) (*Model, error) {

	l := len(obs)
	features := len(obs[0].Features)

	fMap := make(map[string]int, features)

	fpos := 0
	for name := range obs[0].Features {
		fMap[name] = fpos
		fpos++
	}

	problem := C.struct_svm_problem{
		l: C.int(l),
	}

	x := (**C.struct_svm_node)(C.malloc(C.size_t(l) * C.size_t(unsafe.Sizeof(C.struct_svm_node{}))))
	xArr := unsafe.Slice((**C.struct_svm_node)(x), l)
	defer C.free(unsafe.Pointer(x))

	xSpace := (*C.struct_svm_node)(C.malloc(C.size_t(l) * C.size_t(features+1) * C.size_t(unsafe.Sizeof(C.struct_svm_node{}))))
	xSpaceArr := unsafe.Slice((*C.struct_svm_node)(xSpace), l*(features+1))
	defer C.free(unsafe.Pointer(xSpace))

	y := (*C.double)(C.malloc(C.size_t(l) * C.size_t(unsafe.Sizeof(C.double(0)))))
	yArr := unsafe.Slice((*C.double)(y), l)
	defer C.free(unsafe.Pointer(y))

	j := 0
	for i, ob := range obs {
		xArr[i] = &xSpaceArr[j]
		for fName, feature := range ob.Features {

			nd := C.struct_svm_node{}
			xSpaceArr[j] = nd

			xSpaceArr[j].index = C.int(fMap[fName])
			xSpaceArr[j].value = C.double(feature)

			j++
		}
		yArr[i] = C.double(ob.Value)

		xSpaceArr[j].index = -1

		j++
	}

	problem.x = (**C.struct_svm_node)(x)
	problem.y = y

	params := C.struct_svm_parameter{}

	params.svm_type = NU_SVC
	params.gamma = C.double(1.0 / float32(features))
	params.degree = 3
	params.nu = 0.5
	params.C = 1
	params.eps = 1e-3
	params.p = 0.1
	params.shrinking = 1
	params.probability = 1
	params.kernel_type = 1
	params.degree = 3
	params.coef0 = 0
	// params. = 0.1
	params.cache_size = 100

	model := C.svm_train(&problem, &params)

	mdl := &Model{
		model:         model,
		modelType:     int(params.svm_type),
		modelNClasses: 2,
	}

	return mdl, nil
}

func (model *Model) Predict(obs *Observation, predictProbability bool) (float64, []float64) {
	l := 1
	features := len(obs.Features)

	fMap := make(map[string]int, features)

	fpos := 0
	for name := range obs.Features {
		fMap[name] = fpos
		fpos++
	}

	x := (*C.struct_svm_node)(C.malloc(C.size_t(l) * C.size_t(unsafe.Sizeof(C.struct_svm_node{}))))
	xArr := unsafe.Slice((*C.struct_svm_node)(x), l)
	defer C.free(unsafe.Pointer(x))

	xSpace := (*C.struct_svm_node)(C.malloc(C.size_t(l) * C.size_t(features+1) * C.size_t(unsafe.Sizeof(C.struct_svm_node{}))))
	xSpaceArr := unsafe.Slice((*C.struct_svm_node)(xSpace), l*(features+1))
	defer C.free(unsafe.Pointer(xSpace))

	j := 0

	xArr[0] = xSpaceArr[j]
	for fName, feature := range obs.Features {

		nd := C.struct_svm_node{}
		xSpaceArr[j] = nd

		xSpaceArr[j].index = C.int(fMap[fName])
		xSpaceArr[j].value = C.double(feature)

		j++
	}

	xSpaceArr[j].index = -1

	if predictProbability && (model.modelType == C_SVC || model.modelType == NU_SVC || model.modelType == ONE_CLASS) {
		probEstimates := (*C.double)(C.malloc(C.size_t(model.modelNClasses) * C.size_t(unsafe.Sizeof(C.double(0)))))
		defer C.free(unsafe.Pointer(probEstimates))
		predictLabel := C.svm_predict_probability(model.model, xSpace, probEstimates)

		probEstimatesGo := (*[1 << 30]float64)(unsafe.Pointer(probEstimates))[:model.modelNClasses:model.modelNClasses]

		return float64(predictLabel), probEstimatesGo
	} else {
		predictLabel := C.svm_predict(model.model, xSpace)

		return float64(predictLabel), nil
	}

}
