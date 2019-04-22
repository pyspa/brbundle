package brbundle

func (r *Repository) registerBundle(path string, option Option) error {
	return nil
}

func (r *Repository) registerFolder(path string, encrypted bool, option Option) error {
	return fmt.Errorf("Gopher.js doens't support folder bundle")
}

func (r *Repository) initFolderBundleByEnvVar() error {
	return nil
}

func (r *Repository) initExeBundle() error {
	return nil
}
