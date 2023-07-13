package intrface

type Idb interface {
	Begin() (Itx, error)
}
