package jsonnet

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/go-jsonnet"
	jbV1 "github.com/jsonnet-bundler/jsonnet-bundler/spec/v1"
	"github.com/spf13/afero"
)

type FileImporter struct {
	JPaths      []string
	fs          afero.Fs
	projectPath string
	fsCache     map[string]*fsCacheEntry
}

type fsCacheEntry struct {
	exists   bool
	contents jsonnet.Contents
}

// NewImporter - Get the Jsonnet File Import customized to the Jsonnet Bundler type.
func NewImporter(fs afero.Fs, projectPath string, jbFile jbV1.JsonnetFile) *FileImporter {
	libsPath := []string{}
	fileImporter := FileImporter{
		fs:          fs,
		projectPath: projectPath,
	}

	libsPath = append(libsPath, projectPath)

	if jbFile.LegacyImports {
		// Jsonnet-Bundler LegacyImports put the imported folders in the top directory with symlinks.
		libPath := filepath.Join(projectPath, ShareLibsPath)
		libsPath = append(libsPath, libPath)
	} else {
		// Jsonnet-Bundler Imports put complies to the GoMod style of artifact management (vendoring)
		// This means we need to take an extra step to find the top level key for each shared folder.
		libsMap := make(map[string][]string)

		for k := range jbFile.Dependencies {
			libPath := filepath.Join(projectPath, ShareLibsPath, k)
			libPathSplit := strings.Split(libPath, "/")
			libsKey := strings.Join(libPathSplit[:len(libPathSplit)-1], "/")

			if len(libsMap[libsKey]) > 0 {
				libsMap[libsKey] = append(libsMap[libsKey], k)
			} else {
				libsMap[libsKey] = []string{k}
			}
		}

		for k := range libsMap {
			libsPath = append(libsPath, k)
		}
	}

	fileImporter.JPaths = append(fileImporter.JPaths, libsPath...)
	return &fileImporter
}

func (importer *FileImporter) tryPath(dir, importedPath string) (found bool, contents jsonnet.Contents, foundHere string, err error) {
	if importer.fsCache == nil {
		importer.fsCache = make(map[string]*fsCacheEntry)
	}
	var absPath string
	if path.IsAbs(importedPath) {
		absPath = importedPath
	} else {
		absPath = path.Join(dir, importedPath)
	}
	var entry *fsCacheEntry
	if cacheEntry, isCached := importer.fsCache[absPath]; isCached {
		entry = cacheEntry
	} else {
		contentBytes, err := afero.ReadFile(importer.fs, absPath)
		if err != nil {
			if os.IsNotExist(err) {
				entry = &fsCacheEntry{
					exists: false,
				}
			} else {
				return false, jsonnet.Contents{}, "", err
			}
		} else {
			entry = &fsCacheEntry{
				exists:   true,
				contents: jsonnet.MakeContents(string(contentBytes)),
			}
		}
		importer.fsCache[absPath] = entry
	}
	return entry.exists, entry.contents, absPath, nil
}

// Import imports file from the filesystem.
func (importer *FileImporter) Import(importedFrom, importedPath string) (contents jsonnet.Contents, foundAt string, err error) {
	// TODO(sbarzowski) Make sure that dir is absolute and resolving of ""
	// is independent from current CWD. The default path should be saved
	// in the importer.
	// We need to relativize the paths in the error formatter, so that the stack traces
	// don't have ugly absolute paths (less readable and messy with golden tests).
	dir, _ := path.Split(importedFrom)
	found, content, foundHere, err := importer.tryPath(dir, importedPath)
	if err != nil {
		return jsonnet.Contents{}, "", err
	}

	for i := len(importer.JPaths) - 1; !found && i >= 0; i-- {
		found, content, foundHere, err = importer.tryPath(importer.JPaths[i], importedPath)
		if err != nil {
			return jsonnet.Contents{}, "", err
		}
	}

	if !found {
		return jsonnet.Contents{}, "", fmt.Errorf("couldn't open import %#v: no match locally or in the Jsonnet library paths", importedPath)
	}
	return content, foundHere, nil
}
