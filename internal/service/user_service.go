package service

import (
	"context"
	"errors"
	"time"

	"github.com/Karthisgowda/Ainyx/db/sqlc"
	"github.com/Karthisgowda/Ainyx/internal/models"
	"github.com/Karthisgowda/Ainyx/internal/repository"
	"github.com/jackc/pgx/v5"
)

const dateLayout = "2006-01-02"

var ErrFutureDOB = errors.New("dob cannot be in the future")

type UserService interface {
	Create(ctx context.Context, request models.UserRequest) (models.UserResponse, error)
	GetByID(ctx context.Context, id int32) (models.UserResponse, error)
	List(ctx context.Context, limit int32, offset int32) ([]models.UserResponse, error)
	Update(ctx context.Context, id int32, request models.UserRequest) (models.UserResponse, error)
	Delete(ctx context.Context, id int32) error
}

type userService struct {
	repo repository.UserRepository
	now  func() time.Time
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
		now:  time.Now,
	}
}

func (s *userService) Create(ctx context.Context, request models.UserRequest) (models.UserResponse, error) {
	dob, err := parseDOB(request.Dob, s.now())
	if err != nil {
		return models.UserResponse{}, err
	}

	user, err := s.repo.Create(ctx, sqlc.CreateUserParams{Name: request.Name, Dob: dob})
	if err != nil {
		return models.UserResponse{}, err
	}

	return toUserResponse(user, false, s.now()), nil
}

func (s *userService) GetByID(ctx context.Context, id int32) (models.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return models.UserResponse{}, err
	}

	return toUserResponse(user, true, s.now()), nil
}

func (s *userService) List(ctx context.Context, limit int32, offset int32) ([]models.UserResponse, error) {
	users, err := s.repo.List(ctx, sqlc.ListUsersParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}

	response := make([]models.UserResponse, 0, len(users))
	now := s.now()
	for _, user := range users {
		response = append(response, toUserResponse(user, true, now))
	}

	return response, nil
}

func (s *userService) Update(ctx context.Context, id int32, request models.UserRequest) (models.UserResponse, error) {
	dob, err := parseDOB(request.Dob, s.now())
	if err != nil {
		return models.UserResponse{}, err
	}

	user, err := s.repo.Update(ctx, sqlc.UpdateUserParams{ID: id, Name: request.Name, Dob: dob})
	if err != nil {
		return models.UserResponse{}, err
	}

	return toUserResponse(user, false, s.now()), nil
}

func (s *userService) Delete(ctx context.Context, id int32) error {
	rowsAffected, err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func parseDOB(value string, now time.Time) (time.Time, error) {
	dob, err := time.Parse(dateLayout, value)
	if err != nil {
		return time.Time{}, err
	}
	if dob.After(now) {
		return time.Time{}, ErrFutureDOB
	}
	return dob, nil
}

func toUserResponse(user sqlc.User, includeAge bool, now time.Time) models.UserResponse {
	response := models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Format(dateLayout),
	}
	if includeAge {
		age := CalculateAge(user.Dob, now)
		response.Age = &age
	}
	return response
}
