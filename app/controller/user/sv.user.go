package user

import (
	"app/app/enum"
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"
	"errors"
	"math"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Create(ctx context.Context, req request.CreateUser) (*model.User, error) {
	// Check if the user already exists
	ex, err := s.db.NewSelect().Model(&model.User{}).Where("username = ?", req.Username).Exists(ctx)
	if err != nil {
		return nil, err
	}

	if ex {
		return nil, errors.New("username already exists")
	}

	ex, err = s.db.NewSelect().Model(&model.User{}).Where("email = ?", req.Email).Exists(ctx)
	if err != nil {
		return nil, err
	}

	if ex {
		return nil, errors.New("email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create the user
	user := &model.User{
		Username:    req.Username,
		Email:       req.Email,
		Password:    string(hashedPassword),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		DisplayName: req.DisplayName,
		RoleID:      1,
		Status:      enum.STATUS_ACTIVE,
	}
	if _, err := s.db.NewInsert().Model(user).Exec(ctx); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) List(ctx context.Context, limit, page int, search string, roleID string, status string, planType string) ([]model.User, response.Pagination, error) {
	Offset := (page - 1) * limit
	users := []model.User{}
	query := s.db.NewSelect().Model(&users)
	if search != "" {
		query.Where("first_name LIKE ? OR display_name LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if roleID != "" {
		query.Where("role_id = ?", roleID)
	}
	if status != "" {
		query.Where("status = ?", status)
	}
	// if planType != "" {
	// 	query.Where("plan_type = ?", planType)
	// }
	if limit == 0 {
		limit = 10
	}
	if page == 0 {
		page = 1
	}
	query.Limit(limit).Offset(Offset).
		Order("created_at ASC")
	if err := query.Scan(ctx); err != nil {
		return nil, response.Pagination{}, err
	}

	total, err := s.db.NewSelect().Model(&model.User{}).Count(ctx)
	if err != nil {
		return nil, response.Pagination{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	paginate := response.Pagination{
		CurrentPage: page,
		PerPage:     limit,
		TotalPages:  totalPages,
		Total:       total,
	}

	return users, paginate, nil
}

func (s *Service) Get(ctx context.Context, id string) (*model.User, error) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	user := model.User{ID: idInt}
	if err := s.db.NewSelect().Model(&user).WherePK().Scan(ctx); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) Update(ctx context.Context, req request.UpdateUser, id string) (*model.User, error) {

	user := model.User{}
	if err := s.db.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	// Hash the password
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.DisplayName = req.DisplayName

	if _, err := s.db.NewUpdate().Model(&user).Where("id = ?", id).Exec(ctx); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	user := model.User{ID: idInt}
	if _, err := s.db.NewDelete().Model(&user).WherePK().Exec(ctx); err != nil {
		return err
	}

	return nil
}
