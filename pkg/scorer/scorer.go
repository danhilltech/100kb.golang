package scorer

import (
	"encoding/json"
	"math"
	"os"
)

type LogisticModel struct {
	Weights []float64 `json:"weights"`
	Bias    float64   `json:"bias"`
}

func (l *LogisticModel) Save(fileName string) error {
	jsonString, err := json.Marshal(l)
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, jsonString, os.ModePerm)

}

func LoadModel(fileName string) (*LogisticModel, error) {
	raw, err := os.ReadFile((fileName))
	if err != nil {
		return nil, err
	}

	l := LogisticModel{}

	err = json.Unmarshal(raw, &l)
	if err != nil {
		return nil, err
	}

	return &l, nil

}

func (l *LogisticModel) Predict(x [][]float64) []float64 {
	var y_pred []float64

	for i := range x {
		z := 0.0
		for j := range l.Weights {
			z += l.Weights[j] * x[i][j]
		}
		z += l.Bias
		a := Sigmoid(z)

		y_pred = append(y_pred, a)

	}
	return y_pred
}

func Sigmoid(z float64) float64 {
	sigmoid := 1 / (1 + math.Pow(math.E, -1*z))
	return sigmoid
}
