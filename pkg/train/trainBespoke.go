package train

import (
	"fmt"

	"github.com/danhilltech/100kb.golang/pkg/domain"
	"github.com/sjwhitworth/golearn/evaluation"
)

func trainBespoke(goodEntries []Entry, allDomains []*domain.Domain) error {

	instances, _, _, _, err := domainsToFeaturesFloat(goodEntries, allDomains, 1.0)
	if err != nil {
		return err
	}

	fmt.Println("## Bespoke")

	predictions, err := PredictBespoke(instances)
	if err != nil {
		return err
	}

	// var predictions base.FixedDataGrid

	// Prints precision/recall metrics
	confusionMat, err := evaluation.GetConfusionMatrix(instances, predictions)
	if err != nil {
		return fmt.Errorf("Unable to get confusion matrix: %s", err.Error())
	}
	fmt.Println(evaluation.GetSummary(confusionMat))
	return nil
}
