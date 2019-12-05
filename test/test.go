package main

import (
	"fmt"
	"github.com/Myriad-Dreamin/dorm"
	"time"
)

type fmtLogger struct {
	e []interface{}
}

// NewFmtLogger returns a logger that doesn't do anything.
func NewFmtLogger() dorm.Logger { return &fmtLogger{} }

func (l *fmtLogger) Info(a string, e ...interface{}) {
	fmt.Printf("info: %v ", a)
	fmt.Println(append(l.e, e...)...)
}

func (l *fmtLogger) Debug(a string, e ...interface{}) {
	fmt.Printf("debug: %v ", a)
	fmt.Println(append(l.e, e...)...)
}

func (l *fmtLogger) Error(a string, e ...interface{}) {
	fmt.Printf("error: %v ", a)
	fmt.Println(append(l.e, e...)...)
}

func (l *fmtLogger) With(e ...interface{}) dorm.Logger {
	return &fmtLogger{e: append(l.e, e...)}
}

type Group struct {
	ID          uint      `dorm:"id"`
	CreatedAt   time.Time `dorm:"created_at"`
	UpdatedAt   time.Time `dorm:"updated_at"`
	DeletedAt   time.Time `dorm:"deleted_at"`
	Name        string    `gorm:"type:varchar(128);column:name;default:'anonymous group';unique" json:"name"`
	Description string    `gorm:"column:description;type:text"`
	Owner       User      `gorm:"ForeignKey:OwnerID;AssociationForeignKey:ID"`
	OwnerID     uint      `gorm:"column:owner_id" json:"owner_id"`
	UsersBuffer []User    `gorm:"many2many:group_members;association_foreignkey:ID;foreignkey:ID;preload:false"`
}

type User struct {
	ID                  uint      `dorm:"id"`
	CreatedAt           time.Time `dorm:"created_at"`
	UpdatedAt           time.Time `dorm:"updated_at"`
	DeletedAt           time.Time `dorm:"deleted_at"`
	Password            string    `dorm:"password" gorm:"type:varchar(128);column:password" json:"-"`
	Gender              uint8     `dorm:"gender" gorm:"type:varchar(128);column:gender" json:"gender"`
	LastLogin           time.Time `dorm:"last_login" gorm:"column:last_login;default:CURRENT_TIMESTAMP" json:"last_login"`
	UserName            string    `dorm:"user_name" gorm:"type:varchar(30);column:user_name;not null;unique" json:"user_name"` // todo: regex
	NickName            string    `dorm:"nick_name" gorm:"type:varchar(30);column:nick_name;not null" json:"nick_name"`        // todo: regex
	Email               string    `dorm:"email" gorm:"column:email;unique;default:NULL" json:"email" binding:"email"`          // todo: email
	Motto               string    `dorm:"motto" gorm:"column:motto" json:"motto"`
	SolvedProblemsCount int64     `dorm:"solved_problems" gorm:"column:solved_problems" json:"-"`
	TriedProblemsCount  int64     `dorm:"tried_problems" gorm:"column:tried_problems" json:"-"`
}

// TableName specification
func (Group) TableName() string {
	return "group"
}
func (d Group) GetID() uint {
	return d.ID
}


// TableName specification
func (User) TableName() string {
	return "user"
}

func (d User) GetID() uint {
	return d.ID
}

func main() {}
//func main() {
//	logger := NewFmtLogger()
//	db, err := dorm.Open(dsn, logger)
//	if err != nil {
//		logger.Error("error open", "error", err)
//	}
//
//	r, err := db.ManyToManyRelation(&Group{}, &User{},  dorm.RCommon{
//		UPK:       "group_id",
//		VPK:       "user_id",
//		TableName: "group_members",
//	})
//	fmt.Println(r)
//	if err != nil {
//		logger.Error("error open", "error", err)
//	}
//	var result []uint
//	fmt.Println(r.Anchor(&Group{ID:8}).Find(&result))
//	fmt.Println(result)
//
//	result = nil
//	fmt.Println(r.Anchor(&Group{ID:8}).Limit(5).Find(&result))
//	fmt.Println(result)
//
//	result = nil
//	fmt.Println(r.Anchor(&Group{ID:8}).Offset(5).Limit(5).Find(&result))
//	fmt.Println(result)
//
//	result = nil
//	fmt.Println(r.Anchor(&Group{ID:8}).Limit(5).Offset(1).Find(&result))
//	fmt.Println(result)
//
//
//	result = nil
//	fmt.Println(r.Anchor(&Group{ID:8}).Order("user_id desc").Find(&result))
//	fmt.Println(result)
//
//	result = nil
//	fmt.Println(r.Anchor(&Group{ID:8}).Order("user_id desc").Limit(5).Find(&result))
//	fmt.Println(result)
//
//	result = nil
//	fmt.Println(r.Anchor(&Group{ID:8}).Order("user_id desc").Offset(5).Limit(5).Find(&result))
//	fmt.Println(result)
//
//	result = nil
//	fmt.Println(r.Anchor(&Group{ID:8}).Order("user_id desc").Limit(5).Offset(1).Find(&result))
//	fmt.Println(result)
//
//
//	result = nil
//	fmt.Println(r.Anchor(&Group{ID:8}).Order("user_id desc").Where("user_id between ? and ?").Find(&result, 2, 5))
//	fmt.Println(result)
//
//	result = nil
//	fmt.Println(r.Anchor(&Group{ID:8}).Order("user_id desc").Where("user_id between ? and ?").Limit(5).Find(&result, 2, 5))
//	fmt.Println(result)
//
//	result = nil
//	fmt.Println(r.Anchor(&Group{ID:8}).Order("user_id desc").Where("user_id between ? and ?").Offset(5).Limit(5).Find(&result, 2, 5))
//	fmt.Println(result)
//
//	result = nil
//	fmt.Println(r.Anchor(&Group{ID:8}).Order("user_id desc").Where("user_id between ? and ?").Limit(5).Offset(1).Find(&result, 2, 5))
//	fmt.Println(result)
//
//
//	result = nil
//	fmt.Println(r.Anchor(&Group{ID:8}).Where("user_id between ? and ?").Find(&result, 2, 5))
//	fmt.Println(result)
//
//	result = nil
//	fmt.Println(r.Anchor(&Group{ID:8}).Where("user_id between ? and ?").Limit(5).Find(&result, 2, 5))
//	fmt.Println(result)
//
//	result = nil
//	fmt.Println(r.Anchor(&Group{ID:8}).Where("user_id between ? and ?").Offset(5).Limit(5).Find(&result, 2, 5))
//	fmt.Println(result)
//
//	result = nil
//	fmt.Println(r.Anchor(&Group{ID:8}).Where("user_id between ? and ?").Limit(5).Offset(1).Find(&result, 2, 5))
//	fmt.Println(result)
//
//	defer db.Close()
//
//}
