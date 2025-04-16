package blockchain

type TxInput struct {
	ID  []byte
	Out int
	Sig string
}

type TxOutput struct {
	Value  int
	PubKey string
}

func (in *TxInput) CanUnLock(data string) bool {
	return in.Sig == data
}

func (out *TxOutput) CanBeUnLocked(data string) bool {
	return out.PubKey == data
}

type UTXO struct {
	TxID   string
	OutIdx int
	OutPut TxOutput
}
