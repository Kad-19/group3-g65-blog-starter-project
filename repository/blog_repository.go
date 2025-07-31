package repository

import (
    "context"
    "errors"
    "time"
    "g3-g65-bsp/domain"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson/primitive"
)


type mongoBlogRepository struct {
    collection *mongo.Collection
}

// ErrBlogNotFound is returned when a blog is not found in the repository
var ErrBlogNotFound = errors.New("blog not found")

// NewBlogRepository returns a MongoDB implementation of BlogRepository
func NewBlogRepository(collection *mongo.Collection) domain.BlogRepository {
    return &mongoBlogRepository{collection: collection}
}

func (r *mongoBlogRepository) CreateBlog(ctx context.Context, blog *domain.Blog) (string, error) {
    blog.ID = primitive.NewObjectID().Hex()
    blog.CreatedAt = time.Now()
    blog.UpdatedAt = blog.CreatedAt
    _, err := r.collection.InsertOne(ctx, blog)
    if err != nil {
        return "", err
    }
    return blog.ID, nil
}

func (r *mongoBlogRepository) GetBlogByID(ctx context.Context, id string) (*domain.Blog, error) {
    var blog domain.Blog
    filter := bson.M{"id": id}
    err := r.collection.FindOne(ctx, filter).Decode(&blog)
    if err == mongo.ErrNoDocuments {
        return nil, ErrBlogNotFound
    }
    if err != nil {
        return nil, err
    }
    return &blog, nil
}

func (r *mongoBlogRepository) UpdateBlog(ctx context.Context, blog *domain.Blog) error {
    blog.UpdatedAt = time.Now()
    filter := bson.M{"id": blog.ID}
    update := bson.M{"$set": blog}
    result, err := r.collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }
    if result.MatchedCount == 0 {
        return ErrBlogNotFound
    }
    return nil
}

func (r *mongoBlogRepository) DeleteBlog(ctx context.Context, id string) error {
    filter := bson.M{"id": id}
    result, err := r.collection.DeleteOne(ctx, filter)
    if err != nil {
        return err
    }
    if result.DeletedCount == 0 {
        return ErrBlogNotFound
    }
    return nil
}

func (r *mongoBlogRepository) ListBlogs(ctx context.Context, filter map[string]interface{}) ([]*domain.Blog, error) {
    bsonFilter := bson.M(filter)
    cur, err := r.collection.Find(ctx, bsonFilter)
    if err != nil {
        return nil, err
    }
    defer cur.Close(ctx)
    var blogs []*domain.Blog
    for cur.Next(ctx) {
        var blog domain.Blog
        if err := cur.Decode(&blog); err != nil {
            return nil, err
        }
        blogs = append(blogs, &blog)
    }
    if err := cur.Err(); err != nil {
        return nil, err
    }
    return blogs, nil
}
