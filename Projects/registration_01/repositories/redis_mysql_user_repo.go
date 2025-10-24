package repositories

import (
	"context"
	"encoding/json"
	"registration/models"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisMysqlUserRepository struct {
	repo  UserRepository
	redis *redis.Client
	ctx context.Context
}

func NewRedisMysqlUserRepository(repo UserRepository, redis *redis.Client) *RedisMysqlUserRepository {
	return &RedisMysqlUserRepository{
		repo: repo,
		redis: redis,
		ctx: context.Background(),
	}
}

func (r *RedisMysqlUserRepository) Create(user models.User) error {
	return r.repo.Create(user)
}

func (r *RedisMysqlUserRepository) GetByEmail(email string) (*models.User, error) {
	// query cache
	// key-value
	cacheKey := "user:" + email                  // user:labib@gmail.com = {}
	value, err := r.redis.Get(r.ctx, cacheKey).Result() // value is JSON-encoded string
	if err == nil {
		var cacheUser models.User
		if err := json.Unmarshal([]byte(value), &cacheUser); err == nil {
			return &cacheUser, nil
		}
	}

	user, err := r.repo.GetByEmail(email)
	if err != nil || user == nil {
		return nil, err
	}

	// // cache user

	userJson, err := json.Marshal(user)
	if err == nil {
		r.redis.Set(r.ctx, cacheKey, userJson, 10*time.Minute)
	}
	return user, nil
}
