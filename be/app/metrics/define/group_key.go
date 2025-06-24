package define

type GroupKey struct {
	ID          uint64 `json:"id"`
	DisplayName string `json:"name"`
	ParentID    uint64 `json:"parent_id"`
	Order       uint16 `json:"order_index"`
	Remark      string `json:"remark"`
	Mtime       uint32 `json:"mtime"`
	Ctime       uint32 `json:"ctime"`
}

type DataGroupKey struct {
	ID          uint64 `json:"id"`
	DisplayName string `json:"name"`
	Remark      string `json:"remark"`
	Mtime       uint32 `json:"mtime"`
	Ctime       uint32 `json:"ctime"`
}

type GroupKeyTree struct {
	ID          uint64          `json:"id"`
	DisplayName string          `json:"name"`
	Order       uint16          `json:"order_index"`
	Children    []*GroupKeyTree `json:"children"`
}

type GroupKeyTreeReq struct {
	RootID uint64 `json:"root_id" form:"root_id" query:"root_id"`
}
