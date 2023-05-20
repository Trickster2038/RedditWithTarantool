package handlers

import "redditclone/pkg/post"

func (h *PostsHandler) getPostsDependentData(ps []post.Post) error {
	for i := range ps {
		err := h.getPostDependentData(&ps[i])
		if err != nil {
			return err
		}
	}
	post.SortByScoreDesc(ps)
	return nil
}

func (h *PostsHandler) getPostDependentData(ps *post.Post) error {
	cms, err := h.CommentRepo.GetAllCommentsForPost(ps.ID)
	if err != nil {
		return err
	}
	ps.Comments = cms
	vts, err := h.VoteRepo.GetAllVotesForPost(ps.ID)
	if err != nil {
		return err
	}
	ps.Votes = vts
	ps.UpdateStats()
	return nil
}
