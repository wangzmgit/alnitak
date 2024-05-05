package dto

type AddPartitionReq struct {
	Name     string
	ParentId uint //所属分区ID
	Type     int
}
