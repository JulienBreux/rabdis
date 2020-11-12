package redis

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/jpillora/backoff"
	"github.com/julienbreux/rabdis/pkg/logger"
)

var (
	errorConnection = errors.New("redis not connected")
)

const (
	// INF for Infinity
	INF = "+inf"
)

// Redis represents the Redis interface
type Redis interface {
	Connect()
	Disconnect() error

	Increment(string) (int64, error)
	Decrement(string) (int64, error)

	Set(string, string, time.Duration) error
	Get(string) (string, error)
	Del(string) (int64, error)
	Exists(string) (int64, error)

	SearchByKey(string) ([]string, error)

	SetMemberAdd(string, ...string) error
	SetMemberExists(string, string) (bool, error)
	SetMemberDelete(string, string) (int64, error)
	SetMembers(string) ([]string, error)
	SetLength(string) (int64, error)

	HashItemAdd(string, string, string) error
	HashItemExists(string, string) (bool, error)
	HashItemGet(string, string) (string, error)
	HashItemDelete(string, string) (int64, error)
	HashItems(string) (map[string]string, error)
	HashLength(string) (int64, error)

	ScoreItemAdd(string, float64, string) error
	ScoreItemCount(string, string, string) (float64, error)

	PubSubPublish(string, string) error
	PubSubSubscribe(string, chan<- string) error

	FlushAll() error
}

type red struct {
	o *Options

	r *redis.Client
}

// New creates a new Redis instance
func New(opts ...Option) (Redis, error) {
	o, err := newOptions(opts...)
	if err != nil {
		return nil, err
	}

	return &red{
		o: o,
	}, nil
}

// Connect connects Redis
func (r *red) Connect() {
	r.connect()
}

func (r *red) connect() {
	r.o.Logger.Info("redis: connecting")

	b := &backoff.Backoff{}
	attempts := 1
	for {
		cli := redis.NewClient(r.config())
		_, err := cli.Ping().Result()
		if err == nil {
			r.r = cli
			r.o.Logger.Info("redis: connected")
			break
		}
		r.o.Logger.Error(
			"redis: failed to connect",
			logger.F("attempts", attempts),
			logger.E(err),
		)
		attempts++
		time.Sleep(b.Duration())
	}
	go r.watchConnectCloseAndReconnect()
}

func (r *red) watchConnectCloseAndReconnect() {
	for {
		if _, err := r.r.Ping().Result(); err == nil {
			r.o.Logger.Debug("redis: ping ok")
			time.Sleep(r.o.PingDelay)
			continue
		}

		r.r = nil

		r.o.Logger.Warn("redis: disconnected")

		r.connect()
		break
	}
}

func (r *red) config() *redis.Options {
	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.o.Host, r.o.Port),
		Password: r.o.Password,
		DB:       r.o.Database,
	}
}

func (r *red) prefixed(key string) string {
	return r.o.KeyPrefix + key
}

// Disconnect disconnects Redis
func (r *red) Disconnect() error {
	if r.r == nil {
		return nil
	}

	r.o.Logger.Info("redis: disconnecting")

	err := r.r.Close()
	if err == nil {
		r.o.Logger.Info("redis: disconnected")
		return nil
	}

	r.o.Logger.Warn(
		"redis: unable to stop",
		logger.E(err),
	)

	return err
}

// FlushAll flushes all [...] -_-
func (r *red) FlushAll() error {
	if r.r == nil {
		return errorConnection
	}
	_, err := r.r.FlushAll().Result()
	return err
}

// SIMPLE

// Increment increments a value
func (r *red) Increment(key string) (int64, error) {
	if r.r == nil {
		return 0, errorConnection
	}

	return r.r.Incr(r.prefixed(key)).Result()
}

// Decrement decrements a value
func (r *red) Decrement(key string) (int64, error) {
	if r.r == nil {
		return 0, errorConnection
	}

	return r.r.Decr(r.prefixed(key)).Result()
}

// Set sets a value
func (r *red) Set(key, value string, expiration time.Duration) error {
	if r.r == nil {
		return errorConnection
	}

	return r.r.Set(r.prefixed(key), value, expiration).Err()
}

