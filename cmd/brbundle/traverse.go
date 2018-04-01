package main

import (
	"os"
	"path/filepath"

	"github.com/monochromegane/go-gitignore"
)

const ignoreSettingFile = ".bundleignore"

type Entry struct {
	Path string
	Info os.FileInfo
}

func Traverse(srcDirPath string) (entries chan Entry, dirs []Entry, ignored []string) {
	ignoreMatcher, _ := gitignore.NewGitIgnore(filepath.Join(srcDirPath, ignoreSettingFile), srcDirPath)
	var paths []string
	var infos []os.FileInfo
	filepath.Walk(srcDirPath,
		func(path string, info os.FileInfo, err error) error {
			rel, err := filepath.Rel(srcDirPath, path)
			if err != nil {
				return err
			}
			if path == "." {
				return nil
			}
			if ignoreMatcher != nil && ignoreMatcher.Match(rel, info.IsDir()) {
				ignored = append(ignored, rel)
				return nil
			}
			if !info.IsDir() && rel != ignoreSettingFile {
				paths = append(paths, rel)
				infos = append(infos, info)
			}
			return nil
		})
	dirMap := make(map[string]bool)
	entries = make(chan Entry, len(paths))
	for i, path := range paths {
		entries <- Entry{path, infos[i]}
		dir := filepath.Dir(path)
		if dir != "." && !dirMap[dir] {
			dirMap[dir] = true
		}
	}
	for dir := range dirMap {
		stat, _ := os.Stat(filepath.Join(srcDirPath, dir))
		dirs = append(dirs, Entry{dir, stat})
	}
	return
}
