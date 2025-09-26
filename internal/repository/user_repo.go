package repository

import (
	"context"
	"encoding/json"

	"contrack-be/internal/models"
	"contrack-be/internal/services/supabase"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo { return &UserRepo{} }

// List returns users from the "users" table using PostgREST client
func (r *UserRepo) List(ctx context.Context) ([]models.User, error) {
	client := supabase.Client()
	if client == nil {
		return nil, nil
	}

	// Use Postgrest client via client.From(...)
	resp, _, err := client.From("users").Select("id,email", "exact", false).Execute()
	if err != nil {
		return nil, err
	}

	var out []models.User
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}
	return out, nil
}
