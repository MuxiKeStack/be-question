package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MuxiKeStack/be-question/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExist = redis.Nil

type QuestionCache interface {
	Set(ctx context.Context, question domain.Question) error
	Get(ctx context.Context, id int64) (domain.Question, error)
}

type RedisQuestionCache struct {
	cmd redis.Cmdable
}

// Set 这里缓存是给消费者监听到消息发送才会去预缓存的，这需要feed来发消息
func (cache *RedisQuestionCache) Set(ctx context.Context, question domain.Question) error {
	key := cache.key(question.Id)
	val, err := json.Marshal(question)
	if err != nil {
		return err
	}
	return cache.cmd.Set(ctx, key, val, time.Minute*30).Err()
}

func (cache *RedisQuestionCache) Get(ctx context.Context, id int64) (domain.Question, error) {
	key := cache.key(id)
	val, err := cache.cmd.Get(ctx, key).Bytes()
	if err != nil {
		return domain.Question{}, err
	}
	var q domain.Question
	err = json.Unmarshal(val, &q)
	return q, err
}

func (cache *RedisQuestionCache) key(id int64) string {
	return fmt.Sprintf("kstack:question:%d", id)
}

func NewRedisQuestionCache(cmd redis.Cmdable) QuestionCache {
	return &RedisQuestionCache{cmd: cmd}
}
