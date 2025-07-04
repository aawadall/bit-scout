package index

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/aawadall/bit-scout/internal/models"
	"github.com/rs/zerolog/log"
	"go.etcd.io/bbolt"
)

/**
A simple index that persists to disk. in boltDB
*/

// dbOperation represents a database operation to be performed asynchronously
type dbOperation struct {
	opType string
	data   interface{}
}

type PersistedSimpleIndex struct {
	index  *SimpleIndex
	db     *bbolt.DB
	opChan chan dbOperation
	done   chan struct{}
	wg     sync.WaitGroup
	mu     sync.RWMutex
}

func NewPersistedSimpleIndex() *PersistedSimpleIndex {
	return &PersistedSimpleIndex{
		index:  NewSimpleIndex(),
		db:     nil,                          // Will be initialized when database is opened
		opChan: make(chan dbOperation, 1000), // Buffer for async operations
		done:   make(chan struct{}),
	}
}

// OpenDatabase opens the BoltDB database for persistence, creating it if it doesn't exist
func (p *PersistedSimpleIndex) OpenDatabase(dbPath string) error {
	if p.db != nil {
		return fmt.Errorf("database already open")
	}

	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory %s: %w", dir, err)
	}

	// Check if database file exists
	_, err := os.Stat(dbPath)
	dbExists := err == nil

	// Open or create the database
	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Create buckets if they don't exist
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("documents"))
		if err != nil {
			return fmt.Errorf("failed to create documents bucket: %w", err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte("config"))
		if err != nil {
			return fmt.Errorf("failed to create config bucket: %w", err)
		}
		return nil
	})

	if err != nil {
		db.Close()
		return err
	}

	p.db = db

	// Start the async database worker
	p.startAsyncWorker()

	if dbExists {
		log.Info().Msgf("Opened existing persistent database at %s", dbPath)
	} else {
		log.Info().Msgf("Created new persistent database at %s", dbPath)
	}
	return nil
}

// startAsyncWorker starts the goroutine that handles database operations asynchronously
func (p *PersistedSimpleIndex) startAsyncWorker() {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			select {
			case op := <-p.opChan:
				p.processDBOperation(op)
			case <-p.done:
				log.Info().Msg("Async database worker shutting down")
				return
			}
		}
	}()
	log.Info().Msg("Started async database worker")
}

// processDBOperation handles individual database operations
func (p *PersistedSimpleIndex) processDBOperation(op dbOperation) {
	p.mu.RLock()
	db := p.db
	p.mu.RUnlock()

	if db == nil {
		log.Warn().Msg("Database not available for async operation")
		return
	}

	switch op.opType {
	case "add_document":
		if doc, ok := op.data.(models.Document); ok {
			p.asyncAddDocument(doc)
		}
	case "add_documents":
		if docs, ok := op.data.([]models.Document); ok {
			p.asyncAddDocuments(docs)
		}
	case "update_document":
		if data, ok := op.data.(map[string]interface{}); ok {
			if id, ok := data["id"].(string); ok {
				if doc, ok := data["document"].(models.Document); ok {
					p.asyncUpdateDocument(id, doc)
				}
			}
		}
	case "delete_document":
		if id, ok := op.data.(string); ok {
			p.asyncDeleteDocument(id)
		}
	case "delete_documents":
		if ids, ok := op.data.([]string); ok {
			p.asyncDeleteDocuments(ids)
		}
	case "update_documents":
		if docs, ok := op.data.([]models.Document); ok {
			p.asyncUpdateDocuments(docs)
		}
	case "configure":
		if config, ok := op.data.(map[string]interface{}); ok {
			p.asyncConfigure(config)
		}
	default:
		log.Warn().Msgf("Unknown async operation type: %s", op.opType)
	}
}

// asyncAddDocument performs the actual database operation for adding a document
func (p *PersistedSimpleIndex) asyncAddDocument(doc models.Document) {
	p.mu.RLock()
	db := p.db
	p.mu.RUnlock()

	err := db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("documents"))
		docData, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("failed to marshal document: %w", err)
		}
		return bucket.Put([]byte(doc.ID), docData)
	})

	if err != nil {
		log.Error().Err(err).Msgf("Async add document failed for %s", doc.ID)
	} else {
		log.Debug().Msgf("Async added document %s to database", doc.ID)
	}
}

