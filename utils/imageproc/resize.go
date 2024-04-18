package imageproc

import (
	"bytes"
	"io"
	"path/filepath"

	"github.com/disintegration/imaging"
)

func IsContentTypeValid(contentType string, filename string) bool {
	if contentType != "image/jpeg" && contentType != "image/png" {
		return false
	}

	if contentType == "image/jpeg" && filepath.Ext(filename) != ".jpg" {
		return false
	}

	if contentType == "image/png" && filepath.Ext(filename) != ".png" {
		return false
	}

	return true
}

// This function resizes an image to fit in 250x200 pixels bounding box (aspect ratio is maintained).
func ResizeImage(f io.ReadSeeker, name string) (*bytes.Reader, int, error) {
	_, err := f.Seek(0, io.SeekStart)
	if err != nil {
		return nil, 0, err
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, 0, err
	}

	src, err := imaging.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, 0, err
	}
	dst := imaging.Fit(src, 250, 200, imaging.Lanczos)

	var o []byte
	w := bytes.NewBuffer(o)
	err = imaging.Encode(w, dst, imaging.JPEG)
	if err != nil {
		return nil, 0, err
	}

	r := bytes.NewReader(w.Bytes())
	return r, int(r.Size()), nil
}
