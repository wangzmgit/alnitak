package initialize

import (
	"go.uber.org/zap"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
)

func InitDefaultData() {
	initApiData()        // 初始化API数据
	initCasbinRuleData() // 初始化CasbinRule数据
	initMenuData()       // 初始化菜单数据
	initPartitionData()  // 初始化分区数据
	initRoleData()       // 初始化角色数据
	initUserData()       // 初始化用户数据
}

// 初始化API数据
func initApiData() {
	var total int64
	global.Mysql.Model(&model.Api{}).Count(&total)
	if total > 0 {
		return
	}

	zap.L().Info("API数据不存在，添加默认数据", zap.String("module", "initialize"))
	entities := []model.Api{
		{Method: "POST", Path: "/api/v1/api/addApi", Category: "API管理", Desc: "新增API（后台管理）"},
		{Method: "DELETE", Path: "/api/v1/api/deleteApi/:id", Category: "API管理", Desc: "删除API（后台管理）"},
		{Method: "PUT", Path: "/api/v1/api/editApi", Category: "API管理", Desc: "编辑API（后台管理）"},
		{Method: "PUT", Path: "/api/v1/api/editRoleApi", Category: "API管理", Desc: "编辑角色API（后台管理）"},
		{Method: "GET", Path: "/api/v1/api/getAllApiList", Category: "API管理", Desc: "获取全部API列表（后台管理）"},
		{Method: "POST", Path: "/api/v1/api/getApiList", Category: "API管理", Desc: "获取API列表（后台管理）"},
		{Method: "GET", Path: "/api/v1/api/getRoleApi", Category: "API管理", Desc: "获取角色API（后台管理）"},
		{Method: "POST", Path: "/api/v1/archive/article/cancelCollect", Category: "点赞收藏", Desc: "文章取消收藏"},
		{Method: "POST", Path: "/api/v1/archive/article/cancelLike", Category: "点赞收藏", Desc: "文章取消赞"},
		{Method: "POST", Path: "/api/v1/archive/article/collect", Category: "点赞收藏", Desc: "文章收藏"},
		{Method: "GET", Path: "/api/v1/archive/article/hasCollect", Category: "点赞收藏", Desc: "文章是否收藏"},
		{Method: "GET", Path: "/api/v1/archive/article/hasLike", Category: "点赞收藏", Desc: "文章是否点赞"},
		{Method: "POST", Path: "/api/v1/archive/article/like", Category: "点赞收藏", Desc: "文章点赞"},
		{Method: "POST", Path: "/api/v1/archive/video/cancelLike", Category: "点赞收藏", Desc: "视频取消赞"},
		{Method: "POST", Path: "/api/v1/archive/video/collect", Category: "点赞收藏", Desc: "视频收藏"},
		{Method: "GET", Path: "/api/v1/archive/video/getCollectInfo", Category: "点赞收藏", Desc: "获取视频已收藏的文件夹"},
		{Method: "GET", Path: "/api/v1/archive/video/hasCollect", Category: "点赞收藏", Desc: "视频是否收藏"},
		{Method: "GET", Path: "/api/v1/archive/video/hasLike", Category: "点赞收藏", Desc: "视频是否点赞"},
		{Method: "POST", Path: "/api/v1/archive/video/like", Category: "点赞收藏", Desc: "视频点赞"},
		{Method: "DELETE", Path: "/api/v1/article/deleteArticle/:id", Category: "文章", Desc: "删除文章"},
		{Method: "DELETE", Path: "/api/v1/article/deleteArticleManage/:id", Category: "文章", Desc: "删除文章（后台管理）"},
		{Method: "PUT", Path: "/api/v1/article/editArticleInfo", Category: "文章", Desc: "编辑文章信息"},
		{Method: "GET", Path: "/api/v1/article/getAllArticleList", Category: "文章", Desc: "获取所有的文章列表"},
		{Method: "POST", Path: "/api/v1/article/getArticleListManage", Category: "文章", Desc: "获取文章列表（后台管理）"},
		{Method: "GET", Path: "/api/v1/article/getArticleStatus", Category: "文章", Desc: "获取文章状态信息"},
		{Method: "POST", Path: "/api/v1/article/getReviewArticleList", Category: "文章", Desc: "获取文章审核列表（后台管理）"},
		{Method: "GET", Path: "/api/v1/article/getUploadArticle", Category: "文章", Desc: "获取上传的文章"},
		{Method: "POST", Path: "/api/v1/article/uploadArticleInfo", Category: "文章", Desc: "上传文章信息"},
		{Method: "POST", Path: "/api/v1/auth/logout", Category: "Auth", Desc: "退出登录"},
		{Method: "POST", Path: "/api/v1/carousel/addCarousel", Category: "轮播图", Desc: "新增轮播图（后台管理）"},
		{Method: "DELETE", Path: "/api/v1/carousel/deleteCarousel/:id", Category: "轮播图", Desc: "删除轮播图（后台管理）"},
		{Method: "PUT", Path: "/api/v1/carousel/editCarousel", Category: "轮播图", Desc: "编辑轮播图（后台管理）"},
		{Method: "POST", Path: "/api/v1/carousel/getCarouselList", Category: "轮播图", Desc: "获取轮播图列表（后台管理）"},
		{Method: "POST", Path: "/api/v1/collection/addCollection", Category: "收藏夹", Desc: "添加收藏夹"},
		{Method: "DELETE", Path: "/api/v1/collection/deleteCollection/:id", Category: "收藏夹", Desc: "删除收藏夹"},
		{Method: "PUT", Path: "/api/v1/collection/editCollection", Category: "收藏夹", Desc: "编辑收藏夹"},
		{Method: "GET", Path: "/api/v1/collection/getCollectionInfo", Category: "收藏夹", Desc: "获取收藏夹信息"},
		{Method: "GET", Path: "/api/v1/collection/getCollectionList", Category: "收藏夹", Desc: "获取收藏夹列表"},
		{Method: "GET", Path: "/api/v1/collection/getVideoList", Category: "收藏夹", Desc: "获取收藏夹的视频列表"},
		{Method: "POST", Path: "/api/v1/comment/article/addComment", Category: "评论回复", Desc: "发表文章评论或回复"},
		{Method: "GET", Path: "/api/v1/comment/article/getCommentList", Category: "评论回复", Desc: "获取文章评论列表"},
		{Method: "DELETE", Path: "/api/v1/comment/varticle/deleteComment/:id", Category: "评论回复", Desc: "删除文章评论或回复"},
		{Method: "POST", Path: "/api/v1/comment/video/addComment", Category: "评论回复", Desc: "发表视频评论或回复"},
		{Method: "DELETE", Path: "/api/v1/comment/video/deleteComment/:id", Category: "评论回复", Desc: "删除视频评论或回复"},
		{Method: "GET", Path: "/api/v1/comment/video/getCommentList", Category: "评论回复", Desc: "获取视频评论列表"},
		{Method: "POST", Path: "/api/v1/danmaku/sendDanmaku", Category: "弹幕", Desc: "发送弹幕"},
		{Method: "POST", Path: "/api/v1/history/video/addHistory", Category: "历史记录", Desc: "保存视频历史记录"},
		{Method: "GET", Path: "/api/v1/history/video/getHistory", Category: "历史记录", Desc: "获取视频历史记录"},
		{Method: "GET", Path: "/api/v1/history/video/getProgress", Category: "历史记录", Desc: "获取视频播放进度"},
		{Method: "POST", Path: "/api/v1/menu/addMenu", Category: "菜单管理", Desc: "添加菜单（后台管理）"},
		{Method: "DELETE", Path: "/api/v1/menu/deleteMenu/:id", Category: "菜单管理", Desc: "删除菜单（后台管理）"},
		{Method: "PUT", Path: "/api/v1/menu/editMenu", Category: "菜单管理", Desc: "编辑菜单（后台管理）"},
		{Method: "PUT", Path: "/api/v1/menu/editRoleMenu", Category: "菜单管理", Desc: "编辑角色菜单（后台管理）"},
		{Method: "GET", Path: "/api/v1/menu/getMenuTree", Category: "菜单管理", Desc: "获取菜单树（后台管理）"},
		{Method: "GET", Path: "/api/v1/menu/getRoleMenu", Category: "菜单管理", Desc: "获取角色菜单（后台管理）"},
		{Method: "GET", Path: "/api/v1/menu/getUserMenu", Category: "菜单管理", Desc: "获取用户菜单树（后台管理）"},
		{Method: "POST", Path: "/api/v1/message/addAnnounce", Category: "消息", Desc: "添加公告（后台管理）"},
		{Method: "DELETE", Path: "/api/v1/message/deleteAnnounce/:id", Category: "消息", Desc: "删除公告（后台管理）"},
		{Method: "GET", Path: "/api/v1/message/getAtMsg", Category: "消息", Desc: "获取AT通知"},
		{Method: "GET", Path: "/api/v1/message/getLikeMsg", Category: "消息", Desc: "获取点赞通知"},
		{Method: "GET", Path: "/api/v1/message/getReplyMsg", Category: "消息", Desc: "获取回复通知"},
		{Method: "GET", Path: "/api/v1/message/getWhisperDetails", Category: "消息", Desc: "获取私信详情"},
		{Method: "GET", Path: "/api/v1/message/getWhisperList", Category: "消息", Desc: "获取私信列表"},
		{Method: "POST", Path: "/api/v1/message/readWhisper", Category: "消息", Desc: "已读私信"},
		{Method: "POST", Path: "/api/v1/message/sendWhisper", Category: "消息", Desc: "发送私信"},
		{Method: "POST", Path: "/api/v1/partition/addPartition", Category: "分区", Desc: "添加分区（后台管理）"},
		{Method: "DELETE", Path: "/api/v1/partition/deletePartition/:id", Category: "分区", Desc: "删除分区（后台管理）"},
		{Method: "POST", Path: "/api/v1/relation/follow", Category: "关注", Desc: "关注用户"},
		{Method: "GET", Path: "/api/v1/relation/getUserRelation", Category: "关注", Desc: "获取用户关系"},
		{Method: "POST", Path: "/api/v1/relation/unfollow", Category: "关注", Desc: "取关用户"},
		{Method: "DELETE", Path: "/api/v1/resource/deleteResource/:id", Category: "资源", Desc: "删除视频资源"},
		{Method: "PUT", Path: "/api/v1/resource/modifyTitle", Category: "资源", Desc: "修改资源标题"},
		{Method: "GET", Path: "/api/v1/review/getArticleReviewRecord", Category: "审核", Desc: "获取文章审核记录"},
		{Method: "GET", Path: "/api/v1/review/getVideoReviewRecord", Category: "审核", Desc: "获取视频审核记录"},
		{Method: "POST", Path: "/api/v1/review/reviewArticleApproved", Category: "审核", Desc: "文章审核通过（后台管理）"},
		{Method: "POST", Path: "/api/v1/review/reviewArticleFailed", Category: "审核", Desc: "文章审核不通过（后台管理）"},
		{Method: "POST", Path: "/api/v1/review/reviewVideoApproved", Category: "审核", Desc: "视频审核通过（后台管理）"},
		{Method: "POST", Path: "/api/v1/review/reviewVideoFailed", Category: "审核", Desc: "视频审核不通过（后台管理）"},
		{Method: "POST", Path: "/api/v1/role/addRole", Category: "角色", Desc: "添加角色（后台管理）"},
		{Method: "DELETE", Path: "/api/v1/role/deleteRole/:id", Category: "角色", Desc: "删除角色（后台管理）"},
		{Method: "PUT", Path: "/api/v1/role/editRole", Category: "角色", Desc: "编辑角色（后台管理）"},
		{Method: "PUT", Path: "/api/v1/role/editRoleHome", Category: "角色", Desc: "编辑角色首页（后台管理）"},
		{Method: "GET", Path: "/api/v1/role/getAllRoleList", Category: "角色", Desc: "获取全部角色（后台管理）"},
		{Method: "GET", Path: "/api/v1/role/getRoleInfo", Category: "角色", Desc: "获取个人角色信息（后台管理）"},
		{Method: "POST", Path: "/api/v1/role/getRoleList", Category: "角色", Desc: "获取角色列表（后台管理）"},
		{Method: "POST", Path: "/api/v1/upload/image", Category: "上传", Desc: "上传图片"},
		{Method: "POST", Path: "/api/v1/upload/video", Category: "上传", Desc: "上传视频"},
		{Method: "POST", Path: "/api/v1/upload/video/:vid", Category: "上传", Desc: "上传视频分P"},
		{Method: "DELETE", Path: "/api/v1/user/deleteUser/:id", Category: "用户", Desc: "删除用户"},
		{Method: "PUT", Path: "/api/v1/user/editUserInfo", Category: "用户", Desc: "编辑用户信息"},
		{Method: "PUT", Path: "/api/v1/user/editUserInfoManage", Category: "用户", Desc: "编辑用户信息（后台管理）"},
		{Method: "PUT", Path: "/api/v1/user/editUserRole", Category: "用户", Desc: "编辑用户角色"},
		{Method: "GET", Path: "/api/v1/user/getUserInfo", Category: "用户", Desc: "获取用户信息"},
		{Method: "POST", Path: "/api/v1/user/getUserListManage", Category: "用户", Desc: "获取用户列表（后台管理）"},
		{Method: "DELETE", Path: "/api/v1/video/deleteVideo/:id", Category: "视频", Desc: "删除视频"},
		{Method: "DELETE", Path: "/api/v1/video/deleteVideoManage/:id", Category: "视频", Desc: "删除视频（后台管理）"},
		{Method: "PUT", Path: "/api/v1/video/editVideoInfo", Category: "视频", Desc: "编辑视频信息"},
		{Method: "GET", Path: "/api/v1/video/getAllVideoList", Category: "视频", Desc: "获取所有的视频列表"},
		{Method: "POST", Path: "/api/v1/video/getReviewList", Category: "视频", Desc: "获取审核列表（后台管理）"},
		{Method: "GET", Path: "/api/v1/video/getReviewResourceList", Category: "视频", Desc: "获取审核资源列表（后台管理）"},
		{Method: "GET", Path: "/api/v1/video/getUploadVideo", Category: "视频", Desc: "获取上传的视频"},
		{Method: "POST", Path: "/api/v1/video/getVideoListManage", Category: "视频", Desc: "获取视频列表（后台管理）"},
		{Method: "GET", Path: "/api/v1/video/getVideoStatus", Category: "视频", Desc: "获取上传视频状态信息"},
		{Method: "POST", Path: "/api/v1/video/uploadVideoInfo", Category: "视频", Desc: "上传视频信息"},
		{Method: "GET", Path: "/api/v1/video/getResourceQualityManage", Category: "视频", Desc: "获取视频资源支持的分辨率信息（后台管理）"},
		{Method: "GET", Path: "/api/v1/video/getVideoFileManage", Category: "视频", Desc: "获取视频文件URL（后台管理）"},
		{Method: "GET", Path: "/api/v1/config/getEmailConfig", Category: "配置", Desc: "获取邮箱配置（后台管理）"},
		{Method: "POST", Path: "/api/v1/config/setEmailConfig", Category: "配置", Desc: "编辑邮箱配置（后台管理）"},
		{Method: "GET", Path: "/api/v1/config/getStorageConfig", Category: "配置", Desc: "获取存储配置（后台管理）"},
		{Method: "POST", Path: "/api/v1/config/setStorageConfig", Category: "配置", Desc: "编辑存储配置（后台管理）"},
		{Method: "GET", Path: "/api/v1/config/getOtherConfig", Category: "配置", Desc: "获取其他配置（后台管理）"},
		{Method: "POST", Path: "/api/v1/config/setOtherConfig", Category: "配置", Desc: "编辑其他配置（后台管理）"},
	}
	if err := global.Mysql.Create(&entities).Error; err != nil {
		zap.L().Error("API数据初始化失败", zap.String("err", err.Error()), zap.String("module", "initialize"))
	}
}

