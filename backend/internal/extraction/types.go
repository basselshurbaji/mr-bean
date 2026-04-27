package extraction

import "time"

type BeanSummary struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Roaster *string `json:"roaster"`
	Roast   *string `json:"roast"`
}

type GearSummary struct {
	ID     string `json:"id"`
	TypeID string `json:"type_id"`
	Name   string `json:"name"`
}

type ExtractionResponse struct {
	ID          string        `json:"id"`
	UserID      string        `json:"user_id"`
	Bean        BeanSummary   `json:"bean"`
	DoseIn      float64       `json:"dose_in"`
	YieldOut    float64       `json:"yield_out"`
	Time        float64       `json:"time"`
	TargetTime  float64       `json:"target_time"`
	GrindSize   float64       `json:"grind_size"`
	PreInfusion bool          `json:"pre_infusion"`
	TastingNote *string       `json:"tasting_note,omitempty"`
	Gear        []GearSummary `json:"gear"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

func toResponse(e Extraction) ExtractionResponse {
	gear := make([]GearSummary, len(e.Gear))
	for i, g := range e.Gear {
		gear[i] = GearSummary(g)
	}
	return ExtractionResponse{
		ID: e.ID,
		UserID: e.UserID,
		Bean: BeanSummary{
			ID:      e.Bean.ID,
			Name:    e.Bean.Name,
			Roaster: e.Bean.Roaster,
			Roast:   e.Bean.Roast,
		},
		DoseIn:      e.DoseIn,
		YieldOut:    e.YieldOut,
		Time:        e.Time,
		TargetTime:  e.TargetTime,
		GrindSize:   e.GrindSize,
		PreInfusion: e.PreInfusion,
		TastingNote: e.TastingNote,
		Gear:        gear,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
