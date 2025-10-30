package gymModels

type GymRecord struct {
	GymId         int     `json:"gymId"`
	GymName       string  `json:"gymName"`
	State         string  `json:"state"`
	City          string  `json:"city"`
	GymAddress    string  `json:"gymAddress"`
	ContactNo     string  `json:"contactNo"`
	OfficialEmail string  `json:"officialEmail"`
	CreatedBy     int     `json:"createdBy"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     *string `json:"updatedAt"`
	LocationLat   string  `json:"locationLat"`
	LocationLong  string  `json:"locationLong"`
}

type GymRecordWithAdditionalData struct {
	GymId         int     `json:"gymId"`
	GymName       string  `json:"gymName"`
	State         string  `json:"state"`
	City          string  `json:"city"`
	GymAddress    string  `json:"gymAddress"`
	ContactNo     string  `json:"contactNo"`
	OfficialEmail string  `json:"officialEmail"`
	CreatedBy     int     `json:"createdBy"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     *string `json:"updatedAt"`
	LocationLat   string  `json:"locationLat"`
	LocationLong  string  `json:"locationLong"`
	AverageRating int     `json:"averageRating"`
	TotalFeedback int     `json:"totalFeedbacks"`
	TrainersCount int     `json:"trainersCount"`
	StaffCount    int     `json:"staffCount"`
	ManagersCount int     `json:"managersCount"`
	ActivePlans   int     `json:"activePlans"`
}

type UpdateGymRequest struct {
	GymId         int    `json:"gymId"`
	GymName       string `json:"gymName"`
	State         string `json:"state"`
	City          string `json:"city"`
	GymAddress    string `json:"gymAddress"`
	ContactNo     string `json:"contactNo"`
	OfficialEmail string `json:"officialEmail"`
	LocationLat   string `json:"locationLat"`
	LocationLong  string `json:"locationLong"`
}
