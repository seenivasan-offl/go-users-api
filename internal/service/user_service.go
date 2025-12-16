package service

import (
	"context"
	"time"

	"go-users-api/internal/models"
	"go-users-api/internal/repository"
)

type UserService interface {
	Create(ctx context.Context, req models.CreateUserRequest) (models.User, error)
	Get(ctx context.Context, id int64) (models.User, error)
	List(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, id int64, req models.UpdateUserRequest) (models.User, error)
	Delete(ctx context.Context, id int64) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func calculateAge(dob, now time.Time) int {
	years := now.Year() - dob.Year()
	if now.Month() < dob.Month() ||
		(now.Month() == dob.Month() && now.Day() < dob.Day()) {
		years--
	}
	if years < 0 {
		return 0
	}
	return years
}

func (s *userService) Create(ctx context.Context, req models.CreateUserRequest) (models.User, error) {
	dob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		return models.User{}, err
	}

	u, err := s.repo.Create(ctx, req.Name, dob)
	if err != nil {
		return models.User{}, err
	}

	d := u.Dob.Time // pgtype.Date -> time.Time

	return models.User{
		ID:   int64(u.ID),
		Name: u.Name,
		Dob:  d,
	}, nil
}

func (s *userService) Get(ctx context.Context, id int64) (models.User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return models.User{}, err
	}

	d := u.Dob.Time

	return models.User{
		ID:   int64(u.ID),
		Name: u.Name,
		Dob:  d,
		Age:  calculateAge(d, time.Now()),
	}, nil
}

func (s *userService) List(ctx context.Context) ([]models.User, error) {
	us, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	out := make([]models.User, 0, len(us))

	for _, u := range us {
		d := u.Dob.Time
		out = append(out, models.User{
			ID:   int64(u.ID),
			Name: u.Name,
			Dob:  d,
			Age:  calculateAge(d, now),
		})
	}

	return out, nil
}

func (s *userService) Update(ctx context.Context, id int64, req models.UpdateUserRequest) (models.User, error) {
	dob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		return models.User{}, err
	}

	u, err := s.repo.Update(ctx, id, req.Name, dob)
	if err != nil {
		return models.User{}, err
	}

	d := u.Dob.Time

	return models.User{
		ID:   int64(u.ID),
		Name: u.Name,
		Dob:  d,
	}, nil
}

func (s *userService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
