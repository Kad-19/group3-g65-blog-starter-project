

package repository

import (
    "context"
    "errors"
    "time"
    "g3-g65-bsp/domain"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
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

func (r *mongoBlogRepository) ListBlogs(ctx context.Context, filter map[string]any, page, limit int) ([]*domain.Blog, *domain.Pagination, error) {
    var andFilters []bson.M

    if search, ok := filter["search"].(string); ok && search != "" {
        orFilters := []bson.M{
            {"title": bson.M{"$regex": search, "$options": "i"}},
            {"authorusername": bson.M{"$regex": search, "$options": "i"}},
        }
        andFilters = append(andFilters, bson.M{"$or": orFilters})
    }

    // Add other specific filters to the $and array

    if tags, ok := filter["tags"].([]string); ok && len(tags) > 0 {
        andFilters = append(andFilters, bson.M{"tags": bson.M{"$all": tags}})
    }

    if from, ok := filter["created_at_from"].(string); ok && from != "" {
        if fromTime, err := time.Parse(time.RFC3339, from); err == nil {
            andFilters = append(andFilters, bson.M{"createdat": bson.M{"$gte": fromTime}})
        }
    }

    if to, ok := filter["created_at_to"].(string); ok && to != "" {
        if toTime, err := time.Parse(time.RFC3339, to); err == nil {
            andFilters = append(andFilters, bson.M{"createdat": bson.M{"$lte": toTime}})
        }
    }
    if minViews, ok := filter["min_views"].(int); ok && minViews > 0 {
        andFilters = append(andFilters, bson.M{"metrics.viewcount": bson.M{"$gte": minViews}})
    }

    // ... add other filters similarly ...


    bsonFilter := bson.M{}
    if len(andFilters) > 0 {
        bsonFilter["$and"] = andFilters
    }

    // Get total count for pagination
    total, err := r.collection.CountDocuments(ctx, bsonFilter)
    if err != nil {
        return nil, nil, err
    }

    opts := options.Find()
    if limit > 0 {
        opts.SetLimit(int64(limit))
    }
    if page > 1 && limit > 0 {
        opts.SetSkip(int64((page - 1) * limit))
    }

    cur, err := r.collection.Find(ctx, bsonFilter, opts)
    if err != nil {
        return nil, nil, err
    }
    defer cur.Close(ctx)
    var blogs []*domain.Blog
    for cur.Next(ctx) {
        var blog domain.Blog
        if err := cur.Decode(&blog); err != nil {
            return nil, nil, err
        }
        blogs = append(blogs, &blog)
    }
    if err := cur.Err(); err != nil {
        return nil, nil, err
    }

    pagination := &domain.Pagination{
        Total: int(total),
        Page:  page,
        Limit: limit,
    }
    return blogs, pagination, nil
}
// IncrementBlogViewCount atomically increments the view count of a blog post
func (r *mongoBlogRepository) IncrementBlogViewCount(ctx context.Context, id string, blog *domain.Blog) error {
    filter := bson.M{"id": blog.ID}
    update := bson.M{"$inc": bson.M{"metrics.viewcount": 1}}
    result, err := r.collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }
    if result.MatchedCount == 0 {
        return ErrBlogNotFound
    }
    return nil
}

