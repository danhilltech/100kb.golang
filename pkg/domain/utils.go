package domain

import (
	"fmt"
	"text/tabwriter"
)

func (d Domain) Tabulate(w *tabwriter.Writer) {
	fmt.Fprintf(w,
		"%s\t%s\t%d\t%t\t%t\t%t\t%t\t%s\t%t\t%t\t%t\t%0.4f\t%d\t%0.4f\t%d\t%d\t%0.4f\t%0.4f\t%0.4f\t%0.4f\t%d4\n",
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
		d.GetWordCount(),
		d.GetWordsPerByte(),
		d.GetGoodTagCount(),
		d.GetBadTagCount(),
		d.GetWordsPerParagraph(),
		d.GetGoodBadTagRatio(),
		d.GetGoodTagsPerByte(),
		d.GetLatestArticlesPerMonth(),
		d.GetCodeTagCount(),
	)
}

func (d Domain) TabulateHeader(w *tabwriter.Writer) {
	fmt.Fprintf(w, "Domain\tFeed\tArticles\tAbout\tNow\tBlogRoll\tPopular\tTLD\tBlog\tHuman\tNews\tFPR\tWC\tWPB\tGood\tBad\tWPP\tGBR\tGTBP\tAPM\tCodeCn\n")
}
