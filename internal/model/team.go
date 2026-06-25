package model

type TeamRole string

const (
	TeamRoleOwner  TeamRole = "owner"
	TeamRoleAdmin  TeamRole = "admin"
	TeamRoleMember TeamRole = "member"
)

func (r TeamRole) CanInvite() bool {
	return r == TeamRoleOwner || r == TeamRoleAdmin
}

func (r TeamRole) CanUpdate() bool {
	return r == TeamRoleOwner || r == TeamRoleAdmin
}

type Team struct {
	ID        uint64
	Name      string
	CreatedBy uint64
}

type TeamExtended struct {
	ID                    uint64
	Name                  string
	MemberCount           int
	LastWeekDoneTaskCount int
}

type TeamMember struct {
	UserID uint64
	TeamID uint64
	Role   TeamRole
}
