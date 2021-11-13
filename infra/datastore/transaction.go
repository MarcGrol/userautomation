package datastore

type transaction struct {
	operations []datastoreOperation
}

type datastoreOperationKind int

const (
	datastoreOperationUpsert datastoreOperationKind = iota
	datastoreOperationRemove
)

type datastoreOperation struct {
	Kind datastoreOperationKind
	UID  string
	data interface{}
}

func (trx *transaction) put(uid string, item interface{}) error {
	trx.operations = append(trx.operations, datastoreOperation{
		Kind: datastoreOperationUpsert,
		UID:  uid,
		data: item,
	})
	return nil
}

func (trx *transaction) remove(uid string) error {
	trx.operations = append(trx.operations, datastoreOperation{
		Kind: datastoreOperationRemove,
		UID:  uid,
		data: nil,
	})
	return nil
}
