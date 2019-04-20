package brbundle

var DefaultRepository, _ = NewRepository()

func SetDecryptoKeyToEmbeddedBundle(name, key string) error {
	return DefaultRepository.SetDecryptoKeyToEmbeddedBundle(name, key)
}

func SetDecryptoKeyToExeBundle(key string) error {
	return DefaultRepository.SetDecryptoKeyToExeBundle(key)
}

func RegisterBundle(path string, option ...Option) error {
	return DefaultRepository.RegisterBundle(path, option...)
}

func RegisterFolder(path string, option ...Option) error {
	return DefaultRepository.RegisterFolder(path, option...)
}

func RegisterEncryptedFolder(path string, option ...Option) error {
	return DefaultRepository.RegisterEncryptedFolder(path, option...)
}

func Unload(name string) error {
	return DefaultRepository.Unload(name)
}

func Find(path string) (FileEntry, error) {
	return DefaultRepository.Find(path)
}