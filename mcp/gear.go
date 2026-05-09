package main

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ListGearParams struct{}

type CreateGearParams struct {
	TypeID string  `json:"type_id" jsonschema:"Gear type: machine, grinder, scale, portafilter, tamper, distributor, wdt, basket, puckscreen, or other"`
	Name   string  `json:"name" jsonschema:"Gear name"`
	Brand  *string `json:"brand,omitempty" jsonschema:"Brand name"`
	Model  *string `json:"model,omitempty" jsonschema:"Model name"`
	Year   *string `json:"year,omitempty" jsonschema:"Year of manufacture (4-digit)"`
	Notes  *string `json:"notes,omitempty" jsonschema:"Additional notes"`
}

type GetGearParams struct {
	ID string `json:"id" jsonschema:"Gear ID to retrieve"`
}

type UpdateGearParams struct {
	ID     string  `json:"id" jsonschema:"Gear ID to update"`
	TypeID string  `json:"type_id" jsonschema:"Gear type: machine, grinder, scale, portafilter, tamper, distributor, wdt, basket, puckscreen, or other"`
	Name   string  `json:"name" jsonschema:"Gear name"`
	Brand  *string `json:"brand,omitempty" jsonschema:"Brand name"`
	Model  *string `json:"model,omitempty" jsonschema:"Model name"`
	Year   *string `json:"year,omitempty" jsonschema:"Year of manufacture (4-digit)"`
	Notes  *string `json:"notes,omitempty" jsonschema:"Additional notes"`
}

type DeleteGearParams struct {
	ID string `json:"id" jsonschema:"Gear ID to delete"`
}

type ListStationsParams struct{}

type CreateStationParams struct {
	Name    string   `json:"name" jsonschema:"Station name"`
	GearIDs []string `json:"gear_ids,omitempty" jsonschema:"List of gear IDs to associate with this station"`
}

type UpdateStationParams struct {
	ID      string   `json:"id" jsonschema:"Station ID to update"`
	Name    string   `json:"name" jsonschema:"Station name"`
	GearIDs []string `json:"gear_ids,omitempty" jsonschema:"List of gear IDs to associate with this station"`
}

type DeleteStationParams struct {
	ID string `json:"id" jsonschema:"Station ID to delete"`
}

func registerGearTools(server *mcp.Server, client *Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_gear",
		Description: "List all gear items for the authenticated user",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params ListGearParams) (*mcp.CallToolResult, any, error) {
		var result any
		if err := client.get(ctx, "/gear", &result); err != nil {
			return nil, nil, err
		}
		text, err := toJSON(result)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_gear",
		Description: "Create a new gear item (machine, grinder, scale, etc.)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params CreateGearParams) (*mcp.CallToolResult, any, error) {
		var result any
		if err := client.post(ctx, "/gear", params, &result); err != nil {
			return nil, nil, err
		}
		text, err := toJSON(result)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_gear",
		Description: "Get a specific gear item by ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params GetGearParams) (*mcp.CallToolResult, any, error) {
		var result any
		if err := client.get(ctx, fmt.Sprintf("/gear/%s", params.ID), &result); err != nil {
			return nil, nil, err
		}
		text, err := toJSON(result)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_gear",
		Description: "Update an existing gear item by ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params UpdateGearParams) (*mcp.CallToolResult, any, error) {
		body := struct {
			TypeID string  `json:"type_id"`
			Name   string  `json:"name"`
			Brand  *string `json:"brand,omitempty"`
			Model  *string `json:"model,omitempty"`
			Year   *string `json:"year,omitempty"`
			Notes  *string `json:"notes,omitempty"`
		}{
			TypeID: params.TypeID,
			Name:   params.Name,
			Brand:  params.Brand,
			Model:  params.Model,
			Year:   params.Year,
			Notes:  params.Notes,
		}
		var result any
		if err := client.put(ctx, fmt.Sprintf("/gear/%s", params.ID), body, &result); err != nil {
			return nil, nil, err
		}
		text, err := toJSON(result)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_gear",
		Description: "Delete a gear item by ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params DeleteGearParams) (*mcp.CallToolResult, any, error) {
		if err := client.delete(ctx, fmt.Sprintf("/gear/%s", params.ID)); err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "Gear deleted successfully"}}}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_stations",
		Description: "List all espresso stations with their associated gear",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params ListStationsParams) (*mcp.CallToolResult, any, error) {
		var result any
		if err := client.get(ctx, "/stations", &result); err != nil {
			return nil, nil, err
		}
		text, err := toJSON(result)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_station",
		Description: "Create a new espresso station grouping gear items together",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params CreateStationParams) (*mcp.CallToolResult, any, error) {
		var result any
		if err := client.post(ctx, "/stations", params, &result); err != nil {
			return nil, nil, err
		}
		text, err := toJSON(result)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_station",
		Description: "Update an existing station by ID. This is a full replacement — gear_ids must list all gear to keep; omitting it clears all gear associations.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params UpdateStationParams) (*mcp.CallToolResult, any, error) {
		gearIDs := params.GearIDs
		if gearIDs == nil {
			gearIDs = []string{}
		}
		body := struct {
			Name    string   `json:"name"`
			GearIDs []string `json:"gear_ids"`
		}{
			Name:    params.Name,
			GearIDs: gearIDs,
		}
		var result any
		if err := client.put(ctx, fmt.Sprintf("/stations/%s", params.ID), body, &result); err != nil {
			return nil, nil, err
		}
		text, err := toJSON(result)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_station",
		Description: "Delete a station by ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params DeleteStationParams) (*mcp.CallToolResult, any, error) {
		if err := client.delete(ctx, fmt.Sprintf("/stations/%s", params.ID)); err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "Station deleted successfully"}}}, nil, nil
	})
}
