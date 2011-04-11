package goarchive

import (
	. "fmt"
	"os"
	"testing"
)

var tmpDir = "./tmp"

// setup tmpDir for decompressions
func init() {
	if fileExists(tmpDir) {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			Printf("%v\n", err)
			os.Exit(1)
		}
	}
	err := os.Mkdir(tmpDir, 0755)
	if err != nil {
		Printf("%v\n", err)
		os.Exit(1)
	}
}


// Test struct
type testZip struct {
	name       string
	zipFile    string
	dir        string
	file       string
	data       string
	linkSrc    string
	linkTarget string
	longPath   string
}

var gzipTest = &testZip{
	name:       "gzip",
	zipFile:    "testdata/gzip.tar.gz",
	dir:        "directory",
	file:       "small.txt",
	data:       "small2.txt",
	linkSrc:    "src",
	linkTarget: "target",
}

var bzip2Test = &testZip{
	name:       "bzip2",
	zipFile:    "testdata/bzip2.tar.bz2",
	dir:        "directory",
	file:       "small.txt",
	data:       "small2.txt",
	linkSrc:    "src",
	linkTarget: "target",
}

var longLinkTest = &testZip{
	name:       "longlink",
	zipFile:    "testdata/longlink.tar.bz2",
	dir:        "directory",
	file:       "small.txt",
	data:       "small2.txt",
	linkSrc:    "src",
	linkTarget: "target",
	longPath:   "0123456789101112131415161718192021222324252627282930313233343536373839404142434445464748495051525354555657585960616263646566676869707172737475767778798081828384858687888990919293949596979899100",
}

var tests = []*testZip{
	gzipTest,
	bzip2Test,
	longLinkTest,
}

// Loop through each test and test for decompression
// TODO: test each test struct field
func TestDecompress(t *testing.T) {
	for _, zt := range tests {
		zip, err := NewZip(zt.zipFile)
		if err != nil {
			t.Errorf("NewZip %v : Unexpected error: %v", zt.name, err)
		}
		if err := zip.Decompress(tmpDir); err != nil {
			t.Errorf("Decompress %v : Unexpected error: %v", zt.name, err)
		}
	}
}

func TestPeek(t *testing.T) {
	for _, zt := range tests {
		zip, err := NewZip(zt.zipFile)
		if err != nil {
			t.Errorf("NewZip %v : Unexpected error: %v", zt.name, err)
		}
		if dir,_ := zip.Peek(); dir != zt.name {
			t.Errorf("Peek expected %v got %v", zt.name, dir)
		}
	}

}
