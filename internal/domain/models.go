package domain

import (
	"time"
)

// Cat represents a spy cat in the system.
type Cat struct {
	ID                int       `db:"id" json:"id"`
	Name              string    `db:"name" json:"name"`
	YearsOfExperience int       `db:"years_of_experience" json:"years_of_experience"`
	Breed             string    `db:"breed" json:"breed"`
	Salary            float64   `db:"salary" json:"salary"`
	Status            string    `db:"status" json:"status"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

// Mission represents a mission assigned to a spy cat.
type Mission struct {
	ID        int       `db:"id" json:"id"`
	CatID     *int      `db:"cat_id" json:"cat_id"` // Nullable, as a mission can be unassigned
	Completed bool      `db:"completed" json:"completed"`
	Targets   []Target  `db:"-" json:"targets"` // Skip DB mapping for nested slice
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Target represents a target within a mission.
type Target struct {
	ID        int       `db:"id" json:"id"`
	MissionID int       `db:"mission_id" json:"mission_id"`
	Name      string    `db:"name" json:"name"`
	Country   string    `db:"country" json:"country"`
	Notes     string    `db:"notes" json:"notes"`
	Completed bool      `db:"completed" json:"completed"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
