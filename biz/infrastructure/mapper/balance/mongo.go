package balance

import (
	"context"
	"errors"
	"github.com/CloudStriver/cloudmind-trade/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-trade/biz/infrastructure/consts"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const CollectionName = "balance"

var PrefixBalanceCacheKey = "cache:balance:"

var _ IBalanceMongoMapper = (*MongoMapper)(nil)

type (
	IBalanceMongoMapper interface {
		Insert(ctx context.Context, data *Balance) (string, error)              // 插入
		FindOne(ctx context.Context, id string) (*Balance, error)               // 查找
		Update(ctx context.Context, data *Balance) (*mongo.UpdateResult, error) // 修改
		Delete(ctx context.Context, id string) (int64, error)                   // 删除
	}
	Balance struct {
		ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // 用户ID
		Flow   *int64             `bson:"flow,omitempty" json:"flow,omitempty"`
		Memory *int64             `bson:"memory,omitempty" json:"memory,omitempty"`
		Point  *int64             `bson:"point,omitempty" json:"point,omitempty"`
	}

	MongoMapper struct {
		conn *monc.Model
	}
)

func NewMongoMapper(config *config.Config) IBalanceMongoMapper {
	conn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, CollectionName, config.CacheConf)
	return &MongoMapper{
		conn: conn,
	}
}

func (m *MongoMapper) Insert(ctx context.Context, data *Balance) (string, error) {
	if data.ID.IsZero() {
		data.ID = primitive.NewObjectID()
	}

	key := PrefixBalanceCacheKey + data.ID.Hex()
	ID, err := m.conn.InsertOne(ctx, key, data)
	if err != nil {
		return "", err
	}
	return ID.InsertedID.(primitive.ObjectID).Hex(), err
}

func (m *MongoMapper) FindOne(ctx context.Context, id string) (*Balance, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, consts.ErrInvalidObjectId
	}
	var data Balance
	key := PrefixBalanceCacheKey + id
	err = m.conn.FindOne(ctx, key, &data, bson.M{consts.ID: oid})
	switch {
	case err == nil:
		return &data, nil
	case errors.Is(err, monc.ErrNotFound):
		return nil, consts.ErrNotFound
	default:
		return nil, err
	}
}

func (m *MongoMapper) Update(ctx context.Context, data *Balance) (*mongo.UpdateResult, error) {
	key := PrefixBalanceCacheKey + data.ID.Hex()
	b := bson.M{}
	filter := bson.M{consts.ID: data.ID}
	if data.Flow != nil {
		b[consts.Flow] = data.Flow
		filter[consts.Flow] = bson.M{"$gte": data.Flow}
	}
	if data.Memory != nil {
		b[consts.Memory] = data.Memory
		filter[consts.Memory] = bson.M{"$gte": data.Memory}
	}
	if data.Point != nil {
		b[consts.Point] = data.Point
		filter[consts.Point] = bson.M{"$gte": data.Point}
	}
	res, err := m.conn.UpdateOne(ctx, key, filter, bson.M{"$inc": b})
	return res, err
}

func (m *MongoMapper) Delete(ctx context.Context, id string) (int64, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}
	key := PrefixBalanceCacheKey + id
	res, err := m.conn.DeleteOne(ctx, key, bson.M{consts.ID: oid})
	return res, err
}
