package team

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"
	"errors"
	"math"
)

func (s *Service) Create(ctx context.Context, req request.CreateTeam, user string) (*model.Team, error) {
	team := &model.Team{
		TeamName:  req.Name,
		CreatedBy: user,
	}

	if _, err := s.db.NewInsert().Model(team).Exec(ctx); err != nil {
		return nil, err
	}

	return team, nil
}

func (s *Service) AddTeamMember(ctx context.Context, req request.CreateTeamMember) (*model.TeamMember, error) {
	ex, err := s.db.NewSelect().Model(&model.TeamMember{}).Where("team_id = ? AND user_id = ?", req.TeamID, req.UserID).Exists(ctx)
	if err != nil {
		return nil, err
	}

	if ex {
		return nil, errors.New("this user has been added to this team")
	} else {
		team := &model.TeamMember{
			TeamID: req.TeamID,
			UserID: req.UserID,
		}

		if _, err := s.db.NewInsert().Model(team).Exec(ctx); err != nil {
			return nil, err
		}

		return team, nil
	}
}

func (s *Service) RemoveTeamMember(ctx context.Context, req request.CreateTeamMember) error {
	_, err := s.db.NewDelete().Model(&model.TeamMember{}).Where("team_id = ? AND user_id = ?", req.TeamID, req.UserID).Exec(ctx)
	return err
}

func (s *Service) List(ctx context.Context, limit, page int, search string) ([]response.TeamResponse, *response.Pagination, error) {
	m := []response.TeamResponse{}
	query := s.db.NewSelect().TableExpr("teams AS t").
		Column("t.id", "t.team_name", "t.created_by", "t.created_at", "t.updated_at").
		ColumnExpr("json_agg(json_build_object('user_id', u.id, 'display_name', u.display_name)) AS members").
		Join("LEFT JOIN team_members AS tm ON t.id = tm.team_id").
		Join("LEFT JOIN users AS u ON tm.user_id = u.id").
		Group("t.id", "t.team_name", "t.created_by", "t.created_at", "t.updated_at").
		Where("t.deleted_at IS NULL")

	if search != "" {
		query.Where("t.team_name ILIKE ?", "%"+search+"%")
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	paginate := &response.Pagination{
		CurrentPage: page,
		PerPage:     limit,
		TotalPages:  totalPages,
		Total:       total,
	}

	if limit > 0 {
		query.Limit(limit).Offset((page - 1) * limit)
	}

	if err := query.Scan(ctx, &m); err != nil {
		return nil, nil, err
	}

	return m, paginate, nil
}

func (s *Service) Get(ctx context.Context, id string) (*response.TeamResponse, error) {

	m := response.TeamResponse{}
	err := s.db.NewSelect().TableExpr("teams AS t").
		Column("t.id", "t.team_name", "t.created_by", "t.created_at", "t.updated_at").
		ColumnExpr("json_agg(json_build_object('user_id', u.id, 'display_name', u.display_name)) AS members").
		Join("LEFT JOIN team_members AS tm ON t.id = tm.team_id").
		Join("LEFT JOIN users AS u ON tm.user_id = u.id").
		Group("t.id", "t.team_name", "t.created_by", "t.created_at", "t.updated_at").
		Where("t.deleted_at IS NULL").
		Where("t.id = ?", id).Scan(ctx, &m)

	return &m, err
}

func (s *Service) Update(ctx context.Context, req request.UpdateTeam, ID string) (*model.Team, error) {
	ex := &model.Team{}
	err := s.db.NewSelect().Model(ex).Where("id = ?", ID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	team := model.Team{
		ID:        ID,
		TeamName:  req.Name,
		CreatedBy: ex.CreatedBy,
	}
	team.SetCreated(ex.CreatedAt)
	team.SetUpdateNow()

	if _, err := s.db.NewUpdate().Model(&team).WherePK().Exec(ctx); err != nil {
		return nil, err
	}

	return &team, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	// check user has  this team
	exists, err := s.db.NewSelect().Model((*model.TeamMember)(nil)).Where("team_id = ?", id).Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("this team has been used by user")
	}
	_, err = s.db.NewDelete().Model(&model.Team{}).Where("id = ?", id).Exec(ctx)

	return err
}
