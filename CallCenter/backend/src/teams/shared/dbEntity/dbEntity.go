package dbEntity

type Team struct {
	ID          string   `json:"id,omitempty" bson:"id,omitempty"`
	Name        string   `json:"displayName,omitempty" bson:"displayName,omitempty"`
	Owners      []string `json:"owners,omitempty" bson:"owners,omitempty"`
	UpdatedBy   string   `json:"updatedBy,omitempty" bson:"updatedBy,omitempty"`
	UpdatedDate string   `json:"updatedDate,omitempty" bson:"updatedDate,omitempty"`
}
