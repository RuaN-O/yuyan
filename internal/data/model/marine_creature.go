package data

import "time"

// History 历史记录模型，对应数据库中的 history 表
type History struct {
	// history_id BIGINT PRIMARY KEY COMMENT '记录唯一标识（自增）'
	ID int64 `gorm:"column:history_id;primaryKey;autoIncrement;comment:记录唯一标识（自增）"`

	// user_id BIGINT UNIQUE COMMENT '用户id'
	UserID int64 `gorm:"column:user_id;uniqueIndex;comment:用户id"`

	// file_name VARCHAR(255) NOT NULL COMMENT '上传的视频文件名'
	FileName string `gorm:"column:file_name;size:255;not null;comment:上传的视频文件名"`

	// creature_id INT UNIQUE COMMENT '识别结果关联的生物ID'
	CreatureID int `gorm:"column:creature_id;uniqueIndex;comment:识别结果关联的生物ID"`

	// result_text TEXT NULLABLE COMMENT '模型输出的原始结果'
	ResultText string `gorm:"column:result_text;type:text;comment:模型输出的原始结果"`

	// identify_time TIMESTAMP DEFAULT NOW COMMENT '识别时间'
	IdentifyTime time.Time `gorm:"column:identify_time;autoCreateTime;comment:识别时间"`
}

// TableName 指定 GORM 使用的表名
func (History) TableName() string {
	return "history"
}

// MarineCreature 海洋生物模型，对应数据库中的 marine_creature 表
type MarineCreature struct {
	// creature_id INT PRIMARY KEY COMMENT '生物唯一标识（自增）'
	ID int `gorm:"column:creature_id;primaryKey;autoIncrement;comment:生物唯一标识（自增）"`

	// name VARCHAR(100) UNIQUE COMMENT '生物名称（如"蓝鲸"）'
	Name string `gorm:"column:name;size:100;uniqueIndex;comment:生物名称（如蓝鲸）"`

	// description TEXT NOT NULL COMMENT '详细信息（习性/分布等）'
	Description string `gorm:"column:description;type:text;not null;comment:详细信息（习性/分布等）"`

	// image_url VARCHAR(255) NULLABLE COMMENT '图片链接'
	ImageURL string `gorm:"column:image_url;size:255;comment:图片链接"`

	// create_time TIMESTAMP DEFAULT NOW COMMENT '记录创建时间'
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime;comment:记录创建时间"`
}

// TableName 指定 GORM 使用的表名
func (MarineCreature) TableName() string {
	return "marine_creature"
}
