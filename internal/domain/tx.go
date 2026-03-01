package domain

type IsolationLevel int

const (
	IsolationDefault IsolationLevel = iota
	IsolationReadUncommitted
	IsolationReadCommitted
	IsolationRepeatableRead
	IsolationSerializable
)

type TxOptions struct {
	IsolationLevel IsolationLevel
	ReadOnly       bool
}
