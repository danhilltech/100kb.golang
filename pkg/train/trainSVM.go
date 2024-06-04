package train

import (
	"fmt"

	"github.com/danhilltech/100kb.golang/pkg/domain"
	"github.com/danhilltech/100kb.golang/pkg/svm"
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
)

func trainSVM(goodEntries []Entry, allDomains []*domain.Domain) error {
	instances, specs, attrs, attrCount, err := domainsToFeaturesFloat(goodEntries, allDomains, 1.0)
	if err != nil {
		return err
	}

	trainData, testData := base.InstancesTrainTestSplit(instances, 0.4)

	_, trainRows := trainData.Size()

	svmInstances := make([]*svm.Observation, trainRows)

	for row := 0; row < trainRows; row++ {

		byteVal := trainData.Get(specs[0], row)

		str := specs[0].GetAttribute().GetStringFromSysVal(byteVal)

		fltVal := -1
		if str == "good" {
			fltVal = 1
		}

		svmInstances[row] = &svm.Observation{}

		svmInstances[row].Value = float32(fltVal)

		svmInstances[row].Features = make(map[string]float32)

		for i := 1; i < attrCount; i++ {
			byteVal := trainData.Get(specs[i], row)

			fltVal := base.UnpackBytesToFloat(byteVal)

			svmInstances[row].Features[fmt.Sprintf("ft-%d", i)] = float32(fltVal)

		}
	}

	fmt.Println("## SVM")

	model, err := svm.Train(svmInstances)
	if err != nil {
		return err
	}

	predictions := base.GeneratePredictionVector(testData)
	attrSpecs := base.ResolveAttributes(testData, attrs)

	testData.MapOverRows(attrSpecs, func(row [][]byte, srcRowNo int) (bool, error) {
		// var prediction float64 = lr.disturbance

		obs := &svm.Observation{}

		obs.Features = make(map[string]float32)

		for i := 1; i < attrCount; i++ {
			byteVal := testData.Get(attrSpecs[i], srcRowNo)

			fltVal := base.UnpackBytesToFloat(byteVal)

			obs.Features[fmt.Sprintf("ft-%d", i)] = float32(fltVal)

		}

		val, _ := model.Predict(obs, false)

		// predictions.Set(attrSpecs[0], srcRowNo, base.PackFloatToBytes(val))
		if val > 0 {
			base.SetClass(predictions, srcRowNo, "good")
		} else {
			base.SetClass(predictions, srcRowNo, "bad")
		}

		return true, nil
	})

	// var predictions base.FixedDataGrid

	// Prints precision/recall metrics
	confusionMat, err := evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		return fmt.Errorf("Unable to get confusion matrix: %s", err.Error())
	}
	fmt.Println(evaluation.GetSummary(confusionMat))

	return nil
}
