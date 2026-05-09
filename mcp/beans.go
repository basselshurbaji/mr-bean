package main

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ListBeansParams struct{}

type CreateBeanParams struct {
	Name         string  `json:"name" jsonschema:"Bean name"`
	Roaster      *string `json:"roaster,omitempty" jsonschema:"Roaster name"`
	Origin       *string `json:"origin,omitempty" jsonschema:"Bean origin country or region"`
	Process      *string `json:"process,omitempty" jsonschema:"Processing method: washed, natural, honey, anaerobic, or other"`
	RoastLevel   *string `json:"roast_level,omitempty" jsonschema:"Roast level: light, medium_light, medium, medium_dark, or dark"`
	TastingNotes *string `json:"tasting_notes,omitempty" jsonschema:"Tasting notes"`
	Notes        *string `json:"notes,omitempty" jsonschema:"Additional notes"`
}

type UpdateBeanParams struct {
	ID           string  `json:"id" jsonschema:"Bean ID to update"`
	Name         string  `json:"name" jsonschema:"Bean name"`
	Roaster      *string `json:"roaster,omitempty" jsonschema:"Roaster name"`
	Origin       *string `json:"origin,omitempty" jsonschema:"Bean origin country or region"`
	Process      *string `json:"process,omitempty" jsonschema:"Processing method: washed, natural, honey, anaerobic, or other"`
	RoastLevel   *string `json:"roast_level,omitempty" jsonschema:"Roast level: light, medium_light, medium, medium_dark, or dark"`
	TastingNotes *string `json:"tasting_notes,omitempty" jsonschema:"Tasting notes"`
	Notes        *string `json:"notes,omitempty" jsonschema:"Additional notes"`
}

type DeleteBeanParams struct {
	ID string `json:"id" jsonschema:"Bean ID to delete"`
}

func registerBeanTools(server *mcp.Server, client *Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_beans",
		Description: "List all coffee beans for the authenticated user",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params ListBeansParams) (*mcp.CallToolResult, any, error) {
		var result any
		if err := client.get(ctx, "/beans", &result); err != nil {
			return nil, nil, err
		}
		text, err := toJSON(result)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_bean",
		Description: "Create a new coffee bean entry",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params CreateBeanParams) (*mcp.CallToolResult, any, error) {
		var result any
		if err := client.post(ctx, "/beans", params, &result); err != nil {
			return nil, nil, err
		}
		text, err := toJSON(result)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_bean",
		Description: "Update an existing coffee bean by ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params UpdateBeanParams) (*mcp.CallToolResult, any, error) {
		body := struct {
			Name         string  `json:"name"`
			Roaster      *string `json:"roaster,omitempty"`
			Origin       *string `json:"origin,omitempty"`
			Process      *string `json:"process,omitempty"`
			RoastLevel   *string `json:"roast_level,omitempty"`
			TastingNotes *string `json:"tasting_notes,omitempty"`
			Notes        *string `json:"notes,omitempty"`
		}{
			Name:         params.Name,
			Roaster:      params.Roaster,
			Origin:       params.Origin,
			Process:      params.Process,
			RoastLevel:   params.RoastLevel,
			TastingNotes: params.TastingNotes,
			Notes:        params.Notes,
		}
		var result any
		if err := client.put(ctx, fmt.Sprintf("/beans/%s", params.ID), body, &result); err != nil {
			return nil, nil, err
		}
		text, err := toJSON(result)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_bean",
		Description: "Delete a coffee bean by ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params DeleteBeanParams) (*mcp.CallToolResult, any, error) {
		if err := client.delete(ctx, fmt.Sprintf("/beans/%s", params.ID)); err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "Bean deleted successfully"}}}, nil, nil
	})
}
