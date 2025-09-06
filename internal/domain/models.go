package domain

import (
	"time"
)

// Cat represents a spy cat in the system.
type Cat struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	YearsOfExperience int       `json:"years_of_experience"`
	Breed             string    `json:"breed"`
	Salary            float64   `json:"salary"`
	Status            string    `json:"status"` // e.g., 'available', 'on_mission'
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Mission represents a mission assigned to a spy cat.
type Mission struct {
	ID        int       `json:"id"`
	CatID     *int      `json:"cat_id"` // Nullable, as a mission can be unassigned
	Completed bool      `json:"completed"`
	Targets   []Target  `json:"targets"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Target represents a target within a mission.
type Target struct {
	ID        int       `json:"id"`
	MissionID int       `json:"mission_id"`
	Name      string    `json:"name"`
	Country   string    `json:"country"`
	Notes     string    `json:"notes"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
