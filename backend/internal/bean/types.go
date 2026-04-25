package bean

import "time"

type BeanResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Roaster      *string   `json:"roaster,omitempty"`
	Origin       *string   `json:"origin,omitempty"`
	Process      *string   `json:"process,omitempty"`
	RoastLevel   *string   `json:"roast_level,omitempty"`
	TastingNotes *string   `json:"tasting_notes,omitempty"`
	Notes        *string   `json:"notes,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func beanToResponse(b Bean) BeanResponse {
	return BeanResponse{
		ID:           b.ID,
		Name:         b.Name,
		Roaster:      b.Roaster,
		Origin:       b.Origin,
		Process:      b.Process,
		RoastLevel:   b.RoastLevel,
		TastingNotes: b.TastingNotes,
		Notes:        b.Notes,
		CreatedAt:    b.CreatedAt,
		UpdatedAt:    b.UpdatedAt,
	}
}
