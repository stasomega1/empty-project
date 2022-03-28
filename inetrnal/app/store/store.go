package store

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"project/inetrnal/app/model"
)

type StoreI interface {
	SqlRepository() SqlRepositoryI
	RedisRepository() RedisRepositoryI
	RabbitPublisher() RabbitPublisherI
}

type SqlRepositoryI interface {
	GetSomeDataFromDB(ctx context.Context, parameter1 string, parameter2 int) (*model.DbModel, error)
}

type RedisRepositoryI interface {
	SetSomeDBCache(ctx context.Context, parameter1 string, parameter2 int, result *model.DbModel) error
	GetSomeDBCache(ctx context.Context, parameter1 string, parameter2 int) (*model.DbModel, error)
}

type RabbitPublisherI interface {
	SendSomeDbMessage(model model.DbModel) error
}

type Store struct {
	db              *sqlx.DB
	sqlRepository   SqlRepositoryI
	redisRepository RedisRepositoryI
	rabbitPublisher RabbitPublisherI
}

func NewStore(db *sqlx.DB, redis *redis.Client, rabbitPublisher RabbitPublisherI) *Store {
	redisRepository := &RedisRepository{Redis: redis}
	sqlRepository := &SqlRepository{db: db}
	return &Store{
		db:              db,
		redisRepository: redisRepository,
		rabbitPublisher: rabbitPublisher,
		sqlRepository:   sqlRepository,
	}
}

func (s *Store) SqlRepository() SqlRepositoryI {
	return s.sqlRepository
}

func (s *Store) RedisRepository() RedisRepositoryI {
	return s.redisRepository
}

func (s *Store) RabbitPublisher() RabbitPublisherI {
	return s.rabbitPublisher
}
