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

type FilesystemLoader struct {
	root string
}

func NewFilesystemLoader(root string) *FilesystemLoader {
	log.Info().Msgf("NewFilesystemLoader: %s", root)
	return &FilesystemLoader{root: root}
}

func (l *FilesystemLoader) Load(source string) ([]Document, error) {
	log.Info().Msgf("FilesystemLoader.Load: %s", source)
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
			ID:   makeID(path),
			Text: string(content),
			Source: path,
			Meta: getMeta(info, path),
		})

		
		return nil
	})

	return documents, err
}

// UUID
func makeID(path string) string {
	return uuid.New().String()
}

func getMeta(info os.FileInfo, path string) map[string]string {
	return map[string]string{
		"filename": info.Name(),
		"path": path,
		"extension": filepath.Ext(info.Name()),
		"fileSize": strconv.FormatInt(info.Size(), 10),
		"lastModified": info.ModTime().Format(time.RFC3339),
		"isDir": strconv.FormatBool(info.IsDir()),
		"isSymlink": strconv.FormatBool(info.Mode()&os.ModeSymlink != 0),
		"isExecutable": strconv.FormatBool(info.Mode()&0100 != 0),
		"isWritable": strconv.FormatBool(info.Mode()&0200 != 0),
		"isReadable": strconv.FormatBool(info.Mode()&0400 != 0),
		"isHidden": strconv.FormatBool(info.Name()[0] == '.'),
		"isSystem": strconv.FormatBool(info.Mode()&01000 != 0),
		"isArchive": strconv.FormatBool(info.Mode()&02000 != 0),
	}
}