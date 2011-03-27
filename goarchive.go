package goarchive

import (
	"archive/tar"
	"compress/bzip2"
	"compress/gzip"
	"io"
	"os"
	"path"
)

const (
	ISDIR  = 53
	ISFILE = 0
	DMASK  = 0755
	RMASK  = 0000
	FMASK  = 0644
)

// Decompress bzip2 or gzip file to destination directory
func Untar(file string, dest string) (err os.Error) {
	f, err := os.Open(file, os.O_RDONLY, RMASK)
	if err != nil {
		return err
	}
	defer f.Close()
	var cr io.Reader
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
		fpath := path.Join(dest, hdr.Name)
		fmask := uint32(hdr.Mode)
		if hdr.Typeflag == ISDIR {
			if err := os.Mkdir(fpath, uint32(hdr.Mode)); err != nil {
				return err
			}
			continue
		}
		if hdr.Typeflag == ISFILE {
			f, err := os.Open(fpath, os.O_WRONLY|os.O_CREAT, fmask)
			if err != nil {
				return err
			}
			_, err = io.Copy(f, tr)
			if err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}
	return
}
