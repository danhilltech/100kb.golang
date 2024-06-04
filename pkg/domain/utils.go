package domain

import (
	"fmt"
	"text/tabwriter"
)

func (d Domain) Tabulate(w *tabwriter.Writer) {

	fts := d.GetFloatFeatures()

	fmt.Fprintf(w,
		"%s\t%s\t%d\t",
		d.Domain,
		d.FeedURL,
		len(d.Articles),
	)

	for _, f := range fts {
		fmt.Fprintf(w, "%0.4f\t", f)
	}
	fmt.Fprintf(w, "\n")
}

func (d Domain) TabulateHeader(w *tabwriter.Writer) {
	fmt.Fprintf(w, "Domain\tFeed\tArticles\t")

	names := d.GetFloatFeatureNames()

	for _, f := range names {
		fmt.Fprintf(w, "%s\t", f)
	}
	fmt.Fprintf(w, "\n")
}
