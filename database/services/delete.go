package services

func delete(hash string) (err error) {
	db := open()
	defer close()
	err = db.Delete([]byte(hash), nil)
	return
}
