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

func (l *LogisticModel) Predict(x [][]float64, names []string) []float64 {
	var y_pred []float64

	for i := range x {
		z := 0.0
		for j := range l.Weights {
			z += l.Weights[j] * x[i][j]
		}
		z += l.Bias
		a := Sigmoid(z)

		// Now bias it by our values

		if getValByName(x[i], names, "fpr") < 0.01 {
			a = a * 0.5
		}
		if getValByName(x[i], names, "fpr") > 0.025 {
			a = a * 1.2
		}

		// if getValByName(x[i], names, "tti") > 10 {
		// 	a = a * 0.2
		// }
		if getValByName(x[i], names, "totalScriptRequests") > 3 {
			a = a * 0.2
		}
		if getValByName(x[i], names, "loadsGoogleAds") > 0 {
			a = a * 0.1
		}
		if getValByName(x[i], names, "loadsGoogleTagManager") > 0 {
			a = a * 0.5
		}

		if a > 1 {
			a = 1.0
		}

		y_pred = append(y_pred, a)

	}
	return y_pred
}

func getValByName(vals []float64, names []string, target string) float64 {
	for i, name := range names {
		if name == target {
			return vals[i]
		}
	}
	return 0
}

func Sigmoid(z float64) float64 {
	sigmoid := 1 / (1 + math.Pow(math.E, -1*z))
	return sigmoid
}
