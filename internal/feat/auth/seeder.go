package auth

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	hm "github.com/hermesgen/hm"
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
	if err := s.JSONSeeder.Setup(ctx); err != nil {
		return err
	}
	return s.SeedAll(ctx)
}

// SeedAll loads and applies all auth seeds in a single transaction.
func (s *Seeder) SeedAll(ctx context.Context) error {
	s.Log().Info("Seeding GitAuth data...")
	byFeature, err := s.JSONSeeder.LoadJSONSeeds()
	if err != nil {
		return fmt.Errorf("failed to load JSON seeds: %w", err)
	}
	const authFeat = "auth"
	for feature, seeds := range byFeature {
		if feature != authFeat {
			continue
		}

		for _, seed := range seeds {
			applied, err := s.JSONSeeder.SeedApplied(seed.Datetime, seed.Name, feature)
			if err != nil {
				return fmt.Errorf("failed to check if seed was applied: %w", err)
			}
			if applied {
				s.Log().Debugf("Seed already applied: %s-%s [%s]", seed.Datetime, seed.Name, feature)
				continue
			}

			var data SeedData
			err = json.Unmarshal([]byte(seed.Content), &data)
			if err != nil {
				return fmt.Errorf("failed to unmarshal %s seed: %w", feature, err)
			}

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
	ctx, tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error at beginning tx for seedUsers: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	s.Log().Debug("Seeding users: start")
	defer s.Log().Debug("Seeding users: end")
	for i := range data.Users {
		u := &data.Users[i]
		u.GenCreateValues()

		err = s.repo.CreateUser(ctx, u)
		if err != nil {
			return fmt.Errorf("error inserting user: %w", err)
		}
		userRefMap[u.Ref()] = u.GetID()
	}
	return tx.Commit()
}
