package corpus

/*
Implementation of corpus loader for filesystem.
*/

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// Constants for file size normalization
const (
	MEAN_FILESIZE = 1024 * 1024       // 1MB mean file size
	MAX_FILESIZE  = 100 * 1024 * 1024 // 100MB max file size for normalization
	MEAN_TIME     = 1640995200        // Unix timestamp for 2022-01-01 (baseline)
	MAX_TIME      = 31536000          // 1 year in seconds (365 days)
)

type FilesystemLoader struct {
	root string
}

func NewFilesystemLoader(root string) *FilesystemLoader {
	log.Info().Msgf("NewFilesystemLoader: %s", root)
	return &FilesystemLoader{root: root}
}

func (l *FilesystemLoader) Load() ([]Document, error) {
	log.Info().Msgf("FilesystemLoader.Load from %s", l.root)
	documents := []Document{}

	err := filepath.Walk(l.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Error().Msgf("FilesystemLoader.Load: %s", err)
			return err
		}

		if info.IsDir() {
			log.Info().Msgf("FilesystemLoader.Load: skipping directory: %s", path)
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			log.Error().Msgf("FilesystemLoader.Load: %s", err)
			return err
		}

		log.Info().Msgf("FilesystemLoader.Load: adding document: %s", path)

		documents = append(documents, Document{
			ID:     makeID(path),
			Text:   string(content),
			Source: path,
			Meta:   getMeta(info, path, content),
			Vector: getVector(path, info, content),
		})

		return nil
	})

	return documents, err
}

// UUID
func makeID(path string) string {
	return uuid.New().String()
}

func getMeta(info os.FileInfo, path string, content []byte) map[string]string {
	return map[string]string{
		"filename":     info.Name(),
		"path":         path,
		"extension":    filepath.Ext(info.Name()),
		"fileSize":     strconv.FormatInt(int64(len(content)), 10),
		"lastModified": info.ModTime().Format(time.RFC3339),
		"isDir":        strconv.FormatBool(info.IsDir()),
		"isSymlink":    strconv.FormatBool(info.Mode()&os.ModeSymlink != 0),
		"isExecutable": strconv.FormatBool(info.Mode()&0100 != 0),
		"isWritable":   strconv.FormatBool(info.Mode()&0200 != 0),
		"isReadable":   strconv.FormatBool(info.Mode()&0400 != 0),
		"isHidden":     strconv.FormatBool(info.Name()[0] == '.'),
		"isSystem":     strconv.FormatBool(info.Mode()&01000 != 0),
		"isArchive":    strconv.FormatBool(info.Mode()&02000 != 0),
	}
}

func getVector(path string, info os.FileInfo, content []byte) []float64 {
	// get some numeric metadata
	fileSize := float64(len(content))
	lastModified := float64(info.ModTime().Unix())

	// normalize values
	fileSize = (fileSize - MEAN_FILESIZE) / MAX_FILESIZE
	lastModified = (lastModified - MEAN_TIME) / MAX_TIME
	// return a vector of these values
	return []float64{fileSize, lastModified}
}
