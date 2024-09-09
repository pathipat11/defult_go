package enum

import (
	"encoding/json"
)

type RelationshipStatus int

const (
	RELATIONSHIP_STATUS_UNKNOWN RelationshipStatus = iota
	RELATIONSHIP_STATUS_SINGLE
	RELATIONSHIP_STATUS_MARRIED
	RELATIONSHIP_STATUS_DIVORCED
	RELATIONSHIP_STATUS_WIDOWED
	RELATIONSHIP_STATUS_SEPARATED
)

var (
	relationshipStatusName = map[RelationshipStatus]string{
		RELATIONSHIP_STATUS_UNKNOWN:   "Unknown",
		RELATIONSHIP_STATUS_SINGLE:    "Single",
		RELATIONSHIP_STATUS_MARRIED:   "Married",
		RELATIONSHIP_STATUS_DIVORCED:  "Divorced",
		RELATIONSHIP_STATUS_WIDOWED:   "Widowed",
		RELATIONSHIP_STATUS_SEPARATED: "Separated",
	}
	intToRelationshipStatus = map[int]RelationshipStatus{
		0: RELATIONSHIP_STATUS_UNKNOWN,
		1: RELATIONSHIP_STATUS_SINGLE,
		2: RELATIONSHIP_STATUS_MARRIED,
		3: RELATIONSHIP_STATUS_DIVORCED,
		4: RELATIONSHIP_STATUS_WIDOWED,
		5: RELATIONSHIP_STATUS_SEPARATED,
	}
)

func (s RelationshipStatus) String() string {
	return relationshipStatusName[s]
}

func GetRelationshipStatus(value int) RelationshipStatus {
	if status, ok := intToRelationshipStatus[value]; ok {
		return status
	}
	return RELATIONSHIP_STATUS_UNKNOWN
}

func (s *RelationshipStatus) UnmarshalJSON(data []byte) error {
	var value int
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*s = GetRelationshipStatus(value)
	return nil
}

func GetRelationshipStatusFromInt(i int) RelationshipStatus {
	switch i {
	case 0:
		return RELATIONSHIP_STATUS_SINGLE
	case 1:
		return RELATIONSHIP_STATUS_MARRIED
	default:
		return RELATIONSHIP_STATUS_UNKNOWN
	}
}

func (s RelationshipStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
