package data

import "time"

// User 用户模型，对应数据库中的 user 表
type User struct {
	// user_id BIGINT NOT NULL AUTO_INCREMENT COMMENT '用户唯一标识（自增）'
	ID int64 `gorm:"column:user_id;primaryKey;autoIncrement;comment:用户唯一标识（自增）"`

	// username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名（登录账号）'
	Username string `gorm:"column:username;size:50;not null;uniqueIndex;comment:用户名（登录账号）"`

	// password_hash CHAR(64) NOT NULL COMMENT '密码哈希值（SHA-256加密）'
	PasswordHash string `gorm:"column:password_hash;size:100;not null;comment:密码哈希值（bcrypt加密）"`

	// create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间'
	CreateTime time.Time `gorm:"column:create_time;not null;autoCreateTime;comment:注册时间"`

	// last_login TIMESTAMP NULL DEFAULT NULL COMMENT '最后登录时间'
	LastLogin time.Time `gorm:"column:last_login;autoCreateTime;comment:最后登录时间"`
}

// TableName 指定 GORM 使用的表名
// GORM 默认会将结构体名复数化作为表名 (User -> users)，这里我们指定为单数 'user'
func (User) TableName() string {
	return "user"
}
