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
		return "b"
	case LZ4:
		return "l"
	case NoCompression:
		return "-"
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
		return "a"
	case NoEncryption:
		return "-"
	}
	return ""
}
