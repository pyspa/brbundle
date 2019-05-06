package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/monochromegane/go-gitignore"
)

const ignoreSettingFile = ".bundleignore"

type Entry struct {
	Path     string
	DestPath string
	Info     os.FileInfo
}

func splitBuildTag(name, buildTag string) (psuedoName string, match bool) {
	splitNames := strings.Split(name, "__")
	if len(splitNames) == 2 && splitNames[0] != "" {
		targetFileBuildTag := splitNames[1][:len(splitNames[1])-len(filepath.Ext(splitNames[1]))]
		if targetFileBuildTag == buildTag {
			psuedoName = splitNames[0] + filepath.Ext(splitNames[1])
			match = true
		} else {
			psuedoName = ""
			match = false
		}
	} else {
		psuedoName = name
		match = true
	}
	return
}

func cleanPseudoPath(source, buildTag string) string {
	fragments := strings.Split(source, string(os.PathSeparator))
	newFragments := make([]string, len(fragments))
	for i, fragment := range fragments {
		newFragments[i], _ = splitBuildTag(fragment, buildTag)
	}
	return strings.Join(newFragments, string(os.PathSeparator))
}

func Traverse(srcDirPath, buildTag string) (entries chan Entry, dirs []Entry, ignored []string) {
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
			if _, match := splitBuildTag(info.Name(), buildTag); !match {
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
		entries <- Entry{path, cleanPseudoPath(path, buildTag), infos[i]}
		dir := filepath.Dir(path)
		if dir != "." && !dirMap[dir] {
			dirMap[dir] = true
		}
	}
	for dir := range dirMap {
		stat, _ := os.Stat(filepath.Join(srcDirPath, dir))
		dirs = append(dirs, Entry{dir, cleanPseudoPath(dir, buildTag), stat})
	}
	return
}
