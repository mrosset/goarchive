package goarchive

import (
	"archive/tar"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	. "fmt"
	"io"
	"os"
	"path"
)


// Struct used to decompress 
type Zip struct {
	path    string
	Verbose bool
	Debug   bool
}


// Returns a new Zip struct
func NewZip(p string) (*Zip, os.Error) {
	if p == "" {
		return nil, os.NewError("path: empty")
	}
	return &Zip{path: p}, nil
}


// Decompress bzip2 or gzip tarball to destination directory
func (z *Zip) Decompress(dest string) (err os.Error) {
	var cr io.Reader
	f, err := os.Open(z.path, os.O_RDONLY, 0000)
	if err != nil {
		return err
	}
	defer f.Close()
	switch path.Ext(z.path) {
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
		if z.Debug {
			Printf("%v\n", hdr)
		}
		// Switch through header Typeflag and handle tar entry accordingly 
		switch hdr.Typeflag {
		// Handles Directories
		case tar.TypeDir:
			path := path.Join(dest, hdr.Name)
			if z.Verbose {
				Printf("%v\n", hdr.Name)
			}
			if err := mkDir(path, hdr.Mode); err != nil {
				return err
			}
			continue
		// TODO: handle symlinks
		case tar.TypeSymlink, tar.TypeLink:
		case tar.TypeReg, tar.TypeRegA:
			path := path.Join(dest, hdr.Name)
			if z.Verbose {
				Printf("%v\n", hdr.Name)
			}
			if err := writeFile(path, hdr, tr); err != nil {
				return err
			}
			continue
		default:
			// Handles gnu LongLink long file names
			if string(hdr.Typeflag) == "L" {
				lfile := new(bytes.Buffer)
				// Get longlink path from tar file data
				lfile.ReadFrom(tr)
				if z.Verbose {
					Printf("%v\n", lfile.String())
				}
				fpath := path.Join(dest, lfile.String())
				// Read next iteration for file data
				hdr, err := tr.Next()
				if hdr.Typeflag == tar.TypeDir {
					err := mkDir(fpath, hdr.Mode)
					if err != nil {
						return err
					}
					continue
				}
				if err != nil && err != os.EOF {
					return err
				}
				// Write long file data to disk
				if err := writeFile(fpath, hdr, tr); err != nil {
					return err
				}
			}
		}
	}
	return
}


// Make directory with permission
func mkDir(path string, mode int64) (err os.Error) {
	if err = os.Mkdir(path, uint32(mode)); err != nil {
		return err
	}
	return
}


// Write file from tar reader
func writeFile(path string, hdr *tar.Header, tr *tar.Reader) (err os.Error) {
	f, err := os.Open(path, os.O_WRONLY|os.O_CREAT, uint32(hdr.Mode))
	if err != nil {
		return err
	}
	_, err = io.Copy(f, tr)
	f.Close()
	if err != nil {
		return err
	}
	return
}


func (z *Zip) pVerbose(path string) {
	if z.Verbose {
		Printf("%v\n", path)
	}
}
