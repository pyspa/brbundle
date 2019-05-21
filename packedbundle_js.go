package brbundle

func (p packedFileEntry) GetLocalPath() (string, error) {
	return "", errors.New("GetLocalPath() is not supported")
}
