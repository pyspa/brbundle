package brbundle

type Decryptor interface {
	Decrypto() ([]byte, error)
}
