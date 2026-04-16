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
