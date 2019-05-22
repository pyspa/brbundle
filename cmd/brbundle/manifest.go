package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/monochromegane/go-gitignore"
)

type ManifestEntry struct {
	File string `json:"file"`
	Sha1 string `json:"sha1"`
}

func manifest(brotli bool, encryptionKey []byte, buildTag, destDirPath, srcDirPath string, date *time.Time) error {
	ignoreMatcher, _ := gitignore.NewGitIgnore(filepath.Join(srcDirPath, ignoreSettingFile), srcDirPath)
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
			dirs, err := ioutil.ReadDir(path)
			if err != nil {
				return err
			}
			hasFile := false
			for _, dir := range dirs {
				if !dir.IsDir() {
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

	manifestSrc := make(map[string]*ManifestEntry)

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

func packSubBundle(brotli bool, encryptionKey []byte, buildTag, destDirPath, srcDirPath, targetFolder string, date *time.Time, result map[string]*ManifestEntry, lock *sync.RWMutex) error {
	rel, err := filepath.Rel(srcDirPath, targetFolder)
	if err != nil {
		return err
	}
	h := md5.New()
	io.WriteString(h, cleanPath("", rel))
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
	result[cleanPath("", rel)] = &ManifestEntry{
		File: fileName,
		Sha1: fmt.Sprintf("%x", hash.Sum(nil)),
	}

	return nil
}
