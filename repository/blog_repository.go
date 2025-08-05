

package repository

import (
    "context"
    "errors"
    "time"
    "g3-g65-bsp/domain"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Blog model in repository package

type BlogModel struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    AuthorID  primitive.ObjectID `bson:"author_id"`
    AuthorUsername string        `bson:"author_username"`
    Title     string        `bson:"title"`
    Content   string        `bson:"content"`
    Tags      []string      `bson:"tags"`
    Metrics   Metrics       `bson:"metrics"`
    Comments  []Comment     `bson:"comments"`
    CreatedAt time.Time     `bson:"created_at"`
    UpdatedAt time.Time     `bson:"updated_at"`
}

type Metrics struct {
    ViewCount int            `bson:"view_count"`
    Likes     Likes          `bson:"likes"`
}

type Likes struct {
    Count int            `bson:"count"`
    Users []string      `bson:"users"`
}

type Comment struct {
    ID             string        `bson:"id"`
    AuthorID       string        `bson:"author_id"`
    AuthorUsername string        `bson:"author_username"`
    Content        string        `bson:"content"`
    CreatedAt      time.Time     `bson:"created_at"`
}

func (m *BlogModel) ToDomain() *domain.Blog {
    comments := make([]domain.Comment, len(m.Comments))
    for i, c := range m.Comments {
        comments[i] = domain.Comment{
            ID:             c.ID,
            AuthorID:       c.AuthorID,
            AuthorUsername: c.AuthorUsername,
            Content:        c.Content,
            CreatedAt:      c.CreatedAt,
        }
    }
    return &domain.Blog{
        ID:             m.ID.Hex(),
        AuthorID:       m.AuthorID.Hex(),
        AuthorUsername: m.AuthorUsername,
        Title:          m.Title,
        Content:       m.Content,
        Tags:          m.Tags,
        Metrics:       domain.Metrics{
            ViewCount: m.Metrics.ViewCount,
            Likes:     domain.Likes{
                Count: m.Metrics.Likes.Count,
                Users: m.Metrics.Likes.Users,
            },
        },
        Comments:      comments,
        CreatedAt:     m.CreatedAt,
        UpdatedAt:     m.UpdatedAt,
    }
}

func (m *BlogModel) FromDomain(blog *domain.Blog) {
    var err error
    m.ID, err = primitive.ObjectIDFromHex(blog.ID)
    if err != nil {
        m.ID = primitive.NilObjectID
    }
    m.AuthorID, err = primitive.ObjectIDFromHex(blog.AuthorID)
    if err != nil {
        m.AuthorID = primitive.NilObjectID
    }
    m.AuthorUsername = blog.AuthorUsername
    m.Title = blog.Title
    m.Content = blog.Content
    m.Tags = blog.Tags
    m.Metrics = Metrics{
        ViewCount: blog.Metrics.ViewCount,
        Likes: Likes{
            Count: blog.Metrics.Likes.Count,
            Users: blog.Metrics.Likes.Users,
        },
    }
    m.Comments = make([]Comment, len(blog.Comments))
    for i, c := range blog.Comments {
        m.Comments[i] = Comment{
            ID:             c.ID,
            AuthorID:       c.AuthorID,
            AuthorUsername: c.AuthorUsername,
            Content:        c.Content,
            CreatedAt:      c.CreatedAt,
        }
    }
    m.CreatedAt = blog.CreatedAt
    m.UpdatedAt = blog.UpdatedAt
}

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
    var model BlogModel
    model.FromDomain(blog)
    model.ID = primitive.NewObjectID()
    model.CreatedAt = time.Now()
    model.UpdatedAt = model.CreatedAt

    res, err := r.collection.InsertOne(ctx, model)
    if err != nil {
        return "", err
    }

    oid, ok := res.InsertedID.(primitive.ObjectID)
    if !ok {
        return "", errors.New("failed to get inserted ID")
    }
    return oid.Hex(), nil
}

func (r *mongoBlogRepository) GetBlogByID(ctx context.Context, id string) (*domain.Blog, error) {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, ErrBlogNotFound
    }
    filter := bson.M{"_id": oid}
    var model BlogModel
    err = r.collection.FindOne(ctx, filter).Decode(&model)
    if err == mongo.ErrNoDocuments {
        return nil, ErrBlogNotFound
    }
    if err != nil {
        return nil, err
    }
    return model.ToDomain(), nil
}

func (r *mongoBlogRepository) UpdateBlog(ctx context.Context, blog *domain.Blog) error {
    var model BlogModel
    model.FromDomain(blog)
    model.UpdatedAt = time.Now()

    oid, err := primitive.ObjectIDFromHex(blog.ID)
    if err != nil {
        return ErrBlogNotFound
    }

    filter := bson.M{"_id": oid}
    update := bson.M{"$set": model}

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
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return ErrBlogNotFound
    }
    filter := bson.M{"_id": oid}
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
            andFilters = append(andFilters, bson.M{"created_at": bson.M{"$gte": fromTime}})
        }
    }

    if to, ok := filter["created_at_to"].(string); ok && to != "" {
        if toTime, err := time.Parse(time.RFC3339, to); err == nil {
            andFilters = append(andFilters, bson.M{"created_at": bson.M{"$lte": toTime}})
        }
    }
    if minViews, ok := filter["min_views"].(int); ok && minViews > 0 {
        andFilters = append(andFilters, bson.M{"metrics.view_count": bson.M{"$gte": minViews}})
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
        var model BlogModel
        if err := cur.Decode(&model); err != nil {
            return nil, nil, err
        }
        blogs = append(blogs, model.ToDomain())
    }
    if err := cur.Err(); err != nil {
        return nil, nil, err
    }

    pagination := &domain.Pagination{
        Total: int(total),
        Page:  page,
        Limit: limit,
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
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return ErrBlogNotFound
    }
    filter := bson.M{"_id": oid}
    update := bson.M{"$inc": bson.M{"metrics.view_count": 1}}
    result, err := r.collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }
    if result.MatchedCount == 0 {
        return ErrBlogNotFound
    }
    return nil
}

