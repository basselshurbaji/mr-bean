package main

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetMeParams struct{}

type UpdateMeParams struct {
	FirstName *string `json:"first_name,omitempty" jsonschema:"First name (at least one of first_name or last_name is required)"`
	LastName  *string `json:"last_name,omitempty" jsonschema:"Last name (at least one of first_name or last_name is required)"`
}

func registerUserTools(server *mcp.Server, client *Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_me",
		Description: "Get the current authenticated user's profile",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params GetMeParams) (*mcp.CallToolResult, any, error) {
		var result any
		if err := client.get(ctx, "/user/me", &result); err != nil {
			return nil, nil, err
		}
		text, err := toJSON(result)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_me",
		Description: "Update the current authenticated user's profile (first_name and/or last_name)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params UpdateMeParams) (*mcp.CallToolResult, any, error) {
		if params.FirstName == nil && params.LastName == nil {
			return nil, nil, fmt.Errorf("at least one of first_name or last_name must be provided")
		}
		var result any
		if err := client.patch(ctx, "/user/me", params, &result); err != nil {
			return nil, nil, err
		}
		text, err := toJSON(result)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}, nil, nil
	})
}