// 初始化CasbinRule数据
func initCasbinRuleData() {
	var total int64
	global.Mysql.Model(&model.CasbinRule{}).Count(&total)
	if total > 0 {
		return
	}

	zap.L().Info("CasbinRule数据不存在，添加默认数据", zap.String("module", "initialize"))
	entities := []model.CasbinRule{
		{Ptype: "p", V0: "001", V1: "/api/v1/archive/article/cancelCollect", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/archive/article/cancelLike", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/archive/article/collect", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/archive/article/hasCollect", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/archive/article/hasLike", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/archive/article/like", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/archive/video/cancelLike", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/archive/video/collect", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/archive/video/getCollectInfo", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/archive/video/hasCollect", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/archive/video/hasLike", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/archive/video/like", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/article/deleteArticle/:id", V2: "DELETE"},
		{Ptype: "p", V0: "001", V1: "/api/v1/article/editArticleInfo", V2: "PUT"},
		{Ptype: "p", V0: "001", V1: "/api/v1/article/getAllArticleList", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/article/getArticleStatus", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/article/getUploadArticle", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/article/uploadArticleInfo", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/auth/logout", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/collection/addCollection", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/collection/deleteCollection/:id", V2: "DELETE"},
		{Ptype: "p", V0: "001", V1: "/api/v1/collection/editCollection", V2: "PUT"},
		{Ptype: "p", V0: "001", V1: "/api/v1/collection/getCollectionInfo", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/collection/getCollectionList", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/collection/getVideoList", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/comment/article/addComment", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/comment/article/getCommentList", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/comment/varticle/deleteComment/:id", V2: "DELETE"},
		{Ptype: "p", V0: "001", V1: "/api/v1/comment/video/addComment", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/comment/video/deleteComment/:id", V2: "DELETE"},
		{Ptype: "p", V0: "001", V1: "/api/v1/comment/video/getCommentList", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/danmaku/sendDanmaku", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/history/video/addHistory", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/history/video/getHistory", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/history/video/getProgress", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/message/getAtMsg", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/message/getLikeMsg", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/message/getReplyMsg", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/message/getWhisperDetails", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/message/getWhisperList", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/message/readWhisper", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/message/sendWhisper", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/relation/follow", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/relation/getUserRelation", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/relation/unfollow", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/resource/deleteResource/:id", V2: "DELETE"},
		{Ptype: "p", V0: "001", V1: "/api/v1/resource/modifyTitle", V2: "PUT"},
		{Ptype: "p", V0: "001", V1: "/api/v1/review/getArticleReviewRecord", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/review/getVideoReviewRecord", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/upload/image", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/upload/video", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/upload/video/:vid", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/user/deleteUser/:id", V2: "DELETE"},
		{Ptype: "p", V0: "001", V1: "/api/v1/user/editUserInfo", V2: "PUT"},
		{Ptype: "p", V0: "001", V1: "/api/v1/user/editUserRole", V2: "PUT"},
		{Ptype: "p", V0: "001", V1: "/api/v1/user/getUserInfo", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/video/deleteVideo/:id", V2: "DELETE"},
		{Ptype: "p", V0: "001", V1: "/api/v1/video/editVideoInfo", V2: "PUT"},
		{Ptype: "p", V0: "001", V1: "/api/v1/video/getAllVideoList", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/video/getUploadVideo", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/video/getVideoStatus", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/video/uploadVideoInfo", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/api/addApi", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/api/deleteApi/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/api/editApi", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/api/editRoleApi", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/api/getAllApiList", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/api/getApiList", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/api/getRoleApi", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/archive/article/cancelCollect", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/archive/article/cancelLike", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/archive/article/collect", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/archive/article/hasCollect", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/archive/article/hasLike", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/archive/article/like", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/archive/video/cancelLike", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/archive/video/collect", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/archive/video/getCollectInfo", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/archive/video/hasCollect", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/archive/video/hasLike", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/archive/video/like", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/article/deleteArticle/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/article/deleteArticleManage/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/article/editArticleInfo", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/article/getAllArticleList", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/article/getArticleListManage", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/article/getArticleStatus", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/article/getReviewArticleList", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/article/getUploadArticle", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/article/uploadArticleInfo", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/auth/logout", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/carousel/addCarousel", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/carousel/deleteCarousel/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/carousel/editCarousel", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/carousel/getCarouselList", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/collection/addCollection", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/collection/deleteCollection/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/collection/editCollection", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/collection/getCollectionInfo", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/collection/getCollectionList", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/collection/getVideoList", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/comment/article/addComment", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/comment/article/getCommentList", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/comment/varticle/deleteComment/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/comment/video/addComment", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/comment/video/deleteComment/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/comment/video/getCommentList", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/config/getEmailConfig", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/config/getOtherConfig", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/config/getStorageConfig", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/config/setEmailConfig", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/config/setOtherConfig", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/config/setStorageConfig", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/danmaku/sendDanmaku", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/history/video/addHistory", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/history/video/getHistory", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/history/video/getProgress", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/menu/addMenu", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/menu/deleteMenu/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/menu/editMenu", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/menu/editRoleMenu", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/menu/getMenuTree", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/menu/getRoleMenu", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/menu/getUserMenu", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/message/addAnnounce", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/message/deleteAnnounce/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/message/getAtMsg", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/message/getLikeMsg", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/message/getReplyMsg", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/message/getWhisperDetails", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/message/getWhisperList", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/message/readWhisper", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/message/sendWhisper", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/partition/addPartition", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/partition/deletePartition/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/relation/follow", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/relation/getUserRelation", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/relation/unfollow", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/resource/deleteResource/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/resource/modifyTitle", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/review/getArticleReviewRecord", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/review/getVideoReviewRecord", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/review/reviewArticleApproved", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/review/reviewArticleFailed", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/review/reviewVideoApproved", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/review/reviewVideoFailed", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/role/addRole", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/role/deleteRole/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/role/editRole", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/role/editRoleHome", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/role/getAllRoleList", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/role/getRoleInfo", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/role/getRoleList", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/upload/image", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/upload/video", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/upload/video/:vid", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/user/deleteUser/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/user/editUserInfo", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/user/editUserInfoManage", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/user/editUserRole", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/user/getUserInfo", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/user/getUserListManage", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/video/deleteVideo/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/video/deleteVideoManage/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/video/editVideoInfo", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/video/getAllVideoList", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/video/getResourceQualityManage", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/video/getReviewList", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/video/getReviewResourceList", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/video/getUploadVideo", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/video/getVideoFileManage", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/video/getVideoListManage", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/video/getVideoStatus", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/video/uploadVideoInfo", V2: "POST"},
	}
	if err := global.Mysql.Create(&entities).Error; err != nil {
		zap.L().Error("CasbinRule数据初始化失败", zap.String("err", err.Error()), zap.String("module", "initialize"))
	}
}

// 初始化菜单数据
func initMenuData() {
	var total int64
	global.Mysql.Model(&model.Menu{}).Count(&total)
	if total > 0 {
		return
	}

	zap.L().Info("菜单数据不存在，添加默认数据", zap.String("module", "initialize"))

	reviewMenu := model.Menu{Name: "review", Path: "review", Component: "", Desc: "", Sort: 1, ParentId: 0, Title: "内容审核", Icon: "LayersOutline", Hidden: false, KeepAlive: false}
	if err := global.Mysql.Create(&reviewMenu).Error; err != nil {
		zap.L().Error("菜单数据初始化失败", zap.String("err", err.Error()), zap.String("module", "initialize"))
	}
	reviewMenuEntities := []model.Menu{
		{Name: "ReviewVideo", Path: "review/video", Component: "views/review/video/index.vue", Desc: "", Sort: 1, ParentId: reviewMenu.ID, Title: "视频审核", Icon: "FileTrayOutline", Hidden: false, KeepAlive: false},
		{Name: "ReviewArticle", Path: "review/article", Component: "views/review/article/index.vue", Desc: "", Sort: 1, ParentId: reviewMenu.ID, Title: "专栏审核", Icon: "AlbumsOutline", Hidden: false, KeepAlive: false},
	}
	if err := global.Mysql.Create(&reviewMenuEntities).Error; err != nil {
		zap.L().Error("菜单数据初始化失败", zap.String("err", err.Error()), zap.String("module", "initialize"))
	}

	contentMenu := model.Menu{Name: "Content", Path: "content", Component: "", Desc: "", Sort: 1, ParentId: 0, Title: "内容管理", Icon: "ReaderOutline", Hidden: false, KeepAlive: false}
	if err := global.Mysql.Create(&contentMenu).Error; err != nil {
		zap.L().Error("菜单数据初始化失败", zap.String("err", err.Error()), zap.String("module", "initialize"))
	}
	contentMenuEntities := []model.Menu{
		{Name: "ContentVideo", Path: "content/video", Component: "views/content/video/index.vue", Desc: "", Sort: 1, ParentId: contentMenu.ID, Title: "视频管理", Icon: "PlayCircleOutline", Hidden: false, KeepAlive: false},
		{Name: "ContentArticle", Path: "content/article", Component: "views/content/article/index.vue", Desc: "", Sort: 1, ParentId: contentMenu.ID, Title: "专栏管理", Icon: "DocumentTextOutline", Hidden: false, KeepAlive: false},
		{Name: "ContentCarousel", Path: "content/carousel", Component: "views/content/carousel/index.vue", Desc: "", Sort: 1, ParentId: contentMenu.ID, Title: "轮播图管理", Icon: "ImagesOutline", Hidden: false, KeepAlive: false},
		{Name: "ContentPartition", Path: "content/partition", Component: "views/content/partition/index.vue", Desc: "", Sort: 1, ParentId: contentMenu.ID, Title: "分区管理", Icon: "BookmarkOutline", Hidden: false, KeepAlive: false},
		{Name: "ContentAnnounce", Path: "content/announce", Component: "views/content/announce/index.vue", Desc: "", Sort: 1, ParentId: contentMenu.ID, Title: "公告管理", Icon: "TodayOutline", Hidden: false, KeepAlive: false},
	}
	if err := global.Mysql.Create(&contentMenuEntities).Error; err != nil {
		zap.L().Error("菜单数据初始化失败", zap.String("err", err.Error()), zap.String("module", "initialize"))
	}

	systemMenu := model.Menu{Name: "System", Path: "system", Component: "", Desc: "", Sort: 1, ParentId: 0, Title: "系统管理", Icon: "TerminalOutline", Hidden: false, KeepAlive: false}
	if err := global.Mysql.Create(&systemMenu).Error; err != nil {
		zap.L().Error("菜单数据初始化失败", zap.String("err", err.Error()), zap.String("module", "initialize"))
	}
	systemMenuEntities := []model.Menu{
		{Name: "SystemApi", Path: "system/api", Component: "views/system/api/index.vue", Desc: "", Sort: 1, ParentId: systemMenu.ID, Title: "API管理", Icon: "ShieldOutline", Hidden: false, KeepAlive: false},
		{Name: "SystemMenu", Path: "system/menu", Component: "views/system/menu/index.vue", Desc: "", Sort: 1, ParentId: systemMenu.ID, Title: "菜单管理", Icon: "GridOutline", Hidden: false, KeepAlive: false},
		{Name: "SystemRole", Path: "system/role", Component: "views/system/role/index.vue", Desc: "", Sort: 1, ParentId: systemMenu.ID, Title: "角色管理", Icon: "PeopleOutline", Hidden: false, KeepAlive: false},
		{Name: "SysUser", Path: "system/user", Component: "views/system/user/index.vue", Desc: "", Sort: 1, ParentId: systemMenu.ID, Title: "用户管理", Icon: "PersonOutline", Hidden: false, KeepAlive: false},
		{Name: "SysConfig", Path: "system/config", Component: "views/system/config/index.vue", Desc: "", Sort: 1, ParentId: systemMenu.ID, Title: "系统配置", Icon: "BriefcaseOutline", Hidden: false, KeepAlive: false},
	}
	if err := global.Mysql.Create(&systemMenuEntities).Error; err != nil {
		zap.L().Error("菜单数据初始化失败", zap.String("err", err.Error()), zap.String("module", "initialize"))
	}
}

// 初始化分区数据
func initPartitionData() {
	var total int64
	global.Mysql.Model(&model.Partition{}).Count(&total)
	if total > 0 {
		return
	}

	zap.L().Info("分区数据不存在，添加默认数据", zap.String("module", "initialize"))
	videoPartition := model.Partition{Name: "生活", Type: 0, ParentId: 0}
	if err := global.Mysql.Create(&videoPartition).Error; err != nil {
		zap.L().Error("分区数据初始化失败", zap.String("err", err.Error()), zap.String("module", "initialize"))
	}

	if err := global.Mysql.Create(&model.Partition{Name: "日常", Type: 0, ParentId: videoPartition.ID}).Error; err != nil {
		zap.L().Error("分区数据初始化失败", zap.String("err", err.Error()), zap.String("module", "initialize"))
	}

	articlePartition := model.Partition{Name: "生活", Type: 1, ParentId: 0}
	if err := global.Mysql.Create(&articlePartition).Error; err != nil {
		zap.L().Error("分区数据初始化失败", zap.String("err", err.Error()), zap.String("module", "initialize"))
	}

	if err := global.Mysql.Create(&model.Partition{Name: "日常", Type: 1, ParentId: articlePartition.ID}).Error; err != nil {
		zap.L().Error("分区数据初始化失败", zap.String("err", err.Error()), zap.String("module", "initialize"))
	}
}

// 初始化角色数据
func initRoleData() {
	var total int64
	global.Mysql.Model(&model.Role{}).Count(&total)
	if total > 0 {
		return
	}

	zap.L().Info("角色数据不存在，添加默认数据", zap.String("module", "initialize"))
	if err := global.Mysql.Create(&model.Role{Name: "用户", Code: "001", Desc: "", HomePage: ""}).Error; err != nil {
		zap.L().Error("角色数据初始化失败", zap.String("err", err.Error()), zap.String("module", "initialize"))
	}

	admin := model.Role{Name: "超级管理员", Code: "002", Desc: "", HomePage: ""}
	if err := global.Mysql.Create(&admin).Error; err != nil {
		zap.L().Error("角色数据初始化失败", zap.String("err", err.Error()), zap.String("module", "initialize"))
	}

	var menus []model.Menu
	global.Mysql.Model(&model.Menu{}).Find(&menus)
	admin.Menus = menus
	global.Mysql.Model(&admin).Association("Menus").Replace(admin.Menus)
}

// 初始化用户数据
func initUserData() {
	var total int64
	global.Mysql.Model(&model.User{}).Count(&total)
	if total > 0 {
		return
	}

	zap.L().Info("用户数据不存在，添加默认数据", zap.String("module", "initialize"))
	entities := []model.User{
		{Username: "超级管理员", Email: "admin@admin.com", Password: "$2a$10$syHVkGIzJL4H4cwKa0/eOOy7KxakWunVbLQUT.eJk0ayzcuVqr56u", Role: "002"},
	}
	if err := global.Mysql.Create(&entities).Error; err != nil {
		zap.L().Error("用户数据初始化失败", zap.String("err", err.Error()), zap.String("module", "initialize"))
	}
}
