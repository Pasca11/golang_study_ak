package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
)

type Cacher interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}

type cache struct {
	client *redis.Client
}

func (c *cache) Set(key string, value interface{}) error {
	js, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = c.client.Set(key, js, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *cache) Get(key string) (interface{}, error) {
	value, err := c.client.Get(key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	return value, nil
}

func NewCache(client *redis.Client) Cacher {
	return &cache{
		client: client,
	}
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	// Создание клиента Redis
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	err := client.Ping().Err()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	cache := NewCache(client)
	// Установка значения по ключу
	err = cache.Set("some:key", "value")
	if err != nil {
		panic(err)
	}
	// Получение значения по ключу
	value, err := cache.Get("some:key")
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
	user := &User{
		ID:   1,
		Name: "John",
		Age:  30,
	}
	// Установка значения по ключу
	err = cache.Set(fmt.Sprintf("user:%v", user.ID), user)
	if err != nil {
		panic(err)
	}
	// Получение значения по ключу
	value, err = cache.Get("key")
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
}
