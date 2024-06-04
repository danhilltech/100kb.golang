package train

import (
	"fmt"
	"math"

	"github.com/danhilltech/100kb.golang/pkg/domain"
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
)

type LogisticModel struct {
	weights []float64
	bias    float64
}

func (l *LogisticModel) Save(filename string) error {
	return nil
}

func trainLogistic(goodEntries []Entry, allDomains []*domain.Domain) (*LogisticModel, error) {

	instances, attrSpecs, allAttrs, _, err := domainsToFeaturesFloat(goodEntries, allDomains, 1.0)
	if err != nil {
		return nil, err
	}

	fmt.Println("## Logistic")

	trainData, testData := base.InstancesTrainTestSplit(instances, 0.1)

	var x_train [][]float64
	var y_train []float64

	var x_test [][]float64
	var y_test []float64

	trainData.MapOverRows(attrSpecs, func(row [][]byte, srcRowNo int) (bool, error) {
		// var prediction float64 = lr.disturbance

		trainRowBuf := make([]float64, len(allAttrs))

		byteVal := trainData.Get(attrSpecs[0], srcRowNo)

		str := attrSpecs[0].GetAttribute().GetStringFromSysVal(byteVal)

		fltVal := float64(0)
		if str == "good" {
			fltVal = 1
		}

		for i := range allAttrs {
			trainRowBuf[i] = base.UnpackBytesToFloat(row[i])
		}

		y_train = append(y_train, fltVal)

		x_train = append(x_train, trainRowBuf[1:])
		return true, nil
	})

	testData.MapOverRows(attrSpecs, func(row [][]byte, srcRowNo int) (bool, error) {
		// var prediction float64 = lr.disturbance

		trainRowBuf := make([]float64, len(allAttrs))

		byteVal := testData.Get(attrSpecs[0], srcRowNo)

		str := attrSpecs[0].GetAttribute().GetStringFromSysVal(byteVal)

		fltVal := float64(0)
		if str == "good" {
			fltVal = 1
		}

		for i := range allAttrs {
			trainRowBuf[i] = base.UnpackBytesToFloat(row[i])
		}

		y_test = append(y_test, fltVal)

		x_test = append(x_test, trainRowBuf[1:])
		return true, nil
	})

	legacyRun(x_train, y_train, x_test, y_test, 5000, 0.01)

	w, b, err := buildLogisticModel(x_train, y_train, 50000, 1e-4)
	if err != nil {
		return nil, err
	}

	predictions, err := predictLogistic(testData, w, b)
	if err != nil {
		return nil, err
	}

	// var predictions base.FixedDataGrid

	// Prints precision/recall metrics
	confusionMat, err := evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		return nil, fmt.Errorf("Unable to get confusion matrix: %s", err.Error())
	}
	fmt.Println(evaluation.GetSummary(confusionMat))
	return &LogisticModel{weights: w, bias: b}, nil
}

func logloss(yTrue float64, yPred float64) float64 {
	loss := yTrue*math.Log1p(yPred) + (1-yTrue)*math.Log1p(1-yPred)
	return loss
}

func sigmoid(z float64) float64 {
	sigmoid := 1 / (1 + math.Pow(math.E, -1*z))
	return sigmoid
}

func initialize(features int) ([]float64, float64) {
	var w []float64
	b := 0.0
	for i := 0; i < features; i++ {
		w = append(w, 0.0)
	}
	return w, b
}

func Propagate(w []float64, b float64, x [][]float64, y []float64) (float64, []float64, float64) {
	logloss_epoch := 0.0
	m := len(y)

	var activations []float64

	var dw []float64
	db := 0.0

	for i := range x {
		var result float64 = 0.0
		for j := range w {
			result = result + (w[j] * x[i][j])
		}
		z := result + b
		a := sigmoid(z)

		activations = append(activations, a)

		var loss_example float64 = logloss(a, y[i])
		logloss_epoch = logloss_epoch + loss_example

		// Calculating gradients in same loop.
		mod_error := a - y[i]

		db = db + mod_error

	}

	for j := range w {
		error_examples := 0.0

		for i := range x {
			error_i := activations[i] - y[i]

			error_examples += error_i * x[i][j]
		}
		error_examples /= float64(m)
		dw = append(dw, error_examples)

	}

	cost := 0 - (logloss_epoch / float64(len(x))) // Negative logloss

	db /= float64(m)

	return cost, dw, db
}

