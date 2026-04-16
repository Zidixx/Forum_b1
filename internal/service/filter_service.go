package service

import "forum/internal/models"

type FilterService struct {
	postService *PostService
}

func NewFilterService(postService *PostService) *FilterService {
	return &FilterService{postService: postService}
}

func (s *FilterService) FilterByCategory(categoryID int, userID int) ([]models.Post, error) {
	return s.postService.GetByCategory(categoryID, userID)
}

func (s *FilterService) FilterByUser(userID int) ([]models.Post, error) {
	return s.postService.GetByUser(userID)
}

func (s *FilterService) FilterByLiked(userID int) ([]models.Post, error) {
	return s.postService.GetLikedByUser(userID)
}
