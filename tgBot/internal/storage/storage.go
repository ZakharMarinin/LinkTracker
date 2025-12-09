package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"linktracker/internal/domain"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type TempUserLinks struct {
	UserID int64          `json:"user_id"`
	Links  []*domain.Link `json:"links"`
}

type RedisCom struct {
	DB  *redis.Client
	log *slog.Logger
}

func NewRedisCom(db *redis.Client, log *slog.Logger) *RedisCom {
	return &RedisCom{db, log}
}

func (r *RedisCom) SetTempUserState(ctx context.Context, userInfo *domain.UserStateInfo) error {
	err := r.DB.HSet(ctx, fmt.Sprintf("TempUserState:%s", userInfo.UserID), map[string]any{
		"UserID": userInfo.UserID,
		"URL":    userInfo.URL,
		"Desc":   userInfo.Desc,
		"State":  userInfo.State,
	}).Err()
	if err != nil {
		r.log.Error("SetUser", "err", err)
		return err
	}

	err = r.DB.Expire(ctx, fmt.Sprintf("TempUserState:%s", userInfo.UserID), 2*time.Hour).Err()
	if err != nil {
		r.log.Error("SetUser", "err", err)
		return err
	}

	return nil
}

func (r *RedisCom) GetTempUserState(ctx context.Context, userID int64) (*domain.UserStateInfo, error) {
	val, err := r.DB.HGetAll(ctx, fmt.Sprintf("TempUserState:%s", userID)).Result()
	if err != nil {
		r.log.Error("GetTempUserState", "err", err)
		return nil, err
	}

	tempUserState := &domain.UserStateInfo{
		UserID: userID,
		URL:    val["URL"],
		Desc:   val["Desc"],
		State:  val["State"],
	}

	return tempUserState, nil
}

func (r *RedisCom) SaveTempUserLinks(ctx context.Context, tempUserLinks *TempUserLinks) error {
	links, err := json.Marshal(tempUserLinks.Links)
	if err != nil {
		r.log.Error("SaveTempUserLinks", "err", err)
		return err
	}
	err = r.DB.HSet(ctx, fmt.Sprintf("TempUserLinks:%s", tempUserLinks.UserID), map[string]any{
		"user_id": tempUserLinks.UserID,
		"links":   links,
	}).Err()
	if err != nil {
		r.log.Error("SaveTempUserLinks", "err", err)
		return err
	}

	err = r.DB.Expire(ctx, fmt.Sprintf("TempUserLinks:%s", tempUserLinks.UserID), 2*time.Hour).Err()
	if err != nil {
		r.log.Error("SetUser", "err", err)
		return err
	}

	return nil
}

func (r *RedisCom) GetTempUserLinks(ctx context.Context, userID int64) (*TempUserLinks, error) {
	val, err := r.DB.HGetAll(ctx, fmt.Sprintf("TempUserLinks:%s", userID)).Result()
	if err != nil {
		r.log.Error("GetTempUserLinks: error HGetAll", "err", err)
		return nil, err
	}

	var links []*domain.Link
	err = json.Unmarshal([]byte(val["links"]), &links)
	if err != nil {
		r.log.Error("GetTempUserLinks: Cannot parse json", "err", err)
		return nil, err
	}

	tempUserLinks := &TempUserLinks{
		UserID: userID,
		Links:  links,
	}

	return tempUserLinks, nil
}
