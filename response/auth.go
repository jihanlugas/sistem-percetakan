package response

type Init struct {
	User        *User        `json:"user,omitempty"`
	Company     *Company     `json:"company,omitempty"`
	Usercompany *Usercompany `json:"usercompany,omitempty"`
}
