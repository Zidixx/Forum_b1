package service

import (
	"fmt"
	"forum/internal/repository"
)

type ReactionService struct {
	reactionRepo *repository.ReactionRepository
}

func NewReactionService(reactionRepo *repository.ReactionRepository) *ReactionService {
	return &ReactionService{reactionRepo: reactionRepo}
}

func (s *ReactionService) ReactToPost(userID, postID int, reactionType string) error {
	if reactionType != "like" && reactionType != "dislike" {
		return fmt.Errorf("type de réaction invalide")
	}
	return s.reactionRepo.SetPostReaction(userID, postID, reactionType)
}

func (s *ReactionService) ReactToComment(userID, commentID int, reactionType string) error {
	if reactionType != "like" && reactionType != "dislike" {
		return fmt.Errorf("type de réaction invalide")
	}
	return s.reactionRepo.SetCommentReaction(userID, commentID, reactionType)
}

func (s *ReactionService) CountPostReactions(postID int) (int, int, error) {
	return s.reactionRepo.CountPostReactions(postID)
}

func (s *ReactionService) GetPostReaction(userID, postID int) (string, error) {
	return s.reactionRepo.GetPostReaction(userID, postID)
}

func (s *ReactionService) CountCommentReactions(commentID int) (int, int, error) {
	return s.reactionRepo.CountCommentReactions(commentID)
}

func (s *ReactionService) GetCommentReaction(userID, commentID int) (string, error) {
	return s.reactionRepo.GetCommentReaction(userID, commentID)
}
