package output

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"slices"
	"strconv"

	"github.com/danhilltech/100kb.golang/pkg/svm"
)

func (engine *RenderEngine) TrainSVM(filePath string) error {
	// Open the scoring file

	scored, err := readCsvFile(filePath)
	if err != nil {
		return err
	}

	rand.Shuffle(len(scored), func(i, j int) { scored[i], scored[j] = scored[j], scored[i] })

	// 80/20 split

	txn, err := engine.db.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()

	obsArr := make([]*svm.Observation, len(scored))

	featureVals := make(map[string][]float32, 3)

	for i, train := range scored {
		url := train[0]
		score, err := strconv.ParseInt(train[1], 10, 64)
		if err != nil {
			return err
		}

		article, err := engine.articleEngine.FindByURL(txn, url)
		if err != nil {
			return err
		}

		scoreClean := -1
		if score >= 4 {
			scoreClean = 1
		}

		obs := svm.Observation{Ref: url, Value: float32(scoreClean)}

		obs.Features = make(map[string]float32)

		setValue(&obs, "bad_count", float32(math.Log(float64(article.BadCount)+1.0)), featureVals)
		setValue(&obs, "p_count", float32(math.Log(float64(article.PCount)+1.0)), featureVals)
		setValue(&obs, "fpr", float32(article.FirstPersonRatio), featureVals)

		obsArr[i] = &obs
		fmt.Printf("%+v\n", obs.Features)

	}

	// Normalize the features
	for _, obs := range obsArr {
		for name, n := range obs.Features {
			min := slices.Min(featureVals[name])
			max := slices.Max(featureVals[name])

			fmt.Printf("Min/Max for %s: %0.3f/%0.3f\n\n", name, min, max)

			obs.Features[name] = (2 * ((n - min) / max)) - min - 1
		}
		fmt.Printf("%+v\n", obs.Features)
	}

	mid := int(float64(len(obsArr)) * 0.8)
	training := obsArr[:mid]
	test := obsArr[mid:]

	engine.svmModel, err = svm.Train(training)
	if err != nil {
		return err
	}

	for _, t := range test {
		outVal, outProbs := engine.svmModel.Predict(t, true)
		fmt.Print("test:\twanted %0.2f\tgot %0.2f\tprops: %0.4f\t\t%s", t.Value, outVal, outProbs, t.Ref)
	}

	return nil

}

func setValue(obs *svm.Observation, name string, value float32, featureVals map[string][]float32) {
	obs.Features[name] = value
	featureVals[name] = append(featureVals[name], value)
}

func readCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}
