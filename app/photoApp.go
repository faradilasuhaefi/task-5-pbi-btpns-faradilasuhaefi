package app

type PhotoCreate struct {
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
}

type PhotoUpdate struct {
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
}

type PhotoResult struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Caption   string `json:"caption"`
	PhotoUrl  string `json:"photo_url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Email     string `json:"email"`
}
