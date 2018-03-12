package account

// Import a file and return a Account object
func Import(filePath string, name string) (*Account, error) {
	if name == "" {
		name = "Test A"
	}
	// TODO add import
	return &Account{
		Address: "0x0000000000000000000000000000000000000000",
		Name:    name,
	}, nil
}

// Export an account into a file and then return the path of the file
func (account *Account) Export() (string, error) {
	// TODO add export
	return "/home/antho/...", nil
}
