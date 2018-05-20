package services

// Delete a service based on its hash
func Delete(hash string) (err error) {
	db := open()
	defer close()
	err = db.Delete([]byte(hash), nil)
	return
}
