package dto

type CreateTeamRequest struct {
	Name string `json:"name"`
}

type TeamResponse struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type InviteRequest struct {
	UserID uint64 `json:"user_id"`
}

type TeamExtendedResponse struct {
	ID                    uint64 `json:"id"`
	Name                  string `json:"name"`
	MemberCount           int    `json:"member_count"`
	LastWeekDoneTaskCount int    `json:"last_week_done_task_count"`
}

type TopTaskCreatorResponse struct {
	TeamID    uint64 `json:"team_id"`
	UserID    uint64 `json:"user_id"`
	UserName  string `json:"user_name"`
	TaskCount int    `json:"task_count"`
	Rank      int    `json:"rank"`
}
