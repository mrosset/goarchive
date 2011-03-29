package goarchive

import (
	"os"
	"testing"
)

//TODO: change this test to work on more then just my machine
func TestDecompress(t *testing.T) {
	tmpdir := "./tmp"
	fi, err := os.Stat(tmpdir)
	if err == nil {
		if fi.IsDirectory() {
			os.RemoveAll(tmpdir)
		}
	}
	os.Mkdir(tmpdir, 0755)
	zip := NewZip()
	err = zip.Decompress("/home/strings/via/cache/gmp-5.0.1.tar.bz2", "./tmp")
	if err != nil {
		t.Errorf("Error: ", err)
	}
	err = zip.Decompress("/home/strings/via/cache/ppl-0.11.2.tar.gz", "./tmp")
}
