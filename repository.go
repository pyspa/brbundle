package brbundle

import (
	"archive/zip"
	"bytes"
	"fmt"
	"strings"

	"github.com/hashicorp/golang-lru"
)

type BundleType int

const (
	FolderBundleType   BundleType = 0
	PackedBundleType              = 1
	ExeBundleType                 = 2
	EmbeddedBundleType            = 3
)

type Repository struct {
	option  *ROption
	init    bool
	Cache   *lru.TwoQueueCache
	bundles [4][]bundle
}

type ROption struct {
	OmitExeBundle          bool
	OmitEmbeddedBundle     bool
	OmitEnvVarFolderBundle bool
}

func NewRepository(option ...ROption) *Repository {
	rOption := &ROption{}
	if len(option) > 0 {
		rOption = &option[0]
	}
	r := &Repository{
		option: rOption,
	}
	return r
}

func (r *Repository) SetCacheSize(size int) error {
	cache, err := lru.New2Q(size)
	if err != nil {
		return err
	}
	r.Cache = cache
	return nil
}

func (r *Repository) lazyInit() {
	if r.init {
		return
	}
	if !r.option.OmitEmbeddedBundle {
		r.initEmbeddedBundle()
	}
	if !r.option.OmitExeBundle {
		r.initExeBundle()
	}
	if !r.option.OmitEnvVarFolderBundle {
		r.initFolderBundleByEnvVar()
	}
	r.init = true
}

func (r *Repository) setDecryptoKey(name, key string, bundleType BundleType) error {
	r.lazyInit()
	found := false
	for _, bundle := range r.bundles[EmbeddedBundleType] {
		if bundle.getName() == name {
			bundle.setDecryptionKey(key)
			found = true
		}
	}
	if !found {
		return fmt.Errorf("name '%s' is not found", name)
	}
	return nil
}

func (r *Repository) ClearCache() {
	if r.Cache != nil {
		r.Cache.Purge()
	}
}

func (r *Repository) initEmbeddedBundle() {
	r.bundles[EmbeddedBundleType] = make([]bundle, len(embeddedBundles))
	for i, e := range embeddedBundles {
		reader, err := zip.NewReader(bytes.NewReader(e.data), int64(len(e.data)))
		if err != nil {
			panic(err)
		}
		r.bundles[EmbeddedBundleType][i] = newPackedBundle(
			reader, nil, Option{
				Name: e.name,
			})
	}
}

func (r *Repository) SetDecryptoKeyToEmbeddedBundle(name, key string) error {
	return r.setDecryptoKey(name, key, EmbeddedBundleType)
}

func (r *Repository) SetDecryptoKeyToExeBundle(key string) error {
	return r.setDecryptoKey("", key, ExeBundleType)
}

type Option struct {
	DecryptoKey string
	MountPoint  string
	Name        string
	Priority    int
}

func (r *Repository) RegisterBundle(path string, option ...Option) error {
	var bo Option
	if len(option) > 0 {
		bo = option[0]
	}
	if bo.Name == "" {
		bo.Name = path
	}
	return r.registerBundle(path, bo)
}

func (r *Repository) RegisterFolder(path string, option ...Option) error {
	var bo Option
	if len(option) > 0 {
		bo = option[0]
	}
	if bo.Name == "" {
		bo.Name = path
	}
	return r.registerFolder(path, false, bo)
}

func (r *Repository) RegisterEncryptedFolder(path string, option ...Option) error {
	var bo Option
	if len(option) > 0 {
		bo = option[0]
	}
	if bo.Name == "" {
		bo.Name = path
	}
	return r.registerFolder(path, true, bo)
}

func (r *Repository) Unload(name string) error {
	bundles := r.bundles[PackedBundleType]
	for i, bundle := range bundles {
		if bundle.getName() == name {
			r.bundles[PackedBundleType] = append(bundles[:i], bundles[(i+1):]...)
			bundle.close()
			if r.Cache != nil {
				keys := r.Cache.Keys()
				for _, key := range keys {
					b, _ := r.Cache.Get(key)
					if b == bundle {
						r.Cache.Remove(key)
					}
				}
			}
			break
		}
	}
	return fmt.Errorf("PackedBundle '%s' is not found", name)
}

func (r *Repository) Find(path string) (FileEntry, error) {
	r.lazyInit()
	if r.Cache != nil {
		cachedBundle, ok := r.Cache.Get(path)
		if ok {
			return cachedBundle.(bundle).find(path)
		}
	}
	for _, bundles := range r.bundles {
		for _, bundle := range bundles {
			relPath := path
			mountPoint := bundle.getMountPoint()
			if bundle.getMountPoint() != "" {
				if !strings.HasPrefix(path, mountPoint) {
					continue
				}
				relPath = path[len(mountPoint):]
			}
			fileEntry, err := bundle.find(relPath)
			if err != nil {
				continue
			}
			if fileEntry != nil {
				if r.Cache != nil {
					r.Cache.Add(path, bundle)
				}
				return fileEntry, nil
			}
		}
	}
	return nil, fmt.Errorf("Asset '%s' is not in bundles", path)
}
