package domain

import (
	"time"
)

type Blog struct {
    ID        string   
    AuthorID  string   
    Title     string   
    Content   string   
    Tags      []string 
    Metrics   Metrics  
    Comments  []Comment 
    CreatedAt time.Time            
    UpdatedAt time.Time            
}

type Metrics struct {
    ViewCount int    
    Likes     Likes  
}

type Likes struct {
    Count int                  
    Users []string 
}

type Comment struct {
    ID             string
    AuthorID       string    
    AuthorUsername string    
    Content        string   
    CreatedAt      time.Time 
}
