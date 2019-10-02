package dorm

import (
	//"fmt"
	"testing"
	"time"
)

func BenchmarkModelScope_BuildFind(b *testing.B) {

	logger := NewFmtLogger()
	db, err := Open(dsn, logger)
	if err != nil {
		b.Error(err)
		return
	}

	m, err := db.Model(&User{ID:1})
	if err != nil {
		b.Error(err)
		return
	}

	b.ResetTimer()

	for i := 0; i < b.N; i ++ {
		_ = m.Scope().BuildFind()
	}
}

func BenchmarkModelScope_Find(b *testing.B) {

	logger := NewNopLogger()
	db, err := Open(dsn, logger)
	if err != nil {
		b.Error(err)
		return
	}

	m, err := db.Model(&User{ID:1})
	if err != nil {
		b.Error(err)
		return
	}
	t := m.Scope().BuildFind()
	b.ResetTimer()

	for i := 0; i < b.N; i ++ {
		_ = t.Find(&User{})
	}
}

func TestModelScope_Find(t *testing.T) {

	logger := NewFmtLogger()
	db, err := Open(dsn, logger)
	if err != nil {
		t.Error(err)
		return
	}

	m, err := db.Model(&User{ID:1})
	if err != nil {
		t.Error(err)
		return
	}

	g := User{ID:1}
	err = m.Scope().Find(&g)
	if err != nil {
		t.Error(err)
		return
	}
	//fmt.Println(g)


	g = User{ID:264}
	_, err = m.Anchor(&g).Delete()
	if err != nil {
		t.Error(err)
		return
	}
	//fmt.Println(a)

	g = User{ID:263}
	_, err = m.Scope().ID(g.ID).Delete()
	if err != nil {
		t.Error(err)
		return
	}
	//fmt.Println(a)




	var gg []User
	m2, err := db.Model(&User{})
	if err != nil {
		t.Error(err)
		return
	}


	g = User{ID:263}
	_, err = m2.Scope().Where("id between ? and ?").Delete(263, 266)
	if err != nil {
		t.Error(err)
		return
	}
	//fmt.Println(a)


	err = m2.Scope().Find(&gg)
	if err != nil {
		t.Error(err)
		return
	}
	//fmt.Println(len(gg))
	xxx := time.Now()
	ggg := &User{ID:291, CreatedAt:&xxx}
	//a, err = m2.Anchor(ggg).Insert()
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	////fmt.Println("???" , a, ggg)

	_, err = m2.Anchor(ggg).Select("created_at").UpdateFields()
	if err != nil {
		t.Error(err)
		return
	}
	//fmt.Println(a)

	err = m2.Scope().ID(291).Select("id", "created_at").Find(ggg)
	if err != nil {
		t.Error(err)
		return
	}
	//fmt.Println(ggg)

	gg = nil
	err = m2.Scope().Find(&gg)
	if err != nil {
		t.Error(err)
		return
	}
	//fmt.Println(len(gg))
}
