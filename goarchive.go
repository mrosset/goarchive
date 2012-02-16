package goarchive

import (
	"archive/tar"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	. "fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Struct used to decompress 
type Tar struct {
	Verbose bool
	Debug   bool
}

const (
	Bzip = "application/x-bzip2"
	Gz   = "application/x-gzip"
)

func init() {
	log.SetPrefix("goarchive: ")
	log.SetFlags(0)

}
func GetReader(path string) (r io.Reader, err error) {
	magic, err := fileMagic(path)
	if err != nil {
		return nil, err
	}
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	switch magic {
	case Bzip:
		log.Println("using bzip reader for", path)
		r = bzip2.NewReader(fd)
	case Gz:
		log.Println("using gzip reader for", path)
		r, err = gzip.NewReader(fd)
		if err != nil {
			return nil, err
		}
	}
	return r, err
}

// Returns a new Zip struct
func NewTar() *Tar {
	return new(Tar)
}

func (z *Tar) Peek(cr io.Reader) (dir string, err error) {
	tr := tar.NewReader(cr)
	hdr, err := tr.Next()
	if err != nil && err != io.EOF {
		return "", err
	}
	return path.Clean(hdr.Name), nil
}

// Decompress bzip2 or gzip Reader to destination directory
func (z *Tar) Untar(cr io.Reader, dest string) (err error) {
	tr := tar.NewReader(cr)
	for {
		hdr, err := tr.Next()
		if err != nil && err != io.EOF {
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
				if err != nil && err != io.EOF {
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
func mkDir(path string, mode int64) (err error) {
	if fileExists(path) {
		return
	}
	err = os.Mkdir(path, os.FileMode(mode))
	if err != nil {
		return err
	}
	return
}

// Write file from tar reader
func writeFile(path string, hdr *tar.Header, tr *tar.Reader) (err error) {
	f, err := os.Create(path)
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

func (z *Tar) pVerbose(path string) {
	if z.Verbose {
		Printf("%v\n", path)
	}
}

// helper function to test if file/directory exists
func fileExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	if !fi.IsDir() || fi.IsDir() {
		return true
	}
	return false
}

func fileMagic(path string) (string, error) {
	output, err := exec.Command("file", "-b", "-i", path).Output()
	if err != nil {
		return "", err
	}
	mime := strings.Split(string(output), " ")
	return mime[0][:len(mime[0])-1], nil
}
