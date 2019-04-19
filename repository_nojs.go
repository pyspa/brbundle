// +build !js

package brbundle

import (
	"archive/zip"
	"fmt"
	"github.com/shibukawa/zipsection"
	"os"
)

func (r *Repository) registerBundle(path string, option ...Option) error {
	var bo Option
	if len(option) > 0 {
		bo = option[0]
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	i, _ := f.Stat()
	reader, err := zip.NewReader(f, i.Size())
	if err != nil {
		fmt.Println(err)
		return err
	}
	b := newPackedBundle(reader, f, bo)
	r.bundles[PackedBundleType] = append(r.bundles[PackedBundleType], b)
	return nil
}

func (r *Repository) registerFolder(path string, option ...Option) error {
	return nil
}

func (r *Repository) initFolderBundleByEnvVar() {
}

func (r *Repository) initExeBundle() error {
	var err error
	filepath, err := os.Executable()
	if err != nil {
		return err
	}
	reader, closer, err := zipsection.Open(filepath)
	b := newPackedBundle(reader, closer, Option{})
	r.bundles[ExeBundleType] = append(r.bundles[ExeBundleType], b)
	return nil
}
