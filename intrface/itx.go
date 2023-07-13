package intrface

type Itx interface {
	Rollback() error
	Commit() error
}
