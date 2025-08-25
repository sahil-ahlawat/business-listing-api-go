package models

type CreateListingRequest struct {
    FullName string  `json:"full_name"`
    Email    string  `json:"email"`
    Phone    string  `json:"phone"`
    Title    string  `json:"title"`
    Category string  `json:"category"`
    Lat      float64 `json:"lat"`
    Long     float64 `json:"long"`
    Slug     string  `json:"slug"` // 👈 Add this
}