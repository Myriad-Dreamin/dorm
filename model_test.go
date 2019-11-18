package dorm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestModelScope_UpdateFields(t *testing.T) {
	gdb, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer gdb.Close()
	if err := gdb.AutoMigrate(&User{}).Error; err != nil {
		t.Fatal(err)
	}
	fdb, err := FromRaw(gdb.DB(), NewFmtLogger())
	if err != nil {
		t.Fatal(err)
	}

	model, err := fdb.Model(&User{ID:1})
	if err != nil {
		t.Fatal(err)
	}
	aff, err := model.Anchor(&User{ID:1}).Select([]string{"gender"}...).UpdateFields()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(aff)
	aff, err = model.Anchor(&User{ID:1, Gender:1}).Select([]string{"gender"}...).UpdateFields()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(aff)
}



