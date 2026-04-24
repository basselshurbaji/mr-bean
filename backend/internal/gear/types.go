package gear

import "time"

type GearResponse struct {
	ID        string    `json:"id"`
	TypeID    string    `json:"type_id"`
	Name      string    `json:"name"`
	Brand     *string   `json:"brand,omitempty"`
	Model     *string   `json:"model,omitempty"`
	Year      *string   `json:"year,omitempty"`
	Notes     *string   `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StationResponse struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Gear      []GearResponse `json:"gear"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func gearToResponse(g GearItem) GearResponse {
	return GearResponse{
		ID:        g.ID,
		TypeID:    g.TypeID,
		Name:      g.Name,
		Brand:     g.Brand,
		Model:     g.Model,
		Year:      g.Year,
		Notes:     g.Notes,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
}

func stationToResponse(s Station) StationResponse {
	gear := make([]GearResponse, len(s.Gear))
	for i, g := range s.Gear {
		gear[i] = gearToResponse(g)
	}
	return StationResponse{
		ID:        s.ID,
		Name:      s.Name,
		Gear:      gear,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
