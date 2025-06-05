package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/richardwooding/duckdb-mcp/duckdbmanager"
	"github.com/richardwooding/duckdb-mcp/model"
)

var ErrExecutorNotSet = errors.New("executor not set")
var ErrQuerierNotSet = errors.New("querier not set")

type Config struct {
	Transport model.Transport
	Executor  duckdbmanager.DbExecutor
	Querier   duckdbmanager.DbQuerier
}

type Server struct {
	transport model.Transport
	executor  duckdbmanager.DbExecutor
	querier   duckdbmanager.DbQuerier
}

func NewServer(config Config) (*Server, error) {
	if config.Transport == model.UndefinedTransport {
		return nil, model.ErrInvalidTransport
	}
	if config.Executor == nil {
		return nil, ErrExecutorNotSet
	}
	if config.Querier == nil {
		return nil, ErrQuerierNotSet
	}
	return &Server{
		transport: config.Transport,
		executor:  config.Executor,
		querier:   config.Querier,
	}, nil
}

func (s *Server) Run() error {
	// Implementation of the server run logic goes here.
	// This could involve setting up the transport layer,
	// handling incoming requests, and using the executor
	// and querier to process database operations.

	srv := server.NewMCPServer(
		"DuckDB MCP Server. Provides access to a DuckDB database. Uses DuckDB SQL dialect.",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	queryTool := mcp.NewTool("query",
		mcp.WithDescription("Executes a SQL query against the DuckDB database. Uses DuckDB SQL Dialect"),
		mcp.WithString("sql",
			mcp.Required(),
			mcp.Description("The SQL query to execute against the DuckDB database. Use DuckDB SQL dialect.")),
	)

	srv.AddTool(queryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sql, err := request.RequireString("sql")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		rows, err := s.querier.Query(ctx, sql)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		defer rows.Close()
		rawJSON, err := rowsToJSON(rows)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		data, err := json.Marshal(rawJSON)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(string(data)), nil
	})

	execTool := mcp.NewTool("exec",
		mcp.WithDescription("Executes a SQL command against the DuckDB database. Uses DuckDB SQL Dialect"),
		mcp.WithString("sql"))

	srv.AddTool(execTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sql, err := request.RequireString("sql")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := s.executor.Exec(ctx, sql)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		execResult, err := resultToJSON(result)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		data, err := json.Marshal(execResult)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(string(data)), nil
	})

	return nil
}

type execResult struct {
	RowsAffected int64 `json:"rows_affected"`
	LastInsertID int64 `json:"last_insert_id"`
}

func resultToJSON(result sql.Result) (*execResult, error) {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &execResult{
		RowsAffected: rowsAffected,
		LastInsertID: lastInsertID,
	}, nil
}

func rowsToJSON(rows *sql.Rows) ([]map[string]any, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]any, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	var results []map[string]any
	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return nil, err
		}
		rowMap := make(map[string]any)
		for i, colName := range columns {
			rowMap[colName] = *(values[i].(*any))
		}
		results = append(results, rowMap)
	}

	return results, nil
}
