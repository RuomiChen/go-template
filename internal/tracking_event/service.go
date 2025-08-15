package tracking_event

import (
	"encoding/json"
	"mvc/internal/news"
	"mvc/internal/redis"
	"mvc/pkg/utils"
	"strconv"
	"time"
)

type Service interface {
	Record(trackingEvent *TrackingEvent) error
	GetUserTrackingEvents(userID, action string) ([]TrackingEvent, error)
	GetUserTrackingEventsByAction(userID, action string) ([]TrackingWithNews, error)
}

type service struct {
	repo      Repository
	newsSvc   news.Service
	hashStore *utils.RedisHashStore
}

func NewService(repo Repository, newsSvc news.Service, redisService redis.Service) Service {
	hashStore := utils.NewRedisHashStore(redisService, "imghash:", time.Hour*24*7)
	return &service{repo: repo, newsSvc: newsSvc, hashStore: hashStore}
}

func (s *service) Record(trackingEvent *TrackingEvent) error {
	return s.repo.addTrackingEvent(trackingEvent)
}
func (s *service) GetUserTrackingEvents(userID, action string) ([]TrackingEvent, error) {
	return s.repo.GetUserTrackingEvents(userID, action)
}

type TrackingWithNews struct {
	TrackingEvent
	Data *news.News `json:"data,omitempty"`
}

func (s *service) GetUserTrackingEventsByAction(userID, action string) ([]TrackingWithNews, error) {
	if action == "news view" {
		return s.GetUserNewsRecordsWithData(userID)
	}

	// 其他 action 直接返回普通记录
	events, err := s.repo.GetUserTrackingEvents(userID, action)
	if err != nil {
		return nil, err
	}

	// 转换为 TrackingWithNews，Data 为 nil
	result := make([]TrackingWithNews, len(events))
	for i, e := range events {
		result[i] = TrackingWithNews{TrackingEvent: e}
	}
	return result, nil
}
func (s *service) GetUserNewsRecordsWithData(userID string) ([]TrackingWithNews, error) {
	// 1. 获取用户所有新闻记录
	events, err := s.repo.GetUserTrackingEvents(userID, "news view")
	if err != nil {
		return nil, err
	}

	// 2. 从每条记录解析 news id
	idSet := make(map[string]struct{})
	for _, e := range events {
		var extraData map[string]string
		if err := json.Unmarshal([]byte(e.Extra), &extraData); err != nil {
			continue
		}
		if id, ok := extraData["id"]; ok {
			idSet[id] = struct{}{}
		}
	}

	// 3. 构造 id 列表
	ids := make([]string, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}

	// 4. 批量获取新闻详情
	newsList, err := s.newsSvc.GetNewsByIDs(ids)
	if err != nil {
		return nil, err
	}

	// 5. 构造 id -> news map
	newsMap := make(map[string]*news.News)
	for i := range newsList {
		newsMap[strconv.FormatUint(newsList[i].ID, 10)] = &newsList[i]
	}

	// 6. 把 news 放入每条记录
	result := make([]TrackingWithNews, 0, len(events))
	for _, e := range events {
		var extraData map[string]string
		var newsID string
		if err := json.Unmarshal([]byte(e.Extra), &extraData); err == nil {
			newsID = extraData["id"]
		}

		result = append(result, TrackingWithNews{
			TrackingEvent: e,
			Data:          newsMap[newsID],
		})
	}

	return result, nil
}
