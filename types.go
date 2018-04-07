package brbundle

type CompressionType int
type EncryptionType int

const (
	NoCompression CompressionType = iota
	Brotli
	LZ4

	NoEncryption EncryptionType = iota
	AES
	ChaCha20Poly1305
)

func (c CompressionType) String() string {
	switch c {
	case Brotli:
		return "brotli"
	case LZ4:
		return "lz4"
	case NoCompression:
		return "no"
	}
	return ""
}

func (c CompressionType) ConstantName() string {
	switch c {
	case Brotli:
		return "brbundle.Brotli"
	case LZ4:
		return "brbundle.LZ4"
	case NoCompression:
		return "brbundle.NoCompression"
	}
	return ""
}

func (c CompressionType) FunctionName() string {
	switch c {
	case Brotli:
		return "brbundle.BrotliDecompressor()"
	case LZ4:
		return "brbundle.LZ4Decompressor()"
	case NoCompression:
		return "brbundle.NullDecompressor()"
	}
	return ""
}

func (e EncryptionType) String() string {
	switch e {
	case AES:
		return "AES-256-GCM"
	case ChaCha20Poly1305:
		return "ChaCha20-Poly1305"
	case NoEncryption:
		return "no"
	}
	return ""
}

func (e EncryptionType) ConstantName() string {
	switch e {
	case AES:
		return "brbundle.AES"
	case ChaCha20Poly1305:
		return "brbundle.ChaCha20Poly1305"
	case NoEncryption:
		return "brbundle.NoEncryption"
	}
	return ""
}

func (e EncryptionType) FunctionName() string {
	switch e {
	case AES:
		return "brbundle.AESDecryptor()"
	case ChaCha20Poly1305:
		return "brbundle.ChaChaDecryptor()"
	case NoEncryption:
		return "brbundle.NullDecryptor()"
	}
	return ""
}
