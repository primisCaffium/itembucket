package model

import "time"

type Item struct {
	Id           *int64
	BucketId     *int64
	Title        *string
	Description  *string
	CreationDate *time.Time
	DoneDate     *time.Time
}
type Bucket struct {
	Id   *int64
	Name *string
}
