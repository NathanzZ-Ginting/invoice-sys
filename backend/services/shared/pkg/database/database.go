package database

import (
	"fmt"
	"log"
	"os"

	"github.com/supabase-community/supabase-go"
)

// Client wraps the Supabase client
type Client struct {
	Supabase *supabase.Client
}

// NewClient creates a new database client
func NewClient() (*Client, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		return nil, fmt.Errorf("SUPABASE_URL and SUPABASE_KEY must be set")
	}

	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Supabase client: %v", err)
	}

	log.Printf("✅ Database client connected")

	return &Client{
		Supabase: client,
	}, nil
}
