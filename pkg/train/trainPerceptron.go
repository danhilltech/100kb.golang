package train

import (
	"fmt"

	"github.com/danhilltech/100kb.golang/pkg/domain"
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/perceptron"
)

func trainPerceptron(goodEntries []Entry, allDomains []*domain.Domain) error {

	instances, _, _, attrCount, err := domainsToFeaturesFloat(goodEntries, allDomains, 100.0)
	if err != nil {
		return err
	}

	trainData, testData := base.InstancesTrainTestSplit(instances, 0.2)

	fmt.Println("## Perceptron")
	// cls, err := linear_models.NewLogisticRegression("l0", 1.0, 1e-6)
	cls := perceptron.NewAveragePerceptron(attrCount-1, 0.9, 0.2, 0.3)
	// if err != nil {
	// 	return err
	// }

	cls.Fit(trainData)

	//Calculates the Euclidean distance and returns the most popular label
	predictions := cls.Predict(testData)

	// var predictions base.FixedDataGrid

	// Prints precision/recall metrics
	confusionMat, err := evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		panic(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
	}
	fmt.Println(evaluation.GetSummary(confusionMat))

	return nil
}
