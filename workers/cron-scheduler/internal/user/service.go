package user

type Service interface {
	GetUsers() (*[]GetUsersResponse, error)
}

type userService struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &userService{repo: r}
}

func (s *userService) GetUsers() (*[]GetUsersResponse, error) {
	var result []GetUsersResponse
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		result = append(result, GetUsersResponse{
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		})
	}
	return &result, nil
}
