package storage

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	db         *sql.DB
	insertStmt *sql.Stmt
	queryStmt  *sql.Stmt
}

func NewSQLiteStore(dbPath string) *SQLiteStore {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	}

	// Apply SQLite optimizations
	_, err = db.Exec(`
		PRAGMA synchronous = OFF;
		PRAGMA journal_mode = WAL;
		PRAGMA temp_store = MEMORY;
		PRAGMA mmap_size = 30000000000;`)
	if err != nil {
		log.Fatalf("Failed to apply PRAGMA settings: %v", err)
	}

	// Create table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS hashes (hash TEXT PRIMARY KEY)")
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Prepare statements
	insertStmt, err := db.Prepare("INSERT OR IGNORE INTO hashes (hash) VALUES (?)")
	if err != nil {
		log.Fatalf("Failed to prepare insert statement: %v", err)
	}

	queryStmt, err := db.Prepare("SELECT COUNT(1) FROM hashes WHERE hash=?")
	if err != nil {
		log.Fatalf("Failed to prepare query statement: %v", err)
	}

	return &SQLiteStore{
		db:         db,
		insertStmt: insertStmt,
		queryStmt:  queryStmt,
	}
}

func (s *SQLiteStore) Add(hash string) {
	_, err := s.insertStmt.Exec(hash)
	if err != nil {
		log.Printf("Failed to insert hash: %v", err)
	}
}

func (s *SQLiteStore) AddBatch(hashes []string) {
	tx, err := s.db.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	for _, hash := range hashes {
		_, err := tx.Stmt(s.insertStmt).Exec(hash)
		if err != nil {
			log.Printf("Failed to insert hash: %v", err)
			_ = tx.Rollback()
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}
}

func (s *SQLiteStore) Exists(hash string) bool {
	var exists int
	err := s.queryStmt.QueryRow(hash).Scan(&exists)
	if err != nil {
		log.Printf("Failed to check hash existence: %v", err)
		return false
	}
	return exists > 0
}

func (s *SQLiteStore) Close() error {
	if s.insertStmt != nil {
		_ = s.insertStmt.Close()
	}
	if s.queryStmt != nil {
		_ = s.queryStmt.Close()
	}
	return s.db.Close()
}
