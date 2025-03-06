package models

import "time"

type WebViewServer struct {
        ID        string    `bson:"_id,omitempty" json:"_id"`
        CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
        UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
        Name      string    `bson:"name" json:"name"`
        Status    string    `bson:"status" json:"status"`
}
