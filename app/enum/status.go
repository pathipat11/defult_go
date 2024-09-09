package enum

import (
	"encoding/json"
)

type Status int

const (
	STATUS_ACTIVE Status = iota + 1
	STATUS_INACTIVE
	STATUS_OUT
)

var (
	statusName = map[Status]string{
		STATUS_ACTIVE:   "Active",
		STATUS_INACTIVE: "Inactive",
		STATUS_OUT:      "Out",
	}
	intToStatus = map[int]Status{
		1: STATUS_ACTIVE,
		2: STATUS_INACTIVE,
		3: STATUS_OUT,
	}
)

func (s Status) String() string {
	return statusName[s]
}

func (s Status) List() map[Status]string {
	return statusName
}

func GetStatus(value int) Status {
	if status, ok := intToStatus[value]; ok {
		return status
	}
	return STATUS_ACTIVE
}

func (s *Status) UnmarshalJSON(data []byte) error {
	var value int
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*s = GetStatus(value)
	return nil
}

// GetStatusFromString converts a string to a Status enum value
func GetStatusFromInt(i int) Status {
	switch i {
	case 0:
		return STATUS_ACTIVE
	case 1:
		return STATUS_INACTIVE
	default:
		return STATUS_OUT
	}
}

// MarshalJSON customizes the JSON output of the Status type
func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
