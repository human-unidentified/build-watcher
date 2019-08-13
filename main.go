package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

// Для сортировки директорий по дате\времени создания
type byCreateDate []os.FileInfo

func (s byCreateDate) Len() int {
	return len(s)
}
func (s byCreateDate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byCreateDate) Less(i, j int) bool {
	diff := s[i].ModTime().Sub(s[j].ModTime())
	return diff < 0
}

// Получение названия последней директории по указанному пути
func getLastDir(path string) (lastDir string) {
	allFiles, _ := ioutil.ReadDir(path)
	sort.Sort(byCreateDate(allFiles))
	dirs := make([]string, 0)

	for _, f := range allFiles {
		if f.IsDir() {
			dirs = append(dirs, f.Name())
		}
	}

	return dirs[len(dirs)-1]
}

func watchDir(path string) {
	fmt.Println("Start watch folder ", path)
	watcherInterval := time.Second * 5
	for true {
		lastDir := getLastDir(path)
		if lastDir != "" {
			fmt.Println("lastDir=", lastDir)
		}

		fmt.Printf("Sleep %s\n", watcherInterval)
		time.Sleep(watcherInterval)
	}

}

func main() {
	watchDir("\\\\s6\\BuildArchive\\T-FLEX DOCs 17\\DOCsDev\\")
}
