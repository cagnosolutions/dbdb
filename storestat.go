package dbdb

type StoreStat struct {
	Name     string
	Id, Docs uint64
}

type StoreStatSorted []*StoreStat

func (sss StoreStatSorted) Len() int {
	return len(sss)
}

func (sss StoreStatSorted) Less(i, j int) bool {
	return sss[i].Name < sss[j].Name
}

func (sss StoreStatSorted) Swap(i, j int) {
	sss[i], sss[j] = sss[j], sss[i]
}
