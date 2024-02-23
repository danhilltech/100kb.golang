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

		// setValue(&obs, "bad_count", float32(math.Log(float64(article.BadCount)+1.0)), featureVals)
		// setValue(&obs, "p_count", float32(math.Log(float64(article.PCount)+1.0)), featureVals)

		setValue(&obs, "bad_count", float32(math.Log(float64(article.BadCount)+1.0)), featureVals)
		setValue(&obs, "word_count", float32(math.Log(float64(article.WordCount)+1.0)), featureVals)
		setValue(&obs, "bad_ratio", float32(math.Log(float64(article.BadCount)+1.0))/float32(article.WordCount), featureVals)

		setValue(&obs, "fpr", float32(article.FirstPersonRatio), featureVals)
		setValueBool(&obs, "domain_popular", article.DomainIsPopular, featureVals)
		setValueBool(&obs, "page_about", article.PageAbout, featureVals)
		obsArr[i] = &obs

		fmt.Printf("%+v\t%s\n", obs.Features, obs.Ref)

	}

	// Normalize the features
	for _, obs := range obsArr {
		for name, n := range obs.Features {
			min := slices.Min(featureVals[name])
			max := slices.Max(featureVals[name])

			obs.Features[name] = ((n - min) / (max - min) * (2)) + min

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

	correct := 0
	correctRight := 0 // we got it right in the positive
	missedRight := 0  // should have been right and we said wrong

	totalPositive := 0
	for _, t := range test {
		outVal, probs := engine.svmModel.Predict(t, true)
		if outVal == float64(t.Value) {
			correct++
		}
		if outVal == float64(t.Value) && t.Value == 1 {
			correctRight++
		}

		if t.Value == 1 && outVal == -1 {
			missedRight++
		}

		if t.Value == 1 {
			totalPositive++
		}

		fmt.Printf("test:\twanted %0.2f\tgot %0.2f\t%+v\t%s\n", t.Value, outVal, probs, t.Ref)
	}
	fmt.Printf("ACCURACY: %0.2f%%\n", (float64(correct) / float64(len(test)) * 100))
	fmt.Printf("CORRECT ACCURACY: %0.2f%%\n", (float64(correctRight) / float64(totalPositive) * 100))
	fmt.Printf("MISSED ACCURACY: %0.2f%%\n", (float64(missedRight) / float64(totalPositive) * 100))

	return nil

}

func setValue(obs *svm.Observation, name string, value float32, featureVals map[string][]float32) {

	v := value

	if math.IsInf(float64(value), 0) {
		v = 0
	}
	if math.IsNaN(float64(value)) {
		v = 0
	}

	obs.Features[name] = v
	featureVals[name] = append(featureVals[name], v)
}

func setValueBool(obs *svm.Observation, name string, value bool, featureVals map[string][]float32) {

	v := float32(0)

	if value {
		v = 1
	} else {
		v = -1
	}
	setValue(obs, name, v, featureVals)
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
