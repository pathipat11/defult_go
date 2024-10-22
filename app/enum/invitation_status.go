package enum

type InvitationStatus string

const (
	INVITATION_STATUS_PENDING  InvitationStatus = "pending"
	INVITATION_STATUS_ACCEPTED InvitationStatus = "accepted"
	INVITATION_STATUS_REJECTED InvitationStatus = "rejected"
)

func GetInvitationStatus(t InvitationStatus) InvitationStatus {
	switch t {
	case INVITATION_STATUS_ACCEPTED:
		return INVITATION_STATUS_ACCEPTED
	case INVITATION_STATUS_REJECTED:
		return INVITATION_STATUS_REJECTED
	default:
		return INVITATION_STATUS_PENDING
	}
}
