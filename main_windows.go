package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"golang.org/x/sys/windows/svc/eventlog"
)

type fileData struct {
	path string
	info os.FileInfo
}

var (
	app     = "Housekeeper"
	version string
	build   string
)

func main() {
	var (
		versionFlg         = flag.Bool("version", false, "Display application version")
		olderThanFlg       = flag.Int("older-than", 0, "Number of days that a file should be older than in order to be deleted")
		extFlg             = flag.String("ext", "", "File extension to be deleted. Use * to match all files")
		pathFlg            = flag.String("path", "", "Path to search for files to be deleted; DO NOT use trailing slashes")
		recursiveFlg       = flag.Bool("recursive", false, "Search all subfolders as well")
		caseInsensitiveFlg = flag.Bool("case-insensitive", false, "Match files regardless of case")
		testFlg            = flag.Bool("test", false, "Test run")
		removeDirsFlg      = flag.Bool("remove-directories", false, "Remove empty directories?")
		debug              = flag.Bool("debug", false, "Enable debugging?")
		logger             log.Logger
		fileInfo           []fileData
		processed          int64
	)

	flag.Parse()

	if *versionFlg {
		fmt.Println(app + " v" + version + " build " + build)
		os.Exit(0)
	}

	if *olderThanFlg == 0 || *extFlg == "" || *pathFlg == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	processed = 0

	logger = log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller, "app", app, "ext", *extFlg, "path", *pathFlg, "version", "v"+version, "build", build, "older-than", *olderThanFlg, "recursive", *recursiveFlg, "test", *testFlg)

	if *debug {
		logger = level.NewFilter(logger, level.AllowDebug())
	} else {
		logger = level.NewFilter(logger, level.AllowInfo())
	}

	var el *eventlog.Log
	el, err := eventlog.Open(app)
	if err != nil {
		msg := fmt.Sprintf("unable to open event log. %s", err)
		level.Error(logger).Log("msg", msg)

		err := eventlog.InstallAsEventCreate(app, eventlog.Error|eventlog.Warning|eventlog.Info)
		if err != nil {
			msg := fmt.Sprintf("unable to create event source. %s", err)
			level.Error(logger).Log("msg", msg)
			os.Exit(1)
		}

		el, err = eventlog.Open(app)
		if err != nil {
			msg := fmt.Sprintf("unable to open event log after creating event source. %s", err)
			level.Error(logger).Log("msg", msg)
			os.Exit(1)
		}
	}

	msg := fmt.Sprintf("Starting %s\nVersion %s\nBuild %s\next: %s\npath: %s\nolder-than: %d\nrecursive: %t\ntest: %t", app, version, build, *extFlg, *pathFlg, *olderThanFlg, *recursiveFlg, *testFlg)
	el.Info(1, msg)

	if _, err := os.Stat(*pathFlg); os.IsNotExist(err) {
		msg := fmt.Sprintf("path, %s, does not exist", *pathFlg)
		level.Error(logger).Log("msg", msg)
		el.Error(2, msg)
		os.Exit(1)
	}

	d := time.Now().AddDate(0, 0, -*olderThanFlg)
	level.Debug(logger).Log("older-than", d)

	ext := "." + strings.Trim(*extFlg, ".")
	level.Debug(logger).Log("extension", ext)

	msg = fmt.Sprintf("delete files with extension %s older than %d days", ext, *olderThanFlg)
	el.Info(3, msg)

	// Build the list of files differently if we're running a recursive search or not
	if !*recursiveFlg {
		files, err := ioutil.ReadDir(*pathFlg)
		if err != nil {
			msg := fmt.Sprintf("unable to read directory; %s", err)
			level.Error(logger).Log("msg", msg)
			el.Error(4, msg)
		}

		for _, file := range files {
			f := fileData{
				path: *pathFlg + "\\" + file.Name(),
				info: file,
			}

			fileInfo = append(fileInfo, f)
		}
	} else {
		filepath.Walk(*pathFlg, func(path string, fi os.FileInfo, err error) error {
			f := fileData{
				path: path,
				info: fi,
			}

			fileInfo = append(fileInfo, f)

			return nil
		})
	}

	// Now process the file list
	for _, file := range fileInfo {
		if file.info.IsDir() {
			continue
		}

		if file.info.ModTime().After(d) {
			continue
		}

		fileExt := filepath.Ext(file.path)

		if *caseInsensitiveFlg {
			fileExt = strings.ToLower(fileExt)
			ext = strings.ToLower(ext)
		}

		if ext != ".*" && fileExt != ext {
			continue
		}

		processed = processed + 1

		if *testFlg {
			level.Info(logger).Log("file", file.path, "msg", "test: would be deleted")

			msg := fmt.Sprintf("testing: %s would be deleted", file.path)
			el.Info(5, msg)

			continue
		}

		err := os.Remove(file.path)
		if err != nil {
			level.Error(logger).Log("file", file.path, "msg", err)

			msg := fmt.Sprintf("unable to delete %s; %s", file.path, err)
			el.Error(6, msg)
		} else {
			level.Info(logger).Log("file", file.path, "msg", "deleted")

			msg := fmt.Sprintf("%s deleted", file.path)
			el.Info(7, msg)
		}
	}

	if *removeDirsFlg && *recursiveFlg {
		for _, file := range fileInfo {
			if !file.info.IsDir() {
				continue
			}

			if file.path == *pathFlg {
				continue
			}

			empty, err := isDirEmpty(file.path)
			if err != nil {
				level.Error(logger).Log("folder", file.path, "msg", err, "task", "is directory empty?")

				msg := fmt.Sprintf("unable to check whether folder, %s, is empty; %s", file.path, err)
				el.Error(8, msg)

				continue
			}

			if empty == false {
				continue
			}

			processed = processed + 1

			if *testFlg {
				level.Info(logger).Log("folder", file.path, "msg", "test: would be deleted")

				msg := fmt.Sprintf("testing: %s would be deleted", file.path)
				el.Info(5, msg)

				continue
			}

			err = os.Remove(file.path)
			if err != nil {
				level.Error(logger).Log("folder", file.path, "msg", err)

				msg := fmt.Sprintf("unable to delete %s; %s", file.path, err)
				el.Error(6, msg)
			} else {
				level.Info(logger).Log("folder", file.path, "msg", "deleted")

				msg := fmt.Sprintf("%s deleted", file.path)
				el.Info(7, msg)
			}
		}
	}

	if processed == 0 {
		level.Info(logger).Log("msg", "no files found to be deleted")
		el.Info(8, "no files found to be deleted")
	}

	err = el.Close()
	if err != nil {
		msg := fmt.Sprintf("unable to close event log. %s", err)
		level.Error(logger).Log("msg", msg)
		os.Exit(1)
	}
}
