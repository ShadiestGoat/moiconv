package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/integrii/flaggy"
)

const helpTemplate = `{{.CommandName}}{{if .Description}} - {{.Description}}{{end}}{{if .PrependMessage}}
{{.PrependMessage}}{{end}}

  Usage: moiconv [...directories]
{{if (gt (len .Flags) 0)}}
  Flags: {{if .Flags}}{{range .Flags}}
    {{if .ShortName}}-{{.ShortName}} {{else}}   {{end}}{{if .LongName}}--{{.LongName}}{{end}}{{if .Description}}   {{.Spacer}}{{.Description}}{{if .DefaultValue}} (default: {{.DefaultValue}}){{end}}{{end}}{{end}}{{end}}
{{end}}{{if .AppendMessage}}{{.AppendMessage}}
{{end}}{{if .Message}}
{{.Message}}{{end}}
`

var (
	FLAG_OUT string
	FLAG_FMT = "mov"
	FLAG_FLAT = false
	FLAG_RECURSIVE = false
)

func init() {
	flaggy.String(&FLAG_OUT, "o", "output", "The output director")
	flaggy.String(&FLAG_FMT, "t", "format", "The file type of the output(s)")
	flaggy.Bool(&FLAG_FLAT, "f", "flat", "Whether to flatten the output to 1 directory or not")
	flaggy.Bool(&FLAG_RECURSIVE, "r", "recursive", "Whether to go into subdirectories or nt")
	
	flaggy.SetName("MOI Converter")
	flaggy.SetDescription("Convert MOD+MOI data into a more modern type, adding EXIF data :3")
	flaggy.ShowHelpOnUnexpectedDisable()
	flaggy.DefaultParser.DisableShowVersionWithVersion()
	flaggy.DefaultParser.SetHelpTemplate(helpTemplate)

	flaggy.Parse()
}

func main() {
	if len(flaggy.TrailingArguments) == 0 {
		flaggy.ShowHelpAndExit("==> No directories provided, please see usage")
	}
	if FLAG_OUT == "" {
		flaggy.ShowHelpAndExit("==> No output directory provided, please see usage")
	}
	if strings.ToLower(FLAG_FMT) == "mod" {
		panic("I'm not converting mod to mod >:(")
	}

	requireBin("ffmpeg")
	requireBin("exiftool")

	s, err := os.Stat(FLAG_OUT)
	if err != nil {
		panic("Can't look at the output dir: " + err.Error())
	}
	if !s.IsDir() {
		panic("Output must be a directory")
	}

	for _, d := range flaggy.TrailingArguments {
		resolvedDir, err := filepath.EvalSymlinks(d)
		if err != nil {
			panic("Can't resolve path '" + d + "': " + err.Error())
		}

		abs, err := filepath.Abs(resolvedDir)
		if err != nil {
			panic("Can't resolve path '" + d + "': " + err.Error())
		}

		findAllFiles(abs, nil)
	}
}

var scannedDirs = map[string]bool{}

func findAllFiles(topDir string, pathToHere []string) {
	curDir := filepath.Join(append([]string{topDir}, pathToHere...)...)

	if scannedDirs[curDir] {
		return
	}
	scannedDirs[curDir] = true

	files, err := os.ReadDir(curDir)
	if err != nil {
		panic("Can't read dir '" + curDir + "': " + err.Error())
	}

	dstDir := FLAG_OUT
	if !FLAG_FLAT {
		dstDir = filepath.Join(append([]string{dstDir}, pathToHere...)...)
	}

	for _, f := range files {
		n := f.Name()
		if f.IsDir() {
			if FLAG_RECURSIVE {
				findAllFiles(topDir, append(pathToHere, n))
			}
			continue
		}

		if !strings.HasSuffix(n, ".MOD") {
			continue
		}

		moiFile, err := os.OpenFile(filepath.Join(curDir, n[:len(n) - 3] + "MOI"), os.O_RDONLY, 0755)
		if err != nil {
			fmt.Println(filepath.Join(curDir, n) + ": no MOI file :/")
			// Just in case
			moiFile = nil
		}

		dst := n[:len(n) - 4]
		s := 0
		for {
			tmp := dst
			if s != 0 {
				tmp += "-" + strconv.Itoa(s)
			}
			tmp = filepath.Join(dstDir, tmp + "." + FLAG_FMT)

			_, err := os.Stat(tmp)
			if err != nil && errors.Is(err, os.ErrNotExist) {
				dst = tmp
				break
			}
			
			s++
		}

		convertFile(
			filepath.Join(curDir, n),
			dst,
			moiFile,
		)
	}
}

func convertFile(ogFile, dstFile string, moiFile *os.File) {
	fmt.Printf("File:\t\t%v\n", ogFile)

	err := os.MkdirAll(filepath.Dir(dstFile), 0755)
	if err != nil {
		panic("Can't create dir: " + err.Error())
	}

	err = runCmd(
		"ffmpeg",
		"-i", ogFile,
		"-vcodec", "copy", "-acodec", "aac",
		"-v", "quiet",
		dstFile,
	)
	if err != nil {
		panic("Failed to do conversion; See earlier failure message. " + err.Error())
	}
	if moiFile == nil {
		return
	}

	t, d := getTimestampAndDuration(moiFile)
	fmt.Printf("Time:\t\t%v\n", t)
	fmt.Printf("Duration:\t%v\n", d)

	timeFmt := t.Format(time.DateTime)

	err = runCmd(
		"exiftool",
		"-AllDates=" + timeFmt,
		"-Track*Date=" + timeFmt,
		"-Media*Date=" + timeFmt,
		"-overwrite_original",
		dstFile,
	)
	if err != nil {
		panic("Failed to apply metadata; See earlier failure message. " + err.Error())
	}

	fmt.Println("Success!")
	fmt.Println("<<=====||=====>>")
}
