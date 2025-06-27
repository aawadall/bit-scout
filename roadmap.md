# Roadmap — bit-scout

`bit-scout` is being built in 4 public phases:

1. **Alpha** — local single-node prototype
2. **Beta** — clustered architecture with gRPC API
3. **Pre-prod** — fallback, metrics, lightweight UI
4. **v1.0** — stable, pluggable, production-ready release

---

## 🧪 Alpha (Weeks 1–4)

🎯 Goal: Validate core trie + inverted index with token expansion and introduce Boolean search model.

✅ Features
- [ ] Corpus loader 
- [ ] Tokenizer
- [ ] Token expansion (typos, stems, translits)
- [ ] Inverted index
- [ ] Trie for next-token prediction
- [ ] Boolean search model (AND, OR, NOT queries)
- [ ] Basic scorer (weighted source signals)
- [ ] CLI search tool
- [ ] Index snapshot + load from disk
- [ ] Sample test corpus
- [ ] Pluggable persistence interface (filesystem, BoltDB)

---

## 🛠️ Beta (Weeks 5–10)

🎯 Goal: Multi-node architecture with gRPC query serving and expand search models.

✅ Features
- [ ] gRPC server for search queries
- [ ] Static shard config (manual node assignment)
- [ ] Coordinator node for routing
- [ ] Redis optional cache layer
- [ ] Trie hydration via RPC or shared volume
- [ ] Trie pruning and compression
- [ ] Vector Space search model
- [ ] BM25 ranking model
- [ ] Plugin architecture for additional persistence backends

---

## ⚙️ Pre-Prod (Weeks 11–18)

🎯 Goal: Fallback + monitoring + field test readiness

✅ Features
- [ ] ML fallback hook (e.g., cosine similarity service)
- [ ] Epsilon-greedy fallback trigger
- [ ] Logging and telemetry per query
- [ ] Admin dashboard (basic UI or CLI)
- [ ] Config-based A/B ranking logic

---

## 🟢 v1.0 Public

🎯 Goal: Stable deployment + documentation + community-friendly

✅ Features
- [ ] REST/GraphQL wrapper for gRPC
- [ ] Dockerized multi-node cluster
- [ ] Docs + API spec
- [ ] Scalable embedding support (optional)
- [ ] Embeddable as a Go module
- [ ] Mature plugin ecosystem for models and persistence

---

## 🌱 Future / Stretch Ideas

- Learner graph-style re-ranking
- WASM version for edge inference
- Real-time personalization scoring
- Cohere/OpenAI plugin module
- Visual trie explorer
- Hot reloadable expansions / rules

---

_This roadmap is subject to evolution — see issues and milestones for latest plans._
