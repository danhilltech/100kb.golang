package train

import (
	"fmt"

	"github.com/sjwhitworth/golearn/base"
)

func PredictBespoke(what base.FixedDataGrid) (base.FixedDataGrid, error) {

	_, rows := what.Size()

	classAttrs := what.AllClassAttributes()
	fmt.Println(classAttrs)
	if len(classAttrs) != 1 {
		return nil, fmt.Errorf("Only 1 class variable is permitted")
	}
	// classAttrSpecs := base.ResolveAttributes(what, classAttrs)

	allAttrs := base.NonClassAttributes(what)
	attrs := make([]base.Attribute, 0)
	for _, a := range allAttrs {
		if _, ok := a.(*base.FloatAttribute); ok {
			attrs = append(attrs, a)
		}
	}

	cols := len(attrs) + 1

	if rows < cols {
		return nil, fmt.Errorf("not enough data")
	}

	// cls := classAttrs[0]

	ret := base.GeneratePredictionVector(what)
	attrSpecs := base.ResolveAttributes(what, attrs)
	// // clsSpec, err := ret.GetAttribute(cls)
	// if err != nil {
	// 	return nil, err
	// }

	trainRowBuf := make([]float64, len(allAttrs))

	what.MapOverRows(attrSpecs, func(row [][]byte, srcRowNo int) (bool, error) {
		// var prediction float64 = lr.disturbance

		for i := range allAttrs {
			trainRowBuf[i] = base.UnpackBytesToFloat(row[i])
		}

		// fmt.Println(trainRowBuf)
		// ret.SetClass(clsSpec, i, base.PackFloatToBytes(prediction))
		if trainRowBuf[0] > 0.05 {
			base.SetClass(ret, srcRowNo, "good")
		} else {
			base.SetClass(ret, srcRowNo, "bad")
		}
		return true, nil
	})

	return ret, nil

}
