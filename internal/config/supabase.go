package config

import (
	"fmt"
	"os"

	supabase "github.com/supabase-community/supabase-go"
)

var SupabaseClient *supabase.Client

func InitSupabase() error {
	// ambil dari env
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_KEY") // service_role key

	if supabaseURL == "" || supabaseKey == "" {
		return fmt.Errorf("SUPABASE_URL or SUPABASE_SERVICE_KEY not set")
	}

	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to create supabase client: %v", err))
	}
	SupabaseClient = client

	return nil
}
