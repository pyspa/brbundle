package brbundle

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
)

type ManifestEntry struct {
	File string `json:"file"`
	Sha1 string `json:"sha1"`
	Size int    `json:"size"`
}

type Progress struct {
	folders        map[string]string
	shouldDownload chan *ManifestEntry
	downlodFiles   []*ManifestEntry
	deleteFiles    []*ManifestEntry
	keepFiles      []*ManifestEntry
	errors         []error
	wait           chan struct{}
}

func (p Progress) DownloadFiles() []string {
	result := make([]string, len(p.downlodFiles))
	for i, entry := range p.downlodFiles {
		result[i] = entry.File + ".pb"
	}
	return result
}

func (p Progress) DeleteFiles() []string {
	result := make([]string, len(p.deleteFiles))
	for i, entry := range p.deleteFiles {
		result[i] = entry.File + ".pb"
	}
	return result
}

func (p Progress) KeepFiles() []string {
	result := make([]string, len(p.keepFiles))
	for i, entry := range p.keepFiles {
		result[i] = entry.File + ".pb"
	}
	return result
}

func (p *Progress) Wait() error {
	<-p.wait
	return nil
}

func (p *Progress) startDownload(r *Repository, b *manifestBundle, manifestUrl, parentFolderPath string, jsonContent []byte, parallelDownload int) {
	tasks := make(chan *ManifestEntry)
	go func() {
		for _, downloadFile := range p.downlodFiles {
			tasks <- downloadFile
		}
		close(tasks)
	}()

	var wg sync.WaitGroup
	var lock sync.RWMutex

	wg.Add(parallelDownload)
	for i := 0; i < parallelDownload; i++ {
		go func() {
			for downloadFile := range tasks {
				r, err := http.Get(manifestUrl + downloadFile.File + ".pb")
				if err != nil {
					lock.Lock()
					defer lock.Unlock()
					p.errors = append(p.errors, err)
					continue
				}
				defer r.Body.Close()
				f, err := os.Create(filepath.Join(parentFolderPath, downloadFile.File+".pb"))
				if err != nil {
					lock.Lock()
					defer lock.Unlock()
					p.errors = append(p.errors, err)
					continue
				}
				defer f.Close()
				io.Copy(f, r.Body)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	for _, entry := range p.deleteFiles {
		os.Remove(filepath.Join(parentFolderPath, entry.File+".pb"))
	}
	ioutil.WriteFile(filepath.Join(parentFolderPath, "manifest.json"), jsonContent, 0644)
	r.bundles[ManifestBundleType] = append(r.bundles[ManifestBundleType], b)
	close(p.wait)
}

type manifestBundle struct {
	baseBundle
	rootFolder string
	folders    map[string]string
}

func newManifestBundle(parentFolderPath string, o Option, files map[string]*ManifestEntry) *manifestBundle {
	mountPoint := o.MountPoint
	if mountPoint != "" && !strings.HasSuffix(mountPoint, "/") {
		mountPoint = mountPoint + "/"
	}

	bundle := &manifestBundle{
		baseBundle: baseBundle{
			mountPoint:    mountPoint,
			name:          o.Name,
			decryptorType: NotToEncrypto,
		},
		folders: make(map[string]string),
	}
	if o.DecryptoKey != "" {
		bundle.baseBundle.decryptorType = UseAES
		bundle.baseBundle.setDecryptionKey(o.DecryptoKey)
	}
	for folder, file := range files {
		bundle.folders[folder] = filepath.Join(parentFolderPath, file.File+".pb")
	}
	return bundle
}

func (m manifestBundle) find(path string) (FileEntry, error) {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	if path, ok := m.folders[dir]; ok {
		reader, err := zip.OpenReader(path)
		reader.RegisterDecompressor(ZIPMethodSnappy, snappyDecompressor)
		if err != nil {
			return nil, err
		}
		for _, file := range reader.File {
			if file.Name == base {
				entry, err := newPackedFileEntry(file, dir, &m.baseBundle)
				runtime.SetFinalizer(entry, func(*packedFileEntry) {
					reader.Close()
				})
				return entry, err
			}
		}
	}
	return nil, nil
}

func (manifestBundle) close() {
}

func (m manifestBundle) dirs() []string {
	dirNames := make([]string, len(m.folders))
	i := 0
	for name := range m.folders {
		dirNames[i] = m.mountPoint + name
		i++
	}
	sort.Strings(dirNames)
	return dirNames
}

func (m manifestBundle) filesInDir(dirName string) []string {
	if !strings.HasPrefix(dirName, m.mountPoint) {
		return nil
	}
	dirName = dirName[len(m.mountPoint) : len(dirName)-1]
	folder, ok := m.folders[dirName]
	if !ok {
		return nil
	}
	reader, err := zip.OpenReader(folder)
	defer reader.Close()
	if err != nil {
		return nil
	}
	var fileNames []string
	for _, f := range reader.File {
		fileNames = append(fileNames, f.Name)
	}
	sort.Strings(fileNames)

	return fileNames
}
