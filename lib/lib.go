package lib

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type GitObjectType string

const (
	COMMIT           GitObjectType = "commit"
	TREE             GitObjectType = "tree"
	BLOB             GitObjectType = "blob"
	GIT_OBJ_NAME_LEN               = 2
	// TODO: Think about a buffer size
	GIT_OBJ_BUF_SIZE = 4096
)

type GitObject struct {
	Header  GitObjectType
	Size    int
	Content string
}

type Repository struct {
	Path string
}

func readObject(absPath string) GitObject {
	rawData, err := os.ReadFile(absPath)

	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(rawData)

	data, err := zlib.NewReader(reader)

	if err != nil {
		log.Fatal(err)
	}

	buff := make([]byte, GIT_OBJ_BUF_SIZE)
	data.Read(buff)
	objData := string(buff)

	meta := strings.Split(objData, "\000")[0]
	header := GitObjectType(strings.Split(meta, " ")[0])
	contentSize, err := strconv.Atoi(strings.Split(meta, " ")[1])

	if err != nil {
		log.Fatal(err)
	}
	content := strings.Split(objData, "\000")[1]

	return GitObject{Header: header, Size: contentSize, Content: content}
}

func (r *Repository) GetObjects() []GitObject {
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

		if info.Name() == "pack" {
			return filepath.SkipDir
		}

		if !info.IsDir() {
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

// TODO: Handle gpg key in commit body
func parseCommit(commit string) string {
	commitInfo := strings.Split(commit, "\n")

	tree := commitInfo[0]
	parent := commitInfo[1]
	author := commitInfo[2]
	// NB: Git adds new line after commit message
	msg := commitInfo[len(commitInfo)-2]

	commitToPrint := []string{tree, parent, author, msg}
	commitObj := strings.Join(commitToPrint, "\n")

	return commitObj
}

func (r *Repository) Log() {
	objects := r.GetObjects()
	var commits []string

	for _, obj := range objects {
		if obj.Header == COMMIT {
			commit := parseCommit(obj.Content)
			commits = append(commits, commit)
		}
	}

	sort.Slice(commits, func(i, j int) bool {
		author1 := strings.Split(commits[i], "\n")[2]
		author2 := strings.Split(commits[j], "\n")[2]
		commitTime1, err := strconv.Atoi(strings.Split(author1, " ")[len(strings.Split(author1, " "))-2])
		commitTime2, err := strconv.Atoi(strings.Split(author2, " ")[len(strings.Split(author2, " "))-2])
		if err != nil {
			log.Fatal(err)
		}
		return commitTime1 > commitTime2
	})

	for _, c := range commits {
		fmt.Println(c, "\n")
	}
}
