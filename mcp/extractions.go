package main

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ListExtractionsParams struct {
	Limit *int `json:"limit,omitempty" jsonschema:"Number of extractions per page (default: 20)"`
	Page  *int `json:"page,omitempty" jsonschema:"Page number (default: 1)"`
}

type CreateExtractionParams struct {
	BeanID       string   `json:"bean_id" jsonschema:"Bean ID to use for this extraction"`
	DoseIn       float64  `json:"dose_in" jsonschema:"Dose in grams (must be > 0)"`
	YieldOut     float64  `json:"yield_out" jsonschema:"Yield out in grams (must be > 0)"`
	BrewTime     float64  `json:"time" jsonschema:"Extraction time in seconds (must be > 0)"`
	TargetTime   float64  `json:"target_time" jsonschema:"Target extraction time in seconds (must be > 0)"`
	GrindSize    float64  `json:"grind_size" jsonschema:"Grind size setting (must be > 0)"`
	GearIDs      []string `json:"gear_ids,omitempty" jsonschema:"List of gear IDs used in this extraction"`
	PreInfusion  *bool    `json:"pre_infusion,omitempty" jsonschema:"Whether pre-infusion was used"`
	TastingNote  *string  `json:"tasting_note,omitempty" jsonschema:"Tasting note for this extraction"`
}

type GetExtractionParams struct {
	ID string `json:"id" jsonschema:"Extraction ID to retrieve"`
}

type UpdateExtractionParams struct {
	ID           string   `json:"id" jsonschema:"Extraction ID to update"`
	BeanID       string   `json:"bean_id" jsonschema:"Bean ID to use for this extraction"`
	DoseIn       float64  `json:"dose_in" jsonschema:"Dose in grams (must be > 0)"`
	YieldOut     float64  `json:"yield_out" jsonschema:"Yield out in grams (must be > 0)"`
	BrewTime     float64  `json:"time" jsonschema:"Extraction time in seconds (must be > 0)"`
	TargetTime   float64  `json:"target_time" jsonschema:"Target extraction time in seconds (must be > 0)"`
	GrindSize    float64  `json:"grind_size" jsonschema:"Grind size setting (must be > 0)"`
	GearIDs      []string `json:"gear_ids,omitempty" jsonschema:"List of gear IDs used in this extraction"`
	PreInfusion  *bool    `json:"pre_infusion,omitempty" jsonschema:"Whether pre-infusion was used"`
	TastingNote  *string  `json:"tasting_note,omitempty" jsonschema:"Tasting note for this extraction"`
}

type DeleteExtractionParams struct {
	ID string `json:"id" jsonschema:"Extraction ID to delete"`
}

func registerExtractionTools(server *mcp.Server, client *Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_extractions",
		Description: "List extractions for the authenticated user with optional pagination",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params ListExtractionsParams) (*mcp.CallToolResult, any, error) {
		path := "/extractions"
		sep := "?"
		if params.Limit != nil {
			path += fmt.Sprintf("%slimit=%d", sep, *params.Limit)
			sep = "&"
		}
		if params.Page != nil {
			path += fmt.Sprintf("%spage=%d", sep, *params.Page)
		}

		var result any
		if err := client.get(ctx, path, &result); err != nil {
			return nil, nil, err
		}
		text, err := toJSON(result)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_extraction",
		Description: "Record a new espresso extraction",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params CreateExtractionParams) (*mcp.CallToolResult, any, error) {
		body := map[string]any{
			"bean_id":     params.BeanID,
			"dose_in":     params.DoseIn,
			"yield_out":   params.YieldOut,
			"time":        params.BrewTime,
			"target_time": params.TargetTime,
			"grind_size":  params.GrindSize,
		}
		if params.GearIDs != nil {
			body["gear_ids"] = params.GearIDs
		}
		if params.PreInfusion != nil {
			body["pre_infusion"] = *params.PreInfusion
		}
		if params.TastingNote != nil {
			body["tasting_note"] = *params.TastingNote
		}

		var result any
		if err := client.post(ctx, "/extractions", body, &result); err != nil {
			return nil, nil, err
		}
		text, err := toJSON(result)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_extraction",
		Description: "Get a specific extraction by ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params GetExtractionParams) (*mcp.CallToolResult, any, error) {
		var result any
		if err := client.get(ctx, fmt.Sprintf("/extractions/%s", params.ID), &result); err != nil {
			return nil, nil, err
		}
		text, err := toJSON(result)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_extraction",
		Description: "Update an existing extraction by ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params UpdateExtractionParams) (*mcp.CallToolResult, any, error) {
		body := map[string]any{
			"bean_id":     params.BeanID,
			"dose_in":     params.DoseIn,
			"yield_out":   params.YieldOut,
			"time":        params.BrewTime,
			"target_time": params.TargetTime,
			"grind_size":  params.GrindSize,
		}
		if params.GearIDs != nil {
			body["gear_ids"] = params.GearIDs
		}
		if params.PreInfusion != nil {
			body["pre_infusion"] = *params.PreInfusion
		}
		if params.TastingNote != nil {
			body["tasting_note"] = *params.TastingNote
		}

		var result any
		if err := client.put(ctx, fmt.Sprintf("/extractions/%s", params.ID), body, &result); err != nil {
			return nil, nil, err
		}
		text, err := toJSON(result)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_extraction",
		Description: "Delete an extraction by ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params DeleteExtractionParams) (*mcp.CallToolResult, any, error) {
		if err := client.delete(ctx, fmt.Sprintf("/extractions/%s", params.ID)); err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "Extraction deleted successfully"}}}, nil, nil
	})
}
