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

func Unload(name string) error {
	return DefaultRepository.Unload(name)
}
