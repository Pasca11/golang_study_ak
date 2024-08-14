package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type SomeRepository interface {
	GetData() string
}
type SomeRepositoryImpl struct{}

func (r *SomeRepositoryImpl) GetData() string {
	// Здесь происходит запрос к базе данных
	time.Sleep(1 * time.Second)
	return "data"
}

type SomeRepositoryProxy struct {
	repository SomeRepository
	cache      redis.Client
}

func (r *SomeRepositoryProxy) GetData() string {
	// Здесь происходит проверка наличия данных в кэше
	// Если данные есть в кэше, то они возвращаются
	// Если данных нет в кэше, то они запрашиваются у оригинального объекта и сохраняются в кэш
	c, err := r.cache.Get("data").Result()
	if err == nil {
		return c
	}
	data := r.repository.GetData()
	r.cache.Set("data", data, time.Minute)
	return r.repository.GetData()
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()
	proxy := &SomeRepositoryProxy{
		repository: &SomeRepositoryImpl{},
		cache:      *client,
	}
	start := time.Now()
	proxy.GetData()
	end := time.Now()
	fmt.Println("First time", end.Sub(start))

	start = time.Now()
	proxy.GetData()
	end = time.Now()
	fmt.Println("Second time", end.Sub(start))
}
