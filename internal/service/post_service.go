package service

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type PostService struct {
	postRepo     *repository.PostRepository
	catRepo      *repository.CategoryRepository
	reactionRepo *repository.ReactionRepository
}

func NewPostService(postRepo *repository.PostRepository, catRepo *repository.CategoryRepository, reactionRepo *repository.ReactionRepository) *PostService {
	return &PostService{postRepo: postRepo, catRepo: catRepo, reactionRepo: reactionRepo}
}

func (s *PostService) Create(post *models.Post, categoryIDs []int) error {
	return s.postRepo.Create(post, categoryIDs)
}

func (s *PostService) GetByID(id int, currentUserID int) (*models.Post, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	s.enrichPost(post, currentUserID)
	return post, nil
}

func (s *PostService) GetAll(currentUserID int) ([]models.Post, error) {
	posts, err := s.postRepo.FindAll()
	if err != nil {
		return nil, err
	}
	s.enrichPosts(posts, currentUserID)
	return posts, nil
}

func (s *PostService) GetByCategory(categoryID int, currentUserID int) ([]models.Post, error) {
	posts, err := s.postRepo.FindByCategory(categoryID)
	if err != nil {
		return nil, err
	}
	s.enrichPosts(posts, currentUserID)
	return posts, nil
}

func (s *PostService) GetByUser(userID int) ([]models.Post, error) {
	posts, err := s.postRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	s.enrichPosts(posts, userID)
	return posts, nil
}

func (s *PostService) GetLikedByUser(userID int) ([]models.Post, error) {
	posts, err := s.postRepo.FindLikedByUserID(userID)
	if err != nil {
		return nil, err
	}
	s.enrichPosts(posts, userID)
	return posts, nil
}

func (s *PostService) Update(post *models.Post, categoryIDs []int) error {
	return s.postRepo.Update(post, categoryIDs)
}

func (s *PostService) Delete(id int) error {
	return s.postRepo.Delete(id)
}

func (s *PostService) GetCategoryIDs(postID int) ([]int, error) {
	return s.postRepo.GetCategoryIDs(postID)
}

func (s *PostService) Excerpt(content string) string {
	return s.postRepo.Excerpt(content, 200)
}

func (s *PostService) enrichPosts(posts []models.Post, currentUserID int) {
	for i := range posts {
		s.enrichPost(&posts[i], currentUserID)
	}
}

func (s *PostService) enrichPost(post *models.Post, currentUserID int) {
	post.Likes, post.Dislikes, _ = s.reactionRepo.CountPostReactions(post.ID)
	cats, _ := s.catRepo.FindByPostID(post.ID)
	post.Categories = cats

	if currentUserID > 0 {
		post.UserVote, _ = s.reactionRepo.GetPostReaction(currentUserID, post.ID)
	}
}
