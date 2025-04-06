package controller

import "github.com/jackc/pgx"

// Controller example
type Controller struct {
	dbconn *pgx.Conn
}

// NewController example
func NewController(dbconn *pgx.Conn) *Controller {
	return &Controller{dbconn: dbconn}
}
