package crypto

type InvalidStructure struct{}

func (invStr *InvalidStructure) Error() string {
	return "Invalid Structure of received data"
}
