package train

import (
	"fmt"
	"math"

	"github.com/danhilltech/100kb.golang/pkg/scorer"
)

func trainLogistic(goodEntries []Entry) (*scorer.LogisticModel, error) {

	fmt.Println("## Logistic")

	// var x_train [][]float64
	x_train := make([][]float64, len(goodEntries))
	y_train := make([]float64, len(goodEntries))

	for i, e := range goodEntries {
		y_train[i] = float64(e.score - 1)

		fts := e.domain.GetFloatFeatures()

		x_train[i] = fts
	}

	fmt.Println(y_train[0], x_train[0])

	w, b, err := buildLogisticModel(x_train, y_train, 50000, 1e-4)
	if err != nil {
		return nil, err
	}

	model := &scorer.LogisticModel{Weights: w, Bias: b}

	return model, nil
}

func logloss(yTrue float64, yPred float64) float64 {
	loss := yTrue*math.Log1p(yPred) + (1-yTrue)*math.Log1p(1-yPred)
	return loss
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
		a := scorer.Sigmoid(z)

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

func legacyPredict(w []float64, b float64, x [][]float64) []float64 {
	var y_pred []float64

	for i := range x {
		z := 0.0
		for j := range w {
			z += w[j] * x[i][j]
		}
		z += b
		a := scorer.Sigmoid(z)

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
