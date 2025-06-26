# bit-scout

**bit-scout** is a fast, modular search engine focused on typeahead, token-based matching, and trie-based prediction — built in Go with future plans for clustering and optional ML-based fallback.

> 🧪 Alpha: This is a single-node prototype with offline indexing and in-memory querying via CLI.

---

## 🌟 Features

- Trie-based prefix prediction
- Inverted index for token lookup
- Token expansion: stemming, typos, transliteration
- Offline corpus ingestion and index building
- In-memory or BoltDB-backed persistence (WIP)
- Simple CLI search interface

---

## 🧠 Why bit-scout?

Modern search systems (e.g., Lucene, ElasticSearch) are powerful but heavyweight. `bit-scout` explores how far you can go with:

- Fully embedded Go infrastructure
- Trie + inverted index hybrid architecture
- Fast suggestions and typeahead use cases
- Clean separation of scoring, expansion, and retrieval
- Optional fallback for semantic or ML-based suggestions

---

## 🚀 Quick Start (planned)

```bash
# Index builder (offline)
bit-scout index --input corpus.txt --output snapshot.db

# CLI search
bit-scout query "canned sou" --max-results 5
