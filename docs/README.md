# rag-postgres — documentation

PostgreSQL vector store for ai-rag — pgvector + pg_search BM25 hybrid retrieval

## Overview

Package ragpostgres is a PostgreSQL vector store for togo ai-rag, backed by
pgvector (similarity) and — when available — pg_search/ParadeDB BM25 for hybrid
retrieval. Blank-import it and set RAG_STORE=postgres + DATABASE_URL.

## Install

```bash
togo install togo-framework/rag-postgres
```



## Configuration

Environment variables read by this plugin (extracted from the source — see the gateway/provider docs for each value):

| Env var |
|---|
| `DATABASE_URL` |
| `RAG_PG_DIM` |

## Usage

```go
// Registers a Postgres-backed rag.Store (pgvector + pg_search BM25 -> hybrid RRF).
// Install it + set DATABASE_URL; ai-rag uses it automatically.
```

## Links

- Marketplace: https://to-go.dev/marketplace
- Source: https://github.com/togo-framework/rag-postgres
- Full README: ../README.md
