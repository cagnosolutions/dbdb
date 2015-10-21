package dbdb

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func Walk(start string) map[string][]int {
	w := Walker{
		StartDir: start,
		Stores:   make(map[string][]int),
	}
	filepath.Walk(w.StartDir, w.Texas)
	for k := range w.Stores {
		w.Stores[k] = w.Ranger(k)
	}
	return w.Stores
}

type Walker struct {
	StartDir string
	Stores   map[string][]int
}

// walks the db root and gathers all the stores/folders
func (w *Walker) Texas(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.Name() == w.StartDir {
		return nil
	}
	if info.IsDir() {
		w.Stores[info.Name()] = make([]int, 0)
		return filepath.SkipDir
	}
	return nil
}

// takes folder/store as key and walks files/docs...
// returns list of files/docs in this folder/store
func (w *Walker) Ranger(key string) []int {
	var ids []int
	filepath.Walk(w.StartDir+"/"+key, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == key {
			return nil
		}
		if info.IsDir() {
			return filepath.SkipDir
		}
		sid := strings.Split(info.Name(), ".")[0]
		id, _ := strconv.ParseInt(sid, 10, 64)
		ids = append(ids, int(id))
		//ss = append(ss, info.Name())
		return nil
	})
	sort.Ints(ids)
	return ids
}
