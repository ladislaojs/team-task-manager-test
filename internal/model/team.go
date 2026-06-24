package model

type TeamID string

type TeamRole string

const (
	TeamRoleOwner  TeamRole = "owner"
	TeamRoleAdmin  TeamRole = "admin"
	TeamRoleMember TeamRole = "member"
)

func (r TeamRole) CanInvite() bool {
	return r == TeamRoleOwner || r == TeamRoleAdmin
}

type Team struct {
	ID        TeamID
	Name      string
	CreatedBy uint64
}

type TeamMember struct {
	UserID UserID
	TeamID TeamID
	Role   TeamRole
}
