// Package ragpostgres is a PostgreSQL vector store for togo ai-rag, backed by
// pgvector (similarity) and — when available — pg_search/ParadeDB BM25 for hybrid
// retrieval. Blank-import it and set RAG_STORE=postgres + DATABASE_URL.
package ragpostgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	rag "github.com/togo-framework/ai-rag"
	"github.com/togo-framework/togo"
)

func init() {
	rag.RegisterStore("postgres", func(k *togo.Kernel) (rag.Store, error) {
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			return nil, errors.New("rag-postgres: DATABASE_URL not set")
		}
		dim := 1536
		if v := os.Getenv("RAG_PG_DIM"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 {
				dim = n
			}
		}
		ctx := context.Background()
		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			return nil, err
		}
		s := &store{pool: pool, dim: dim, table: "rag_chunks"}
		if err := s.ensure(ctx); err != nil {
			pool.Close()
			return nil, err
		}
		return s, nil
	})
}

type store struct {
	pool  *pgxpool.Pool
	dim   int
	table string
	bm25  bool
}

func (s *store) ensure(ctx context.Context) error {
	for _, q := range []string{
		"CREATE EXTENSION IF NOT EXISTS vector",
		fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id text PRIMARY KEY, doc_id text NOT NULL, text text NOT NULL, embedding vector(%d))", s.table, s.dim),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s_embedding_idx ON %s USING hnsw (embedding vector_cosine_ops)", s.table, s.table),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s_doc_idx ON %s (doc_id)", s.table, s.table),
	} {
		if _, err := s.pool.Exec(ctx, q); err != nil {
			return fmt.Errorf("rag-postgres ensure: %q: %w", q, err)
		}
	}
	// Best-effort BM25 (pg_search / ParadeDB) — enables hybrid retrieval when present.
	if _, err := s.pool.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS pg_search"); err == nil {
		if _, err := s.pool.Exec(ctx, fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s_bm25_idx ON %s USING bm25 (id, text) WITH (key_field='id')", s.table, s.table)); err == nil {
			s.bm25 = true
		}
	}
	return nil
}

func vec(v []float32) string {
	parts := make([]string, len(v))
	for i, f := range v {
		parts[i] = strconv.FormatFloat(float64(f), 'f', -1, 32)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func (s *store) Upsert(ctx context.Context, chunks []rag.Chunk) error {
	q := fmt.Sprintf("INSERT INTO %s (id, doc_id, text, embedding) VALUES ($1,$2,$3,$4::vector) ON CONFLICT (id) DO UPDATE SET doc_id=EXCLUDED.doc_id, text=EXCLUDED.text, embedding=EXCLUDED.embedding", s.table)
	for _, c := range chunks {
		if _, err := s.pool.Exec(ctx, q, c.ID, c.DocID, c.Text, vec(c.Vector)); err != nil {
			return err
		}
	}
	return nil
}

// Search runs pgvector cosine similarity (the rag.Store interface only passes a
// vector). For hybrid vector+BM25 use HybridSearch.
func (s *store) Search(ctx context.Context, vector []float32, topK int) ([]rag.Chunk, error) {
	if topK <= 0 {
		topK = 5
	}
	rows, err := s.pool.Query(ctx,
		fmt.Sprintf("SELECT id, doc_id, text, 1 - (embedding <=> $1::vector) AS score FROM %s ORDER BY embedding <=> $1::vector LIMIT $2", s.table),
		vec(vector), topK)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []rag.Chunk
	for rows.Next() {
		var c rag.Chunk
		if err := rows.Scan(&c.ID, &c.DocID, &c.Text, &c.Score); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *store) Delete(ctx context.Context, docID string) error {
	_, err := s.pool.Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE doc_id=$1", s.table), docID)
	return err
}

// HybridSearch combines pgvector similarity with pg_search BM25 keyword ranking
// via reciprocal-rank fusion. Falls back to vector-only when pg_search is absent.
func (s *store) HybridSearch(ctx context.Context, vector []float32, query string, topK int) ([]rag.Chunk, error) {
	if topK <= 0 {
		topK = 5
	}
	if !s.bm25 || strings.TrimSpace(query) == "" {
		return s.Search(ctx, vector, topK)
	}
	const k = 60
	sql := fmt.Sprintf(`
WITH v AS (SELECT id, row_number() OVER (ORDER BY embedding <=> $1::vector) AS r FROM %[1]s),
     b AS (SELECT id, row_number() OVER (ORDER BY paradedb.score(id) DESC) AS r FROM %[1]s WHERE text @@@ $2)
SELECT c.id, c.doc_id, c.text,
       COALESCE(1.0/($3 + v.r),0) + COALESCE(1.0/($3 + b.r),0) AS score
FROM %[1]s c
LEFT JOIN v ON v.id=c.id
LEFT JOIN b ON b.id=c.id
WHERE v.id IS NOT NULL OR b.id IS NOT NULL
ORDER BY score DESC LIMIT $4`, s.table)
	rows, err := s.pool.Query(ctx, sql, vec(vector), query, k, topK)
	if err != nil {
		return s.Search(ctx, vector, topK)
	}
	defer rows.Close()
	var out []rag.Chunk
	for rows.Next() {
		var c rag.Chunk
		if err := rows.Scan(&c.ID, &c.DocID, &c.Text, &c.Score); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}
