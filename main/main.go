package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"sort"
	"todobucket/model"
	"todobucket/persistance"
	"todobucket/utils"
)

func main() {
	addItemFlag := flag.String("add", "", "Specify a title for adding a new item into the general bucket.")
	addItemTodayFlag := flag.String("addToday", "", "Specify a title for adding a new item into the today bucket.")
	listItemsFlag := flag.String("list", "", "Specify 'general' or 'today' to list all active items.")
	markItemDoneFlag := flag.Int64("toggle", -1, "Specify the item id to mark it as done.")
	moveItemToTodayFlag := flag.Int64("moveToday", -1, "Specify the item id to move it from general bucket to today.")
	moveItemToGeneralFlag := flag.Int64("moveGeneral", -1, "Specify the item id to move it from today bucket to general.")
	emptyTodayFlag := flag.Bool("emptyToday", false, "Moves all items in today list to general.")
	cleanupFlag := flag.Bool("cleanup", false, "Deletes all checked items.")
	editItemFlag := flag.Int64("edit", -1, "Specify the id of the item to edit. This argument needs to be followed by -text argument.")
	textFlag := flag.String("text", "", "Specify the text for editing an item. This argument needs to be following the -edit argument.")
	printStorageFlag := flag.Bool("printStorage", false, "Prints the storage path.")

	flag.Parse()
	isListArg := len(*listItemsFlag) > 0

	dirname, err := os.UserHomeDir()
	utils.Panic(err)
	storageFilePath := path.Join(dirname, ".ibstorage.json")
	tool := NewTool(&storageFilePath)
	switch {
	case isListArg:
		tool.ListItem(persistance.BucketKey(*listItemsFlag))
	case len(*addItemFlag) > 0:
		tool.AddItem(*addItemFlag, persistance.BucketKeyGeneral)
	case len(*addItemTodayFlag) > 0:
		tool.AddItem(*addItemTodayFlag, persistance.BucketKeyToday)
	case *markItemDoneFlag > -1:
		tool.ToggleItemCheck(markItemDoneFlag)
	case *moveItemToTodayFlag > -1:
		tool.MoveItemToTodayList(moveItemToTodayFlag)
	case *moveItemToGeneralFlag > -1:
		tool.MoveItemToGeneralList(moveItemToGeneralFlag)
	case *emptyTodayFlag:
		tool.EmptyToday()
	case *cleanupFlag:
		tool.CleanupCheckedItems()
	case *editItemFlag > -1:
		if len(*textFlag) == 0 {
			panic("Missing -text argument")
		}
		tool.EditItem(editItemFlag, textFlag)
	case *printStorageFlag:
		fmt.Printf("IB storage file: %s, you can backup this.\n", storageFilePath)
		os.Exit(0)
	}

	if !isListArg {
		tool.Storage.Save(tool.StorageFile)

		fmt.Printf("GENERAL:\n")
		tool.ListItem(persistance.BucketKeyGeneral)

		fmt.Printf("\nTODAY:\n")
		tool.ListItem(persistance.BucketKeyToday)
		fmt.Printf("\n")
	}
}

type Tool struct {
	StorageFile *string
	Storage     *persistance.Storage
}

func NewTool(storageFile *string) *Tool {
	return &Tool{
		StorageFile: storageFile,
		Storage:     persistance.NewStorage(storageFile),
	}
}

func (o *Tool) AddItem(title string, bucketKey persistance.BucketKey) *model.Item {
	item := o.Storage.CreateItem(title, bucketKey)
	return item
}

func (o *Tool) EditItem(itemId *int64, text *string) {
	item, _ := o.Storage.FindItem(itemId)
	item.Title = text
	o.Storage.EditItem(itemId, item)
}

func (o *Tool) ListItem(bucketKey persistance.BucketKey) []model.Item {
	list := o.Storage.ListItem()
	listOfGivenBucket := make([]model.Item, 0)
	bucketId := o.Storage.FindBucketByKey(bucketKey).Id
	doneList := make([]model.Item, 0)
	pendingList := make([]model.Item, 0)
	for _, cur := range list {
		if *cur.BucketId == *bucketId {
			if cur.DoneDate != nil {
				doneList = append(doneList, cur)
			} else {
				pendingList = append(pendingList, cur)
			}
		}
	}

	sort.Slice(doneList, func(i, j int) bool {
		return *doneList[i].Id < *doneList[j].Id
	})
	sort.Slice(pendingList, func(i, j int) bool {
		return *pendingList[i].Id < *pendingList[j].Id
	})

	listOfGivenBucket = pendingList
	listOfGivenBucket = append(listOfGivenBucket, doneList...)

	for _, cur := range listOfGivenBucket {
		doneChar := " "
		if cur.DoneDate != nil {
			doneChar = "x"
		}
		id := fmt.Sprintf("%d", *cur.Id)
		fmt.Printf("%s-[%s]-%s\n", id, doneChar, *cur.Title)
	}
	return listOfGivenBucket
}

func (o *Tool) ToggleItemCheck(itemId *int64) {
	o.Storage.ToggleDone(itemId)
}
func (o *Tool) MoveItemToTodayList(itemId *int64) {
	o.doMoveItemToBucket(itemId, o.Storage.FindBucketByKey(persistance.BucketKeyToday).Id)
}
func (o *Tool) MoveItemToGeneralList(itemId *int64) {
	o.doMoveItemToBucket(itemId, o.Storage.FindBucketByKey(persistance.BucketKeyGeneral).Id)
}
func (o *Tool) doMoveItemToBucket(itemId *int64, bucketId *int64) {
	item, idx := o.Storage.FindItem(itemId)
	if idx == nil {
		panic(fmt.Sprintf("Item id '%d' doesn't exist", *itemId))
	}
	item.BucketId = bucketId
	o.Storage.ItemList[*idx] = *item
}
func (o *Tool) EmptyToday() {
	todayBucketId := *o.Storage.FindBucketByKey(persistance.BucketKeyToday).Id
	generalBucketId := *o.Storage.FindBucketByKey(persistance.BucketKeyGeneral).Id
	for idx, cur := range o.Storage.ItemList {
		if *cur.BucketId == todayBucketId {
			o.Storage.ItemList[idx].BucketId = &generalBucketId
		}
	}
}
func (o *Tool) CleanupCheckedItems() {
	result := make([]model.Item, 0)
	for _, cur := range o.Storage.ItemList {
		if cur.DoneDate == nil {
			result = append(result, cur)
		}
	}
	o.Storage.ItemList = result
}
