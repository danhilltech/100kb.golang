package domain

import (
	"fmt"
	"text/tabwriter"
)

func (d Domain) Tabulate(w *tabwriter.Writer) {
	fmt.Fprintf(w,
		"%s\t%s\t%d\t%t\t%t\t%t\t%t\t%s\t%t\t%t\t%t\t%0.4f\n",
		d.Domain,
		d.FeedURL,
		len(d.Articles),
		d.PageAbout,
		d.PageNow,
		d.PageBlogRoll,
		d.DomainIsPopular,
		d.DomainTLD,
		d.URLBlog,
		d.URLHumanName,
		d.URLNews,
		d.GetFPR(),
	)
}

func (d Domain) TabulateHeader(w *tabwriter.Writer) {
	fmt.Fprintf(w, "Domain\tFeed\tArticles\tAbout\tNow\tBlogRoll\tPopular\tTLD\tBlog\tHuman\tNews\tFPR\n")
}
