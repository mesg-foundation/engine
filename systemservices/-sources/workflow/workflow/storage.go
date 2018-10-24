package workflow

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Storage is a database storage for keeping workflows.
type Storage interface {
	// Save saves a workflow to database.
	// Save ensures that ID and Name fields individually unique.
	// Otherwise it should return an error.
	Save(w *WorkflowDocument) error

	// Delete deletes workflow from database by unique id/name of workflow document.
	Delete(id string) error

	// All returns all the workflows from database.
	All() (defs []*WorkflowDocument, err error)
}

// MongoStorage implements Storage.
type MongoStorage struct {
	// workflowsC is workflows collection.
	workflowsC string

	//dbName is database name.
	dbName string

	// uniqueKeys are individually set unique key indexes.
	uniqueKeys []string

	ss *mgo.Session
}

// NewMongoStorage returns a new MongoDB storage with given mongodb address and database.
func NewMongoStorage(addr, dbName string) (Storage, error) {
	s := &MongoStorage{
		workflowsC: "workflow",
		dbName:     dbName,
		uniqueKeys: []string{"id", "name"},
	}

	// dial to mongo.
	var err error
	s.ss, err = mgo.Dial(addr)
	if err != nil {
		return nil, err
	}

	// set consistency mode.
	s.ss.SetMode(mgo.Strong, true)

	// create unique key indexes.
	ss := s.ss.Clone()
	defer ss.Close()
	for _, key := range s.uniqueKeys {
		index := mgo.Index{
			Key:    []string{key},
			Unique: true,
		}
		if err := ss.DB(s.dbName).C(s.workflowsC).EnsureIndex(index); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// Save saves a workflow to database.
func (s *MongoStorage) Save(w *WorkflowDocument) error {
	ss := s.ss.Clone()
	defer ss.Close()
	return ss.DB(s.dbName).C(s.workflowsC).Insert(w)
}

// Delete deletes workflow from database by unique id/name of workflow document.
func (s *MongoStorage) Delete(id string) error {
	ss := s.ss.Clone()
	defer ss.Close()
	return ss.DB(s.dbName).C(s.workflowsC).Remove(bson.M{"$or": []bson.M{
		bson.M{"id": id},
		bson.M{"name": id},
	}})
}

// All returns all the workflows from database.
func (s *MongoStorage) All() ([]*WorkflowDocument, error) {
	ss := s.ss.Clone()
	defer ss.Close()
	var workflows []*WorkflowDocument
	return workflows, ss.DB(s.dbName).C(s.workflowsC).Find(bson.M{}).All(&workflows)
}
