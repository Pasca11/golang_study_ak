package cache

import (
	"context"
	"errors"
	"geoservice/internal/repository"
	"geoservice/models"
	"github.com/go-redis/redis"
)

type ProxyCache struct {
	DataBase repository.UserRepository
	Client   *redis.Client
}

func NewCache(ctx context.Context) (*ProxyCache, error) {
	database, err := repository.StartPostgresDataBase(ctx)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	newCache := &ProxyCache{
		DataBase: database,
		Client:   client,
	}

	return newCache, nil
}

func (pc *ProxyCache) Create(ctx context.Context, user models.User) error {
	return pc.DataBase.Create(ctx, user)
}

func (pc *ProxyCache) GetByID(ctx context.Context, id string) (models.User, error) {
	return pc.DataBase.GetByID(ctx, id)
}

func (pc *ProxyCache) Update(ctx context.Context, user models.User) error {
	return pc.DataBase.Update(ctx, user)
}

func (pc *ProxyCache) Delete(ctx context.Context, id string) error {
	return pc.DataBase.Delete(ctx, id)
}

func (pc *ProxyCache) List(ctx context.Context) ([]models.User, error) {
	return pc.DataBase.List(ctx)
}

func (pc *ProxyCache) GetByName(ctx context.Context, name string) (models.User, bool, error) {
	user := models.User{}
	result, err := pc.Client.Get(name).Result()

	if errors.Is(err, redis.Nil) {
		user, exist, err := pc.DataBase.GetByName(ctx, name)
		if !exist {
			return user, false, err
		}

		err = pc.Client.Set(user.Username, user.Password, 0).Err()
		if err != nil {
			return user, true, err
		}

		return user, true, nil
	} else if err != nil {
		return user, false, err
	}

	user.Username = name
	user.Password = result
	return user, true, nil
}
