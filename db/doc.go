// db holds all the database interactions for the API server.
//
// The goal of this package is to make things as database agnostic as possible in the case that
// switching databases or database drivers should is needed. In those cases, it should be as simple
// as creating a new implementation of db.Database and updating daemon/daemon.go to use the new implementation.
package db