// asyncAddDocuments performs the actual database operation for adding multiple documents
func (p *PersistedSimpleIndex) asyncAddDocuments(docs []models.Document) {
	p.mu.RLock()
	db := p.db
	p.mu.RUnlock()

	err := db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("documents"))
		for _, doc := range docs {
			docData, err := json.Marshal(doc)
			if err != nil {
				return fmt.Errorf("failed to marshal document %s: %w", doc.ID, err)
			}
			if err := bucket.Put([]byte(doc.ID), docData); err != nil {
				return fmt.Errorf("failed to store document %s: %w", doc.ID, err)
			}
		}
		return nil
	})

	if err != nil {
		log.Error().Err(err).Msgf("Async add documents failed for %d documents", len(docs))
	} else {
		log.Debug().Msgf("Async added %d documents to database", len(docs))
	}
}

// asyncUpdateDocument performs the actual database operation for updating a document
func (p *PersistedSimpleIndex) asyncUpdateDocument(id string, doc models.Document) {
	p.mu.RLock()
	db := p.db
	p.mu.RUnlock()

	err := db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("documents"))
		docData, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("failed to marshal document: %w", err)
		}
		return bucket.Put([]byte(id), docData)
	})

	if err != nil {
		log.Error().Err(err).Msgf("Async update document failed for %s", id)
	} else {
		log.Debug().Msgf("Async updated document %s in database", id)
	}
}

// asyncDeleteDocument performs the actual database operation for deleting a document
func (p *PersistedSimpleIndex) asyncDeleteDocument(id string) {
	p.mu.RLock()
	db := p.db
	p.mu.RUnlock()

	err := db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("documents"))
		return bucket.Delete([]byte(id))
	})

	if err != nil {
		log.Error().Err(err).Msgf("Async delete document failed for %s", id)
	} else {
		log.Debug().Msgf("Async deleted document %s from database", id)
	}
}

// asyncDeleteDocuments performs the actual database operation for deleting multiple documents
func (p *PersistedSimpleIndex) asyncDeleteDocuments(ids []string) {
	p.mu.RLock()
	db := p.db
	p.mu.RUnlock()

	err := db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("documents"))
		for _, id := range ids {
			if err := bucket.Delete([]byte(id)); err != nil {
				return fmt.Errorf("failed to delete document %s: %w", id, err)
			}
		}
		return nil
	})

	if err != nil {
		log.Error().Err(err).Msgf("Async delete documents failed for %d documents", len(ids))
	} else {
		log.Debug().Msgf("Async deleted %d documents from database", len(ids))
	}
}

// asyncUpdateDocuments performs the actual database operation for updating multiple documents
func (p *PersistedSimpleIndex) asyncUpdateDocuments(docs []models.Document) {
	p.mu.RLock()
	db := p.db
	p.mu.RUnlock()

	err := db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("documents"))
		for _, doc := range docs {
			docData, err := json.Marshal(doc)
			if err != nil {
				return fmt.Errorf("failed to marshal document %s: %w", doc.ID, err)
			}
			if err := bucket.Put([]byte(doc.ID), docData); err != nil {
				return fmt.Errorf("failed to update document %s: %w", doc.ID, err)
			}
		}
		return nil
	})

	if err != nil {
		log.Error().Err(err).Msgf("Async update documents failed for %d documents", len(docs))
	} else {
		log.Debug().Msgf("Async updated %d documents in database", len(docs))
	}
}

// asyncConfigure performs the actual database operation for configuration
func (p *PersistedSimpleIndex) asyncConfigure(config map[string]interface{}) {
	p.mu.RLock()
	db := p.db
	p.mu.RUnlock()

	err := db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("config"))
		configData, err := json.Marshal(config)
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}
		return bucket.Put([]byte("index_config"), configData)
	})

	if err != nil {
		log.Error().Err(err).Msg("Async configure failed")
	} else {
		log.Debug().Msg("Async configured database")
	}
}

