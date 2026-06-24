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
