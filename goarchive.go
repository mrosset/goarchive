package goarchive

import (
	"archive/tar"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
)


var (
	Debug = false
)


const (
	DMASK = 0755
	RMASK = 0000
	FMASK = 0644
)


// Decompress bzip2 or gzip file to destination directory
func Untar(file string, dest string) (err os.Error) {
	var cr io.Reader
	f, err := os.Open(file, os.O_RDONLY, RMASK)
	if err != nil {
		return err
	}
	defer f.Close()
	switch path.Ext(file) {
	case ".bz2":
		cr = bzip2.NewReader(f)
	case ".gz":
		cr, err = gzip.NewReader(f)
		if err != nil {
			return err
		}
	default:
		return os.NewError("unable to determine filetype")
	}
	tr := tar.NewReader(cr)
	for {
		hdr, err := tr.Next()
		if err != nil && err != os.EOF {
			return err
		}
		if hdr == nil {
			break
		}
		if Debug {
			fmt.Printf("%v %v\n", hdr.Name, hdr.Typeflag)
		}
		fpath := path.Join(dest, hdr.Name)
		fmask := uint32(hdr.Mode)
		if hdr.Typeflag == tar.TypeDir {
			if err := os.Mkdir(fpath, uint32(hdr.Mode)); err != nil {
				return err
			}
			continue
		}
		f, err := os.Open(fpath, os.O_WRONLY|os.O_CREAT, fmask)
		if err != nil {
			return err
		}
		_, err = io.Copy(f, tr)
		f.Close()
		if err != nil {
			return err
		}
	}
	return
}
