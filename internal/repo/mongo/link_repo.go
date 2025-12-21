package mongo

import (
	"context"
	"time"

	"github.com/oanhngy/url-shortener/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodrv "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LinkRepo struct {
	col *mongodrv.Collection // MongoDB collection for links
}

type linkDoc struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	ShortCode  string             `bson:"short_code"` //unique
	LongURL    string             `bson:"long_url"`
	ClickCount int                `bson:"click_count"`
	CreatedAt  time.Time          `bson:"created_at"`
}

func NewMongoRepo(db *mongodrv.Database) *LinkRepo {
	return &LinkRepo{
		col: db.Collection("links"),
	}
}

func (r *LinkRepo) EnsureIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	index := mongodrv.IndexModel{
		Keys: bson.D{{Key: "short_code", Value: 1}},
		Options: options.Index().
			SetUnique(true).
			SetName("u_short_code"),
	}

	_, err := r.col.Indexes().CreateOne(ctx, index)
	return err
}

// SAVE NEW LINK
func (r *LinkRepo) Save(link *model.Link) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := linkDoc{
		ShortCode:  link.ShortCode,
		LongURL:    link.LongURL,
		ClickCount: link.ClickCount,
		CreatedAt:  link.CreatedAt,
	}

	res, err := r.col.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		link.ID = oid.Hex()
	}

	return nil
}

// FIND BY CODE
func (r *LinkRepo) FindByCode(code string) (*model.Link, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var doc linkDoc
	err := r.col.FindOne(ctx, bson.M{"short_code": code}).Decode(&doc)
	if err != nil {
		return nil, err //404 nếu k thấy
	}

	return &model.Link{
		ID:         doc.ID.Hex(),
		LongURL:    doc.LongURL,
		ShortCode:  doc.ShortCode,
		CreatedAt:  doc.CreatedAt,
		ClickCount: doc.ClickCount,
	}, nil
}

// EXIST, chech tồn tại chưa, tránh collision
func (r *LinkRepo) Exists(code string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	n, err := r.col.CountDocuments(ctx, bson.M{"short_code": code})
	if err != nil {
		return true
	}
	return n > 0
}

// FIND ALL
func (r *LinkRepo) FindAll() ([]model.Link, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := r.col.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	out := make([]model.Link, 0)

	for cur.Next(ctx) {
		var doc linkDoc
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}

		out = append(out, model.Link{
			ID:         doc.ID.Hex(),
			LongURL:    doc.LongURL,
			ShortCode:  doc.ShortCode,
			CreatedAt:  doc.CreatedAt,
			ClickCount: doc.ClickCount,
		})
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

// INCREMENT CLICK
func (r *LinkRepo) IncrementClick(code string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.col.UpdateOne(
		ctx,
		bson.M{"short_code": code},
		bson.M{"$inc": bson.M{"click_count": 1}},
	)

	return err
}
