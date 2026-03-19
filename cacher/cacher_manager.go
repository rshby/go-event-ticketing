package cacher

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rshby/go-event-ticketing/config"
	"github.com/rshby/go-event-ticketing/tracing"
)

// CacheManager interface yang fleksibel
type CacheManager interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string, dest any) error
	Delete(ctx context.Context, key ...string) error
	Compress(value any) ([]byte, error)
	Decompress(data []byte, dest any) error
	Raw() *redis.Client
}

type cacheManager struct {
	client *redis.Client
}

// NewCacheManager creates new instance of CacheManager
func NewCacheManager(client *redis.Client) CacheManager {
	return &cacheManager{
		client: client,
	}
}

func (c *cacheManager) Compress(value any) ([]byte, error) {
	// 1. Convert struct/data ke JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	// 2. Compress JSON menggunakan GZIP
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(jsonData); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (c *cacheManager) Decompress(data []byte, dest any) error {
	// 1. Decompress data GZIP dari Redis
	b := bytes.NewBuffer(data)
	gz, err := gzip.NewReader(b)
	if err != nil {
		return err
	}
	defer gz.Close()

	decompressedData, err := io.ReadAll(gz)
	if err != nil {
		return err
	}

	// 2. Unmarshal JSON kembali ke Struct tujuan (dest)
	return json.Unmarshal(decompressedData, dest)
}

func (c *cacheManager) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	ctx, span := tracing.Start(ctx)
	defer span.End()

	// Lakukan kompresi otomatis
	compressedData, err := c.Compress(value)
	if err != nil {
		return err
	}

	if ttl == 0 {
		ttl = config.DefaultRedisTTL
	}

	// Set byte yang sudah dikompres ke Redis
	return c.client.Set(ctx, key, compressedData, ttl).Err()
}

func (c *cacheManager) Get(ctx context.Context, key string, dest any) error {
	ctx, span := tracing.Start(ctx)
	defer span.End()

	// Ambil raw byte dari Redis
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return err // Akan me-return redis.Nil jika tidak ketemu
	}

	// Dekompresi dan mapping ke struct tujuan
	return c.Decompress(data, dest)
}

func (c *cacheManager) Delete(ctx context.Context, key ...string) error {
	ctx, span := tracing.Start(ctx)
	defer span.End()

	return c.client.Del(ctx, key...).Err()
}

func (c *cacheManager) Raw() *redis.Client {
	return c.client
}
