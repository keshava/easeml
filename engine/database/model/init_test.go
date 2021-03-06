package model

import (
	"github.com/ds3lab/easeml/engine/utils"
	"testing"

	"github.com/ds3lab/easeml/engine/database"
	"github.com/ds3lab/easeml/engine/database/model/types"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
)

func TestClear(t *testing.T) {
	assert := assert.New(t)
	var MongoInstance = utils.GetEnvVariableOrDefault("EASEML_DATABASE_ADDRESS","localhost")
	connection, err := database.Connect(MongoInstance, "testdb")
	assert.Nil(err)

	// Create a temp database.
	var collectionInfo = mgo.CollectionInfo{}
	err = connection.Session.DB("testdb").C("testcol").Create(&collectionInfo)
	assert.Nil(err)

	// Run the Clear function.
	var context = Context{Session: connection.Session, DBName: connection.DBName}
	err = context.Clear("testdb")
	assert.Nil(err)

	// Verify that the database has been dropped.
	names, err := connection.Session.DatabaseNames()
	assert.Nil(err)
	assert.NotContains(names, "testdb")
}

func TestInitialize(t *testing.T) {
	assert := assert.New(t)
	var MongoInstance = utils.GetEnvVariableOrDefault("EASEML_DATABASE_ADDRESS","localhost")
	connection, err := database.Connect(MongoInstance, "testdb")
	assert.Nil(err)

	// Run the Initialize function.
	var context = Context{Session: connection.Session, DBName: connection.DBName}
	err = context.Initialize("testdb")
	assert.Nil(err)

	// Verify that all the collections have been created.
	names, err := connection.Session.DB("testdb").CollectionNames()
	assert.Nil(err)
	assert.ElementsMatch(names, []string{"users", "processes", "datasets", "modules", "jobs", "tasks"})

	// Verify that the root user has been created.
	n, err := connection.Session.DB("testdb").C("users").Find(bson.M{"id": types.UserRoot}).Count()
	assert.Nil(err)
	assert.Equal(1, n)

	// Drop the test database.
	err = connection.Session.DB("testdb").DropDatabase()
	assert.Nil(err)
}
