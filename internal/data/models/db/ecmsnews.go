package db

import (
	"context"
	"fmt"
)

type Ecmsnews struct {
	Id       int32 `gorm:"primary_key;column:id"`
	ClassId  int32 `gorm:"column:classid"`
	NewsTime int32 `gorm:"column:newstime"`
}

func (Ecmsnews) TableName() string {
	return "phome_ecms_news"
}

func (m *DB) GetFollowNewsByClassId(ctx context.Context, followClassIds []int32, page int32) (index []string, err error) {
	db := m.Source.WithContext(ctx)
	var news []Ecmsnews
	offset := (page - 1) * 12
	if err = db.Select([]string{"id", "classid"}).Where("classid in ?", followClassIds).Order("newstime desc").Offset(int(offset)).Limit(12).Find(&news).Error; err != nil {
		return nil, err
	}
	for _, n := range news {
		index = append(index, fmt.Sprint(n.ClassId, ":", n.Id))
	}
	return index, nil
}
