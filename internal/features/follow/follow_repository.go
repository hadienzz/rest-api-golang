package follow

import "gorm.io/gorm"

type FollowersRepository interface {
	AddFollower(follow *Follow) (*Follow, error)
	UnfollowMerchant(follow *Follow) (*Follow, error)
	GetFollowStatus(follow *Follow) (*Follow, error)
}

type followerRepository struct {
	db *gorm.DB
}

func NewFollowersRepository(db *gorm.DB) FollowersRepository {
	return &followerRepository{
		db: db,
	}
}

func (r *followerRepository) AddFollower(follow *Follow) (*Follow, error) {
	if err := r.db.Create(follow).Error; err != nil {
		return nil, err
	}

	return follow, nil
}

func (r *followerRepository) UnfollowMerchant(follow *Follow) (*Follow, error) {
	result := r.db.
		Where("merchant_id = ? AND user_id = ?", follow.MerchantID, follow.UserID).
		Delete(&Follow{})

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return follow, nil
}

func (r *followerRepository) GetFollowStatus(follow *Follow) (*Follow, error) {
	var req Follow
	result := r.db.Where("merchant_id = ? AND user_id = ?", follow.MerchantID, follow.UserID).Find(&req)

	if result.Error != nil {
		return nil, result.Error
	}

	return &req, nil
}
