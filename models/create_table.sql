DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `user_id` bigint(20) NOT NULL,
    `username` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
    `password` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
    `email` varchar(64) COLLATE utf8mb4_general_ci,
    `gender` tinyint(4) NOT NULL DEFAULT '0',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE
            CURRENT_TIMESTAMP,
        PRIMARY KEY (`id`),
        UNIQUE KEY `idx_username` (`username`) USING BTREE,
        UNIQUE KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

DROP TABLE IF EXISTS `community`;
CREATE TABLE `community` (
     `id` int(11) NOT NULL AUTO_INCREMENT,
     `community_id` int(10) unsigned NOT NULL,
     `community_name` varchar(128) COLLATE utf8mb4_general_ci NOT NULL,
     `introduction` varchar(256) COLLATE utf8mb4_general_ci NOT NULL,
     `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
     `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
     PRIMARY KEY (`id`),
     UNIQUE KEY `idx_community_id` (`community_id`),
     UNIQUE KEY `idx_community_name` (`community_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `community` VALUES ('1', '1', 'Go', 'Golang', '2016-11-01 08:10:10', '2016-11-01 08:10:10');
INSERT INTO `community` VALUES ('2', '2', 'leetcode', '刷题刷题刷题', '2020-01-01 08:00:00', '2020-01-01 08:00:00');
INSERT INTO `community` VALUES ('3', '3', 'CS:GO', 'Rush B。。。', '2018-08-07 08:30:00', '2018-08-07 08:30:00');
INSERT INTO `community` VALUES ('4', '4', 'LOL', '欢迎来到英雄联盟!', '2016-01-01 08:00:00', '2016-01-01 08:00:00');

CREATE TABLE `post` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `post_id` bigint(20) NOT NULL COMMENT '帖子id',
    `title` varchar(128) COLLATE utf8mb4_general_ci NOT NULL COMMENT '标题',
    `content` varchar(8192) COLLATE utf8mb4_general_ci NOT NULL COMMENT '内容',
    `author_id` bigint(20) NOT NULL COMMENT '作者的用户id',
    `community_id` bigint(20) NOT NULL COMMENT '所属社区',
    `like_count` bigint(20) NOT NULL DEFAULT '0' COMMENT '点赞数',
    `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '帖子状态',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_post_id` (`post_id`),
    KEY `idx_author_id` (`author_id`),
    KEY `idx_community_id` (`community_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `comment_id` bigint(20) NOT NULL COMMENT '评论id',
    `post_id` bigint(20) NOT NULL COMMENT '评论所属帖子id',
    `author_id` bigint(20) NOT NULL COMMENT '评论作者用户id',
    `author_name` varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '评论作者名',
    `content` varchar(2048) COLLATE utf8mb4_general_ci NOT NULL COMMENT '评论内容',
    `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '1正常，0删除',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '逻辑删除时间，0表示未删除',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_comment_id` (`comment_id`),
    KEY `idx_post_id` (`post_id`),
    KEY `idx_author_id` (`author_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

DROP TABLE IF EXISTS `comment_relation`;
CREATE TABLE `comment_relation` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `post_id` bigint(20) NOT NULL COMMENT '评论所属帖子id',
    `comment_id` bigint(20) NOT NULL COMMENT '当前评论id',
    `parent_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '父评论id，0表示一级评论',
    `reply_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '被回复的评论id',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '逻辑删除时间，0表示未删除',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_comment_id` (`comment_id`),
    KEY `idx_post_parent` (`post_id`, `parent_id`),
    KEY `idx_reply_id` (`reply_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

DROP TABLE IF EXISTS `post_like`;
CREATE TABLE `post_like` (
    `post_id` bigint(20) NOT NULL COMMENT '帖子id',
    `user_id` bigint(20) NOT NULL COMMENT '用户id',
    `liked` tinyint(4) NOT NULL DEFAULT '1' COMMENT '1点赞，0取消点赞',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`post_id`, `user_id`),
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

DROP TABLE IF EXISTS `like_event_failed`;
CREATE TABLE `like_event_failed` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `event_id` bigint(20) NOT NULL COMMENT '点赞事件id',
    `post_id` bigint(20) NOT NULL COMMENT '帖子id',
    `user_id` bigint(20) NOT NULL COMMENT '用户id',
    `liked` tinyint(4) NOT NULL COMMENT '1点赞，0取消点赞',
    `delta` bigint(20) NOT NULL COMMENT '点赞数变化',
    `retry_count` int(11) NOT NULL DEFAULT '0' COMMENT '重试次数',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_event_id` (`event_id`),
    KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

DROP TABLE IF EXISTS `user_checkin_detail`;
CREATE TABLE `user_checkin_detail` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `checkin_id` bigint(20) NOT NULL COMMENT '签到记录id',
    `user_id` bigint(20) NOT NULL COMMENT '用户id',
    `sign_date` date NOT NULL COMMENT '签到日期',
    `sign_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '签到时间',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_checkin_id` (`checkin_id`),
    UNIQUE KEY `idx_user_sign_date` (`user_id`, `sign_date`),
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

DROP TABLE IF EXISTS `user_checkin_count`;
CREATE TABLE `user_checkin_count` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `user_id` bigint(20) NOT NULL COMMENT '用户id',
    `total_count` int(11) NOT NULL DEFAULT '0' COMMENT '累计签到次数',
    `continuous_count` int(11) NOT NULL DEFAULT '0' COMMENT '连续签到次数',
    `last_sign_date` date NOT NULL COMMENT '上一次签到日期',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
