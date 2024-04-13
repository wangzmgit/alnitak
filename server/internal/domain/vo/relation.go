package vo

type RelationResp struct {
	Uid       uint         `json:"-"`
	TargetUid uint         `json:"-"`
	Relation  int          `json:"relation"`
	User      UserInfoResp `json:"user" gorm:"-"`
}
