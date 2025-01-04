package main

import (
	"itembucket/persistance"
	"log"
	"os"
	"path"
	"testing"
)

func setupTestStorage(t *testing.T) (*Tool, string) {
	t.Helper()
	storagePath := path.Join(os.TempDir(), "test_ibstorage.json")
	tool := NewTool(&storagePath)

	err := os.Remove(storagePath)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("Error setting up test storage: %s", err)
	}

	return tool, storagePath
}

func TestAddItem(t *testing.T) {
	tool, _ := setupTestStorage(t)

	title := "Test item"
	item := tool.AddItem(title, persistance.BucketKeyGeneral)

	if item == nil {
		t.Fatal("Failed to add item")
	}
	if *item.Title != title || *item.BucketId != *tool.Storage.FindBucketByKey(persistance.BucketKeyGeneral).Id {
		t.Fatalf("Unexpected item details: %+v", item)
	}
}

func TestListItems(t *testing.T) {
	tool, _ := setupTestStorage(t)

	titleG1 := "General 1"
	titleG2 := "General 2"
	titleG3 := "General 3"
	titleT1 := "Today 1"
	titleT2 := "Today 2"
	tool.AddItem(titleG1, persistance.BucketKeyGeneral)
	g2 := tool.AddItem(titleG2, persistance.BucketKeyGeneral)
	g3 := tool.AddItem(titleG3, persistance.BucketKeyGeneral)

	tool.AddItem(titleT1, persistance.BucketKeyToday)
	tool.AddItem(titleT2, persistance.BucketKeyToday)

	tool.ToggleItemCheck(g2.Id)
	tool.ToggleItemCheck(g3.Id)

	generalItems := tool.ListItem(persistance.BucketKeyGeneral)
	if len(generalItems) != 3 || (*generalItems[0].Title != titleG1 && *generalItems[1].Title != titleG2 && *generalItems[2].Title != titleG3) {
		t.Fatalf("Unexpected general items: %+v", generalItems)
	}

	todayItems := tool.ListItem(persistance.BucketKeyToday)
	if len(todayItems) != 2 || (*todayItems[0].Title != titleT1 && *todayItems[1].Title != titleT2) {
		t.Fatalf("Unexpected today items: %+v", todayItems)
	}
}

func TestToggleItemCheck(t *testing.T) {
	tool, _ := setupTestStorage(t)

	item := tool.AddItem("Toggle me", persistance.BucketKeyGeneral)
	tool.ToggleItemCheck(item.Id)

	toggledItem, _ := tool.Storage.FindItem(item.Id)
	if toggledItem.DoneDate == nil {
		t.Fatal("Item was not marked as done")
	}

	tool.ToggleItemCheck(item.Id)
	toggledItem, _ = tool.Storage.FindItem(item.Id)
	if toggledItem.DoneDate != nil {
		t.Fatal("Item was not toggled back to pending")
	}
}

func TestMoveItemBetweenBuckets(t *testing.T) {
	tool, _ := setupTestStorage(t)

	item := tool.AddItem("Move me", persistance.BucketKeyGeneral)
	tool.MoveItemToTodayList(item.Id)

	updatedItem, _ := tool.Storage.FindItem(item.Id)
	if *updatedItem.BucketId != *tool.Storage.FindBucketByKey(persistance.BucketKeyToday).Id {
		t.Fatalf("Item was not moved to the 'Today' bucket: %+v", updatedItem)
	}

	tool.MoveItemToGeneralList(item.Id)
	updatedItem, _ = tool.Storage.FindItem(item.Id)
	if *updatedItem.BucketId != *tool.Storage.FindBucketByKey(persistance.BucketKeyGeneral).Id {
		t.Fatalf("Item was not moved to the 'General' bucket: %+v", updatedItem)
	}
}

func TestEditItem(t *testing.T) {
	tool, _ := setupTestStorage(t)

	item := tool.AddItem("Edit me", persistance.BucketKeyGeneral)
	newText := "New text"
	tool.EditItem(item.Id, &newText)
	editedItem, _ := tool.Storage.FindItem(item.Id)
	if *editedItem.Title != newText {
		t.Fatalf("Item was not updated: %+v", editedItem)
	}
}

func TestCleanupCheckedItems(t *testing.T) {
	tool, _ := setupTestStorage(t)

	checkedItem := tool.AddItem("Done item", persistance.BucketKeyGeneral)
	tool.ToggleItemCheck(checkedItem.Id)

	tool.AddItem("Pending item", persistance.BucketKeyGeneral)
	tool.CleanupCheckedItems()

	items := tool.Storage.ListItem()
	if len(items) != 1 {
		t.Fatalf("Expected only 1 item after cleanup, got %v", len(items))
	}
	if *items[0].Title != "Pending item" {
		t.Fatalf("Unexpected item remaining after cleanup: %+v", items)
	}
}

func TestEmptyToday(t *testing.T) {
	tool, _ := setupTestStorage(t)

	tool.AddItem("General item", persistance.BucketKeyGeneral)
	todayItem := tool.AddItem("Today item", persistance.BucketKeyToday)
	tool.EmptyToday()

	generalItems := tool.ListItem(persistance.BucketKeyGeneral)
	if len(generalItems) != 2 {
		t.Fatalf("Expected 2 items in general bucket, got %v", len(generalItems))
	}
	found := false
	for _, item := range generalItems {
		if *item.Title == *todayItem.Title {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("Item was not moved to the general bucket")
	}
}

func TestCompactIds(t *testing.T) {
	tool, _ := setupTestStorage(t)

	tool.AddItem("First item", persistance.BucketKeyGeneral)
	tool.AddItem("Second item", persistance.BucketKeyGeneral)
	tool.AddItem("Third item", persistance.BucketKeyGeneral)

	tool.ToggleItemCheck(tool.Storage.ItemList[1].Id)
	tool.CleanupCheckedItems()

	tool.CompactIds()

	for i, item := range tool.Storage.ItemList {
		expectedId := int64(i + 1)
		if *item.Id != expectedId {
			t.Fatalf("Expected item ID %v, got %v", expectedId, *item.Id)
		}
	}
}

func TestSaveAndLoad(t *testing.T) {
	tool, _ := setupTestStorage(t)

	tool.AddItem("First item", persistance.BucketKeyGeneral)
	tool.AddItem("Second item", persistance.BucketKeyGeneral)
	tool.AddItem("Third item", persistance.BucketKeyGeneral)

	err := os.MkdirAll(*tool.StorageFile, 0755)
	if err != nil {
		log.Fatalf("Failed to create directories: %v", err)
	}
	tool.Save()

	tool2, _ := setupTestStorage(t)
	entryList := tool2.ListItem(persistance.BucketKeyGeneral)
	if len(entryList) != 3 {
		t.Fatalf("Expected 3 items in general bucket, got %v", len(entryList))
	}
}

func TestAutoBucketList(t *testing.T) {
	tool, _ := setupTestStorage(t)
	gi := tool.AddItem("General item", persistance.BucketKeyGeneral)
	ti := tool.AddItem("Today item", persistance.BucketKeyToday)

	general, today := tool.PrintGeneralAndOrTodayDependingOnItemCurrentBucket(gi.Id)
	if !general || today {
		t.Fatal("Expected general true and today false")
	}

	general, today = tool.PrintGeneralAndOrTodayDependingOnItemCurrentBucket(ti.Id)
	if general || !today {
		t.Fatal("Expected today true and general false")
	}
}
