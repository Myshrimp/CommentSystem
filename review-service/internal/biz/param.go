package biz

type AuditParam struct {
	ReviewID int64
	Status   int32
	Reason   string
}

type AuditAppealParam struct {
	ReviewID int64
	Status   int32
	Reason   string
}

type AppealParam struct {
	ReviewID int64
	Content  string
}

type ReplyParam struct {
	ReviewID int64
	StoreID   int64
	Content  string
}