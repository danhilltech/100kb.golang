package article

type ByRecency []*Article

func (a ByRecency) Len() int           { return len(a) }
func (a ByRecency) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRecency) Less(i, j int) bool { return a[i].PublishedAt < a[j].PublishedAt }
