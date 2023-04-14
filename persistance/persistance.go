package persistance

import (
	"fmt"
	"os"
	"time"
	"todobucket/model"
	"todobucket/utils"
)

type BucketKey string

const (
	BucketKeyGeneral BucketKey = "general"
	BucketKeyToday   BucketKey = "today"
)

type Storage struct {
	ItemSequence *Sequence
	ItemList     []model.Item
	ItemDoneList []model.Item
	BucketList   []model.Bucket
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
	o.BucketList = []model.Bucket{
		{
			Id:   utils.PInt64(1),
			Name: utils.PStr(string(BucketKeyGeneral)),
		},
		{
			Id:   utils.PInt64(2),
			Name: utils.PStr(string(BucketKeyToday)),
		},
	}
}
func (o *Storage) Save(file *string) {
	content := *o
	marshalled := utils.Marshal(&content)
	utils.WriteToFile(*file, string(marshalled))
}
func (o *Storage) load(file string) {
	exists := utils.FileExists(file)
	if exists {
		data, err := os.ReadFile(file)
		utils.Panic(err)
		utils.Unmarshal(data, &o)
	}
}
func (o *Storage) CreateItem(title string, bucketKey BucketKey) *model.Item {
	item := model.Item{
		Id:           o.ItemSequence.Next(),
		BucketId:     o.FindBucketByKey(bucketKey).Id,
		Title:        &title,
		CreationDate: utils.PTime(time.Now()),
	}
	o.ItemList = append(o.ItemList, item)
	return &item
}
func (o *Storage) FindBucketByKey(key BucketKey) *model.Bucket {
	searchedBucket := string(key)
	for _, cur := range o.BucketList {
		if *cur.Name == searchedBucket {
			return &cur
		}
	}
	panic(fmt.Sprintf("Bucket key %s not supported", key))
}
func (o *Storage) ListItem() []model.Item {
	return o.ItemList
}

func (o *Storage) FindItem(itemId *int64) (*model.Item, *int) {
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
		item.DoneDate = utils.PTime(time.Now())
	}
	o.ItemList[*idx] = *item
}

func (o *Storage) EditItem(itemId *int64, item *model.Item) {
	existingItem, idx := o.FindItem(itemId)
	if existingItem == nil {
		panic(fmt.Sprintf("Item id '%d' not found.", *itemId))
	}
	o.ItemList[*idx] = *item
}

func (o *Storage) DeleteItem(itemId *int64) {
	_, idx := o.FindItem(itemId)
	if idx == nil {
		return
	}

	result := make([]model.Item, 0)
	for curIdx, cur := range o.ItemList {
		if *idx != curIdx {
			result = append(result, cur)
		}
	}
}
