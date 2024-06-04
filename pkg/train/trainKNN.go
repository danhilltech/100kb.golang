package train

import (
	"fmt"

	"github.com/danhilltech/100kb.golang/pkg/domain"
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/knn"
)

func trainKNN(goodEntries []Entry, allDomains []*domain.Domain) error {

	instances, _, _, _, err := domainsToFeaturesFloat(goodEntries, allDomains, 100.0)
	if err != nil {
		return err
	}

	trainData, testData := base.InstancesTrainTestSplit(instances, 0.25)

	fmt.Println("## KNN")
	cls := knn.NewKnnClassifier("cosine", "linear", 3)
	// cls := trees.NewID3DecisionTree(0.6)

	// cls := ensemble.NewRandomForest(40, 3)

	// Create a 60-40 training-test split

	err = cls.Fit(trainData)
	if err != nil {
		return err
	}

	//Calculates the Euclidean distance and returns the most popular label
	predictions, err := cls.Predict(testData)
	if err != nil {
		panic(err)
	}

	// var predictions base.FixedDataGrid

	// Prints precision/recall metrics
	confusionMat, err := evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		panic(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
	}
	fmt.Println(evaluation.GetSummary(confusionMat))

	return nil
}
