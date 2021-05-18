package Helper

type ProdsRequest struct {
	//size参数，指定生成n个产品的列表
	Size int `form:"size"`
}