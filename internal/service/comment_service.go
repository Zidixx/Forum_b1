package service

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type CommentService struct {
	commentRepo  *repository.CommentRepository
	reactionRepo *repository.ReactionRepository
}

func NewCommentService(commentRepo *repository.CommentRepository, reactionRepo *repository.ReactionRepository) *CommentService {
	return &CommentService{commentRepo: commentRepo, reactionRepo: reactionRepo}
}

func (s *CommentService) Create(comment *models.Comment) error {
	return s.commentRepo.Create(comment)
}

func (s *CommentService) GetByPostID(postID int, currentUserID int) ([]models.Comment, error) {
	comments, err := s.commentRepo.FindByPostID(postID)
	if err != nil {
		return nil, err
	}

	for i := range comments {
		comments[i].Likes, comments[i].Dislikes, _ = s.reactionRepo.CountCommentReactions(comments[i].ID)
		if currentUserID > 0 {
			comments[i].UserVote, _ = s.reactionRepo.GetCommentReaction(currentUserID, comments[i].ID)
		}
	}

	return comments, nil
}

func (s *CommentService) GetByID(id int) (*models.Comment, error) {
	return s.commentRepo.FindByID(id)
}

func (s *CommentService) Update(comment *models.Comment) error {
	return s.commentRepo.Update(comment)
}

func (s *CommentService) Delete(id int) error {
	return s.commentRepo.Delete(id)
}
