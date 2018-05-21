package services

// Delete a service based on its hash
func Delete(hash string) (err error) {
	db, err := open()
	defer close()
	if err != nil {
		return
	}
	err = db.Delete([]byte(hash), nil)
	return
}