// Configure sets the index configuration and persists it asynchronously
func (p *PersistedSimpleIndex) Configure(config map[string]interface{}) error {
	// Configure the in-memory index
	if err := p.index.Configure(config); err != nil {
		return err
	}

	// Queue async database operation if database is open
	p.mu.RLock()
	if p.db != nil {
		select {
		case p.opChan <- dbOperation{opType: "configure", data: config}:
			log.Debug().Msg("Queued async configure operation")
		default:
			log.Warn().Msg("Async operation queue full, configure operation dropped")
		}
	}
	p.mu.RUnlock()

	return nil
}

// ShowConfig returns the current index configuration (memory-only operation)
func (p *PersistedSimpleIndex) ShowConfig() (map[string]interface{}, error) {
	return p.index.ShowConfig()
}

// AddDocument adds a single document to the index and persists it asynchronously
func (p *PersistedSimpleIndex) AddDocument(doc models.Document) error {
	// Add to in-memory index
	if err := p.index.AddDocument(doc); err != nil {
		return err
	}

	// Queue async database operation if database is open
	p.mu.RLock()
	if p.db != nil {
		select {
		case p.opChan <- dbOperation{opType: "add_document", data: doc}:
			log.Debug().Msgf("Queued async add document operation for %s", doc.ID)
		default:
			log.Warn().Msgf("Async operation queue full, add document operation dropped for %s", doc.ID)
		}
	}
	p.mu.RUnlock()

	return nil
}

// AddDocuments adds multiple documents to the index and persists them asynchronously
func (p *PersistedSimpleIndex) AddDocuments(docs []models.Document) error {
	// Add to in-memory index
	if err := p.index.AddDocuments(docs); err != nil {
		return err
	}

	// Queue async database operation if database is open
	p.mu.RLock()
	if p.db != nil {
		select {
		case p.opChan <- dbOperation{opType: "add_documents", data: docs}:
			log.Debug().Msgf("Queued async add documents operation for %d documents", len(docs))
		default:
			log.Warn().Msgf("Async operation queue full, add documents operation dropped for %d documents", len(docs))
		}
	}
	p.mu.RUnlock()

	return nil
}

// Search performs search using only the in-memory index (no database access)
func (p *PersistedSimpleIndex) Search(query string) ([]models.Document, error) {
	// Search operations work purely from memory for maximum performance
	return p.index.Search(query)
}

// DeleteDocument removes a document from the index and database asynchronously
func (p *PersistedSimpleIndex) DeleteDocument(id string) error {
	// Delete from in-memory index
	if err := p.index.DeleteDocument(id); err != nil {
		return err
	}

	// Queue async database operation if database is open
	p.mu.RLock()
	if p.db != nil {
		select {
		case p.opChan <- dbOperation{opType: "delete_document", data: id}:
			log.Debug().Msgf("Queued async delete document operation for %s", id)
		default:
			log.Warn().Msgf("Async operation queue full, delete document operation dropped for %s", id)
		}
	}
	p.mu.RUnlock()

	return nil
}

// DeleteDocuments removes multiple documents from the index and database asynchronously
func (p *PersistedSimpleIndex) DeleteDocuments(ids []string) error {
	// Delete from in-memory index
	if err := p.index.DeleteDocuments(ids); err != nil {
		return err
	}

	// Queue async database operation if database is open
	p.mu.RLock()
	if p.db != nil {
		select {
		case p.opChan <- dbOperation{opType: "delete_documents", data: ids}:
			log.Debug().Msgf("Queued async delete documents operation for %d documents", len(ids))
		default:
			log.Warn().Msgf("Async operation queue full, delete documents operation dropped for %d documents", len(ids))
		}
	}
	p.mu.RUnlock()

	return nil
}

