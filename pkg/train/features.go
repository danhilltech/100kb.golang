package train

import (
	"github.com/danhilltech/100kb.golang/pkg/domain"
	"github.com/sjwhitworth/golearn/base"
)

/*
func domainsToFeaturesCategorical(goodEntries []Entry, allDomains []*domain.Domain) (*base.DenseInstances, error) {

	// Training

	attrCount := 14

	attrs := make([]base.Attribute, attrCount)

	n := 0
	attrs[n] = base.NewCategoricalAttribute()
	n++

	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("fpr")

	n++
	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("wordCount")

	n++
	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("wordsPerByte")

	n++
	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("wordsPerP")

	n++
	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("goodPcnt")

	n++
	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("urlHuman")

	n++
	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("urlNews")

	n++
	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("pageNow")

	n++
	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("pageAbout")

	n++
	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("pageBlogRoll")

	n++
	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("pageWriting")

	n++
	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("goodTagsPerByte")

	n++
	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("articlesPerMonth")

	instances := base.NewDenseInstances()

	// Add the attributes
	newSpecs := make([]base.AttributeSpec, len(attrs))
	for i, a := range attrs {
		newSpecs[i] = instances.AddAttribute(a)
	}

	instances.Extend(len(goodEntries))

	// 1 title begins with a number
	// 2 number of paragraphs with more than 40 words
	// 3 average sentence length
	// 4 number of code tags
	// 5 bad keyword density ("how to", "github")
	// 6 identify self help
	// 7 youtube/podcasts
	// https://webring.xxiivv.com/#vitbaisa
	// https://frankmeeuwsen.com/blogroll/
	// title uniqueness/levenstien
	// GetCodeTagCount

	for i, train := range goodEntries {

		var domain *domain.Domain

		for _, d := range allDomains {
			if train.domain == d.Domain {
				domain = d
			}
		}

		n = 0

		if train.score == 2 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("good"))
		}
		if train.score == 1 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("bad"))
		}
		n++

		if domain.GetFPR() > 0.08 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("5"))
		} else if domain.GetFPR() > 0.04 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("4"))
		} else if domain.GetFPR() > 0.02 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("3"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}
		n++

		if domain.GetWordCount() > 1200 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("2"))
		} else if domain.GetWordCount() > 300 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}
		n++

		if domain.GetWordsPerByte() > 0.05 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("2"))
		} else if domain.GetWordsPerByte() > 0.01 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}
		n++

		if domain.GetWordsPerParagraph() > 200 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("2"))
		} else if domain.GetWordsPerParagraph() > 40 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}
		n++

		if domain.GetGoodBadTagRatio() > 0.95 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else if domain.GetGoodBadTagRatio() > 0.8 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}
		n++

		if domain.URLHumanName {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}
		n++

		if domain.URLNews {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}
		n++

		if domain.PageNow {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}
		n++

		if domain.PageAbout {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}
		n++

		if domain.PageBlogRoll {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}
		n++

		if domain.PageWriting {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}
		n++

		if domain.GetGoodTagsPerByte() > 0.002 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("2"))
		} else if domain.GetGoodTagsPerByte() > 0.001 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}
		n++

		if domain.GetLatestArticlesPerMonth() > 10 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("2"))
		} else if domain.GetLatestArticlesPerMonth() > 2 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}

	}

	instances.AddClassAttribute(attrs[0])

	return instances, nil

}
*/

func domainsToFeaturesFloat(goodEntries []Entry, allDomains []*domain.Domain, scaleUp float64) (*base.DenseInstances, []base.AttributeSpec, []base.Attribute, int, error) {
	// Training

	names := allDomains[0].GetFloatFeatureNames()

	attrCount := 1 + len(names)

	attrs := make([]base.Attribute, attrCount)

	n := 0
	attrs[n] = base.NewCategoricalAttribute()
	n++

	for _, name := range names {
		attrs[n] = base.NewFloatAttribute(name)
		n++
	}

	instances := base.NewDenseInstances()

	// Add the attributes
	newSpecs := make([]base.AttributeSpec, len(attrs))
	for i, a := range attrs {
		newSpecs[i] = instances.AddAttribute(a)
	}

	instances.Extend(len(goodEntries))

	for i, train := range goodEntries {

		var domain *domain.Domain

		for _, d := range allDomains {
			if train.domain == d.Domain {
				domain = d
			}
		}

		n = 0

		if train.score == 2 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("good"))
		}
		if train.score == 1 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("bad"))
		}
		n++

		fts := domain.GetFloatFeatures()

		for _, f := range fts {
			instances.Set(newSpecs[n], i, base.PackFloatToBytes(f))
			n++
		}

	}

	maxFloats := make([]float64, attrCount)
	minFloats := make([]float64, attrCount)
	for i := 1; i < attrCount; i++ {
		for row := 0; row < len(goodEntries); row++ {
			byteVal := instances.Get(newSpecs[i], row)

			fltVal := base.UnpackBytesToFloat(byteVal)

			if fltVal > maxFloats[i] {
				maxFloats[i] = fltVal
			}
			if fltVal < minFloats[i] {
				minFloats[i] = fltVal
			}

		}
	}

	for row := 0; row < len(goodEntries); row++ {
		for i := 1; i < attrCount; i++ {
			byteVal := instances.Get(newSpecs[i], row)

			fltVal := base.UnpackBytesToFloat(byteVal)

			rng := maxFloats[i] - minFloats[i]

			nrmVal := float64(0)

			if rng != 0 {

				nrmVal = ((fltVal - minFloats[i]) / rng) * scaleUp
			}

			instances.Set(newSpecs[i], row, base.PackFloatToBytes(nrmVal))

		}

	}

	instances.AddClassAttribute(attrs[0])

	return instances, newSpecs, attrs, attrCount, nil

}
