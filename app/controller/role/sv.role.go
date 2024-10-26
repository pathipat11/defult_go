package role

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"
	"errors"
	"math"
	"strconv"
)

func (s *Service) Create(ctx context.Context, req request.CreateRole) (*model.Role, error) {
	role := &model.Role{
		RoleName:    req.Name,
		Description: req.Description,
	}

	if _, err := s.db.NewInsert().Model(role).Exec(ctx); err != nil {
		return nil, err
	}

	// change role on user in request
	if req.UserID != nil {
		for _, userID := range req.UserID {
			user := model.User{ID: userID}
			user.RoleID = role.ID
			if _, err := s.db.NewUpdate().Model(&user).Column("role_id").Where("id = ?", user.ID).Exec(ctx); err != nil {
				return nil, err
			}
		}
	}

	return role, nil
}

func (s *Service) List(ctx context.Context, limit, page int) ([]model.Role, *response.Pagination, error) {
	var roles []model.Role
	query := s.db.NewSelect().Model(&roles)

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

	if err := query.Scan(ctx); err != nil {
		return nil, nil, err
	}

	return roles, paginate, nil
}

func (s *Service) Get(ctx context.Context, id string) (*model.Role, error) {
	roleID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	role := model.Role{ID: roleID}
	if err := s.db.NewSelect().Model(&role).WherePK().Scan(ctx); err != nil {
		return nil, err
	}

	return &role, nil
}

func (s *Service) Update(ctx context.Context, req request.UpdateRole, ID int) (*model.Role, error) {
	ex := &model.Role{}
	err := s.db.NewSelect().Model(ex).Where("id = ?", ID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	role := model.Role{
		ID:          int64(ID),
		RoleName:    req.Name,
		Description: req.Description,
	}
	role.SetCreated(ex.CreatedAt)
	role.SetUpdateNow()

	if _, err := s.db.NewUpdate().Model(&role).WherePK().Exec(ctx); err != nil {
		return nil, err
	}

	return &role, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	// check uset has  this role
	exists, err := s.db.NewSelect().Model((*model.User)(nil)).Where("role_id = ?", id).Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("this role has been used by user")
	}
	_, err = s.db.NewDelete().Model(&model.Role{}).Where("id = ?", id).Exec(ctx)

	return err
}
