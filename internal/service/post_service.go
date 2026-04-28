package service

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type PostService struct {
	postRepo     *repository.PostRepository
	catRepo      *repository.CategoryRepository
	reactionRepo *repository.ReactionRepository
	repostRepo   *repository.RepostRepository
}

func NewPostService(postRepo *repository.PostRepository, catRepo *repository.CategoryRepository, reactionRepo *repository.ReactionRepository, repostRepo *repository.RepostRepository) *PostService {
	return &PostService{postRepo: postRepo, catRepo: catRepo, reactionRepo: reactionRepo, repostRepo: repostRepo}
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

func (s *PostService) GetAllSorted(sort string, currentUserID int) ([]models.Post, error) {
	posts, err := s.postRepo.FindAllSorted(sort)
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

func (s *PostService) GetTrending(limit int, currentUserID int) ([]models.Post, error) {
	posts, err := s.postRepo.FindTrending(limit)
	if err != nil {
		return nil, err
	}
	s.enrichPosts(posts, currentUserID)
	return posts, nil
}

func (s *PostService) Search(query string, currentUserID int) ([]models.Post, error) {
	posts, err := s.postRepo.Search(query)
	if err != nil {
		return nil, err
	}
	s.enrichPosts(posts, currentUserID)
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
	post.CommentCount, _ = s.postRepo.CountComments(post.ID)

	if s.repostRepo != nil {
		post.Reposts, _ = s.repostRepo.Count(post.ID)
	}

	if currentUserID > 0 {
		post.UserVote, _ = s.reactionRepo.GetPostReaction(currentUserID, post.ID)
		if s.repostRepo != nil {
			post.UserReposted, _ = s.repostRepo.HasReposted(currentUserID, post.ID)
		}
	}
}
