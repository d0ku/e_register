package core

import "testing"

func TestDBConnection(t *testing.T) {
	Initialize("postgres", "test_database", "../page/")
}
