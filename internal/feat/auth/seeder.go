package auth

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hermesgen/hm"
)

type Seeder struct {
	*hm.JSONSeeder
	repo Repo
}

type SeedData struct {
	Users []User `json:"users"`
}

func NewSeeder(assetsFS embed.FS, engine string, repo Repo, params hm.XParams) *Seeder {
	return &Seeder{
		JSONSeeder: hm.NewJSONSeeder("auth", assetsFS, engine, params),
		repo:       repo,
	}
}

func (s *Seeder) Setup(ctx context.Context) error {
	return s.JSONSeeder.Setup(ctx)
}

func (s *Seeder) Start(ctx context.Context) error {
	return s.SeedAll(ctx)
}

// SeedAll loads and applies all auth seeds in a single transaction.
func (s *Seeder) SeedAll(ctx context.Context) error {
	s.Log().Info("Seeding Auth data...")
	byFeature, err := s.JSONSeeder.LoadJSONSeeds()
	if err != nil {
		return fmt.Errorf("failed to load JSON seeds: %w", err)
	}

	s.Log().Info("Loaded seed files", "feature_count", len(byFeature))
	for feat, seeds := range byFeature {
		s.Log().Info("Feature seeds found", "feature", feat, "count", len(seeds))
	}

	const authFeat = "auth"
	for feature, seeds := range byFeature {
		if feature != authFeat {
			s.Log().Info("Skipping non-auth feature", "feature", feature)
			continue
		}

		for _, seed := range seeds {
			s.Log().Info("Processing auth seed", "name", seed.Name, "datetime", seed.Datetime)
			applied, err := s.JSONSeeder.SeedApplied(seed.Datetime, seed.Name, feature)
			if err != nil {
				return fmt.Errorf("failed to check if seed was applied: %w", err)
			}
			if applied {
				s.Log().Info("Seed already applied, skipping", "name", seed.Name)
				continue
			}

			s.Log().Info("Unmarshaling seed data", "content_length", len(seed.Content))
			var data SeedData
			err = json.Unmarshal([]byte(seed.Content), &data)
			if err != nil {
				return fmt.Errorf("failed to unmarshal %s seed: %w", feature, err)
			}

			s.Log().Info("Seed data unmarshaled", "users_count", len(data.Users))

			err = s.seedData(ctx, &data)
			if err != nil {
				return err
			}

			err = s.JSONSeeder.ApplyJSONSeed(seed.Datetime, seed.Name, feature, seed.Content)
			if err != nil {
				s.Log().Errorf("error recording JSON seed: %v", err)
			}
		}
	}
	return nil
}

// seedData applies a single SeedData in a transaction.
func (s *Seeder) seedData(ctx context.Context, data *SeedData) error {
	userRefMap := make(map[string]uuid.UUID)

	err := s.seedUsers(ctx, data, userRefMap)
	if err != nil {
		return err
	}

	return nil
}

// --- Helper functions for each entity type ---

func (s *Seeder) seedUsers(ctx context.Context, data *SeedData, userRefMap map[string]uuid.UUID) error {
	s.Log().Info("Seeding users", "user_count", len(data.Users))
	if len(data.Users) == 0 {
		s.Log().Info("No users in seed data to insert")
		return nil
	}

	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error at beginning tx for seedUsers: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	for i := range data.Users {
		u := &data.Users[i]
		s.Log().Info("Creating user", "username", u.Username, "ref", u.Ref())
		u.GenCreateValues()

		err = s.repo.CreateUser(ctx, u)
		if err != nil {
			return fmt.Errorf("error inserting user %s: %w", u.Username, err)
		}
		userRefMap[u.Ref()] = u.GetID()
		s.Log().Info("User created successfully", "username", u.Username)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing user transaction: %w", err)
	}
	s.Log().Info("All users committed successfully")
	return nil
}
