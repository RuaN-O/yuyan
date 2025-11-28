-- 创建用户表
CREATE TABLE `user` (
    `user_id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '用户唯一标识（自增）',
    `username` VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名（登录账号）',
    `password_hash` CHAR(64) NOT NULL COMMENT '密码哈希值（SHA-256加密）',
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
    `last_login` TIMESTAMP NULL DEFAULT NULL COMMENT '最后登录时间',
    PRIMARY KEY (`user_id`),
    INDEX `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 创建海洋生物表
CREATE TABLE `marine_creature` (
    `creature_id` INT NOT NULL AUTO_INCREMENT COMMENT '生物唯一标识（自增）',
    `name` VARCHAR(100) NOT NULL UNIQUE COMMENT '生物名称（如"蓝鲸"）',
    `description` TEXT NOT NULL COMMENT '详细信息（习性/分布等）',
    `image_url` VARCHAR(255) NULL DEFAULT NULL COMMENT '图片链接',
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
    PRIMARY KEY (`creature_id`),
    INDEX `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='海洋生物表';

-- 创建用户历史记录表
CREATE TABLE `history` (
    `history_id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '记录唯一标识（自增）',
    `user_id` BIGINT NOT NULL COMMENT '用户id',
    `file_name` VARCHAR(255) NOT NULL COMMENT '上传的视频文件名',
    `creature_id` INT NOT NULL COMMENT '识别结果关联的生物ID',
    `result_text` TEXT NULL DEFAULT NULL COMMENT '模型输出的原始结果',
    `identify_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '识别时间',
    PRIMARY KEY (`history_id`),
    UNIQUE KEY `uk_user_id` (`user_id`),
    UNIQUE KEY `uk_creature_id` (`creature_id`),
    INDEX `idx_identify_time` (`identify_time`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户历史记录表';
