package user

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Create(ctx context.Context, req request.CreateUser) (*model.User, bool, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return nil, false, err
	}

	m := &model.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(bytes),
	}

	_, err = s.db.NewInsert().Model(m).Exec(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("email already exists")
		}
	}

	return m, false, err
}

func (s *Service) Update(ctx context.Context, req request.UpdateUser, id request.GetByIDUser) (*model.User, bool, error) {
	ex, err := s.db.NewSelect().Table("users").Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}

	if !ex {
		return nil, true, errors.New("user not found")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return nil, false, err
	}
	m := &model.User{
		ID:        id.ID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(bytes),
	}
	m.SetUpdateNow()
	_, err = s.db.NewUpdate().Model(m).
		Set("first_name = ?first_name").
		Set("last_name = ?last_name").
		Set("email = ?email").
		Set("password = ?password").
		Set("updated_at = ?updated_at").
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("email already exists")
		}
	}
	return m, false, err
}

func (s *Service) List(ctx context.Context, req request.ListUser) ([]response.ListUser, int, error) {
	offset := (req.Page - 1) * req.Size

	m := []response.ListUser{}
	query := s.db.NewSelect().
		TableExpr("users AS u").
		Column("u.id", "u.first_name", "u.last_name", "u.email", "u.created_at", "u.updated_at").
		Where("deleted_at IS NULL")

	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		if req.SearchBy != "" {
			query.Where(fmt.Sprintf("LOWER(u.%s) LIKE ?", req.SearchBy), search)
		} else {
			query.Where("LOWER(first_name) LIKE ?", search)
		}
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("u.%s %s", req.SortBy, req.OrderBy)

	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}
	return m, count, err
}

func (s *Service) Get(ctx context.Context, id request.GetByIDUser) (*response.ListUser, error) {
	m := response.ListUser{}
	err := s.db.NewSelect().
		TableExpr("users AS u").
		Column("u.id", "u.first_name", "u.last_name", "u.email", "u.created_at", "u.updated_at").
		Where("id = ?", id.ID).Where("deleted_at IS NULL").Scan(ctx, &m)
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id request.GetByIDUser) error {
	ex, err := s.db.NewSelect().Table("users").Where("id = ?", id.ID).Where("deleted_at IS NULL").Exists(ctx)
	if err != nil {
		return err
	}

	if !ex {
		return errors.New("user not found")
	}

	_, err = s.db.NewDelete().Model((*model.User)(nil)).Where("id = ?", id.ID).Exec(ctx)
	return err
}
