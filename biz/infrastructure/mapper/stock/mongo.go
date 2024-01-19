package stock

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

const CollectionName = "Stock"

var PrefixUserCacheKey = "cache:Stock:"

var _ IStockMongoMapper = (*MongoMapper)(nil)

type (
	IStockMongoMapper interface {
		Insert(ctx context.Context, data *Stock) (string, error)                                // 插入
		FindOne(ctx context.Context, id string) (*Stock, error)                                 // 查找
		Update(ctx context.Context, data *Stock, oldAmount *int64) (*mongo.UpdateResult, error) // 修改
		Delete(ctx context.Context, id string) (int64, error)                                   // 删除
		FindMany(ctx context.Context, ids []string) ([]*Stock, error)
	}

	Stock struct {
		ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // 商品ID
		Amount *int64             `bson:"amount,omitempty" json:"amount,omitempty"`
	}

	MongoMapper struct {
		conn *monc.Model
	}
)

func (m *MongoMapper) FindMany(ctx context.Context, ids []string) ([]*Stock, error) {
	var data []*Stock
	var oids []primitive.ObjectID
	for _, id := range ids {
		oid, _ := primitive.ObjectIDFromHex(id)
		oids = append(oids, oid)
	}
	err := m.conn.Find(ctx, &data, bson.M{consts.ID: bson.M{"$in": oids}})
	switch {
	case err == nil:
		return data, nil
	case errors.Is(err, monc.ErrNotFound):
		return nil, consts.ErrNotFound
	default:
		return nil, err
	}
}

func NewMongoMapper(config *config.Config) IStockMongoMapper {
	conn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, CollectionName, config.CacheConf)
	return &MongoMapper{
		conn: conn,
	}
}

func (m *MongoMapper) Insert(ctx context.Context, data *Stock) (string, error) {
	if data.ID.IsZero() {
		data.ID = primitive.NewObjectID()
	}

	key := PrefixUserCacheKey + data.ID.Hex()
	ID, err := m.conn.InsertOne(ctx, key, data)
	if err != nil {
		return "", err
	}
	return ID.InsertedID.(primitive.ObjectID).Hex(), err
}

func (m *MongoMapper) FindOne(ctx context.Context, id string) (*Stock, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, consts.ErrInvalidObjectId
	}
	var data Stock
	key := PrefixUserCacheKey + id
	err = m.conn.FindOne(ctx, key, &data, bson.M{"_id": oid})
	switch {
	case err == nil:
		return &data, nil
	case errors.Is(err, monc.ErrNotFound):
		return nil, consts.ErrNotFound
	default:
		return nil, err
	}
}

func (m *MongoMapper) Update(ctx context.Context, data *Stock, oldAmount *int64) (*mongo.UpdateResult, error) {
	key := PrefixUserCacheKey + data.ID.Hex()
	res, err := m.conn.UpdateOne(ctx, key, bson.M{consts.ID: data.ID, consts.Amount: oldAmount}, bson.M{"$set": data})
	return res, err
}

func (m *MongoMapper) Delete(ctx context.Context, id string) (int64, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}
	key := PrefixUserCacheKey + id
	res, err := m.conn.DeleteOne(ctx, key, bson.M{"_id": oid})
	return res, err
}
