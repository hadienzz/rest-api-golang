package follow

type FollowService interface {
	FollowMerchant(request *FollowRequest) (*FollowResponse, error)
	UnfollowMerchant(request *FollowRequest) (*FollowResponse, error)
	GetMerchantFollowStatus(request *FollowRequest) (*FollowResponse, error)
}

type followService struct {
	repo FollowersRepository
}

func NewFollowService(repo FollowersRepository) FollowService {
	return &followService{
		repo: repo,
	}
}

func (s *followService) FollowMerchant(request *FollowRequest) (*FollowResponse, error) {

	follow := &Follow{
		UserID:     request.UserID,
		MerchantID: request.MerchantID,
	}

	createdFollow, err := s.repo.AddFollower(follow)

	if err != nil {
		return nil, err
	}

	return &FollowResponse{
		IsFollowing: true,
		MerchantID:  createdFollow.MerchantID,
		FollowedAt:  createdFollow.CreatedAt.Time,
	}, nil
}

func (s *followService) UnfollowMerchant(request *FollowRequest) (*FollowResponse, error) {
	unfollow := &Follow{
		UserID:     request.UserID,
		MerchantID: request.MerchantID,
	}
	createdUnfollow, err := s.repo.UnfollowMerchant(unfollow)

	if err != nil {
		return nil, err
	}

	return &FollowResponse{
		IsFollowing: false,
		MerchantID:  createdUnfollow.MerchantID,
	}, nil
}

func (s *followService) GetMerchantFollowStatus(request *FollowRequest) (*FollowResponse, error) {
	follow, err := s.repo.GetFollowStatus(&Follow{
		UserID:     request.UserID,
		MerchantID: request.MerchantID,
	})

	if err != nil {
		return nil, err
	}

	return &FollowResponse{
		IsFollowing: follow != nil,
		MerchantID:  request.MerchantID,
	}, nil
}
