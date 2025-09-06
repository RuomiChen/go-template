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
	// 1. 获取用户所有新闻浏览记录
	events, err := s.repo.GetUserTrackingEvents(userID, "news view")
	if err != nil {
		return nil, err
	}

	// 2. 收集所有 newsID
	idSet := make(map[uint64]struct{})
	for _, e := range events {
		var extraData map[string]string
		if err := json.Unmarshal([]byte(e.Extra), &extraData); err != nil {
			continue
		}
		if idStr, ok := extraData["id"]; ok {
			if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
				idSet[id] = struct{}{}
			}
		}
	}

	// 3. 转换成切片
	ids := make([]uint64, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}

	// 4. 批量获取新闻详情
	newsList, err := s.newsSvc.GetNewsByIDs(ids)
	if err != nil {
		return nil, err
	}

	// 5. 构造 id -> *News map
	newsMap := make(map[uint64]*news.News)
	for _, n := range newsList {
		if n != nil {
			newsMap[n.ID] = n
		}
	}

	// 6. 拼装结果，填充 Data
	result := make([]TrackingWithNews, 0, len(events))
	for _, e := range events {
		var extraData map[string]string
		var newsID uint64
		if err := json.Unmarshal([]byte(e.Extra), &extraData); err == nil {
			if idStr, ok := extraData["id"]; ok {
				if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
					newsID = id
				}
			}
		}

		result = append(result, TrackingWithNews{
			TrackingEvent: e,
			Data:          newsMap[newsID], // 可能为 nil，如果找不到
		})
	}

	return result, nil
}
