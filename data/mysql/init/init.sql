-- 授权 root 用户可以远程链接
ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'AlnitakRoot@123';
flush privileges;

use alnitak;

-- 创建API表
CREATE TABLE IF NOT EXISTS `api`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `method` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '请求方式',
  `path` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '访问路径',
  `category` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '所属类别',
  `desc` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '说明',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_api_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- 插入初始API数据
INSERT INTO `api` VALUES (1, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/api/addApi', 'API管理', '新增API（后台管理）');
INSERT INTO `api` VALUES (2, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'DELETE', '/api/v1/api/deleteApi/:id', 'API管理', '删除API（后台管理）');
INSERT INTO `api` VALUES (3, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'PUT', '/api/v1/api/editApi', 'API管理', '编辑API（后台管理）');
INSERT INTO `api` VALUES (4, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'PUT', '/api/v1/api/editRoleApi', 'API管理', '编辑角色API（后台管理）');
INSERT INTO `api` VALUES (5, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/api/getAllApiList', 'API管理', '获取全部API列表（后台管理）');
INSERT INTO `api` VALUES (6, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/api/getApiList', 'API管理', '获取API列表（后台管理）');
INSERT INTO `api` VALUES (7, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/api/getRoleApi', 'API管理', '获取角色API（后台管理）');
INSERT INTO `api` VALUES (8, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/archive/article/cancelCollect', '点赞收藏', '文章取消收藏');
INSERT INTO `api` VALUES (9, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/archive/article/cancelLike', '点赞收藏', '文章取消赞');
INSERT INTO `api` VALUES (10, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/archive/article/collect', '点赞收藏', '文章收藏');
INSERT INTO `api` VALUES (11, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/archive/article/hasCollect', '点赞收藏', '文章是否收藏');
INSERT INTO `api` VALUES (12, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/archive/article/hasLike', '点赞收藏', '文章是否点赞');
INSERT INTO `api` VALUES (13, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/archive/article/like', '点赞收藏', '文章点赞');
INSERT INTO `api` VALUES (14, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/archive/video/cancelLike', '点赞收藏', '视频取消赞');
INSERT INTO `api` VALUES (15, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/archive/video/collect', '点赞收藏', '视频收藏');
INSERT INTO `api` VALUES (16, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/archive/video/getCollectInfo', '点赞收藏', '获取视频已收藏的文件夹');
INSERT INTO `api` VALUES (17, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/archive/video/hasCollect', '点赞收藏', '视频是否收藏');
INSERT INTO `api` VALUES (18, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/archive/video/hasLike', '点赞收藏', '视频是否点赞');
INSERT INTO `api` VALUES (19, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/archive/video/like', '点赞收藏', '视频点赞');
INSERT INTO `api` VALUES (20, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'DELETE', '/api/v1/article/deleteArticle/:id', '文章', '删除文章');
INSERT INTO `api` VALUES (21, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'DELETE', '/api/v1/article/deleteArticleManage/:id', '文章', '删除文章（后台管理）');
INSERT INTO `api` VALUES (22, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'PUT', '/api/v1/article/editArticleInfo', '文章', '编辑文章信息');
INSERT INTO `api` VALUES (23, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/article/getAllArticleList', '文章', '获取所有的文章列表');
INSERT INTO `api` VALUES (24, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/article/getArticleListManage', '文章', '获取文章列表（后台管理）');
INSERT INTO `api` VALUES (25, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/article/getArticleStatus', '文章', '获取文章状态信息');
INSERT INTO `api` VALUES (26, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/article/getReviewArticleList', '文章', '获取文章审核列表（后台管理）');
INSERT INTO `api` VALUES (27, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/article/getUploadArticle', '文章', '获取上传的文章');
INSERT INTO `api` VALUES (28, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/article/uploadArticleInfo', '文章', '上传文章信息');
INSERT INTO `api` VALUES (29, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/auth/logout', 'Auth', '退出登录');
INSERT INTO `api` VALUES (30, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/carousel/addCarousel', '轮播图', '新增轮播图（后台管理）');
INSERT INTO `api` VALUES (31, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'DELETE', '/api/v1/carousel/deleteCarousel/:id', '轮播图', '删除轮播图（后台管理）');
INSERT INTO `api` VALUES (32, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'PUT', '/api/v1/carousel/editCarousel', '轮播图', '编辑轮播图（后台管理）');
INSERT INTO `api` VALUES (33, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/carousel/getCarouselList', '轮播图', '获取轮播图列表（后台管理）');
INSERT INTO `api` VALUES (34, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/collection/addCollection', '收藏夹', '添加收藏夹');
INSERT INTO `api` VALUES (35, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'DELETE', '/api/v1/collection/deleteCollection/:id', '收藏夹', '删除收藏夹');
INSERT INTO `api` VALUES (36, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'PUT', '/api/v1/collection/editCollection', '收藏夹', '编辑收藏夹');
INSERT INTO `api` VALUES (37, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/collection/getCollectionInfo', '收藏夹', '获取收藏夹信息');
INSERT INTO `api` VALUES (38, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/collection/getCollectionList', '收藏夹', '获取收藏夹列表');
INSERT INTO `api` VALUES (39, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/collection/getVideoList', '收藏夹', '获取收藏夹的视频列表');
INSERT INTO `api` VALUES (40, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/comment/article/addComment', '评论回复', '发表文章评论或回复');
INSERT INTO `api` VALUES (41, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/comment/article/getCommentList', '评论回复', '获取文章评论列表');
INSERT INTO `api` VALUES (42, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'DELETE', '/api/v1/comment/varticle/deleteComment/:id', '评论回复', '删除文章评论或回复');
INSERT INTO `api` VALUES (43, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/comment/video/addComment', '评论回复', '发表视频评论或回复');
INSERT INTO `api` VALUES (44, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'DELETE', '/api/v1/comment/video/deleteComment/:id', '评论回复', '删除视频评论或回复');
INSERT INTO `api` VALUES (45, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/comment/video/getCommentList', '评论回复', '获取视频评论列表');
INSERT INTO `api` VALUES (46, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/danmaku/sendDanmaku', '弹幕', '发送弹幕');
INSERT INTO `api` VALUES (47, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/history/video/addHistory', '历史记录', '保存视频历史记录');
INSERT INTO `api` VALUES (48, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/history/video/getHistory', '历史记录', '获取视频历史记录');
INSERT INTO `api` VALUES (49, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/history/video/getProgress', '历史记录', '获取视频播放进度');
INSERT INTO `api` VALUES (50, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/menu/addMenu', '菜单管理', '添加菜单（后台管理）');
INSERT INTO `api` VALUES (51, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'DELETE', '/api/v1/menu/deleteMenu/:id', '菜单管理', '删除菜单（后台管理）');
INSERT INTO `api` VALUES (52, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'PUT', '/api/v1/menu/editMenu', '菜单管理', '编辑菜单（后台管理）');
INSERT INTO `api` VALUES (53, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'PUT', '/api/v1/menu/editRoleMenu', '菜单管理', '编辑角色菜单（后台管理）');
INSERT INTO `api` VALUES (54, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/menu/getMenuTree', '菜单管理', '获取菜单树（后台管理）');
INSERT INTO `api` VALUES (55, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/menu/getRoleMenu', '菜单管理', '获取角色菜单（后台管理）');
INSERT INTO `api` VALUES (56, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/menu/getUserMenu', '菜单管理', '获取用户菜单树（后台管理）');
INSERT INTO `api` VALUES (57, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/message/addAnnounce', '消息', '添加公告（后台管理）');
INSERT INTO `api` VALUES (58, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'DELETE', '/api/v1/message/deleteAnnounce/:id', '消息', '删除公告（后台管理）');
INSERT INTO `api` VALUES (59, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/message/getAtMsg', '消息', '获取AT通知');
INSERT INTO `api` VALUES (60, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/message/getLikeMsg', '消息', '获取点赞通知');
INSERT INTO `api` VALUES (61, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/message/getReplyMsg', '消息', '获取回复通知');
INSERT INTO `api` VALUES (62, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/message/getWhisperDetails', '消息', '获取私信详情');
INSERT INTO `api` VALUES (63, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/message/getWhisperList', '消息', '获取私信列表');
INSERT INTO `api` VALUES (64, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/message/readWhisper', '消息', '已读私信');
INSERT INTO `api` VALUES (65, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/message/sendWhisper', '消息', '发送私信');
INSERT INTO `api` VALUES (66, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/partition/addPartition', '分区', '添加分区（后台管理）');
INSERT INTO `api` VALUES (67, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'DELETE', '/api/v1/partition/deletePartition/:id', '分区', '删除分区（后台管理）');
INSERT INTO `api` VALUES (68, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/relation/follow', '关注', '关注用户');
INSERT INTO `api` VALUES (69, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/relation/getUserRelation', '关注', '获取用户关系');
INSERT INTO `api` VALUES (70, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/relation/unfollow', '关注', '取关用户');
INSERT INTO `api` VALUES (71, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'DELETE', '/api/v1/resource/deleteResource/:id', '资源', '删除视频资源');
INSERT INTO `api` VALUES (72, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'PUT', '/api/v1/resource/modifyTitle', '资源', '修改资源标题');
INSERT INTO `api` VALUES (73, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/review/getArticleReviewRecord', '审核', '获取文章审核记录');
INSERT INTO `api` VALUES (74, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/review/getVideoReviewRecord', '审核', '获取视频审核记录');
INSERT INTO `api` VALUES (75, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/review/reviewArticleApproved', '审核', '文章审核通过（后台管理）');
INSERT INTO `api` VALUES (76, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/review/reviewArticleFailed', '审核', '文章审核不通过（后台管理）');
INSERT INTO `api` VALUES (77, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/review/reviewVideoApproved', '审核', '视频审核通过（后台管理）');
INSERT INTO `api` VALUES (78, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/review/reviewVideoFailed', '审核', '视频审核不通过（后台管理）');
INSERT INTO `api` VALUES (79, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/role/addRole', '角色', '添加角色（后台管理）');
INSERT INTO `api` VALUES (80, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'DELETE', '/api/v1/role/deleteRole/:id', '角色', '删除角色（后台管理）');
INSERT INTO `api` VALUES (81, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'PUT', '/api/v1/role/editRole', '角色', '编辑角色（后台管理）');
INSERT INTO `api` VALUES (82, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'PUT', '/api/v1/role/editRoleHome', '角色', '编辑角色首页（后台管理）');
INSERT INTO `api` VALUES (83, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/role/getAllRoleList', '角色', '获取全部角色（后台管理）');
INSERT INTO `api` VALUES (84, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/role/getRoleInfo', '角色', '获取个人角色信息（后台管理）');
INSERT INTO `api` VALUES (85, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/role/getRoleList', '角色', '获取角色列表（后台管理）');
INSERT INTO `api` VALUES (86, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/upload/image', '上传', '上传图片');
INSERT INTO `api` VALUES (87, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/upload/video', '上传', '上传视频');
INSERT INTO `api` VALUES (88, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/upload/video/:vid', '上传', '上传视频分P');
INSERT INTO `api` VALUES (89, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'DELETE', '/api/v1/user/deleteUser/:id', '用户', '删除用户');
INSERT INTO `api` VALUES (90, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'PUT', '/api/v1/user/editUserInfo', '用户', '编辑用户信息');
INSERT INTO `api` VALUES (91, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'PUT', '/api/v1/user/editUserInfoManage', '用户', '编辑用户信息（后台管理）');
INSERT INTO `api` VALUES (92, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'PUT', '/api/v1/user/editUserRole', '用户', '编辑用户角色');
INSERT INTO `api` VALUES (93, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/user/getUserInfo', '用户', '获取用户信息');
INSERT INTO `api` VALUES (94, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/user/getUserListManage', '用户', '获取用户列表（后台管理）');
INSERT INTO `api` VALUES (95, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'DELETE', '/api/v1/video/deleteVideo/:id', '视频', '删除视频');
INSERT INTO `api` VALUES (96, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'DELETE', '/api/v1/video/deleteVideoManage/:id', '视频', '删除视频（后台管理）');
INSERT INTO `api` VALUES (97, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'PUT', '/api/v1/video/editVideoInfo', '视频', '编辑视频信息');
INSERT INTO `api` VALUES (98, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/video/getAllVideoList', '视频', '获取所有的视频列表');
INSERT INTO `api` VALUES (99, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/video/getReviewList', '视频', '获取审核列表（后台管理）');
INSERT INTO `api` VALUES (100, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/video/getReviewResourceList', '视频', '获取审核资源列表（后台管理）');
INSERT INTO `api` VALUES (101, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/video/getUploadVideo', '视频', '获取上传的视频');
INSERT INTO `api` VALUES (102, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/video/getVideoListManage', '视频', '获取视频列表（后台管理）');
INSERT INTO `api` VALUES (103, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'GET', '/api/v1/video/getVideoStatus', '视频', '获取上传视频状态信息');
INSERT INTO `api` VALUES (104, '2024-05-18 11:22:33.000', '2024-05-18 11:22:33.000', NULL, 'POST', '/api/v1/video/uploadVideoInfo', '视频', '上传视频信息');
INSERT INTO `api` VALUES (105, '2024-05-20 23:07:37.810', '2024-05-20 23:07:37.810', NULL, 'GET', '/api/v1/video/getResourceQualityManage', '视频', '获取视频资源支持的分辨率信息（后台管理）');
INSERT INTO `api` VALUES (106, '2024-05-20 23:19:27.882', '2024-05-20 23:19:27.882', NULL, 'GET', '/api/v1/video/getVideoFileManage', '视频', '获取视频文件URL（后台管理）');
INSERT INTO `api` VALUES (107, '2024-06-08 16:26:36.000', '2024-06-08 16:26:36.000', NULL, 'GET', '/api/v1/config/getEmailConfig', '配置', '获取邮箱配置（后台管理）');
INSERT INTO `api` VALUES (108, '2024-06-08 16:26:36.000', '2024-06-08 16:26:36.000', NULL, 'POST', '/api/v1/config/setEmailConfig', '配置', '编辑邮箱配置（后台管理）');
INSERT INTO `api` VALUES (109, '2024-06-08 16:26:36.000', '2024-06-08 16:26:36.000', NULL, 'GET', '/api/v1/config/getStorageConfig', '配置', '获取存储配置（后台管理）');
INSERT INTO `api` VALUES (110, '2024-06-08 16:26:36.000', '2024-06-08 16:26:36.000', NULL, 'POST', '/api/v1/config/setStorageConfig', '配置', '编辑存储配置（后台管理）');
INSERT INTO `api` VALUES (111, '2024-06-08 16:26:36.000', '2024-06-08 16:26:36.000', NULL, 'GET', '/api/v1/config/getOtherConfig', '配置', '获取其他配置（后台管理）');
INSERT INTO `api` VALUES (112, '2024-06-08 16:26:36.000', '2024-06-08 16:26:36.000', NULL, 'POST', '/api/v1/config/setOtherConfig', '配置', '编辑其他配置（后台管理）');


-- 创建CasbinRule表
CREATE TABLE IF NOT EXISTS `casbin_rule`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v0` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v1` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v2` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v3` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v4` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v5` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_casbin_rule`(`ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- 插入初始CasbinRule数据
INSERT INTO `casbin_rule` VALUES (1, 'p', '001', '/api/v1/archive/article/cancelCollect', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (2, 'p', '001', '/api/v1/archive/article/cancelLike', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (3, 'p', '001', '/api/v1/archive/article/collect', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (4, 'p', '001', '/api/v1/archive/article/hasCollect', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (5, 'p', '001', '/api/v1/archive/article/hasLike', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (6, 'p', '001', '/api/v1/archive/article/like', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (7, 'p', '001', '/api/v1/archive/video/cancelLike', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (8, 'p', '001', '/api/v1/archive/video/collect', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (9, 'p', '001', '/api/v1/archive/video/getCollectInfo', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (10, 'p', '001', '/api/v1/archive/video/hasCollect', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (11, 'p', '001', '/api/v1/archive/video/hasLike', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (12, 'p', '001', '/api/v1/archive/video/like', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (13, 'p', '001', '/api/v1/article/deleteArticle/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (14, 'p', '001', '/api/v1/article/editArticleInfo', 'PUT', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (15, 'p', '001', '/api/v1/article/getAllArticleList', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (16, 'p', '001', '/api/v1/article/getArticleStatus', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (17, 'p', '001', '/api/v1/article/getUploadArticle', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (18, 'p', '001', '/api/v1/article/uploadArticleInfo', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (19, 'p', '001', '/api/v1/auth/logout', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (20, 'p', '001', '/api/v1/collection/addCollection', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (21, 'p', '001', '/api/v1/collection/deleteCollection/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (22, 'p', '001', '/api/v1/collection/editCollection', 'PUT', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (23, 'p', '001', '/api/v1/collection/getCollectionInfo', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (24, 'p', '001', '/api/v1/collection/getCollectionList', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (25, 'p', '001', '/api/v1/collection/getVideoList', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (26, 'p', '001', '/api/v1/comment/article/addComment', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (27, 'p', '001', '/api/v1/comment/article/getCommentList', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (28, 'p', '001', '/api/v1/comment/varticle/deleteComment/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (29, 'p', '001', '/api/v1/comment/video/addComment', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (30, 'p', '001', '/api/v1/comment/video/deleteComment/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (31, 'p', '001', '/api/v1/comment/video/getCommentList', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (32, 'p', '001', '/api/v1/danmaku/sendDanmaku', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (33, 'p', '001', '/api/v1/history/video/addHistory', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (34, 'p', '001', '/api/v1/history/video/getHistory', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (35, 'p', '001', '/api/v1/history/video/getProgress', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (36, 'p', '001', '/api/v1/message/getAtMsg', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (37, 'p', '001', '/api/v1/message/getLikeMsg', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (38, 'p', '001', '/api/v1/message/getReplyMsg', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (39, 'p', '001', '/api/v1/message/getWhisperDetails', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (40, 'p', '001', '/api/v1/message/getWhisperList', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (41, 'p', '001', '/api/v1/message/readWhisper', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (42, 'p', '001', '/api/v1/message/sendWhisper', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (43, 'p', '001', '/api/v1/relation/follow', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (44, 'p', '001', '/api/v1/relation/getUserRelation', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (45, 'p', '001', '/api/v1/relation/unfollow', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (46, 'p', '001', '/api/v1/resource/deleteResource/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (47, 'p', '001', '/api/v1/resource/modifyTitle', 'PUT', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (48, 'p', '001', '/api/v1/review/getArticleReviewRecord', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (49, 'p', '001', '/api/v1/review/getVideoReviewRecord', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (50, 'p', '001', '/api/v1/upload/image', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (51, 'p', '001', '/api/v1/upload/video', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (52, 'p', '001', '/api/v1/upload/video/:vid', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (53, 'p', '001', '/api/v1/user/deleteUser/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (54, 'p', '001', '/api/v1/user/editUserInfo', 'PUT', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (55, 'p', '001', '/api/v1/user/editUserRole', 'PUT', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (56, 'p', '001', '/api/v1/user/getUserInfo', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (57, 'p', '001', '/api/v1/video/deleteVideo/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (58, 'p', '001', '/api/v1/video/editVideoInfo', 'PUT', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (59, 'p', '001', '/api/v1/video/getAllVideoList', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (60, 'p', '001', '/api/v1/video/getUploadVideo', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (61, 'p', '001', '/api/v1/video/getVideoStatus', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (62, 'p', '001', '/api/v1/video/uploadVideoInfo', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (63, 'p', '002', '/api/v1/api/addApi', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (64, 'p', '002', '/api/v1/api/deleteApi/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (65, 'p', '002', '/api/v1/api/editApi', 'PUT', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (66, 'p', '002', '/api/v1/api/editRoleApi', 'PUT', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (67, 'p', '002', '/api/v1/api/getAllApiList', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (68, 'p', '002', '/api/v1/api/getApiList', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (69, 'p', '002', '/api/v1/api/getRoleApi', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (70, 'p', '002', '/api/v1/archive/article/cancelCollect', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (71, 'p', '002', '/api/v1/archive/article/cancelLike', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (72, 'p', '002', '/api/v1/archive/article/collect', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (73, 'p', '002', '/api/v1/archive/article/hasCollect', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (74, 'p', '002', '/api/v1/archive/article/hasLike', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (75, 'p', '002', '/api/v1/archive/article/like', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (76, 'p', '002', '/api/v1/archive/video/cancelLike', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (77, 'p', '002', '/api/v1/archive/video/collect', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (78, 'p', '002', '/api/v1/archive/video/getCollectInfo', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (79, 'p', '002', '/api/v1/archive/video/hasCollect', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (80, 'p', '002', '/api/v1/archive/video/hasLike', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (81, 'p', '002', '/api/v1/archive/video/like', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (82, 'p', '002', '/api/v1/article/deleteArticle/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (83, 'p', '002', '/api/v1/article/deleteArticleManage/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (84, 'p', '002', '/api/v1/article/editArticleInfo', 'PUT', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (85, 'p', '002', '/api/v1/article/getAllArticleList', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (86, 'p', '002', '/api/v1/article/getArticleListManage', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (87, 'p', '002', '/api/v1/article/getArticleStatus', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (88, 'p', '002', '/api/v1/article/getReviewArticleList', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (89, 'p', '002', '/api/v1/article/getUploadArticle', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (90, 'p', '002', '/api/v1/article/uploadArticleInfo', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (91, 'p', '002', '/api/v1/auth/logout', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (92, 'p', '002', '/api/v1/carousel/addCarousel', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (93, 'p', '002', '/api/v1/carousel/deleteCarousel/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (94, 'p', '002', '/api/v1/carousel/editCarousel', 'PUT', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (95, 'p', '002', '/api/v1/carousel/getCarouselList', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (96, 'p', '002', '/api/v1/collection/addCollection', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (97, 'p', '002', '/api/v1/collection/deleteCollection/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (98, 'p', '002', '/api/v1/collection/editCollection', 'PUT', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (99, 'p', '002', '/api/v1/collection/getCollectionInfo', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (100, 'p', '002', '/api/v1/collection/getCollectionList', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (101, 'p', '002', '/api/v1/collection/getVideoList', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (102, 'p', '002', '/api/v1/comment/article/addComment', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (103, 'p', '002', '/api/v1/comment/article/getCommentList', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (104, 'p', '002', '/api/v1/comment/varticle/deleteComment/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (105, 'p', '002', '/api/v1/comment/video/addComment', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (106, 'p', '002', '/api/v1/comment/video/deleteComment/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (107, 'p', '002', '/api/v1/comment/video/getCommentList', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (108, 'p', '002', '/api/v1/danmaku/sendDanmaku', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (109, 'p', '002', '/api/v1/history/video/addHistory', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (110, 'p', '002', '/api/v1/history/video/getHistory', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (111, 'p', '002', '/api/v1/history/video/getProgress', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (112, 'p', '002', '/api/v1/menu/addMenu', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (113, 'p', '002', '/api/v1/menu/deleteMenu/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (114, 'p', '002', '/api/v1/menu/editMenu', 'PUT', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (115, 'p', '002', '/api/v1/menu/editRoleMenu', 'PUT', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (116, 'p', '002', '/api/v1/menu/getMenuTree', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (117, 'p', '002', '/api/v1/menu/getRoleMenu', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (118, 'p', '002', '/api/v1/menu/getUserMenu', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (119, 'p', '002', '/api/v1/message/addAnnounce', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (120, 'p', '002', '/api/v1/message/deleteAnnounce/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (121, 'p', '002', '/api/v1/message/getAtMsg', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (122, 'p', '002', '/api/v1/message/getLikeMsg', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (123, 'p', '002', '/api/v1/message/getReplyMsg', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (124, 'p', '002', '/api/v1/message/getWhisperDetails', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (125, 'p', '002', '/api/v1/message/getWhisperList', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (126, 'p', '002', '/api/v1/message/readWhisper', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (127, 'p', '002', '/api/v1/message/sendWhisper', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (128, 'p', '002', '/api/v1/partition/addPartition', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (129, 'p', '002', '/api/v1/partition/deletePartition/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (130, 'p', '002', '/api/v1/relation/follow', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (131, 'p', '002', '/api/v1/relation/getUserRelation', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (132, 'p', '002', '/api/v1/relation/unfollow', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (133, 'p', '002', '/api/v1/resource/deleteResource/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (134, 'p', '002', '/api/v1/resource/modifyTitle', 'PUT', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (135, 'p', '002', '/api/v1/review/getArticleReviewRecord', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (136, 'p', '002', '/api/v1/review/getVideoReviewRecord', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (137, 'p', '002', '/api/v1/review/reviewArticleApproved', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (138, 'p', '002', '/api/v1/review/reviewArticleFailed', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (139, 'p', '002', '/api/v1/review/reviewVideoApproved', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (140, 'p', '002', '/api/v1/review/reviewVideoFailed', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (141, 'p', '002', '/api/v1/role/addRole', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (142, 'p', '002', '/api/v1/role/deleteRole/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (143, 'p', '002', '/api/v1/role/editRole', 'PUT', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (144, 'p', '002', '/api/v1/role/editRoleHome', 'PUT', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (145, 'p', '002', '/api/v1/role/getAllRoleList', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (146, 'p', '002', '/api/v1/role/getRoleInfo', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (147, 'p', '002', '/api/v1/role/getRoleList', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (148, 'p', '002', '/api/v1/upload/image', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (149, 'p', '002', '/api/v1/upload/video', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (150, 'p', '002', '/api/v1/upload/video/:vid', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (151, 'p', '002', '/api/v1/user/deleteUser/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (152, 'p', '002', '/api/v1/user/editUserInfo', 'PUT', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (153, 'p', '002', '/api/v1/user/editUserInfoManage', 'PUT', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (154, 'p', '002', '/api/v1/user/editUserRole', 'PUT', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (155, 'p', '002', '/api/v1/user/getUserInfo', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (156, 'p', '002', '/api/v1/user/getUserListManage', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (157, 'p', '002', '/api/v1/video/deleteVideo/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (158, 'p', '002', '/api/v1/video/deleteVideoManage/:id', 'DELETE', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (159, 'p', '002', '/api/v1/video/editVideoInfo', 'PUT', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (160, 'p', '002', '/api/v1/video/getAllVideoList', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (161, 'p', '002', '/api/v1/video/getResourceQualityManage', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (162, 'p', '002', '/api/v1/video/getReviewList', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (163, 'p', '002', '/api/v1/video/getReviewResourceList', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (164, 'p', '002', '/api/v1/video/getUploadVideo', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (165, 'p', '002', '/api/v1/video/getVideoFileManage', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (166, 'p', '002', '/api/v1/video/getVideoListManage', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (167, 'p', '002', '/api/v1/video/getVideoStatus', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (168, 'p', '002', '/api/v1/video/uploadVideoInfo', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (169, 'p', '002', '/api/v1/config/getEmailConfig', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (170, 'p', '002', '/api/v1/config/setEmailConfig', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (171, 'p', '002', '/api/v1/config/getStorageConfig', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (172, 'p', '002', '/api/v1/config/setStorageConfig', 'POST', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (173, 'p', '002', '/api/v1/config/getOtherConfig', 'GET', NULL, NULL, NULL);
INSERT INTO `casbin_rule` VALUES (174, 'p', '002', '/api/v1/config/setOtherConfig', 'POST', NULL, NULL, NULL);

-- 创建菜单表
CREATE TABLE IF NOT EXISTS  `menu`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '菜单名称',
  `path` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '菜单访问路径',
  `component` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '前端组件路径',
  `desc` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '介绍',
  `sort` int NULL DEFAULT 1 COMMENT '菜单排序',
  `parent_id` bigint UNSIGNED NULL DEFAULT 0 COMMENT '父菜单ID',
  `title` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '菜单标题',
  `icon` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '菜单图标',
  `hidden` tinyint(1) NULL DEFAULT 0 COMMENT '在菜单中隐藏',
  `keep_alive` tinyint(1) NULL DEFAULT 0 COMMENT '缓存页面',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_menu_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- 插入初始菜单数据
INSERT INTO `menu` VALUES (1, '2024-05-05 13:45:25.291', '2024-05-05 13:45:25.291', NULL, 'review', 'review', NULL, NULL, 1, 0, '内容审核', 'LayersOutline', 0, 0);
INSERT INTO `menu` VALUES (2, '2024-04-13 21:28:51.838', '2024-05-05 13:47:42.897', NULL, 'ReviewVideo', 'review/video', 'views/review/video/index.vue', NULL, 1, 1, '视频审核', 'FileTrayOutline', 0, 0);
INSERT INTO `menu` VALUES (3, '2024-05-05 13:44:36.819', '2024-05-05 13:48:00.048', NULL, 'ReviewArticle', 'review/article', 'views/review/article/index.vue', NULL, 1, 1, '专栏审核', 'AlbumsOutline', 0, 0);
INSERT INTO `menu` VALUES (4, '2024-04-13 21:15:24.522', '2024-04-13 21:15:51.101', NULL, 'Content', 'content', NULL, NULL, 1, 0, '内容管理', 'ReaderOutline', 0, 0);
INSERT INTO `menu` VALUES (5, '2024-04-13 21:17:11.536', '2024-04-13 21:17:29.466', NULL, 'ContentVideo', 'content/video', 'views/content/video/index.vue', NULL, 1, 4, '视频管理', 'PlayCircleOutline', 0, 0);
INSERT INTO `menu` VALUES (6, '2024-05-05 13:43:27.079', '2024-05-05 13:43:27.079', NULL, 'ContentArticle', 'content/article', 'views/content/article/index.vue', NULL, 1, 4, '专栏管理', 'DocumentTextOutline', 0, 0);
INSERT INTO `menu` VALUES (7, '2024-04-13 21:30:14.985', '2024-04-13 21:30:26.532', NULL, 'ContentCarousel', 'content/carousel', 'views/content/carousel/index.vue', NULL, 1, 4, '轮播图管理', 'ImagesOutline', 0, 0);
INSERT INTO `menu` VALUES (8, '2024-05-03 13:02:41.455', '2024-05-03 13:03:07.324', NULL, 'ContentPartition', 'content/partition', 'views/content/partition/index.vue', NULL, 1, 4, '分区管理', 'BookmarkOutline', 0, 0);
INSERT INTO `menu` VALUES (9, '2023-11-24 22:04:14.000', '2024-04-13 15:21:54.739', NULL, 'System', 'system', NULL, NULL, 1, 0, '系统管理', 'TerminalOutline', 0, 0);
INSERT INTO `menu` VALUES (10, '2023-11-24 22:05:27.000', '2023-11-27 23:47:00.933', NULL, 'SystemApi', 'system/api', 'views/system/api/index.vue', NULL, 1, 9, 'API管理', 'ShieldOutline', 0, 0);
INSERT INTO `menu` VALUES (11, '2023-11-27 23:43:00.083', '2023-11-27 23:47:13.533', NULL, 'SystemMenu', 'system/menu', 'views/system/menu/index.vue', NULL, 1, 9, '菜单管理', 'GridOutline', 0, 0);
INSERT INTO `menu` VALUES (12, '2023-11-27 23:44:10.955', '2023-11-27 23:47:22.606', NULL, 'SystemRole', 'system/role', 'views/system/role/index.vue', NULL, 1, 9, '角色管理', 'PeopleOutline', 0, 0);
INSERT INTO `menu` VALUES (13, '2024-04-13 15:40:49.675', '2024-04-13 15:41:15.863', NULL, 'SysUser', 'system/user', 'views/system/user/index.vue', NULL, 1, 9, '用户管理', 'PersonOutline', 0, 0);
INSERT INTO `menu` VALUES (14, '2024-05-18 15:06:40.539', '2024-05-18 15:06:40.539', NULL, 'ContentAnnounce', 'content/announce', 'views/content/announce/index.vue', NULL, 1, 4, '公告管理', 'TodayOutline', 0, 0);
INSERT INTO `menu` VALUES (15, '2024-06-08 14:17:20.000', '2024-06-08 14:17:20.000', NULL, 'SysConfig', 'system/config', 'views/system/config/index.vue', NULL, 1, 9, '系统配置', 'BriefcaseOutline', 0, 0);

-- 创建角色表
CREATE TABLE IF NOT EXISTS `role`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '角色名',
  `code` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '角色代码',
  `desc` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '介绍',
  `home_page` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '角色首页',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name`) USING BTREE,
  UNIQUE INDEX `code`(`code`) USING BTREE,
  INDEX `idx_role_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- 插入初始角色数据
INSERT INTO `role` VALUES (1, '2023-11-24 14:32:33.000', '2023-11-27 23:48:47.500', NULL, '用户', '001', NULL, NULL);
INSERT INTO `role` VALUES (2, '2023-11-24 14:32:54.000', '2024-05-18 15:06:57.535', NULL, '超级管理员', '002', NULL, 'SystemRole');

-- 创建角色菜单表
CREATE TABLE IF NOT EXISTS `role_menu`  (
  `menu_id` bigint UNSIGNED NOT NULL,
  `role_id` bigint UNSIGNED NOT NULL,
  PRIMARY KEY (`menu_id`, `role_id`) USING BTREE,
  INDEX `fk_role_menu_role`(`role_id`) USING BTREE,
  CONSTRAINT `fk_role_menu_menu` FOREIGN KEY (`menu_id`) REFERENCES `menu` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_role_menu_role` FOREIGN KEY (`role_id`) REFERENCES `role` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- 插入初始角色菜单数据
INSERT INTO `role_menu` VALUES (1, 2);
INSERT INTO `role_menu` VALUES (2, 2);
INSERT INTO `role_menu` VALUES (3, 2);
INSERT INTO `role_menu` VALUES (4, 2);
INSERT INTO `role_menu` VALUES (5, 2);
INSERT INTO `role_menu` VALUES (6, 2);
INSERT INTO `role_menu` VALUES (7, 2);
INSERT INTO `role_menu` VALUES (8, 2);
INSERT INTO `role_menu` VALUES (9, 2);
INSERT INTO `role_menu` VALUES (10, 2);
INSERT INTO `role_menu` VALUES (11, 2);
INSERT INTO `role_menu` VALUES (12, 2);
INSERT INTO `role_menu` VALUES (13, 2);
INSERT INTO `role_menu` VALUES (14, 2);
INSERT INTO `role_menu` VALUES (15, 2);

-- 创建用户表
CREATE TABLE IF NOT EXISTS `user`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `username` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名',
  `email` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '邮箱',
  `password` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '头像',
  `space_cover` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '空间封面',
  `gender` tinyint NULL DEFAULT 0 COMMENT '性别:0未知、1男、3女',
  `birthday` datetime(3) NULL DEFAULT '1970-01-01 00:00:00.000' COMMENT '生日',
  `sign` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '这个人很懒，什么都没有留下' COMMENT '个性签名',
  `client_ip` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '客户端IP',
  `role` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '001' COMMENT '角色ID',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_deleted_at`(`deleted_at`) USING BTREE,
  INDEX `idx_user_username`(`username`) USING BTREE,
  INDEX `idx_user_email`(`email`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- 插入初始用户数据
INSERT INTO `user` VALUES (1, '2023-11-13 23:12:06.565', '2024-04-04 16:04:23.559', NULL, '超级管理员', 'admin@admin.com', '$2a$10$syHVkGIzJL4H4cwKa0/eOOy7KxakWunVbLQUT.eJk0ayzcuVqr56u', NULL, NULL, 0, '1970-01-01 00:00:00.000', '这个人很懒，什么都没有留下', NULL, '002');
