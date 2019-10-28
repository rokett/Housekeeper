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
		versionFlg    = flag.Bool("version", false, "Display application version")
		olderThanFlg  = flag.Int("older-than", 0, "Number of days that a file should be older than in order to be deleted")
		extFlg        = flag.String("ext", "", "File extension to be deleted. Use * to match all files")
		pathFlg       = flag.String("path", "", "Path to search for files to be deleted")
		recursiveFlg  = flag.Bool("recursive", false, "Search all subfolders as well")
		testFlg       = flag.Bool("test", false, "Test run")
		removeDirsFlg = flag.Bool("remove-directories", false, "Remove empty directories?")
		debug         = flag.Bool("debug", false, "Enable debugging?")
		logger        log.Logger
		fileInfo      []fileData
		processed     int64
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

	if _, err := os.Stat(*pathFlg); os.IsNotExist(err) {
		msg := fmt.Sprintf("path, %s, does not exist", *pathFlg)
		level.Error(logger).Log("msg", msg)
		os.Exit(1)
	}

	d := time.Now().AddDate(0, 0, -*olderThanFlg)
	level.Debug(logger).Log("older-than", d)

	ext := "." + strings.Trim(*extFlg, ".")
	level.Debug(logger).Log("extension", ext)

	// Build the list of files differently if we're running a recursive search or not
	if !*recursiveFlg {
		files, err := ioutil.ReadDir(*pathFlg)
		if err != nil {
			msg := fmt.Sprintf("unable to read directory; %s", err)
			level.Error(logger).Log("msg", msg)
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

			continue
		}

		err := os.Remove(file.path)
		if err != nil {
			level.Error(logger).Log("file", file.path, "msg", err)
		} else {
			level.Info(logger).Log("file", file.path, "msg", "deleted")
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

			// There is a good chance that a subfolder has already been deleted as the slice of folders is listed from the root down.
			// So we may well have deleted the parent of a subfolder, because all of its subfolders were empty, before we get to check the subfolder.
			// Checking whether the folder we want to act on already exists or not, removes the possibility of an error.
			if _, err := os.Stat(file.path); os.IsNotExist(err) {
				continue
			}

			empty, err := isDirEmpty(file.path)
			if err != nil {
				level.Error(logger).Log("folder", file.path, "msg", err, "task", "is directory empty?")

				continue
			}

			if empty == false {
				continue
			}

			processed = processed + 1

			if *testFlg {
				level.Info(logger).Log("folder", file.path, "msg", "test: would be deleted")

				continue
			}

			err = os.RemoveAll(file.path)
			if err != nil {
				level.Error(logger).Log("folder", file.path, "msg", err)
			} else {
				level.Info(logger).Log("folder", file.path, "msg", "deleted")
			}
		}
	}

	if processed == 0 {
		level.Info(logger).Log("msg", "no files found to be deleted")
	}
}
