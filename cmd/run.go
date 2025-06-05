package cmd

import (
	"errors"
	"github.com/richardwooding/duckdb-mcp/model"
)

var NoDatasourceForInMemory = errors.New("cannot specify a datasource when using in-memory database")

type RunCmd struct {
	Transport  string `name:"transport" default:"stdio" enum:"stdio,http-with-sse" help:"Transport to use for the MCP server."`
	InMemory   bool   `name:"in-memory" default:"false" help:"Use in-memory database."`
	Datasource string `name:"datasource" default:"" help:"Data source for the database. Required if not using in-memory."`
}

func (c *RunCmd) Run(g *model.Globals) error {
	transport, err := model.ParseTransport(c.Transport)
	if err != nil {
		return err
	}

	if !c.InMemory && c.Datasource != "" {
		return NoDatasourceForInMemory
	}

	return nil
}
