package user

import (
	"app/app/enum"
	"app/app/model"
	"app/app/response"
	"app/internal/logger"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-pg/pg/v10"
)

func (s *Service) Create(ctx context.Context, user model.User) (*model.User, error) {
	// Check if a user with the same citizen_id already exists
	existingUserByCitizenID, err := s.GetUserByCitizenID(ctx, user.CitizenID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if existingUserByCitizenID != nil {
		return nil, errors.New("user with the same citizen_id already exists")
	}

	// Check if a user with the same email already exists
	existingUserByEmail, err := s.GetUserByEmail(ctx, user.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if existingUserByEmail != nil {
		return nil, errors.New("user with the same email already exists")
	}

	// Check if a user with the same username already exists
	existingUserByUsername, err := s.GetUserByUsername(ctx, user.Username)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if existingUserByUsername != nil {
		return nil, errors.New("user with the same username already exists")
	}

	// Insert the user to get the ID
	if _, err := s.db.NewInsert().Model(&user).Exec(ctx); err != nil {
		logger.Infof("[error]: %v", err)
		return nil, err
	}

	return &user, nil
}

func (s *Service) GetUserByCitizenID(ctx context.Context, citizenID string) (*model.User, error) {
	var user model.User
	err := s.db.NewSelect().Model(&user).Where("citizen_id = ?", citizenID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := s.db.NewSelect().Model(&user).Where("email = ?", email).Scan(ctx)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil // No user found with this email
		}
		return nil, err
	}
	return &user, nil
}

func (s *Service) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := s.db.NewSelect().Model(&user).Where("username = ?", username).Scan(ctx)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil // No user found with this email
		}
		return nil, err
	}
	return &user, nil
}

func (s *Service) List(ctx context.Context, limit, page int, search string) (*response.UserListResponse, error) {
	var users []response.UserResponse

	// Calculate the offset for pagination
	offset := (page - 1) * limit

	// Build the query with optional search and sorting
	query := s.db.NewSelect().
		Model((*model.User)(nil)).
		ColumnExpr("id, username, firstname, lastname, nickname, email, role_id, CASE status " +
			"WHEN 1 THEN 'Active' " +
			"WHEN 2 THEN 'Inactive' " +
			"WHEN 3 THEN 'Out' " +
			"ELSE 'Unknown' END AS status").
		Order("id ASC")

	// Apply search filter if search string is provided
	if search != "" {
		searchParam := fmt.Sprintf("%%%s%%", strings.ToLower(search))
		query.Where("LOWER(firstname) ILIKE ?", searchParam)
	}

	// Count total users for pagination (with the search filter applied)
	totalCount, err := query.Count(ctx)
	if err != nil {
		return nil, err
	}
	if totalCount == 0 {
		return &response.UserListResponse{
			Users:      users,
			Pagination: response.Pagination{},
		}, nil
	}

	// Apply limit and offset for pagination
	if limit > 0 {
		query.Limit(limit)
	}
	if limit > 0 && offset >= 0 {
		query.Offset(offset)
	}

	// Execute the query and map results to the response struct directly
	err = query.Scan(ctx, &users)
	if err != nil {
		return nil, err
	}

	// Calculate pagination details, ensuring no division by zero
	perPage := limit
	if limit <= 0 {
		perPage = totalCount
	}

	totalPages := 1
	if perPage > 0 {
		totalPages = (totalCount + perPage - 1) / perPage
	}

	pagination := response.Pagination{
		CurrentPage: page,
		PerPage:     perPage,
		TotalPages:  totalPages,
		Total:       totalCount,
	}

	return &response.UserListResponse{
		Users:      users,
		Pagination: pagination,
	}, nil
}

func (s *Service) ListSingle(ctx context.Context, userID uint) (*response.UserListResponse, error) {
	var user response.UserResponse

	// Build the query to fetch a single user by ID
	query := s.db.NewSelect().
		Model((*model.User)(nil)).
		ColumnExpr("id, username, firstname, lastname, nickname, email, role_id, CASE status "+
			"WHEN 1 THEN 'Active' "+
			"WHEN 2 THEN 'Inactive' "+
			"WHEN 3 THEN 'Out' "+
			"ELSE 'Unknown' END AS status").
		Where("id = ?", userID).
		Order("id ASC").
		Limit(1).
		Offset(0)

	// Execute the query and map the result directly to the `user` struct
	err := query.Scan(ctx, &user)
	if err != nil {
		return nil, err
	}

	// Return the response with the single user in a slice
	return &response.UserListResponse{
		Users: []response.UserResponse{user},
	}, nil
}

func (s *Service) Update(ctx context.Context, user model.User, userID uint) (*model.User, error) {
	// Fetch the existing user data from the database
	existingUser, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if a user with the same email already exists (excluding the current user)
	if user.Email != "" && user.Email != existingUser.Email {
		existingUserByEmail, err := s.GetUserByEmail(ctx, user.Email)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		if existingUserByEmail != nil {
			return nil, errors.New("user with the same email already exists")
		}
	}

	// Only update fields that are not empty in the request
	if user.Firstname != "" {
		existingUser.Firstname = user.Firstname
	}
	if user.Lastname != "" {
		existingUser.Lastname = user.Lastname
	}
	if user.Nickname != "" {
		existingUser.Nickname = user.Nickname
	}
	if !user.Birthdate.IsZero() {
		existingUser.Birthdate = user.Birthdate
	}
	if user.Gender != enum.GENDER_UNKNOWN { // Adjusted to use GENDER_UNKNOWN
		existingUser.Gender = user.Gender
	}
	if user.Nationality != "" {
		existingUser.Nationality = user.Nationality
	}
	if user.RelationshipStatus != enum.RELATIONSHIP_STATUS_UNKNOWN { // Adjusted to use RELATIONSHIP_STATUS_UNKNOWN
		existingUser.RelationshipStatus = user.RelationshipStatus
	}
	if user.Address1 != "" {
		existingUser.Address1 = user.Address1
	}
	if user.Address2 != "" {
		existingUser.Address2 = user.Address2
	}
	if user.MobileNo != "" {
		existingUser.MobileNo = user.MobileNo
	}
	if user.Email != "" {
		existingUser.Email = user.Email
	}
	if user.RoleID != 0 {
		existingUser.RoleID = user.RoleID
	}

	// Update the user in the database
	if _, err := s.db.NewUpdate().Model(existingUser).Where("id = ?", userID).Exec(ctx); err != nil {
		logger.Infof("[error]: %v", err)
		return nil, err
	}

	return existingUser, nil
}

func (s *Service) GetUserByID(ctx context.Context, userID uint) (*model.User, error) {
	var user model.User
	if err := s.db.NewSelect().Model(&user).Where("id = ?", userID).Scan(ctx); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) SoftDelete(ctx context.Context, userID uint) error {
	_, err := s.db.NewUpdate().
		Model((*model.User)(nil)).
		Set("deleted_at = NOW()").
		Where("id = ?", userID).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// func (s *Service) Delete(ctx context.Context, userID uint) error {
// 	_, err := s.db.NewDelete().
// 		Model((*model.User)(nil)).
// 		Where("id = ?", userID).
// 		Exec(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
