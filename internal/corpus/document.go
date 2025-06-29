package corpus

import (
	"fmt"
	"strconv"
)

// Document represents a single document loaded from a corpus source.
type Document struct {
	ID     string
	Text   string
	Source string            // Source of the document (e.g., file path, URL)
	Vector []float64         // Vector representation of the document
	Meta   map[string]string // Optional metadata (e.g., filename, tags)
}

// Print the document
func (d *Document) Print() {
	printedBlob := "\n\n\n"
	printedBlob += "Document ID: " + d.ID + "\n"
	//printedBlob += "Document Text: " + d.Text + "\n"
	printedBlob += "Document Source: " + d.Source + "\n"
	// print vector
	printedBlob += "Document Vector:\n"
	for _, value := range d.Vector {
		printedBlob += "  " + strconv.FormatFloat(value, 'f', -1, 64) + "\n"
	}
	printedBlob += "Document Metadata:\n"
	for key, value := range d.Meta {
		printedBlob += "  " + key + ": " + value + "\n"
	}
	fmt.Print(printedBlob)
}
