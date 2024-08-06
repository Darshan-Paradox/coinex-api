package coinbase

type InvalidStructure struct{}

func (invStr *InvalidStructure) Error() string {
	return "Invalid Structure of received data"
}
