package domain

import (
	"context"
	"douyin-app/repository"
)

type VideoDO struct {
	Id            int64   `json:"id"`
	Author        *Author `json:"author"`
	PlayUrl       string  `json:"play_url"`
	CoverUrl      string  `json:"cover_url"`
	FavoriteCount int64   `json:"favorite_count"`
	CommentCount  int64   `json:"comment_count"`
	IsFavorite    bool    `json:"is_favorite"`
}

func FillVideoList(ctx context.Context, videoList []*repository.Video, audienceId int64) ([]*VideoDO, error) {
	res := make([]*VideoDO, len(videoList))

	authorIds := GetUserIdsFromVideoList(videoList)
	authors, err := GetAuthorsFromIds(ctx, authorIds)
	if err != nil {
		return nil, err
	}

	// TODO audienceId 用来查关注和点赞

	for i := range videoList {
		res[i] = &VideoDO{
			Id:            videoList[i].VideoId,
			Author:        authors[i],
			PlayUrl:       videoList[i].PlayUrl,
			CoverUrl:      videoList[i].CoverUrl,
			FavoriteCount: videoList[i].FavoriteCount, //TODO
			CommentCount:  videoList[i].CommentCount,  //TODO
			IsFavorite:    false,                      //TODO
		}
	}

	return res, nil
}

func GetUserIdsFromVideoList(videoList []*repository.Video) []int64 {
	ids := make([]int64, len(videoList))
	for i := range videoList {
		ids[i] = videoList[i].UserId
	}
	return ids
}

func GetAuthorsFromIds(ctx context.Context, ids []int64) ([]*Author, error) {
	users, err := repository.GetUserRepository().FindByUserIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	authors := make([]*Author, len(ids))
	for i := range users {
		authors[i] = FillAuthor(users[i])
	}

	return authors, nil
}
