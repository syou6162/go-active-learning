package db_test

import (
	"testing"

	"github.com/syou6162/go-active-learning/lib/db"
	"github.com/syou6162/go-active-learning/lib/example"
)

func TestCreateDBConnection(t *testing.T) {
	conn, err := db.CreateDBConnection()
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()
}

func TestInsertExampleFromScanner(t *testing.T) {
	conn, err := db.CreateDBConnection()
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	_, err = db.DeleteAllExamples(conn)
	if err != nil {
		t.Error(err)
	}

	_, err = db.InsertOrUpdateExample(conn, &example.Example{Url: "http://hoge.com", Label: example.NEGATIVE})
	if err != nil {
		t.Error(err)
	}

	examples, err := db.ReadExamples(conn)
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 1 {
		t.Errorf("len(examples) == %d, want 1", len(examples))
	}

	// same url
	_, err = db.InsertOrUpdateExample(conn, &example.Example{Url: "http://hoge.com", Label: example.NEGATIVE})
	if err != nil {
		t.Error(err)
	}

	examples, err = db.ReadExamples(conn)
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 1 {
		t.Errorf("len(examples) == %d, want 1", len(examples))
	}

	// different url
	_, err = db.InsertOrUpdateExample(conn, &example.Example{Url: "http://another.com", Label: example.NEGATIVE})
	if err != nil {
		t.Error(err)
	}

	examples, err = db.ReadExamples(conn)
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 2 {
		t.Errorf("len(examples) == %d, want 2", len(examples))
	}
}
