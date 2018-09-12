package services

// Delete deletes a service from the database based on the id.
func Delete(id string) error {
	db, err := open()
	defer close()
	if err != nil {
		return err
	}
	return db.Delete([]byte(id), nil)
}
