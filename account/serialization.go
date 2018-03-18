package account

// Import a file and return a Account object
func Import(filePath string, name string) (account *Account, err error) {
	if name == "" {
		name = "Test A"
	}
	// TODO add import
	account = &Account{
		Address: "0x0000000000000000000000000000000000000000",
		Name:    name,
	}
	return
}

// Export an account into a file and then return the path of the file
func (account *Account) Export() (path string, err error) {
	// TODO add export
	path = "/home/antho/..."
	return
}
