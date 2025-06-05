package main

import (
	"github.com/alecthomas/kong"
	"github.com/richardwooding/duckdb-mcp/cmd"
	"github.com/richardwooding/duckdb-mcp/model"
)

type CLI struct {
	model.Globals
	Run cmd.RunCmd `cmd:"" help:"Run MCP Server"`
}

func main() {
	cli := CLI{
		Globals: model.Globals{
			Version: model.VersionFlag("0.1.0"),
		},
	}
	ctx := kong.Parse(&cli,
		kong.Name("duckdb-mcp"),
		kong.Description("A MCP server for a DuckDB database. Tools understand DuckDB SQL Dialect."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": "0.1.0",
		})
	err := ctx.Run(&cli.Globals)
	ctx.FatalIfErrorf(err)
}
