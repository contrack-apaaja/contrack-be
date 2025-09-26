package supabase

import (
	"fmt"
	"os"
	"sync"

	supabaseClient "github.com/supabase-community/supabase-go"
)

var (
	client     *supabaseClient.Client
	configured bool
	mu         sync.RWMutex
)

// Init initializes the supabase-go client using SUPABASE_URL and SUPABASE_KEY
func Init() error {
	mu.Lock()
	defer mu.Unlock()

	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")

	if url == "" || key == "" {
		configured = false
		return fmt.Errorf("SUPABASE_URL or SUPABASE_KEY not set")
	}

	c, err := supabaseClient.NewClient(url, key, nil)
	if err != nil {
		configured = false
		return err
	}
	client = c
	configured = true
	return nil
}

func IsConfigured() bool {
	mu.RLock()
	defer mu.RUnlock()
	return configured
}

// Client returns the initialized supabase client or nil
func Client() *supabaseClient.Client {
	mu.RLock()
	defer mu.RUnlock()
	return client
}
