package output

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
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
	mid := int(float64(len(scored)) * 0.8)
	training := scored[:mid]
	// test := scored[mid:]

	txn, err := engine.db.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()

	obsArr := make([]*svm.Observation, len(training))

	mins := make(map[string]float32)
	maxs := make(map[string]float32)

	for i, train := range training {
		url := train[0]
		score, err := strconv.ParseInt(train[1], 10, 64)
		if err != nil {
			return err
		}

		article, err := engine.articleEngine.FindByURL(txn, url)
		if err != nil {
			return err
		}

		obs := svm.Observation{Ref: url, Value: float32(score)}

		obs.Features = make(map[string]float32)

		setValue(&obs, "bad_count", float32(article.BadCount), mins, maxs)
		setValue(&obs, "p_count", float32(article.BadCount), mins, maxs)

		obsArr[i] = &obs

	}

	// Normalize the features
	for _, obs := range obsArr {
		for name, n := range obs.Features {
			obs.Features[name] = (1 - -1)*(n-mins[name])/maxs[name] - mins[name] + -1
		}
	}

	fmt.Println(obsArr)

	engine.svmModel, err = svm.Train(obsArr)
	if err != nil {
		return err
	}

	return nil

}

func setValue(obs *svm.Observation, name string, value float32, mins, maxs map[string]float32) {
	obs.Features[name] = value
	if obs.Features[name] < mins[name] {
		mins[name] = obs.Features[name]
	}
	if obs.Features[name] > maxs[name] {
		maxs[name] = obs.Features[name]
	}
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
