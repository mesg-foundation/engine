package services

// Delete deletes a service from the database based on the hash.
func Delete(hash string) error {
	db, err := open()
	defer close()
	if err != nil {
		return err
	}
	return db.Delete([]byte(hash), nil)
}
