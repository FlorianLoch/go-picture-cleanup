package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	// Remove timestamp prefix, we do not want this
	log.SetFlags(0)

	noDryMode := flag.Bool("d", false, "If set, this tool will actually delete pictures. Otherwise it will run in dry mode")

	flag.Parse()

	inputPath := flag.Arg(0)
	if inputPath == "" {
		log.Fatal("No input path given")
	}

	absPath, err := filepath.Abs(inputPath)
	if err != nil {
		log.Fatalf("Could not make given path %q absolute: %v", inputPath, err)
	}

	file, err := os.Open(absPath)
	if err != nil {
		log.Fatalf("Could not open given input path %q: %v", absPath, err)
	}

	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Fatalf("Could not stat given input path %q: %v", absPath, err)
	}

	if !stat.IsDir() {
		log.Fatalf("Given input path %q is required to point to a directory.", absPath)
	}

	if !*noDryMode {
		log.Println("Running in dry mode, no files will actually be deleted. Run with '-d' to do so.")
	}

	log.Printf("Going to check %q and subdirectories for unnecessary JPEGs...", absPath)

	var deleter Deleter = &DryRunner{}

	if *noDryMode {
		deleter = &Remover{}
	}

	filesDeleted, bytesDeleted, err := checkDir(absPath, deleter)
	if err != nil {
		log.Fatalf("Failed checking directory: %v", err)
	}

	if *noDryMode {
		log.Printf("Successfully deleted %d files and freed %.2f Mebibytes.", filesDeleted, float64(bytesDeleted)/1024/1024)
	} else {
		log.Printf("Would have deleted %d files and freed %.2f Mebibytes, but was running in dry mode. To actually delete files run with '-d'.", filesDeleted, float64(bytesDeleted)/1024/1024)
	}
}

func checkDir(curDir string, deleter Deleter) (int64, int64, error) {
	var (
		rawLookupTbl = make(map[string]struct{})
		jpegList     []string

		filesDeleted int64
		bytesDeleted int64
	)

	entries, err := os.ReadDir(curDir)
	if err != nil {
		return -1, -1, fmt.Errorf("reading dir %q: %w", curDir, err)
	}

	// 1. step: find JPEGs and RAWs
	for _, entry := range entries {
		name := entry.Name()

		if entry.IsDir() {
			subFilesDeleted, subBytesDeleted, err := checkDir(path.Join(curDir, name), deleter)
			if err != nil {
				return -1, -1, err
			}

			filesDeleted += subFilesDeleted
			bytesDeleted += subBytesDeleted
		}

		if isRAW(name) {
			rawLookupTbl[loweredNoExt(name)] = struct{}{}
		}

		if isJPEG(name) {
			jpegList = append(jpegList, name)
		}
	}

	// 2. step: Find a RAW for every, if there is one delete the JPEG
	for _, jpeg := range jpegList {
		if _, ok := rawLookupTbl[loweredNoExt(jpeg)]; ok {
			fullPath := path.Join(curDir, jpeg)

			stats, err := os.Stat(fullPath)
			if err != nil {
				return -1, -1, fmt.Errorf("stating file %q: %w", fullPath, err)
			}

			if err := deleter.Delete(fullPath); err != nil {
				return -1, -1, fmt.Errorf("deleting file %q: %w", fullPath, err)
			}

			filesDeleted++
			bytesDeleted += stats.Size()
		}
	}

	return filesDeleted, bytesDeleted, nil
}

func loweredNoExt(name string) string {
	lowered := strings.ToLower(name)

	return strings.TrimSuffix(lowered, filepath.Ext(lowered))
}

func isJPEG(name string) bool {
	return hasSuffix(name, ".jpg") || hasSuffix(name, ".jpeg")
}

func isRAW(name string) bool {
	return hasSuffix(name, ".arw") || hasSuffix(name, ".dng")
}

func hasSuffix(name, suffix string) bool {
	return strings.HasSuffix(strings.ToLower(name), strings.ToLower(suffix))
}

type Deleter interface {
	Delete(path string) error
}

type DryRunner struct{}

func (d *DryRunner) Delete(path string) error {
	log.Println(path)

	return nil
}

type Remover struct{}

func (r *Remover) Delete(path string) error {
	return os.Remove(path)
}