// UpdateDocument updates a document in the index and database asynchronously
func (p *PersistedSimpleIndex) UpdateDocument(id string, doc models.Document) error {
	// Update in-memory index
	if err := p.index.UpdateDocument(id, doc); err != nil {
		return err
	}

	// Queue async database operation if database is open
	p.mu.RLock()
	if p.db != nil {
		data := map[string]interface{}{
			"id":       id,
			"document": doc,
		}
		select {
		case p.opChan <- dbOperation{opType: "update_document", data: data}:
			log.Debug().Msgf("Queued async update document operation for %s", id)
		default:
			log.Warn().Msgf("Async operation queue full, update document operation dropped for %s", id)
		}
	}
	p.mu.RUnlock()

	return nil
}

// UpdateDocuments updates multiple documents in the index and database asynchronously
func (p *PersistedSimpleIndex) UpdateDocuments(docs []models.Document) error {
	// Update in-memory index
	if err := p.index.UpdateDocuments(docs); err != nil {
		return err
	}

	// Queue async database operation if database is open
	p.mu.RLock()
	if p.db != nil {
		select {
		case p.opChan <- dbOperation{opType: "update_documents", data: docs}:
			log.Debug().Msgf("Queued async update documents operation for %d documents", len(docs))
		default:
			log.Warn().Msgf("Async operation queue full, update documents operation dropped for %d documents", len(docs))
		}
	}
	p.mu.RUnlock()

	return nil
}

// Close closes the database connection and shuts down the async worker
func (p *PersistedSimpleIndex) Close() error {
	// Signal the async worker to shut down
	close(p.done)

	// Wait for the async worker to finish
	p.wg.Wait()

	// Close the database
	p.mu.Lock()
	if p.db != nil {
		if err := p.db.Close(); err != nil {
			p.mu.Unlock()
			return fmt.Errorf("failed to close database: %w", err)
		}
		p.db = nil
		log.Info().Msg("PersistedSimpleIndex database closed")
	}
	p.mu.Unlock()

	return p.index.Close()
}

// Flush ensures all data is written to disk
func (p *PersistedSimpleIndex) Flush() error {
	if p.db != nil {
		return p.db.Sync()
	}
	return p.index.Flush()
}

// Optimize optimizes the index for faster search
func (p *PersistedSimpleIndex) Optimize() error {
	return p.index.Optimize()
}

// Count returns the number of documents in the index (memory-only operation)
func (p *PersistedSimpleIndex) Count() (int, error) {
	return p.index.Count()
}

// Size returns the approximate size of the index in bytes (memory-only operation)
func (p *PersistedSimpleIndex) Size() (int, error) {
	return p.index.Size()
}

// LoadDocumentsFromDatabase loads all documents from the database into memory (synchronous read operation)
func (p *PersistedSimpleIndex) LoadDocumentsFromDatabase() error {
	p.mu.RLock()
	db := p.db
	p.mu.RUnlock()

	if db == nil {
		return fmt.Errorf("database not open")
	}

	// Clear the in-memory index first to avoid duplicates
	p.index = NewSimpleIndex()

	var documents []models.Document

	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("documents"))
		if bucket == nil {
			return fmt.Errorf("documents bucket not found")
		}

		return bucket.ForEach(func(k, v []byte) error {
			var doc models.Document
			if err := json.Unmarshal(v, &doc); err != nil {
				return fmt.Errorf("failed to unmarshal document %s: %w", string(k), err)
			}
			documents = append(documents, doc)
			return nil
		})
	})

	if err != nil {
		return err
	}

	// Add all documents to the in-memory index at once
	if err := p.index.AddDocuments(documents); err != nil {
		return fmt.Errorf("failed to add documents to memory index: %w", err)
	}

	log.Info().Msgf("Loaded %d documents from database into memory", len(documents))
	return nil
}