// Get gets a value
func (r *red) Get(key string) (string, error) {
	if r.r == nil {
		return "", errorConnection
	}

	return r.r.Get(r.prefixed(key)).Result()
}

// Del deletes a value
func (r *red) Del(key string) (int64, error) {
	if r.r == nil {
		return 0, errorConnection
	}

	return r.r.Del(r.prefixed(key)).Result()
}

// Exists checks if a key exists
func (r *red) Exists(key string) (int64, error) {
	if r.r == nil {
		return 0, errorConnection
	}

	return r.r.Exists(r.prefixed(key)).Result()
}

// SearchByKey searchs by key the values
func (r *red) SearchByKey(pattern string) ([]string, error) {
	if r.r == nil {
		return []string{}, errorConnection
	}

	return r.r.Keys(pattern).Val(), nil
}

// SETS

// SetMemberAdd adds item to set
func (r *red) SetMemberAdd(key string, member ...string) error {
	if r.r == nil {
		return errorConnection
	}

	return r.r.SAdd(r.prefixed(key), member).Err()
}

// SetMemberExists checks if a member exists
func (r *red) SetMemberExists(key, member string) (bool, error) {
	if r.r == nil {
		return false, errorConnection
	}

	return r.r.SIsMember(r.prefixed(key), member).Result()
}

// SetMemberDelete deletes an member from a set
func (r *red) SetMemberDelete(key, member string) (int64, error) {
	return r.r.SRem(key, member).Result()
}

// SetMembers gets all members from a set
func (r *red) SetMembers(key string) ([]string, error) {
	return r.r.SMembers(r.prefixed(key)).Result()
}

// SetLength returns the length of a hash
func (r *red) SetLength(key string) (int64, error) {
	if r.r == nil {
		return 0, errorConnection
	}

	results, err := r.r.SMembers(r.prefixed(key)).Result()
	if err != nil {
		return 0, err
	}

	return int64(len(results)), nil
}

// HASH

// HashItemAdd adds item to hash
func (r *red) HashItemAdd(hash, item, value string) error {
	if r.r == nil {
		return errorConnection
	}

	return r.r.HSet(r.prefixed(hash), item, value).Err()
}

// HashItemExists checks if an item hash exists
func (r *red) HashItemExists(hash, item string) (bool, error) {
	if r.r == nil {
		return false, errorConnection
	}

	return r.r.HExists(r.prefixed(hash), item).Result()
}

// HashItemGet gets an item into a hash
func (r *red) HashItemGet(hash, item string) (string, error) {
	if r.r == nil {
		return "", errorConnection
	}

	return r.r.HGet(r.prefixed(hash), item).Result()
}

// HashItemDelete deletes an item from a hash
func (r *red) HashItemDelete(hash, item string) (int64, error) {
	return r.r.HDel(r.prefixed(hash), item).Result()
}

// HashItems returns all items from a hash
func (r *red) HashItems(hash string) (map[string]string, error) {
	if r.r == nil {
		return nil, errorConnection
	}

	return r.r.HGetAll(r.prefixed(hash)).Result()
}

// HashLength returns the length of a hash
func (r *red) HashLength(hash string) (int64, error) {
	if r.r == nil {
		return 0, errorConnection
	}

	return r.r.HLen(r.prefixed(hash)).Result()
}

// SCORE

// ScoreItemAdd adds item to score key
func (r *red) ScoreItemAdd(key string, score float64, member string) error {
	if r.r == nil {
		return errorConnection
	}

	return r.r.ZAdd(r.prefixed(key), redis.Z{
		Score:  score,
		Member: member,
	}).Err()
}

// ScoreItemCount counts item from score key
func (r *red) ScoreItemCount(key string, minScore, maxScore string) (float64, error) {
	if r.r == nil {
		return 0, errorConnection
	}

	opt := redis.ZRangeBy{
		Min: minScore,
		Max: maxScore,
	}

	cmd := r.r.ZRangeByScore(r.prefixed(key), opt)
	if err := cmd.Err(); err != nil {
		return 0, err
	}

	return float64(len(cmd.Val())), nil
}
