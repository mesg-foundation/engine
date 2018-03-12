package account

// List all available accounts on this computer
func List() []*Account {
	// TODO add real list
	return []*Account{
		&Account{Name: "Test1", Address: "0x0000000000000000000000000000000000000000"},
		&Account{Name: "Test2", Address: "0x0000000000000000000000000000000000000001"},
	}
}