// LoadConfigFromDatabase loads configuration from the database into memory
func (p *PersistedSimpleIndex) LoadConfigFromDatabase() error {
	p.mu.RLock()
	db := p.db
	p.mu.RUnlock()

	if db == nil {
		return fmt.Errorf("database not open")
	}

	var config map[string]interface{}

	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("config"))
		if bucket == nil {
			return fmt.Errorf("config bucket not found")
		}

		configData := bucket.Get([]byte("index_config"))
		if configData == nil {
			return fmt.Errorf("no configuration found in database")
		}

		return json.Unmarshal(configData, &config)
	})

	if err != nil {
		return err
	}

	// Apply the configuration to the in-memory index
	if err := p.index.Configure(config); err != nil {
		return fmt.Errorf("failed to apply configuration to memory index: %w", err)
	}

	log.Info().Msg("Loaded configuration from database into memory")
	return nil
}

// LoadFromDatabase loads both documents and configuration from the database into memory
func (p *PersistedSimpleIndex) LoadAllFromDatabase() error {
	// Load configuration first
	if err := p.LoadConfigFromDatabase(); err != nil {
		log.Warn().Err(err).Msg("Failed to load configuration from database, continuing with documents only")
		// Continue even if config loading fails
	}

	// Load documents
	if err := p.LoadDocumentsFromDatabase(); err != nil {
		return fmt.Errorf("failed to load documents from database: %w", err)
	}

	log.Info().Msg("Successfully loaded all data from database into memory")
	return nil
}

// NewPersistedSimpleIndexWithDatabase creates a new index and opens the database (creates if doesn't exist)
func NewPersistedSimpleIndexWithDatabase(dbPath string) (*PersistedSimpleIndex, error) {
	index := NewPersistedSimpleIndex()

	if err := index.OpenDatabase(dbPath); err != nil {
		return nil, fmt.Errorf("failed to open/create database: %w", err)
	}

	return index, nil
}

// NewPersistedSimpleIndexWithDatabaseAndLoad creates a new index, opens the database, and loads existing data
func NewPersistedSimpleIndexWithDatabaseAndLoad(dbPath string) (*PersistedSimpleIndex, error) {
	index, err := NewPersistedSimpleIndexWithDatabase(dbPath)
	if err != nil {
		return nil, err
	}

	// Check if database has any data before trying to load
	isEmpty, err := index.IsDatabaseEmpty()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to check if database is empty, attempting to load anyway")
	} else if isEmpty {
		log.Info().Msg("Database is empty, starting with fresh index")
		return index, nil
	}

	// Try to load existing data from database
	if err := index.LoadAllFromDatabase(); err != nil {
		log.Warn().Err(err).Msg("Failed to load existing data from database, starting with empty index")
		// Continue with empty index if loading fails
	} else {
		log.Info().Msg("Successfully loaded existing data from database")
	}

	return index, nil
}

// IsDatabaseEmpty checks if the database has any documents
func (p *PersistedSimpleIndex) IsDatabaseEmpty() (bool, error) {
	p.mu.RLock()
	db := p.db
	p.mu.RUnlock()

	if db == nil {
		return true, fmt.Errorf("database not open")
	}

	var isEmpty bool
	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("documents"))
		if bucket == nil {
			isEmpty = true
			return nil
		}

		cursor := bucket.Cursor()
		key, _ := cursor.First()
		isEmpty = key == nil
		return nil
	})

	return isEmpty, err
}

// GetDatabaseStats returns statistics about the database
func (p *PersistedSimpleIndex) GetDatabaseStats() (map[string]interface{}, error) {
	p.mu.RLock()
	db := p.db
	p.mu.RUnlock()

	if db == nil {
		return nil, fmt.Errorf("database not open")
	}

	stats := make(map[string]interface{})

	err := db.View(func(tx *bbolt.Tx) error {
		// Count documents
		docBucket := tx.Bucket([]byte("documents"))
		if docBucket != nil {
			docCount := 0
			cursor := docBucket.Cursor()
			for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
				docCount++
			}
			stats["document_count"] = docCount
		} else {
			stats["document_count"] = 0
		}

		// Check if config exists
		configBucket := tx.Bucket([]byte("config"))
		if configBucket != nil {
			configData := configBucket.Get([]byte("index_config"))
			stats["has_config"] = configData != nil
		} else {
			stats["has_config"] = false
		}

		return nil
	})

	return stats, err
}
