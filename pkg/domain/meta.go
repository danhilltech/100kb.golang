package domain

func (d *Domain) GetFPR() float64 {
	fpr := 0.0
	if d.Articles == nil || len(d.Articles) == 0 {
		return 0.0
	}

	for _, a := range d.Articles {
		fpr += a.FirstPersonRatio
	}

	return fpr / float64(len(d.Articles))
}
