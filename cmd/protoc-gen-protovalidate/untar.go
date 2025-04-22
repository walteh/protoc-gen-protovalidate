package protocgenprotovalidate

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"strings"
)

// Untar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
func Untar(r io.Reader) (map[string]string, error) {

	files := make(map[string]string)

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return files, nil

		// return any other error
		case err != nil:
			return nil, err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		// target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			// if _, err := fs.Stat(fsys, header.Name); err != nil {
			// 	if err := fs.Dir(fsys, header.Name, 0755); err != nil {
			// 		return err
			// 	}
			// }

		// if it's a file create it
		case tar.TypeReg:

			content, err := io.ReadAll(tr)
			if err != nil {
				return nil, err
			}

			// remove the first path segment
			filePath := strings.SplitN(header.Name, "/", 2)[1]

			files[filePath] = string(content)

			// fmt.Fprintf(os.Stderr, "filePath: %s\n", filePath)

		}
	}

	// return files, nil
}
