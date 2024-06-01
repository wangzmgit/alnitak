package vo

import (
	"time"
)

const (
	ARTICLE_FIELD        = "`id`,`uid`,`title`,`cover`,`tags`,`content`,`copyright`,`status`,`created_at`,`clicks`"
	ARTICLE_STATUS_FIELD = "`id`,`title`,`cover`,`tags`,`content`,`copyright`,`partition_id`,`status`,`created_at`,`clicks`"
	UPLOAD_ARTICLE_FIELD = "`id`,`title`,`cover`,`tags`,`content_desc` as `content`,`copyright`,`status`,`created_at`,`clicks`"
	ARTICLE_LIST_FIELD   = "`id`,`uid`,`title`,`cover`,`tags`,`content_desc` as `content`,`copyright`,`status`,`created_at`,`clicks`"
)

type ArticleStatusResp struct {
	ID          uint   `json:"aid"`
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	Content     string `json:"content"`
	Copyright   bool   `json:"copyright"`
	Status      int    `json:"status"`
	PartitionId uint   `json:"partitionId"`
	Tags        string `json:"tags"`
}

type UploadArticleResp struct {
	ID        uint      `json:"aid"`
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	Content   string    `json:"content"`
	Status    int       `json:"status"`
	Copyright bool      `json:"copyright"`
	CreatedAt time.Time `json:"createdAt"`
	Tags      string    `json:"tags"`
	Clicks    int64     `json:"clicks"`
}

type ArticleResp struct {
	ID          uint         `json:"aid"`
	Uid         uint         `json:"uid"`
	Title       string       `json:"title"`
	Cover       string       `json:"cover"`
	Content     string       `json:"content"`
	Status      int          `json:"status"`
	Copyright   bool         `json:"copyright"`
	CreatedAt   time.Time    `json:"createdAt"`
	Tags        string       `json:"tags"`
	Clicks      int64        `json:"clicks"`
	PartitionId uint         `json:"partitionId"`
	Author      UserInfoResp `json:"author" gorm:"-"`
}

type ReviewArticleListResp struct {
	ID          uint         `json:"aid"`
	Uid         uint         `json:"uid"`
	Title       string       `json:"title"`
	Cover       string       `json:"cover"`
	Content     string       `json:"content"`
	CreatedAt   time.Time    `json:"createdAt"`
	Copyright   bool         `json:"copyright"`
	Tags        string       `json:"tags"`
	PartitionId uint         `json:"partitionId"`
	Author      UserInfoResp `json:"author" gorm:"-"`
}

type AllArticleResp struct {
	ID    uint   `json:"aid"`
	Title string `json:"title"`
}
