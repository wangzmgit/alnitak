package initialize

import (
	"strings"

	"go.uber.org/zap"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
)

func InitDefaultData() {
	initApiData()        // 初始化API数据
	initCasbinRuleData() // 初始化CasbinRule数据
	initMenuData()       // 初始化菜单数据
	syncMenuData()       // 增量补全菜单（含菜单表、role_menu）
	initPartitionData()  // 初始化分区数据
	initRoleData()       // 初始化角色数据
	initUserData()       // 初始化用户数据
	// 每次启动将 authApiDesc 中尚未入库的接口写入 API 表，并补全 001/002 的 Casbin 规则（须在 InitCasbin 之前执行）
	SyncApiData()
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
		{Method: "POST", Path: "/api/v1/archive/article/share", Category: "点赞收藏", Desc: "文章分享计数"},
		{Method: "POST", Path: "/api/v1/archive/video/cancelLike", Category: "点赞收藏", Desc: "视频取消赞"},
		{Method: "POST", Path: "/api/v1/archive/video/collect", Category: "点赞收藏", Desc: "视频收藏"},
		{Method: "GET", Path: "/api/v1/archive/video/getCollectInfo", Category: "点赞收藏", Desc: "获取视频已收藏的文件夹"},
		{Method: "GET", Path: "/api/v1/archive/video/hasCollect", Category: "点赞收藏", Desc: "视频是否收藏"},
		{Method: "GET", Path: "/api/v1/archive/video/hasLike", Category: "点赞收藏", Desc: "视频是否点赞"},
		{Method: "POST", Path: "/api/v1/archive/video/like", Category: "点赞收藏", Desc: "视频点赞"},
		{Method: "POST", Path: "/api/v1/archive/video/share", Category: "点赞收藏", Desc: "视频分享计数"},
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
		{Method: "POST", Path: "/api/v1/client/log", Category: "客户端", Desc: "客户端日志上报"},
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
		{Method: "DELETE", Path: "/api/v1/comment/article/deleteComment/:id", Category: "评论回复", Desc: "删除文章评论或回复"},
		{Method: "POST", Path: "/api/v1/comment/article/like/:id", Category: "评论回复", Desc: "点赞文章评论"},
		{Method: "DELETE", Path: "/api/v1/comment/article/like/:id", Category: "评论回复", Desc: "取消点赞文章评论"},
		{Method: "POST", Path: "/api/v1/comment/article/dislike/:id", Category: "评论回复", Desc: "点踩文章评论"},
		{Method: "DELETE", Path: "/api/v1/comment/article/dislike/:id", Category: "评论回复", Desc: "取消点踩文章评论"},
		{Method: "POST", Path: "/api/v1/comment/video/addComment", Category: "评论回复", Desc: "发表视频评论或回复"},
		{Method: "DELETE", Path: "/api/v1/comment/video/deleteComment/:id", Category: "评论回复", Desc: "删除视频评论或回复"},
		{Method: "GET", Path: "/api/v1/comment/video/getCommentList", Category: "评论回复", Desc: "获取视频评论列表"},
		{Method: "POST", Path: "/api/v1/comment/video/like/:id", Category: "评论回复", Desc: "点赞视频评论"},
		{Method: "DELETE", Path: "/api/v1/comment/video/like/:id", Category: "评论回复", Desc: "取消点赞视频评论"},
		{Method: "POST", Path: "/api/v1/comment/video/dislike/:id", Category: "评论回复", Desc: "点踩视频评论"},
		{Method: "DELETE", Path: "/api/v1/comment/video/dislike/:id", Category: "评论回复", Desc: "取消点踩视频评论"},
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
		{Method: "GET", Path: "/api/v1/message/readStatus", Category: "消息", Desc: "获取公告/点赞/回复/@已读进度"},
		{Method: "POST", Path: "/api/v1/message/readStatus", Category: "消息", Desc: "上报公告/点赞/回复/@已读进度"},
		{Method: "POST", Path: "/api/v1/message/sendWhisper", Category: "消息", Desc: "发送私信"},
		{Method: "POST", Path: "/api/v1/partition/addPartition", Category: "分区", Desc: "添加分区（后台管理）"},
		{Method: "DELETE", Path: "/api/v1/partition/deletePartition/:id", Category: "分区", Desc: "删除分区（后台管理）"},
		{Method: "POST", Path: "/api/v1/relation/follow", Category: "关注", Desc: "关注用户"},
		{Method: "GET", Path: "/api/v1/relation/getUserRelation", Category: "关注", Desc: "获取用户关系"},
		{Method: "POST", Path: "/api/v1/relation/unfollow", Category: "关注", Desc: "取关用户"},
		{Method: "DELETE", Path: "/api/v1/resource/deleteResource/:id", Category: "资源", Desc: "删除视频资源（可通过deleteDanmaku参数选择是否同时删除弹幕）"},
		{Method: "PUT", Path: "/api/v1/resource/modifyTitle", Category: "资源", Desc: "修改资源标题"},
		{Method: "PUT", Path: "/api/v1/resource/reorder", Category: "资源", Desc: "资源排序"},
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
		{Method: "POST", Path: "/api/v1/upload/checkVideo", Category: "上传", Desc: "获取视频上传进度"},
		{Method: "POST", Path: "/api/v1/upload/chunkVideo", Category: "上传", Desc: "上传视频文件分片"},
		{Method: "POST", Path: "/api/v1/upload/mergeVideo", Category: "上传", Desc: "合并视频文件分片"},
		{Method: "PUT", Path: "/api/v1/user/banUser", Category: "用户", Desc: "封禁用户（后台管理）"},
		{Method: "PUT", Path: "/api/v1/user/unBanUser", Category: "用户", Desc: "解封用户（后台管理）"},
		{Method: "GET", Path: "/api/v1/user/getUserBanRecord", Category: "用户", Desc: "获取封禁记录（后台管理）"},
		{Method: "DELETE", Path: "/api/v1/user/deleteUser/:id", Category: "用户", Desc: "删除用户（后台管理）"},
		{Method: "PUT", Path: "/api/v1/user/editUserInfo", Category: "用户", Desc: "编辑用户信息"},
		{Method: "PUT", Path: "/api/v1/user/editUserInfoManage", Category: "用户", Desc: "编辑用户信息（后台管理）"},
		{Method: "PUT", Path: "/api/v1/user/editUserRole", Category: "用户", Desc: "编辑用户角色（后台管理）"},
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
		{Method: "POST", Path: "/api/v1/video/getFailedVideoList", Category: "视频", Desc: "获取处理失败的视频列表（后台管理）"},
		{Method: "POST", Path: "/api/v1/video/getProcessingVideoList", Category: "视频", Desc: "获取处理中视频列表（后台管理）"},
		{Method: "GET", Path: "/api/v1/video/getVideoStatus", Category: "视频", Desc: "获取上传视频状态信息"},
		{Method: "POST", Path: "/api/v1/video/uploadVideoInfo", Category: "视频", Desc: "上传视频信息"},
		{Method: "GET", Path: "/api/v1/video/getResourceQualityManage", Category: "视频", Desc: "获取视频资源支持的分辨率信息（后台管理）"},
		{Method: "GET", Path: "/api/v1/video/getVideoFileManage", Category: "视频", Desc: "获取视频文件URL（后台管理）"},
		{Method: "POST", Path: "/api/v1/pgc/getReviewList", Category: "PGC", Desc: "获取PGC待审列表（后台管理）"},
		{Method: "POST", Path: "/api/v1/pgc/reviewApproved", Category: "PGC", Desc: "PGC审核通过（后台管理）"},
		{Method: "POST", Path: "/api/v1/pgc/reviewFailed", Category: "PGC", Desc: "PGC审核驳回（后台管理）"},
		{Method: "POST", Path: "/api/v1/pgc/getManageList", Category: "PGC", Desc: "获取PGC管理列表（后台管理）"},
		{Method: "POST", Path: "/api/v1/pgc/adminUpdateStatus", Category: "PGC", Desc: "管理员修改PGC状态（后台管理）"},
		{Method: "DELETE", Path: "/api/v1/pgc/adminDelete/:pgc_id", Category: "PGC", Desc: "管理员删除PGC（后台管理）"},
		{Method: "GET", Path: "/api/v1/config/getEmailConfig", Category: "配置", Desc: "获取邮箱配置（后台管理）"},
		{Method: "POST", Path: "/api/v1/config/setEmailConfig", Category: "配置", Desc: "编辑邮箱配置（后台管理）"},
		{Method: "GET", Path: "/api/v1/config/getStorageConfig", Category: "配置", Desc: "获取存储配置（后台管理）"},
		{Method: "POST", Path: "/api/v1/config/setStorageConfig", Category: "配置", Desc: "编辑存储配置（后台管理）"},
		{Method: "GET", Path: "/api/v1/config/getOtherConfig", Category: "配置", Desc: "获取其他配置（后台管理）"},
		{Method: "POST", Path: "/api/v1/config/setOtherConfig", Category: "配置", Desc: "编辑其他配置（后台管理）"},
		{Method: "GET", Path: "/api/v1/config/getCleanupPreview", Category: "配置", Desc: "获取资源清理预览（后台管理）"},
		{Method: "POST", Path: "/api/v1/config/executeCleanup", Category: "配置", Desc: "执行资源清理（后台管理）"},
		// 用户认证相关
		{Method: "GET", Path: "/api/v1/auth/type/list", Category: "用户认证", Desc: "获取认证类型列表"},
		{Method: "GET", Path: "/api/v1/auth/user/list", Category: "用户认证", Desc: "获取用户认证列表"},
		{Method: "GET", Path: "/api/v1/auth/user/primary", Category: "用户认证", Desc: "获取用户主要认证"},
		{Method: "GET", Path: "/api/v1/auth/user/:uid/auth", Category: "用户认证", Desc: "获取指定用户的认证信息"},
		{Method: "POST", Path: "/api/v1/auth/type/add", Category: "用户认证", Desc: "添加认证类型（需登录）"},
		{Method: "PUT", Path: "/api/v1/auth/type/edit", Category: "用户认证", Desc: "编辑认证类型（需登录）"},
		{Method: "DELETE", Path: "/api/v1/auth/type/:id", Category: "用户认证", Desc: "删除认证类型（需登录）"},
		{Method: "GET", Path: "/api/v1/auth/type/all", Category: "用户认证", Desc: "获取所有认证类型（需登录）"},
		{Method: "GET", Path: "/api/v1/auth/type/:id", Category: "用户认证", Desc: "获取认证类型详情（需登录）"},
		{Method: "POST", Path: "/api/v1/auth/user/add", Category: "用户认证", Desc: "添加用户认证（需登录）"},
		{Method: "PUT", Path: "/api/v1/auth/user/edit", Category: "用户认证", Desc: "编辑用户认证（需登录）"},
		{Method: "DELETE", Path: "/api/v1/auth/user", Category: "用户认证", Desc: "删除用户认证（需登录）"},
		{Method: "GET", Path: "/api/v1/auth/user/all", Category: "用户认证", Desc: "获取用户认证列表（管理用，需登录）"},
		{Method: "GET", Path: "/api/v1/auth/user/:id", Category: "用户认证", Desc: "获取用户认证详情（需登录）"},
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
		{Ptype: "p", V0: "001", V1: "/api/v1/archive/article/share", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/archive/video/cancelLike", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/archive/video/collect", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/archive/video/getCollectInfo", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/archive/video/hasCollect", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/archive/video/hasLike", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/archive/video/like", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/archive/video/share", V2: "POST"},
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
		{Ptype: "p", V0: "001", V1: "/api/v1/comment/article/deleteComment/:id", V2: "DELETE"},
		{Ptype: "p", V0: "001", V1: "/api/v1/comment/article/like/:id", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/comment/article/like/:id", V2: "DELETE"},
		{Ptype: "p", V0: "001", V1: "/api/v1/comment/article/dislike/:id", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/comment/article/dislike/:id", V2: "DELETE"},
		{Ptype: "p", V0: "001", V1: "/api/v1/comment/video/addComment", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/comment/video/deleteComment/:id", V2: "DELETE"},
		{Ptype: "p", V0: "001", V1: "/api/v1/comment/video/getCommentList", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/comment/video/like/:id", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/comment/video/like/:id", V2: "DELETE"},
		{Ptype: "p", V0: "001", V1: "/api/v1/comment/video/dislike/:id", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/comment/video/dislike/:id", V2: "DELETE"},
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
		{Ptype: "p", V0: "001", V1: "/api/v1/message/readStatus", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/message/readStatus", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/message/sendWhisper", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/relation/follow", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/relation/getUserRelation", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/relation/unfollow", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/resource/deleteResource/:id", V2: "DELETE"},
		{Ptype: "p", V0: "001", V1: "/api/v1/resource/modifyTitle", V2: "PUT"},
		{Ptype: "p", V0: "001", V1: "/api/v1/resource/reorder", V2: "PUT"},
		{Ptype: "p", V0: "001", V1: "/api/v1/review/getArticleReviewRecord", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/review/getVideoReviewRecord", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/upload/image", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/upload/video", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/upload/video/:vid", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/upload/checkVideo", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/upload/chunkVideo", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/upload/mergeVideo", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/user/editUserInfo", V2: "PUT"},
		{Ptype: "p", V0: "001", V1: "/api/v1/user/getUserInfo", V2: "GET"},
		// 用户认证相关
		{Ptype: "p", V0: "001", V1: "/api/v1/auth/type/add", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/auth/type/edit", V2: "PUT"},
		{Ptype: "p", V0: "001", V1: "/api/v1/auth/type/:id", V2: "DELETE"},
		{Ptype: "p", V0: "001", V1: "/api/v1/auth/type/all", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/auth/type/:id", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/auth/user/add", V2: "POST"},
		{Ptype: "p", V0: "001", V1: "/api/v1/auth/user/edit", V2: "PUT"},
		{Ptype: "p", V0: "001", V1: "/api/v1/auth/user", V2: "DELETE"},
		{Ptype: "p", V0: "001", V1: "/api/v1/auth/user/all", V2: "GET"},
		{Ptype: "p", V0: "001", V1: "/api/v1/auth/user/:id", V2: "GET"},
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
		{Ptype: "p", V0: "002", V1: "/api/v1/archive/article/share", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/archive/video/cancelLike", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/archive/video/collect", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/archive/video/getCollectInfo", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/archive/video/hasCollect", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/archive/video/hasLike", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/archive/video/like", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/archive/video/share", V2: "POST"},
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
		{Ptype: "p", V0: "002", V1: "/api/v1/comment/article/deleteComment/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/comment/article/like/:id", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/comment/article/like/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/comment/article/dislike/:id", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/comment/article/dislike/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/comment/video/addComment", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/comment/video/deleteComment/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/comment/video/getCommentList", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/comment/video/like/:id", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/comment/video/like/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/comment/video/dislike/:id", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/comment/video/dislike/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/config/getEmailConfig", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/config/getOtherConfig", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/config/getStorageConfig", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/config/setEmailConfig", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/config/setOtherConfig", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/config/setStorageConfig", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/config/getCleanupPreview", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/config/executeCleanup", V2: "POST"},
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
		{Ptype: "p", V0: "002", V1: "/api/v1/message/readStatus", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/message/readStatus", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/message/sendWhisper", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/partition/addPartition", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/partition/deletePartition/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/relation/follow", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/relation/getUserRelation", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/relation/unfollow", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/resource/deleteResource/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/resource/modifyTitle", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/resource/reorder", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/review/getArticleReviewRecord", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/review/getVideoReviewRecord", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/review/reviewArticleApproved", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/review/reviewArticleFailed", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/review/reviewVideoApproved", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/review/reviewVideoFailed", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/pgc/getReviewList", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/pgc/reviewApproved", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/pgc/reviewFailed", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/pgc/getManageList", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/pgc/adminUpdateStatus", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/pgc/adminDelete/:pgc_id", V2: "DELETE"},
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
		{Ptype: "p", V0: "002", V1: "/api/v1/upload/checkVideo", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/upload/chunkVideo", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/upload/mergeVideo", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/user/banUser", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/user/unBanUser", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/user/getUserBanRecord", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/user/deleteUser/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/user/editUserInfo", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/user/editUserInfoManage", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/user/editUserRole", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/user/getUserInfo", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/user/getUserListManage", V2: "POST"},
		// 用户认证相关
		{Ptype: "p", V0: "002", V1: "/api/v1/auth/type/add", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/auth/type/edit", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/auth/type/:id", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/auth/type/all", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/auth/type/:id", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/auth/user/add", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/auth/user/edit", V2: "PUT"},
		{Ptype: "p", V0: "002", V1: "/api/v1/auth/user", V2: "DELETE"},
		{Ptype: "p", V0: "002", V1: "/api/v1/auth/user/all", V2: "GET"},
		{Ptype: "p", V0: "002", V1: "/api/v1/auth/user/:id", V2: "GET"},
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
		{Ptype: "p", V0: "002", V1: "/api/v1/video/getFailedVideoList", V2: "POST"},
		{Ptype: "p", V0: "002", V1: "/api/v1/video/getProcessingVideoList", V2: "POST"},
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
		{Name: "ReviewArticle", Path: "review/article", Component: "views/review/article/index.vue", Desc: "", Sort: 2, ParentId: reviewMenu.ID, Title: "专栏审核", Icon: "AlbumsOutline", Hidden: false, KeepAlive: false},
		{Name: "ReviewPlaylist", Path: "review/playlist", Component: "views/review/playlist/index.vue", Desc: "", Sort: 3, ParentId: reviewMenu.ID, Title: "合集审核", Icon: "ListOutline", Hidden: false, KeepAlive: false},
		{Name: "ReviewPGC", Path: "review/pgc", Component: "views/review/pgc/index.vue", Desc: "", Sort: 4, ParentId: reviewMenu.ID, Title: "PGC审核", Icon: "FilmOutline", Hidden: false, KeepAlive: false},
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
		{Name: "ContentPlaylist", Path: "content/playlist", Component: "views/content/playlist/index.vue", Desc: "", Sort: 2, ParentId: contentMenu.ID, Title: "合集管理", Icon: "ListOutline", Hidden: false, KeepAlive: false},
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

// ensureMenuExists 检查菜单是否存在，不存在则创建，并自动分配给拥有同级菜单的角色
// siblingName 为同级已有菜单的 Name，用于查找需要关联新菜单的角色
func ensureMenuExists(parentName string, siblingName string, menu model.Menu) {
	// 检查菜单是否已存在
	var existingMenu model.Menu
	menuExists := global.Mysql.Where("name = ?", menu.Name).First(&existingMenu).Error == nil

	if !menuExists {
		// 菜单不存在，创建
		var parent model.Menu
		if err := global.Mysql.Where("name = ?", parentName).First(&parent).Error; err != nil {
			zap.L().Error("补全菜单失败：父菜单不存在", zap.String("parent", parentName), zap.String("module", "initialize"))
			return
		}

		menu.ParentId = parent.ID
		if err := global.Mysql.Create(&menu).Error; err != nil {
			zap.L().Error("补全菜单失败", zap.String("name", menu.Name), zap.String("err", err.Error()), zap.String("module", "initialize"))
			return
		}
		existingMenu = menu
		zap.L().Info("自动补全菜单", zap.String("name", menu.Name), zap.String("title", menu.Title), zap.String("module", "initialize"))
	}

	// 检查角色关联是否存在，不存在则补全
	var siblingMenu model.Menu
	if err := global.Mysql.Where("name = ?", siblingName).First(&siblingMenu).Error; err != nil {
		return
	}

	// 查找关联了 sibling 菜单的所有角色ID
	var siblingRoleIds []uint
	global.Mysql.Table("role_menu").Where("menu_id = ?", siblingMenu.ID).Pluck("role_id", &siblingRoleIds)
	if len(siblingRoleIds) == 0 {
		return
	}

	// 查找已关联了新菜单的角色ID
	var existingRoleIds []uint
	global.Mysql.Table("role_menu").Where("menu_id = ?", existingMenu.ID).Pluck("role_id", &existingRoleIds)
	existingSet := make(map[uint]bool)
	for _, id := range existingRoleIds {
		existingSet[id] = true
	}

	// 补全缺失的角色关联
	added := 0
	for _, roleId := range siblingRoleIds {
		if !existingSet[roleId] {
			global.Mysql.Exec("INSERT INTO role_menu (role_id, menu_id) VALUES (?, ?)", roleId, existingMenu.ID)
			added++
		}
	}
	if added > 0 {
		zap.L().Info("自动分配菜单给角色", zap.String("menu", existingMenu.Name), zap.Int("count", added), zap.String("module", "initialize"))
	}
}

// syncMenuData 增量补全菜单（每次启动都会检查）
func syncMenuData() {
	// siblingName="ReviewVideo" 表示：把新菜单分配给所有已拥有"视频审核"菜单的角色
	ensureMenuExists("review", "ReviewVideo", model.Menu{
		Name: "ReviewPlaylist", Path: "review/playlist", Component: "views/review/playlist/index.vue",
		Desc: "", Sort: 3, Title: "合集审核", Icon: "ListOutline", Hidden: false, KeepAlive: false,
	})

	ensureMenuExists("review", "ReviewVideo", model.Menu{
		Name: "ReviewPGC", Path: "review/pgc", Component: "views/review/pgc/index.vue",
		Desc: "", Sort: 4, Title: "PGC审核", Icon: "FilmOutline", Hidden: false, KeepAlive: false,
	})

	// 合集管理（内容管理下）
	ensureMenuExists("Content", "ContentVideo", model.Menu{
		Name: "ContentPlaylist", Path: "content/playlist", Component: "views/content/playlist/index.vue",
		Desc: "", Sort: 2, Title: "合集管理", Icon: "ListOutline", Hidden: false, KeepAlive: false,
	})

	// PGC管理（内容管理下）
	ensureMenuExists("Content", "ContentVideo", model.Menu{
		Name: "ContentPGC", Path: "content/pgc", Component: "views/content/pgc/index.vue",
		Desc: "", Sort: 2, Title: "PGC管理", Icon: "FilmOutline", Hidden: false, KeepAlive: false,
	})

	// 同步内容管理子菜单排序：视频(1) → PGC(2) → 其余(3+)
	contentSortOrder := map[string]uint{
		"ContentVideo":     1,
		"ContentPGC":       2,
		"ContentArticle":   3,
		"ContentCarousel":  4,
		"ContentPartition": 5,
		"ContentAnnounce":  6,
		"ContentPlaylist":  7,
	}
	for name, sort := range contentSortOrder {
		global.Mysql.Model(&model.Menu{}).Where("name = ?", name).Update("sort", sort)
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

	admin := model.Role{Name: "超级管理员", Code: "002", Desc: "", HomePage: "ReviewPGC"}
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

// authApiDesc 需要登录权限的API描述表
// key: "METHOD|PATH", value: Desc说明
// 新增需要鉴权的路由时，在此表中添加对应条目即可自动同步到数据库
var authApiDesc = map[string]string{
	// API管理
	"POST|/api/v1/api/addApi":          "新增API（后台管理）",
	"DELETE|/api/v1/api/deleteApi/:id": "删除API（后台管理）",
	"PUT|/api/v1/api/editApi":          "编辑API（后台管理）",
	"PUT|/api/v1/api/editRoleApi":      "编辑角色API（后台管理）",
	"GET|/api/v1/api/getAllApiList":    "获取全部API列表（后台管理）",
	"POST|/api/v1/api/getApiList":      "获取API列表（后台管理）",
	"GET|/api/v1/api/getRoleApi":       "获取角色API（后台管理）",
	// 点赞收藏
	"POST|/api/v1/archive/article/cancelCollect": "文章取消收藏",
	"POST|/api/v1/archive/article/cancelLike":    "文章取消赞",
	"POST|/api/v1/archive/article/collect":       "文章收藏",
	"GET|/api/v1/archive/article/hasCollect":     "文章是否收藏",
	"GET|/api/v1/archive/article/hasLike":        "文章是否点赞",
	"POST|/api/v1/archive/article/like":          "文章点赞",
	"POST|/api/v1/archive/article/share":         "文章分享计数",
	"POST|/api/v1/archive/video/cancelLike":      "视频取消赞",
	"POST|/api/v1/archive/video/collect":         "视频收藏",
	"GET|/api/v1/archive/video/getCollectInfo":   "获取视频已收藏的文件夹",
	"GET|/api/v1/archive/video/hasCollect":       "视频是否收藏",
	"GET|/api/v1/archive/video/hasLike":          "视频是否点赞",
	"POST|/api/v1/archive/video/like":            "视频点赞",
	"POST|/api/v1/archive/video/share":           "视频分享计数",
	// 文章
	"DELETE|/api/v1/article/deleteArticle/:id":       "删除文章",
	"DELETE|/api/v1/article/deleteArticleManage/:id": "删除文章（后台管理）",
	"PUT|/api/v1/article/editArticleInfo":            "编辑文章信息",
	"GET|/api/v1/article/getAllArticleList":          "获取所有的文章列表",
	"POST|/api/v1/article/getArticleListManage":      "获取文章列表（后台管理）",
	"GET|/api/v1/article/getArticleStatus":           "获取文章状态信息",
	"POST|/api/v1/article/getReviewArticleList":      "获取文章审核列表（后台管理）",
	"GET|/api/v1/article/getUploadArticle":           "获取上传的文章",
	"POST|/api/v1/article/uploadArticleInfo":         "上传文章信息",
	// Auth
	"POST|/api/v1/auth/logout": "退出登录",
	// 客户端
	"POST|/api/v1/client/log": "客户端日志上报",
	// 轮播图
	"POST|/api/v1/carousel/addCarousel":          "新增轮播图（后台管理）",
	"DELETE|/api/v1/carousel/deleteCarousel/:id": "删除轮播图（后台管理）",
	"PUT|/api/v1/carousel/editCarousel":          "编辑轮播图（后台管理）",
	"POST|/api/v1/carousel/getCarouselList":      "获取轮播图列表（后台管理）",
	// 收藏夹
	"POST|/api/v1/collection/addCollection":          "添加收藏夹",
	"DELETE|/api/v1/collection/deleteCollection/:id": "删除收藏夹",
	"PUT|/api/v1/collection/editCollection":          "编辑收藏夹",
	"GET|/api/v1/collection/getCollectionInfo":       "获取收藏夹信息",
	"GET|/api/v1/collection/getCollectionList":       "获取收藏夹列表",
	"GET|/api/v1/collection/getVideoList":            "获取收藏夹的视频列表",
	// 评论回复
	"POST|/api/v1/comment/article/addComment":          "发表文章评论或回复",
	"GET|/api/v1/comment/article/getCommentList":       "获取文章评论列表",
	"DELETE|/api/v1/comment/article/deleteComment/:id": "删除文章评论或回复",
	"POST|/api/v1/comment/article/like/:id":           "点赞文章评论",
	"DELETE|/api/v1/comment/article/like/:id":          "取消点赞文章评论",
	"POST|/api/v1/comment/article/dislike/:id":         "点踩文章评论",
	"DELETE|/api/v1/comment/article/dislike/:id":       "取消点踩文章评论",
	"POST|/api/v1/comment/video/addComment":            "发表视频评论或回复",
	"DELETE|/api/v1/comment/video/deleteComment/:id":   "删除视频评论或回复",
	"GET|/api/v1/comment/video/getCommentList":         "获取视频评论列表",
	"POST|/api/v1/comment/video/like/:id":            "点赞视频评论",
	"DELETE|/api/v1/comment/video/like/:id":           "取消点赞视频评论",
	"POST|/api/v1/comment/video/dislike/:id":          "点踩视频评论",
	"DELETE|/api/v1/comment/video/dislike/:id":        "取消点踩视频评论",
	// 弹幕
	"POST|/api/v1/danmaku/sendDanmaku": "发送弹幕",
	// 历史记录
	"POST|/api/v1/history/video/addHistory": "保存视频历史记录",
	"GET|/api/v1/history/video/getHistory":  "获取视频历史记录",
	"GET|/api/v1/history/video/getProgress": "获取视频播放进度",
	// 菜单管理
	"POST|/api/v1/menu/addMenu":          "添加菜单（后台管理）",
	"DELETE|/api/v1/menu/deleteMenu/:id": "删除菜单（后台管理）",
	"PUT|/api/v1/menu/editMenu":          "编辑菜单（后台管理）",
	"PUT|/api/v1/menu/editRoleMenu":      "编辑角色菜单（后台管理）",
	"GET|/api/v1/menu/getMenuTree":       "获取菜单树（后台管理）",
	"GET|/api/v1/menu/getRoleMenu":       "获取角色菜单（后台管理）",
	"GET|/api/v1/menu/getUserMenu":       "获取用户菜单树（后台管理）",
	// 消息
	"POST|/api/v1/message/addAnnounce":          "添加公告（后台管理）",
	"DELETE|/api/v1/message/deleteAnnounce/:id": "删除公告（后台管理）",
	"GET|/api/v1/message/getAtMsg":              "获取AT通知",
	"GET|/api/v1/message/getLikeMsg":            "获取点赞通知",
	"GET|/api/v1/message/getReplyMsg":           "获取回复通知",
	"GET|/api/v1/message/getWhisperDetails":     "获取私信详情",
	"GET|/api/v1/message/getWhisperList":        "获取私信列表",
	"POST|/api/v1/message/readWhisper":          "已读私信",
	"GET|/api/v1/message/readStatus":            "获取公告/点赞/回复/@已读进度",
	"POST|/api/v1/message/readStatus":           "上报公告/点赞/回复/@已读进度",
	"POST|/api/v1/message/sendWhisper":          "发送私信",
	// 分区
	"POST|/api/v1/partition/addPartition":          "添加分区（后台管理）",
	"DELETE|/api/v1/partition/deletePartition/:id": "删除分区（后台管理）",
	// 关注
	"POST|/api/v1/relation/follow":         "关注用户",
	"GET|/api/v1/relation/getUserRelation": "获取用户关系",
	"POST|/api/v1/relation/unfollow":       "取关用户",
	// 资源
	"DELETE|/api/v1/resource/deleteResource/:id": "删除视频资源",
	"PUT|/api/v1/resource/modifyTitle":           "修改资源标题",
	"PUT|/api/v1/resource/reorder":               "资源排序",
	"POST|/api/v1/resource/replaceResource":      "替换视频资源",
	// 审核
	"GET|/api/v1/review/getArticleReviewRecord": "获取文章审核记录",
	"GET|/api/v1/review/getVideoReviewRecord":   "获取视频审核记录",
	"POST|/api/v1/review/reviewArticleApproved": "文章审核通过（后台管理）",
	"POST|/api/v1/review/reviewArticleFailed":   "文章审核不通过（后台管理）",
	"POST|/api/v1/review/reviewVideoApproved":   "视频审核通过（后台管理）",
	"POST|/api/v1/review/reviewVideoFailed":     "视频审核不通过（后台管理）",
	// 角色
	"POST|/api/v1/role/addRole":          "添加角色（后台管理）",
	"DELETE|/api/v1/role/deleteRole/:id": "删除角色（后台管理）",
	"PUT|/api/v1/role/editRole":          "编辑角色（后台管理）",
	"PUT|/api/v1/role/editRoleHome":      "编辑角色首页（后台管理）",
	"GET|/api/v1/role/getAllRoleList":    "获取全部角色（后台管理）",
	"GET|/api/v1/role/getRoleInfo":       "获取个人角色信息（后台管理）",
	"POST|/api/v1/role/getRoleList":      "获取角色列表（后台管理）",
	// 上传
	"POST|/api/v1/upload/image":      "上传图片",
	"POST|/api/v1/upload/video":      "上传视频",
	"POST|/api/v1/upload/video/:vid": "上传视频分P",
	"POST|/api/v1/upload/checkVideo": "获取视频上传进度",
	"POST|/api/v1/upload/chunkVideo": "上传视频文件分片",
	"POST|/api/v1/upload/mergeVideo": "合并视频文件分片",
	// 用户
	"PUT|/api/v1/user/banUser":            "封禁用户（后台管理）",
	"PUT|/api/v1/user/unBanUser":          "解封用户（后台管理）",
	"GET|/api/v1/user/getUserBanRecord":   "获取封禁记录（后台管理）",
	"DELETE|/api/v1/user/deleteUser/:id":  "删除用户（后台管理）",
	"PUT|/api/v1/user/editUserInfo":       "编辑用户信息",
	"PUT|/api/v1/user/editUserInfoManage": "编辑用户信息（后台管理）",
	"PUT|/api/v1/user/editUserRole":       "编辑用户角色（后台管理）",
	"GET|/api/v1/user/getUserInfo":        "获取用户信息",
	"POST|/api/v1/user/getUserListManage": "获取用户列表（后台管理）",
	// 视频
	"DELETE|/api/v1/video/deleteVideo/:id":       "删除视频",
	"DELETE|/api/v1/video/deleteVideoManage/:id": "删除视频（后台管理）",
	"PUT|/api/v1/video/editVideoInfo":            "编辑视频信息",
	"GET|/api/v1/video/getAllVideoList":          "获取所有的视频列表",
	"POST|/api/v1/video/getReviewList":           "获取审核列表（后台管理）",
	"GET|/api/v1/video/getReviewResourceList":    "获取审核资源列表（后台管理）",
	"GET|/api/v1/video/getUploadVideo":           "获取上传的视频",
	"POST|/api/v1/video/getVideoListManage":      "获取视频列表（后台管理）",
	"POST|/api/v1/video/getFailedVideoList":      "获取处理失败的视频列表（后台管理）",
	"POST|/api/v1/video/getProcessingVideoList":  "获取处理中视频列表（后台管理）",
	"GET|/api/v1/video/getVideoStatus":           "获取上传视频状态信息",
	"POST|/api/v1/video/uploadVideoInfo":         "上传视频信息",
	"GET|/api/v1/video/getResourceQualityManage": "获取视频资源支持的分辨率信息（后台管理）",
	"GET|/api/v1/video/getVideoFileManage":       "获取视频文件URL（后台管理）",
	"POST|/api/v1/video/reTranscodeVideo":        "重新转码视频（后台管理）",
	"POST|/api/v1/video/subtitle/upload":         "上传分P字幕",
	"PUT|/api/v1/video/subtitle/:id":             "更新分P字幕",
	"DELETE|/api/v1/video/subtitle/:id":          "删除分P字幕",
	// 配置
	"GET|/api/v1/config/getEmailConfig":    "获取邮箱配置（后台管理）",
	"POST|/api/v1/config/setEmailConfig":   "编辑邮箱配置（后台管理）",
	"GET|/api/v1/config/getStorageConfig":  "获取存储配置（后台管理）",
	"POST|/api/v1/config/setStorageConfig": "编辑存储配置（后台管理）",
	"GET|/api/v1/config/getOtherConfig":    "获取其他配置（后台管理）",
	"POST|/api/v1/config/setOtherConfig":   "编辑其他配置（后台管理）",
	// 合集
	"POST|/api/v1/playlist/add":                        "创建合集",
	"PUT|/api/v1/playlist/edit":                        "编辑合集",
	"DELETE|/api/v1/playlist/del/:id":                  "删除合集",
	"GET|/api/v1/playlist/myList":                      "获取自己的合集列表",
	"POST|/api/v1/playlist/video/add":                  "添加视频到合集",
	"POST|/api/v1/playlist/video/del":                  "从合集移除视频",
	"POST|/api/v1/playlist/video/sort":                 "调整合集视频排序",
	"POST|/api/v1/playlist/getReviewPlaylistList":      "获取待审核合集列表（后台管理）",
	"POST|/api/v1/playlist/reviewApproved":             "合集审核通过（后台管理）",
	"POST|/api/v1/playlist/reviewFailed":               "合集审核不通过（后台管理）",
	"GET|/api/v1/playlist/getPlaylistReviewRecord":     "获取合集审核记录",
	"POST|/api/v1/playlist/getPlaylistListManage":      "获取全站合集列表（后台管理）",
	"DELETE|/api/v1/playlist/deletePlaylistManage/:id": "删除合集（后台管理）",
	// PGC
	"POST|/api/v1/pgc/create":                 "创建PGC内容",
	"PUT|/api/v1/pgc/update":                  "更新PGC内容",
	"DELETE|/api/v1/pgc/:pgc_id":              "删除PGC内容",
	"POST|/api/v1/pgc/:pgc_id/episodes/add":   "添加PGC剧集",
	"DELETE|/api/v1/pgc/:pgc_id/episodes/:id": "删除PGC剧集",
	"POST|/api/v1/pgc/getReviewList":          "获取PGC待审列表（后台管理）",
	"POST|/api/v1/pgc/reviewApproved":         "PGC审核通过（后台管理）",
	"POST|/api/v1/pgc/reviewFailed":           "PGC审核驳回（后台管理）",
	"POST|/api/v1/pgc/getManageList":          "获取PGC管理列表（后台管理）",
	"POST|/api/v1/pgc/adminUpdateStatus":      "管理员修改PGC状态（后台管理）",
	"DELETE|/api/v1/pgc/adminDelete/:pgc_id":  "管理员删除PGC（后台管理）",
}

// SyncApiData 自动同步需要登录权限的路由到API表
// 只同步 authApiDesc 中定义的鉴权路由，公开路由不写入
func SyncApiData() {
	// 1. 从数据库加载已有的 API 记录
	var existingApis []model.Api
	global.Mysql.Find(&existingApis)

	existSet := make(map[string]bool, len(existingApis))
	for _, api := range existingApis {
		existSet[api.Method+"|"+api.Path] = true
	}

	// 2. 遍历鉴权路由表，找出数据库中缺少的
	var newApis []model.Api
	for key, desc := range authApiDesc {
		if existSet[key] {
			continue
		}

		// 解析 METHOD 和 PATH
		parts := strings.SplitN(key, "|", 2)
		if len(parts) != 2 {
			continue
		}

		newApis = append(newApis, model.Api{
			Method:   parts[0],
			Path:     parts[1],
			Category: inferCategory(parts[1]),
			Desc:     desc,
		})
	}

	// 3. 批量写入 API 表
	if len(newApis) > 0 {
		if err := global.Mysql.Create(&newApis).Error; err != nil {
			zap.L().Error("自动同步API数据失败", zap.String("err", err.Error()), zap.String("module", "initialize"))
		} else {
			zap.L().Info("自动同步API数据完成",
				zap.Int("newCount", len(newApis)),
				zap.String("module", "initialize"))
			for _, api := range newApis {
				zap.L().Info("新增API",
					zap.String("method", api.Method),
					zap.String("path", api.Path),
					zap.String("category", api.Category),
					zap.String("desc", api.Desc))
			}
			// 4. 为角色 001、002 写入新 API 的 Casbin 权限，并重载策略
			var newRules []model.CasbinRule
			for _, api := range newApis {
				newRules = append(newRules,
					model.CasbinRule{Ptype: "p", V0: "001", V1: api.Path, V2: api.Method},
					model.CasbinRule{Ptype: "p", V0: "002", V1: api.Path, V2: api.Method},
				)
			}
			if err := global.Mysql.Create(&newRules).Error; err != nil {
				zap.L().Error("自动同步Casbin权限失败", zap.String("err", err.Error()), zap.String("module", "initialize"))
			} else if global.Casbin != nil {
				_ = global.Casbin.ReloadPolicy()
				zap.L().Info("已为角色001/002分配新API权限并重载Casbin", zap.String("module", "initialize"))
			}
		}
	} else {
		zap.L().Info("API数据已是最新，无需同步", zap.String("module", "initialize"))
	}
}

// inferCategory 根据路由路径推断API分类
// 例如 /api/v1/playlist/add → "合集"
func inferCategory(path string) string {
	// 去掉 /api/v1/ 前缀，取第一段作为分类依据
	trimmed := strings.TrimPrefix(path, "/api/v1/")
	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) == 0 {
		return "其他"
	}

	categoryMap := map[string]string{
		"api":        "API管理",
		"archive":    "点赞收藏",
		"article":    "文章",
		"auth":       "Auth",
		"carousel":   "轮播图",
		"collection": "收藏夹",
		"comment":    "评论回复",
		"config":     "配置",
		"danmaku":    "弹幕",
		"history":    "历史记录",
		"menu":       "菜单管理",
		"message":    "消息",
		"partition":  "分区",
		"pgc":        "PGC",
		"playlist":   "合集",
		"relation":   "关注",
		"resource":   "资源",
		"review":     "审核",
		"role":       "角色",
		"upload":     "上传",
		"user":       "用户",
		"verify":     "验证",
		"video":      "视频",
		"online":     "在线",
	}

	if category, ok := categoryMap[parts[0]]; ok {
		return category
	}
	return "其他"
}