func Optimize(w []float64, b float64, x [][]float64, y []float64, num_iterations int, learning_rate float64) ([]float64, float64) {

	for i := 0; i < num_iterations; i++ {
		// var cost float64
		var dw []float64
		var db float64

		_, dw, db = Propagate(w, b, x, y)

		for j := range w {
			w[j] = w[j] - (learning_rate * dw[j])
		}
		b = b - (learning_rate * db)
		// if i%100 == 0 {
		// 	fmt.Print("Cost after Iteration  ")
		// 	fmt.Println(cost)
		// }
	}
	return w, b
}

func Predict(w []float64, b float64, x [][]float64) []float64 {
	var y_pred []float64

	for i := range x {
		z := 0.0
		for j := range w {
			z += w[j] * x[i][j]
		}
		z += b
		a := sigmoid(z)

		y_pred = append(y_pred, a)

	}
	return y_pred
}

func legacyPredict(w []float64, b float64, x [][]float64) []float64 {
	var y_pred []float64

	for i := range x {
		z := 0.0
		for j := range w {
			z += w[j] * x[i][j]
		}
		z += b
		a := sigmoid(z)

		if a >= 0.5 {
			y_pred = append(y_pred, 1)
		} else {
			y_pred = append(y_pred, 0)
		}
	}
	return y_pred
}

func buildLogisticModel(x_train [][]float64, y_train []float64, num_iterations int, learning_rate float64) ([]float64, float64, error) {

	w, b := initialize(len(x_train[0]))
	_, _, _ = Propagate(w, b, x_train, y_train) // The _ s replace the cost, dw, and db. We don't need them in the main loop
	w, b = Optimize(w, b, x_train, y_train, num_iterations, learning_rate)

	return w, b, nil
}

func predictLogistic(what base.FixedDataGrid, w []float64, b float64) (base.FixedDataGrid, error) {

	_, rows := what.Size()

	classAttrs := what.AllClassAttributes()
	// fmt.Println(classAttrs)
	if len(classAttrs) != 1 {
		return nil, fmt.Errorf("Only 1 class variable is permitted")
	}
	// classAttrSpecs := base.ResolveAttributes(what, classAttrs)

	allAttrs := base.NonClassAttributes(what)
	attrs := make([]base.Attribute, 0)
	for _, a := range allAttrs {
		if _, ok := a.(*base.FloatAttribute); ok {
			attrs = append(attrs, a)
		}
	}

	cols := len(attrs) + 1

	if rows < cols {
		return nil, fmt.Errorf("not enough data")
	}

	// cls := classAttrs[0]

	ret := base.GeneratePredictionVector(what)
	attrSpecs := base.ResolveAttributes(what, attrs)
	// // clsSpec, err := ret.GetAttribute(cls)
	// if err != nil {
	// 	return nil, err
	// }

	what.MapOverRows(attrSpecs, func(row [][]byte, srcRowNo int) (bool, error) {
		// var prediction float64 = lr.disturbance

		trainRowBuf := make([]float64, len(allAttrs))

		for i := range allAttrs {
			trainRowBuf[i] = base.UnpackBytesToFloat(row[i])
		}

		predict := Predict(w, b, [][]float64{trainRowBuf})

		// fmt.Println(trainRowBuf, predict)

		// fmt.Println(trainRowBuf)
		// ret.SetClass(clsSpec, i, base.PackFloatToBytes(prediction))
		if predict[0] >= 0.5 {
			base.SetClass(ret, srcRowNo, "good")
		} else {
			base.SetClass(ret, srcRowNo, "bad")
		}
		return true, nil
	})

	return ret, nil

}

func legacyRun(x_train [][]float64, y_train []float64, x_test [][]float64, y_test []float64, num_iterations int, learning_rate float64) {
	w, b := initialize(len(x_train[0]))
	_, _, _ = Propagate(w, b, x_train, y_train) // The _ s replace the cost, dw, and db. We don't need them in the main loop
	w, b = Optimize(w, b, x_train, y_train, num_iterations, learning_rate)

	y_pred_train := legacyPredict(w, b, x_train)
	y_pred_test := legacyPredict(w, b, x_test)

	train_accuracy := 0.0
	test_accuracy := 0.0

	for i := range y_pred_train {
		if y_pred_train[i] == y_train[i] {
			train_accuracy += 1.0
		}

	}
	for i := range y_pred_test {
		if y_pred_test[i] == y_test[i] {
			test_accuracy += 1.0
		}
	}
	train_accuracy /= float64(len(y_pred_train))
	test_accuracy /= float64(len(y_pred_test))

	fmt.Print("Training Accuracy is   ")
	fmt.Println(train_accuracy)
	fmt.Print("Test Accuracy is   ")
	fmt.Println(test_accuracy)

}
