package brbundle_test

import (
	"flag"
	"github.com/ToQoz/gopwt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()
	gopwt.Empower()
	os.Exit(m.Run())
}
