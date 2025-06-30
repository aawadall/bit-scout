package index

import (
	"github.com/aawadall/bit-scout/internal/models"
)
/* Index Interface */

type Index interface {
	// Adds document to current index
	AddDocument(models.Document) error
}