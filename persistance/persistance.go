package persistance

import (
	"fmt"
	"itembucket/common"
	"os"
	"time"
)

type BucketKey string

const (
	BucketKeyGeneral BucketKey = "general"
	BucketKeyToday   BucketKey = "today"
)

type Storage struct {
	ItemSequence *Sequence
	ItemList     []Item
	ItemDoneList []Item
	BucketList   []Bucket
}

func NewStorage(file *string) *Storage {
	storage := &Storage{
		ItemSequence: NewSequence(nil),
	}
	storage.initBucketList()
	storage.load(*file)
	return storage
}

func (o *Storage) initBucketList() {
	o.BucketList = []Bucket{
		{
			Id:   common.PInt64(1),
			Name: common.PStr(string(BucketKeyGeneral)),
		},
		{
			Id:   common.PInt64(2),
			Name: common.PStr(string(BucketKeyToday)),
		},
	}
}

func (o *Storage) Save(file *string) {
	content := *o
	marshalled := common.Marshal(&content)
	common.WriteToFile(*file, string(marshalled))
}

func (o *Storage) load(file string) {
	exists := common.FileExists(file)
	if exists {
		data, err := os.ReadFile(file)
		common.Panic(err)
		common.Unmarshal(data, &o)
	}
}

func (o *Storage) CreateItem(title string, bucketKey BucketKey) *Item {
	id := o.ItemSequence.Next()
	item := Item{
		Id:           id,
		OrderIdx:     id,
		BucketId:     o.FindBucketByKey(bucketKey).Id,
		Title:        &title,
		CreationDate: common.PTime(time.Now()),
	}
	o.ItemList = append(o.ItemList, item)
	return &item
}

func (o *Storage) FindBucketByKey(key BucketKey) *Bucket {
	searchedBucket := string(key)
	for _, cur := range o.BucketList {
		if *cur.Name == searchedBucket {
			return &cur
		}
	}
	panic(fmt.Sprintf("Bucket key %s not supported", key))
}

func (o *Storage) FindBucketKeyById(id *int64) *BucketKey {
	bucketId := *id - 1
	response := BucketKey(*o.BucketList[bucketId].Name)
	return &response
}

func (o *Storage) ListItem() []Item {
	return o.ItemList
}

func (o *Storage) FindItem(itemId *int64) (*Item, *int) {
	for idx, cur := range o.ItemList {
		if *cur.Id == *itemId {
			return &cur, &idx
		}
	}
	return nil, nil
}

func (o *Storage) ToggleDone(itemId *int64) {
	item, idx := o.FindItem(itemId)
	if item == nil {
		panic(fmt.Sprintf("Item id '%d' not found.", *itemId))
	}

	if item.DoneDate != nil {
		item.DoneDate = nil
	} else {
		item.DoneDate = common.PTime(time.Now())
	}
	o.ItemList[*idx] = *item
}

func (o *Storage) EditItem(itemId *int64, item *Item) {
	existingItem, idx := o.FindItem(itemId)
	if existingItem == nil {
		panic(fmt.Sprintf("Item id '%d' not found.", *itemId))
	}
	o.ItemList[*idx] = *item
}
