-- 创建简历上传状态跟踪表
CREATE TABLE IF NOT EXISTS `resume_upload_status` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `upload_id` varchar(36) NOT NULL COMMENT '上传任务ID(UUID)',
  `user_id` int unsigned NOT NULL COMMENT '用户ID',
  `file_path` varchar(500) NOT NULL COMMENT '文件路径',
  `resume_id` bigint unsigned DEFAULT NULL COMMENT '解析完成后的简历ID',
  `status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '状态',
  `progress` int DEFAULT 0 COMMENT '进度百分比',
  `stage` varchar(100) DEFAULT NULL COMMENT '当前阶段描述',
  `error_msg` text COMMENT '错误信息',
  `extract_duration` bigint DEFAULT NULL COMMENT 'PDF提取耗时(ms)',
  `analyze_duration` bigint DEFAULT NULL COMMENT 'AI分析耗时(ms)',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_upload_id` (`upload_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_resume_id` (`resume_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='简历上传状态跟踪表';
