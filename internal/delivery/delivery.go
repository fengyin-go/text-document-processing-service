package delivery

type Transaction interface {
	Commit() error
	Rollback() error
}

type Publisher interface {
	Publish(string) error
}

type Audit interface {
	Record(string)
}

func Finalize(tx Transaction, operationErr error) error {
	if operationErr == nil {
		return tx.Commit()
	}
	return tx.Rollback()
}
