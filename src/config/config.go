package config

import (
	"ComicFileTools/src/tools"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type RunMod string

type Config struct {
	RunMod       RunMod // 运行模式
	TestFileName string // 测试匹配模式
	Pattern      string // 匹配正则表达式
	InputDir     string // 输入文件夹
	OutputDir    string // 输出文件夹
}

func (c *Config) Print() {
	fmt.Println("[Config]")
	fmt.Println("Run mode: " + c.RunMod)
	fmt.Println("Test file name: " + c.TestFileName)
	fmt.Println("Pattern: " + c.Pattern)
	fmt.Println("Input directory: " + c.InputDir)
	fmt.Println("Output directory: " + c.OutputDir)
}

const (
	Test = RunMod("test")
	Move = RunMod("move")
	Copy = RunMod("copy")
)

var Conf = &Config{}
var ValidMods = []RunMod{Test, Move, Copy}

func validateTestConf() {
	// pattern
	validatePattern()

	// test fileName
	if Conf.TestFileName == "" {
		errorAndExit("test file name can not be empty")
	}
}

func validateRunConf() {
	// pattern
	validatePattern()

	// int dir
	inputDir := Conf.InputDir
	if inputDir == "" {
		errorAndExit("input directory can not be empty")
	}
	if exists, _ := tools.FileExists(inputDir); !exists {
		errorAndExit("input directory '%v' not exists", inputDir)
	}
	if isDir, err := tools.IsDirectory(inputDir); err != nil || !isDir {
		errorAndExit("input directory '%v' is not a directory", inputDir)
	}

	// output dir
	outputDir := Conf.OutputDir
	if outputDir == "" {
		errorAndExit("output directory can not be empty")
	}
	if outputDir == inputDir {
		errorAndExit("output and input directory can not be same")
	}
	if strings.Contains(tools.PathWithSeparator(outputDir), tools.PathWithSeparator(inputDir)) {
		errorAndExit("output directory can not be within input directory")
	}
	if _, err := tools.CreateDirectoryIfNotExists(outputDir); err != nil {
		errorAndExit("create output directory error: %v", err)
	}
}

func validatePattern() {
	if Conf.Pattern == "" {
		errorAndExit("pattern can not be empty")
	}
	if _, err := regexp.Compile(Conf.Pattern); err != nil {
		errorAndExit("invalid regexp pattern '%v'", Conf.Pattern)
	}
}

func errorAndExit(s string, any ...any) {
	_, _ = fmt.Fprintf(os.Stderr, "error: "+s+"\n", any...)
	os.Exit(1)
}

func init() {
	var mod string
	flag.StringVar(&mod, "m", "", "run mode")
	flag.StringVar(&Conf.TestFileName, "t", "", "test filename")
	flag.StringVar(&Conf.Pattern, "p", "", "match pattern")
	flag.StringVar(&Conf.InputDir, "i", "", "input directory")
	flag.StringVar(&Conf.OutputDir, "o", "", "output directory")
	flag.Parse()

	// get run mode
	if mod == "" {
		errorAndExit("run mode can not be empty")
	}
	for _, vm := range ValidMods {
		if RunMod(strings.ToLower(mod)) == vm {
			Conf.RunMod = vm
			break
		}
	}
	if Conf.RunMod == "" {
		errorAndExit("invalid run mod '%v'", mod)
	}

	// validate
	switch Conf.RunMod {
	case Test:
		validateTestConf()
	case Move:
		fallthrough
	case Copy:
		validateRunConf()
	}

	Conf.Print()
}
