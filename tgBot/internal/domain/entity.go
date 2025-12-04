package domain

const (
	WaitingCommand     = "waiting_command"
	WaitingURl         = "waiting_url"
	WaitingDescription = "waiting_description"
	WaitingDelete      = "waiting_delete"
)

type UpdatedLink struct {
	ChatIDs []int64 `json:"chat_ids"`
	Link    *Link   `json:"link"`
}

type Link struct {
	URL        string `json:"url"`
	Domain     string `json:"domain"`
	Author     string `json:"author"`
	Repository string `json:"repository"`
	Desc       string `json:"desc"`
	ChatID     int64  `json:"chat_id"`
}

type UserStateInfo struct {
	UserID int64  `json:"userID"`
	URL    string `json:"url"`
	Desc   string `json:"desc"`
	State  string `json:"state"`
}

type DeleteLinkRequest struct {
	ChatID int64  `json:"chat_id"`
	Alias  string `json:"alias"`
}
