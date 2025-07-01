# bit-scout

**bit-scout** is a fast, modular search engine focused on typeahead, token-based matching, and trie-based prediction — built in Go with future plans for clustering, advanced ranking models, and optional ML-based fallback.

> 🧪 Alpha: This is a single-node prototype with document loading and basic vectorization. Search functionality is in development.

---

## 🌟 Features

### ✅ Implemented
- **Document Loading System**: Pluggable corpus loader interface with filesystem implementation
- **Document Model**: Rich document representation with metadata and vectorization
- **Loader Registry**: Plugin architecture for different document sources
- **Basic Vectorization**: Simple vector generation based on file metadata
- **CLI Interface**: Interactive search interface with document loading and display
- **Simple Index**: In-memory index with basic search functionality
- **Advanced Query System**: Boolean query parser with dimension-based filtering
- **Search Functionality**: Both simple text search and advanced boolean queries

### 🚧 In Development
- Trie-based prefix prediction
- Token expansion: stemming, typos, transliteration
- Advanced search models (Vector Space, BM25)
- Inverted index optimization
- Offline corpus ingestion and index building
- In-memory or pluggable persistence (filesystem, BoltDB, and more via plugin architecture)

---

## 🧠 Why bit-scout?

Modern search systems (e.g., Lucene, ElasticSearch) are powerful but heavyweight. `bit-scout` explores how far you can go with:

- Fully embedded Go infrastructure
- Trie + inverted index hybrid architecture
- Fast suggestions and typeahead use cases
- Clean separation of scoring, expansion, and retrieval
- Modular search models (Boolean, Vector, BM25)
- Optional fallback for semantic or ML-based suggestions
- Pluggable persistence for different storage backends

---

## 🚀 Quick Start

```bash
# Build and run the current implementation
go run cmd/bitscout/main.go

# This will load all documents from the current directory
# and display their metadata and vector representations
```

### Current Functionality
The application currently:
1. Loads documents from the filesystem (excluding directories)
2. Extracts rich metadata (file size, permissions, timestamps, etc.)
3. Generates simple vectors based on file characteristics
4. Builds an in-memory index for fast searching
5. Provides interactive search interface with:
   - Simple text search across content and metadata
   - Advanced boolean queries with dimension filtering
   - Support for operators: =, !=, <, <=, >, >=, contains
   - AND logic for combining conditions

### Search Examples
```bash
# Simple text search
> README
> go
> main

# Advanced boolean queries
> fileExtension=go
> fileSize<1000
> filename contains README
> fileExtension=go and fileSize<1000
> fileExtension!=md
```

### Planned Features
```bash
# Index persistence (planned)
bit-scout index --input corpus.txt --output snapshot.db

# Batch search (planned)
bit-scout query "fileExtension=go" --max-results 5
```

## 📚 Documentation

- **[Architecture](docs/architecture.md)**: High-level system design and components
- **[Implementation Status](docs/implementation-status.md)**: Detailed technical overview of current implementation
- **[Roadmap](roadmap.md)**: Development phases and feature planning
