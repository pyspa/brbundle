package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"go.pyspa.org/brbundle"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

func manifest(brotli bool, encryptionKey []byte, buildTag, destDirPath, srcDirPath string, date *time.Time) error {
	ignoreMatcher := findGitIgnore(srcDirPath)
	manifestFile, err := os.Create(filepath.Join(destDirPath, "manifest.json"))
	if err != nil {
		return err
	}
	defer manifestFile.Close()
	var targetFolders []string

	filepath.Walk(srcDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcDirPath, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if info.IsDir() {
			if ignoreMatcher != nil && ignoreMatcher.Match(path, true) {
				return filepath.SkipDir
			}
			entries, err := ioutil.ReadDir(path)
			if err != nil {
				return err
			}
			hasFile := false
			for _, entry := range entries {
				if !entry.IsDir() && !ignoreMatcher.Match(filepath.Join(path, entry.Name()), false) {
					hasFile = true
					break
				}
			}
			if hasFile {
				targetFolders = append(targetFolders, path)
			}
		}
		return nil
	})

	targetFolderChan := make(chan string, 100)
	var wg sync.WaitGroup
	wg.Add(len(targetFolders))
	go func() {
		for _, folder := range targetFolders {
			targetFolderChan <- folder
		}
	}()

	manifestSrc := make(map[string]*brbundle.ManifestEntry)

	var lock sync.RWMutex

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for targetFolder := range targetFolderChan {
				packSubBundle(brotli, encryptionKey, buildTag, destDirPath, srcDirPath, targetFolder, date, manifestSrc, &lock)
				wg.Done()
			}
		}()
	}
	wg.Wait()
	close(targetFolderChan)
	encoder := json.NewEncoder(manifestFile)
	encoder.SetIndent("", "    ")
	encoder.Encode(manifestSrc)
	return nil
}

func packSubBundle(brotli bool, encryptionKey []byte, buildTag, destDirPath, srcDirPath, targetFolder string, date *time.Time, result map[string]*brbundle.ManifestEntry, lock *sync.RWMutex) error {
	rel, err := filepath.Rel(srcDirPath, targetFolder)
	if err != nil {
		return err
	}
	h := md5.New()
	io.WriteString(h, cleanPath("", rel))
	fmt.Println("@@", cleanPath("", rel))
	fileName := fmt.Sprintf("%x", h.Sum(nil))

	out, err := os.Create(filepath.Join(destDirPath, fileName+".pb"))
	if err != nil {
		return err
	}
	hash := sha1.New()
	mode := "Packed Bundle for manifest (" + cleanPath("", rel) + ")"
	packedBundleShallow(brotli, encryptionKey, buildTag, io.MultiWriter(out, hash), targetFolder, "", mode, date)
	out.Close()

	lock.Lock()
	defer lock.Unlock()
	result[cleanPath("", rel)] = &brbundle.ManifestEntry{
		File: fileName,
		Sha1: fmt.Sprintf("%x", hash.Sum(nil)),
	}

	return nil
}
