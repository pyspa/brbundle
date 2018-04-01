package brbundle

type Decompressor interface {
	Decompress() ([]byte, error)
}
