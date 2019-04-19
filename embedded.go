package brbundle

type embededBundle struct {
	data []byte
	name string
}

var embeddedBundles []*embededBundle

func RegisterEmbeddedBundle(data []byte, name string) {
	embeddedBundles = append(embeddedBundles, &embededBundle{data, name})
}
