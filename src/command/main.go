package main

import (
	"ComicFileTools/src/config"
	"ComicFileTools/src/tools"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type context struct {
	meta         *MetaGetter
	gMap         map[string]string
	successCount int
	skipCount    int
	errorCount   int
}

func (c *context) success() {
	c.successCount++
}

func (c *context) skip() {
	c.skipCount++
}

func (c *context) error() {
	c.errorCount++
}

func main() {
	// mod
	switch config.Conf.RunMod {
	case config.Test:
		fmt.Printf("[%v]\n", strings.ToUpper(string(config.Test)))
		ExecuteTest()
	case config.Run:
		fmt.Printf("[%v]\n", strings.ToUpper(string(config.Run)))
		ExecuteRun()
	}
}

func ExecuteTest() {
	meta := NewMetaGetter(config.Conf.Pattern)
	group, name, err := meta.Match(config.Conf.TestFileName)
	if err != nil {
		panic(err)
	}

	fmt.Printf("filename: %v\n", config.Conf.TestFileName)
	fmt.Printf("group: %v\n", group)
	fmt.Printf("name: %v\n", name)
}

func ExecuteRun() {
	// context
	c := &context{
		meta: NewMetaGetter(config.Conf.Pattern),
		gMap: make(map[string]string),
	}

	// process
	fmt.Printf("----------START----------\n")
	processDir(c, tools.PathWithSeparator(config.Conf.InputDir))
	fmt.Printf("----------END----------\n")

	// show
	fmt.Printf("success count: %v\n", c.successCount)
	fmt.Printf("skip count: %v\n", c.skipCount)
	fmt.Printf("error count: %v\n", c.errorCount)
}

func processDir(c *context, currentDir string) {
	var files []fs.FileInfo
	var err error
	if files, err = ioutil.ReadDir(currentDir); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[error] list directory error, directory: '%v', error: %v\n", currentDir, err)
		c.error()
		return
	}

	for _, f := range files {
		if f.IsDir() {
			processDir(c, path.Join(currentDir, f.Name()))
		} else {
			processFile(c, currentDir, f)
		}
	}
}

func processFile(c *context, currentDir string, f fs.FileInfo) {
	fn := f.Name()
	group, name, err := c.meta.Match(fn)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[warning] skip file for file name not matching: '%v'\n", fn)
		c.skip()
		return
	}

	// prepare group directory
	gDir, exists := c.gMap[group]
	if !exists {
		// init group dir
		p := path.Join(config.Conf.OutputDir, group)
		gDir, err = tools.CreateDirectoryIfNotExists(p)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "[error] creating group directory error, directory: '%v', error: %v\n", p, err)
			c.error()
			return
		}

		c.gMap[group] = gDir
	}

	oldPath := path.Join(currentDir, fn)
	newPath := path.Join(gDir, name)
	if err = os.Rename(oldPath, newPath); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[error] move file '%v' to '%v' error: %v\n", oldPath, newPath, err)
		c.error()
		return
	}

	fmt.Printf("[success] move file '%v' to '%v'\n", oldPath, newPath)
	c.success()
}
