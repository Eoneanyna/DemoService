package db

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type News struct {
	Id         int32     `gorm:"primary_key;column:id"`
	Title      string    `gorm:"column:title;comment:标题"`
	Content    string    `gorm:"column:content;comment:内容"`
	ViewCount  int64     `gorm:"column:view_count;comment:点击量"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;not null;comment:创建时间"`
}

func (News) TableName() string {
	return "news"
}

func (m *DB) GetNewsById(ctx context.Context, Id int32) (news []News, err error) {
	db := m.Source.WithContext(ctx)

	if err = db.Select("*").Where("id = ?", Id).Find(&news).Error; err != nil {
		return nil, err
	}

	return news, nil
}

func (m *DB) GetNews(ctx context.Context, page int32, pageSize int) (news []News, err error) {
	db := m.Source.WithContext(ctx)
	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * int32(pageSize)
	if err = db.Select("*").Order("create_time desc").Offset(int(offset)).Limit(pageSize).Find(&news).Error; err != nil {
		return nil, err
	}

	return news, nil
}

func (m *DB) CreateNews(ctx context.Context, news News) (n News, err error) {
	db := m.Source.WithContext(ctx)

	news.CreateTime = time.Now()
	if err = db.Create(&news).Error; err != nil {
		log.Errorf("数据%+v创建失败，err = %v", news, err)
		return News{}, err
	}

	return news, nil
}
