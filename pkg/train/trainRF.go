package train

/*
func trainRF(goodEntries []Entry, allDomains []*domain.Domain) error {

	instances, err := domainsToFeaturesCategorical(goodEntries, allDomains)
if err != nil {
		return err
	}

	// fmt.Println("Running Chi Merge...")
	// filt := filters.NewChiMergeFilter(instances, 0.99)
	// for _, a := range base.NonClassFloatAttributes(instances) {
	// 	filt.AddAttribute(a)
	// }
	// fmt.Println("Training chi merge...")
	// filt.Train()
	// fmt.Println("Filtering with chi merge...")
	// instf := base.NewLazilyFilteredInstances(instances, filt)

	trainData, testData := base.InstancesTrainTestSplit(instances, 0.4)
	fmt.Println("## RF")

	cls := ensemble.NewRandomForest(80, 5)

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
*/
