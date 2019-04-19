package brbundle

type CompressionType int
type EncryptionType int

const (
	NoCompression CompressionType = iota
	Brotli
	LZ4

	NoEncryption EncryptionType = iota
	AES
)

const (
	UseBrotli     = "b"
	UseLZ4        = "l"
	NotToCompress = "-"

	UseAES        = "a"
	NotToEncrypto = "-"
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

func (c CompressionType) Flag() string {
	switch c {
	case Brotli:
		return UseBrotli
	case LZ4:
		return UseLZ4
	case NoCompression:
		return NotToCompress
	}
	return ""
}

func (e EncryptionType) String() string {
	switch e {
	case AES:
		return "AES-256-GCM"
	case NoEncryption:
		return "no"
	}
	return ""
}

func (e EncryptionType) Flag() string {
	switch e {
	case AES:
		return UseAES
	case NoEncryption:
		return NotToEncrypto
	}
	return ""
}
