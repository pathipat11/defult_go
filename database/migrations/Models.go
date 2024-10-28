package migrations

import "app/app/model"

func Models() []any {
	return []any{
		// (*model.ActivateLog)(nil),
		// (*model.AdminAction)(nil),
		// (*model.Ad)(nil),
		// (*model.Invitation)(nil),
		// (*model.Notification)(nil),
		// (*model.Permission)(nil),
		// (*model.RolePermission)(nil),
		// (*model.Role)(nil),
		// (*model.Setting)(nil),
		// (*model.SprintPlan)(nil),
		(*model.SprintTask)(nil),
		// (*model.Subscription)(nil),
		// (*model.TeamMember)(nil),
		// (*model.Team)(nil),
		// (*model.Transaction)(nil),
		// (*model.User)(nil),
	}
}

func RawBeforeQueryMigrate() []string {
	return []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
	}
}

func RawAfterQueryMigrate() []string {
	return []string{}
}
