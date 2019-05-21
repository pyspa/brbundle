package brbundle

// DefaultRepository is a default repository instance
var DefaultRepository = NewRepository()

// SetDecryptoKeyToEmbeddedBundle registers decrypto key for embedded encrypted assets.
func SetDecryptoKeyToEmbeddedBundle(name, key string) error {
	return DefaultRepository.SetDecryptoKeyToEmbeddedBundle(name, key)
}

// SetDecryptoKeyToExeBundle registers decrypto key for bundled assets appended to executable.
func SetDecryptoKeyToExeBundle(key string) error {
	return DefaultRepository.SetDecryptoKeyToExeBundle(key)
}

// RegisterBundle registers single packed bundle file to repository
func RegisterBundle(path string, option ...Option) error {
	return DefaultRepository.RegisterBundle(path, option...)
}

// RegisterFolder registers folder to repository
func RegisterFolder(path string, option ...Option) error {
	return DefaultRepository.RegisterFolder(path, option...)
}

// RegisterEncryptedFolder registers folder to repository with decryption key
func RegisterEncryptedFolder(path, key string, option ...Option) error {
	return DefaultRepository.RegisterEncryptedFolder(path, key, option...)
}

// Unload removes assets from default repository. The name is specified by option when registration.
func Unload(name string) error {
	return DefaultRepository.Unload(name)
}

// Find returns assets in default asset repository.
func Find(path string) (FileEntry, error) {
	return DefaultRepository.Find(path)
}
