package main

import (
	"ComicFileTools/src/config"
	"ComicFileTools/src/tools"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type context struct {
	fileOpMod    fileOpMod
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

type fileOpMod int

const (
	move fileOpMod = 1
	copy fileOpMod = 2
)

func main() {
	// mod
	switch config.Conf.RunMod {
	case config.Test:
		fmt.Printf("[%v]\n", strings.ToUpper(string(config.Test)))
		ExecuteTest()
	case config.Move:
		fmt.Printf("[%v]\n", strings.ToUpper(string(config.Move)))
		ExecuteRun(move)
	case config.Copy:
		fmt.Printf("[%v]\n", strings.ToUpper(string(config.Copy)))
		ExecuteRun(copy)
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

func ExecuteRun(mod fileOpMod) {
	// context
	c := &context{
		fileOpMod: mod,
		meta:      NewMetaGetter(config.Conf.Pattern),
		gMap:      make(map[string]string),
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
	// extract meta
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

	if c.fileOpMod == move {
		moveFile(c, oldPath, newPath)
	}

	if c.fileOpMod == copy {
		copyFile(c, oldPath, newPath)
	}
}

func moveFile(c *context, s string, d string) {
	// move mod
	var dstEx bool
	var err error
	if dstEx, err = tools.FileExists(d); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[error] check target file error: '%v', error: %v\n", d, err)
		c.error()
		return
	}
	if dstEx {
		_, _ = fmt.Fprintf(os.Stderr, "[warning] file already exists, skip: '%v'\n", d)
		c.skip()
		return
	}

	if err := os.Rename(s, d); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[error] move file '%v' to '%v' error: %v\n", s, d, err)
		c.error()
		return
	}

	fmt.Printf("[success] move file '%v' to '%v' done\n", s, d)
	c.success()
}

func copyFile(c *context, s string, d string) {
	// copy mod
	var src *os.File
	var dst *os.File
	var err error

	// prepare source file
	if src, err = os.Open(s); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[error] open source file error: '%v', error: %v\n", s, err)
		c.error()
		return
	}
	defer func(src *os.File) {
		_ = src.Close()
	}(src)

	// prepare dest file
	var dstEx bool
	if dstEx, err = tools.FileExists(d); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[error] check new file error: '%v', error: %v\n", d, err)
		c.error()
		return
	}
	if dstEx {
		_, _ = fmt.Fprintf(os.Stderr, "[warning] file already exists, skip: '%v'\n", d)
		c.skip()
		return
	}
	if dst, err = os.Create(d); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[error] create new file error: '%v', error: %v\n", d, err)
		c.error()
		return
	}
	defer func(dst *os.File) {
		_ = dst.Close()
	}(dst)

	// copy
	var written int64
	if written, err = io.Copy(dst, src); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[error] copy file '%v' to '%v' error: %v\n", s, d, err)
		c.error()
		return
	}

	fmt.Printf("[success] copy file '%v' to '%v' done, %v bytes copied\n", s, d, written)
	c.success()
}
