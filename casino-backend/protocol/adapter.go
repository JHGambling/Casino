package protocol

type CasinoAdapter interface {
	Table(id string) (Table, error)
}
