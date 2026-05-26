package tools

type Approver interface {
	Approve(msg string) bool
}

type ToolChain struct {
	approver Approver
}

func NewToolChain(approver Approver) *ToolChain {
	return &ToolChain{
		approver: approver,
	}
}
