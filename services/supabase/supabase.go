package supabase

import (
	"net/http"
	"os"
	"sync"
)

var (
	client     *http.Client
	baseURL    string
	apiKey     string
	configured bool
	mu         sync.RWMutex
)

// Init initializes a minimal Supabase HTTP client using env vars SUPABASE_URL and SUPABASE_KEY
func Init() {
	mu.Lock()
	defer mu.Unlock()

	baseURL = os.Getenv("SUPABASE_URL")
	apiKey = os.Getenv("SUPABASE_KEY")

	if baseURL != "" && apiKey != "" {
		client = &http.Client{}
		configured = true
	} else {
		configured = false
	}
}

// IsConfigured returns true when supabase client was initialized with env vars
func IsConfigured() bool {
	mu.RLock()
	defer mu.RUnlock()
	return configured
}
