<!-- togo-header -->
<div align="center">
  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />
  <h1>togo-framework/rag-postgres</h1>
  <p>
    <a href="https://to-go.dev/marketplace"><img src="https://img.shields.io/badge/marketplace-to--go.dev-1FC7DC" alt="marketplace" /></a>
    <a href="https://pkg.go.dev/github.com/togo-framework/rag-postgres"><img src="https://pkg.go.dev/badge/github.com/togo-framework/rag-postgres.svg" alt="pkg.go.dev" /></a>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT" />
  </p>
  <p><strong>Part of the <a href="https://to-go.dev">togo</a> framework.</strong></p>
</div>

## Install

```bash
togo install togo-framework/rag-postgres
```

<!-- /togo-header -->

# rag-postgres — PostgreSQL vector store for togo ai-rag

A `rag.Store` driver backed by **PostgreSQL**: **pgvector** for similarity +
**pg_search** (ParadeDB BM25) for keyword ranking → **hybrid retrieval**.
Pairs with `togo-postgres` (ships pgvector + ParadeDB).

```bash
togo install togo-framework/rag-postgres
```

Set `RAG_STORE=postgres` + `DATABASE_URL` (optionally `RAG_PG_DIM`, default 1536).
On boot it ensures `CREATE EXTENSION vector`, the `rag_chunks` table, an HNSW
cosine index, and — if `pg_search` is present — a BM25 index for hybrid search.

- `Search(vector, topK)` — pgvector cosine (the standard ai-rag path).
- `HybridSearch(vector, query, topK)` — vector + BM25 via reciprocal-rank fusion
  (call it directly via a type assertion; falls back to vector-only without pg_search).

MIT

<!-- togo-sponsors -->
---

<div align="center">
  <h3>Premium sponsors</h3>
  <p>
    <a href="https://id8media.com"><strong>ID8 Media</strong></a> &nbsp;·&nbsp;
    <a href="https://one-studio.co"><strong>One Studio</strong></a>
  </p>
  <p><sub>Support togo — <a href="https://github.com/sponsors/fadymondy">become a sponsor</a>.</sub></p>
</div>
<!-- /togo-sponsors -->
