package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type URLRepository struct {
	client *redis.Client
}

func NewURLRepository(addr string) *URLRepository {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &URLRepository{client: client}
}

func (r *URLRepository) Save(ctx context.Context, code string, encryptedURL []byte) error {
	key := fmt.Sprintf("url:%s", code)
	return r.client.Set(ctx, key, encryptedURL, 24*time.Hour).Err()
}

func (r *URLRepository) Get(ctx context.Context, code string) ([]byte, error) {
	key := fmt.Sprintf("url:%s", code)
	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (r *URLRepository) Delete(ctx context.Context, code string) error {
	key := fmt.Sprintf("url:%s", code)
	return r.client.Del(ctx, key).Err()
}
