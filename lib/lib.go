package lib

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type GitObjectType string

const (
	COMMIT           GitObjectType = "commit"
	TREE             GitObjectType = "tree"
	BLOB             GitObjectType = "blob"
	GIT_OBJ_NAME_LEN               = 2
)

// TODO: Add commit, tree
type GitObject struct {
	Header  GitObjectType
	Size    int
	Content string
}

type Repository struct {
	Path string
}

func parseBlob(blob string) GitObject {
	meta := strings.Split(blob, "\000")[0]
	header := GitObjectType(strings.Split(meta, " ")[0])
	contentSize, err := strconv.Atoi(strings.Split(meta, " ")[1])
	if err != nil {
		log.Fatal(err)
	}
	content := strings.Split(blob, "\000")[1]

	return GitObject{Header: header, Size: contentSize, Content: content[:10]}
}

func readObject(absPath string) GitObject {
	object := GitObject{}
	raw_data, err := os.ReadFile(absPath)

	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(raw_data)

	data, err := zlib.NewReader(reader)

	if err != nil {
		log.Fatal(err)
	}

	// TODO: Think about an appropriate buffer size
	buff := make([]byte, 500)

	data.Read(buff)

	objHeader := GitObjectType(strings.Split(string(buff), " ")[0])

	if err != nil {
		log.Fatal(err)
	}

	if objHeader == BLOB {
		object = parseBlob(string(buff))
	}

	return object
}

func (r *Repository) ListObjects() []GitObject {
	objects := make([]GitObject, 1)

	info, err := os.Stat(r.Path)

	if errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}

	if !info.IsDir() {
		log.Fatal("not a git repository")
	}

	err = filepath.Walk(r.Path+"/objects/", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (info.Name() != "pack" || info.Name() != "info") {
			if err != nil {
				log.Fatal(err)
			}

			absPath, err := filepath.Abs(path)

			if err != nil {
				log.Fatal(err)
			}

			object := readObject(absPath)
			objects = append(objects, object)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return objects
}
