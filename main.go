package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/radurobot/go-markdown-crawler/internal/crawler"
	"github.com/radurobot/go-markdown-crawler/internal/storage"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "go-markdown-crawler",
		Short: "go-markdown-crawler - Crawls websites, converts them to markdown and stores the results",
		Run:   run,
	}

	// Define flags
	crawler.SetupFlags(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, _ []string) {
	options := crawler.ParseFlags(cmd)

	// Initialize hash storage
	var hashStore storage.HashStore
	if options.UseInMemory {
		hashStore = storage.NewInMemoryStore()
	} else {
		dbPath := fmt.Sprintf("hashes-%d.db", os.Getpid())
		hashStore = storage.NewSQLiteStore(dbPath)
		defer hashStore.Close()
		defer os.Remove(dbPath)
		defer os.Remove(dbPath + "-shm")
		defer os.Remove(dbPath + "-wal")
	}
	crawler.Crawl(options, hashStore)
}
