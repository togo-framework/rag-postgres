# rag-postgres — documentation

  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />

## Overview

Package ragpostgres is a PostgreSQL vector store for togo ai-rag, backed by
pgvector (similarity) and — when available — pg_search/ParadeDB BM25 for hybrid
retrieval. Blank-import it and set RAG_STORE=postgres + DATABASE_URL.

## Install

```bash
togo install togo-framework/rag-postgres
```



## Configuration

Environment variables read by this plugin (extracted from the source):

| Env var | Notes |
|---|---|
| `DATABASE_URL` | _see provider docs_ |
| `G` | _see provider docs_ |
| `RAG_PG_DIM` | _see provider docs_ |

## Usage

```go
// Registers a Postgres-backed rag.Store (pgvector + pg_search BM25 -> hybrid RRF).
// Just install it and set DATABASE_URL; ai-rag picks it up automatically.
```

## Links

- Marketplace: https://to-go.dev/marketplace
- Source: https://github.com/togo-framework/rag-postgres
- README: ../README.md
