package dto

type AddCarouselReq struct {
	Img         string
	Url         string
	Title       string
	Color       string
	Use         bool
	PartitionId uint
}

type EditCarouselReq struct {
	ID          uint
	Img         string
	Url         string
	Title       string
	Color       string
	Use         bool
	PartitionId uint
}

type CarouselListReq struct {
	Page     int
	PageSize int
}
