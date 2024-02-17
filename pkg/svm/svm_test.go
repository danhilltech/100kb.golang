package svm

import (
	"testing"
)

func TestBasic(t *testing.T) {

	obs := []*Observation{}

	obs = append(obs, &Observation{
		Ref:   "1",
		Value: 1,
		Features: map[string]float32{
			"a": 0.99,
			"b": 0.89,
		},
	})

	obs = append(obs, &Observation{
		Ref:   "2",
		Value: -1,
		Features: map[string]float32{
			"a": 0.01,
			"b": 0.08,
		},
	})
	obs = append(obs, &Observation{
		Ref:   "2",
		Value: -1,
		Features: map[string]float32{
			"a": 0.09,
			"b": 0.08,
		},
	})

	model, err := Train(obs)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	model.Predict(&Observation{
		Ref:   "2",
		Value: -1,
		Features: map[string]float32{
			"a": 0.99,
			"b": 0.99,
		},
	}, true)
}
