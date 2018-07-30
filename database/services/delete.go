package services

// Delete a service based on its hash
func Delete(hash string) error {
	db, err := open()
	defer close()
	if err != nil {
		return err
	}
	return db.Delete([]byte(hash), nil)
}
