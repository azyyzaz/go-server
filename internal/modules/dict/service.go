package dict

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"go-server/internal/errcode"

	rdb "github.com/redis/go-redis/v9"
)

const cacheKeyPrefix = "dict:lookup:"

type Service interface {
	ListTypes(ctx context.Context, q ListDictTypesQuery) ([]DictTypeResponse, error)
	GetType(ctx context.Context, id uint) (DictTypeResponse, error)
	CreateType(ctx context.Context, req CreateDictTypeRequest) (DictTypeResponse, error)
	UpdateType(ctx context.Context, id uint, req UpdateDictTypeRequest) (DictTypeResponse, error)
	DeleteType(ctx context.Context, id uint) error

	ListData(ctx context.Context, q ListDictDataQuery) ([]DictDataResponse, error)
	GetData(ctx context.Context, id uint) (DictDataResponse, error)
	CreateData(ctx context.Context, req CreateDictDataRequest) (DictDataResponse, error)
	UpdateData(ctx context.Context, id uint, req UpdateDictDataRequest) (DictDataResponse, error)
	DeleteData(ctx context.Context, id uint) error
	LookupByTypeCode(ctx context.Context, typeCode string) ([]DictDataResponse, error)
}

type service struct {
	repo  Repository
	redis *rdb.Client
}

func NewService(repo Repository, redisClient *rdb.Client) Service {
	return &service{repo: repo, redis: redisClient}
}

func (s *service) ListTypes(ctx context.Context, q ListDictTypesQuery) ([]DictTypeResponse, error) {
	items, err := s.repo.ListTypes(ctx, q)
	if err != nil {
		return nil, err
	}
	result := make([]DictTypeResponse, 0, len(items))
	for _, item := range items {
		result = append(result, toTypeResponse(item))
	}
	return result, nil
}

func (s *service) GetType(ctx context.Context, id uint) (DictTypeResponse, error) {
	item, err := s.repo.GetTypeByID(ctx, id)
	if err != nil {
		if err == ErrDictTypeNotFound {
			return DictTypeResponse{}, errcode.ErrDictTypeNotFound.AsError()
		}
		return DictTypeResponse{}, err
	}
	return toTypeResponse(item), nil
}

func (s *service) CreateType(ctx context.Context, req CreateDictTypeRequest) (DictTypeResponse, error) {
	item := DictType{
		Name:      strings.TrimSpace(req.Name),
		Code:      strings.TrimSpace(req.Code),
		Remark:    strings.TrimSpace(req.Remark),
		Status:    req.Status,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	created, err := s.repo.CreateType(ctx, item)
	if err != nil {
		if err == ErrDictDuplicated {
			return DictTypeResponse{}, errcode.ErrDictTypeCodeExists.AsError()
		}
		return DictTypeResponse{}, err
	}
	return toTypeResponse(created), nil
}

func (s *service) UpdateType(ctx context.Context, id uint, req UpdateDictTypeRequest) (DictTypeResponse, error) {
	item, err := s.repo.GetTypeByID(ctx, id)
	if err != nil {
		if err == ErrDictTypeNotFound {
			return DictTypeResponse{}, errcode.ErrDictTypeNotFound.AsError()
		}
		return DictTypeResponse{}, err
	}
	oldCode := item.Code

	item.Name = strings.TrimSpace(req.Name)
	item.Code = strings.TrimSpace(req.Code)
	item.Remark = strings.TrimSpace(req.Remark)
	item.Status = req.Status
	item.UpdatedAt = time.Now().UTC()

	updated, err := s.repo.UpdateType(ctx, item)
	if err != nil {
		if err == ErrDictDuplicated {
			return DictTypeResponse{}, errcode.ErrDictTypeCodeExists.AsError()
		}
		return DictTypeResponse{}, err
	}
	s.invalidateCache(ctx, oldCode)
	s.invalidateCache(ctx, updated.Code)
	return toTypeResponse(updated), nil
}

func (s *service) DeleteType(ctx context.Context, id uint) error {
	item, err := s.repo.GetTypeByID(ctx, id)
	if err != nil {
		if err == ErrDictTypeNotFound {
			return errcode.ErrDictTypeNotFound.AsError()
		}
		return err
	}

	total, err := s.repo.CountTypeItems(ctx, id)
	if err != nil {
		return err
	}
	if total > 0 {
		return errcode.ErrDictTypeHasItems.AsError()
	}

	if err := s.repo.DeleteType(ctx, id); err != nil {
		if err == ErrDictTypeNotFound {
			return errcode.ErrDictTypeNotFound.AsError()
		}
		return err
	}
	s.invalidateCache(ctx, item.Code)
	return nil
}

func (s *service) ListData(ctx context.Context, q ListDictDataQuery) ([]DictDataResponse, error) {
	items, err := s.repo.ListData(ctx, q)
	if err != nil {
		return nil, err
	}
	return toDataResponses(items), nil
}

func (s *service) GetData(ctx context.Context, id uint) (DictDataResponse, error) {
	item, err := s.repo.GetDataByID(ctx, id)
	if err != nil {
		if err == ErrDictDataNotFound {
			return DictDataResponse{}, errcode.ErrDictDataNotFound.AsError()
		}
		return DictDataResponse{}, err
	}
	dictType, err := s.repo.GetTypeByID(ctx, item.TypeID)
	if err != nil {
		if err == ErrDictTypeNotFound {
			return DictDataResponse{}, errcode.ErrDictTypeNotFound.AsError()
		}
		return DictDataResponse{}, err
	}
	return DictDataResponse{
		ID:        item.ID,
		TypeID:    item.TypeID,
		TypeCode:  dictType.Code,
		Label:     item.Label,
		Value:     item.Value,
		Sort:      item.Sort,
		Status:    item.Status,
		Remark:    item.Remark,
		CreatedAt: item.CreatedAt,
	}, nil
}

func (s *service) CreateData(ctx context.Context, req CreateDictDataRequest) (DictDataResponse, error) {
	dictType, err := s.repo.GetTypeByID(ctx, req.TypeID)
	if err != nil {
		if err == ErrDictTypeNotFound {
			return DictDataResponse{}, errcode.ErrDictTypeNotFound.AsError()
		}
		return DictDataResponse{}, err
	}

	item := DictData{
		TypeID:    req.TypeID,
		Label:     strings.TrimSpace(req.Label),
		Value:     strings.TrimSpace(req.Value),
		Sort:      req.Sort,
		Status:    req.Status,
		Remark:    strings.TrimSpace(req.Remark),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	created, err := s.repo.CreateData(ctx, item)
	if err != nil {
		if err == ErrDictDuplicated {
			return DictDataResponse{}, errcode.ErrDictDataLabelExists.AsError()
		}
		if err == ErrDictTypeNotFound {
			return DictDataResponse{}, errcode.ErrDictTypeNotFound.AsError()
		}
		return DictDataResponse{}, err
	}
	s.invalidateCache(ctx, dictType.Code)
	return DictDataResponse{
		ID:        created.ID,
		TypeID:    created.TypeID,
		TypeCode:  dictType.Code,
		Label:     created.Label,
		Value:     created.Value,
		Sort:      created.Sort,
		Status:    created.Status,
		Remark:    created.Remark,
		CreatedAt: created.CreatedAt,
	}, nil
}

func (s *service) UpdateData(ctx context.Context, id uint, req UpdateDictDataRequest) (DictDataResponse, error) {
	item, err := s.repo.GetDataByID(ctx, id)
	if err != nil {
		if err == ErrDictDataNotFound {
			return DictDataResponse{}, errcode.ErrDictDataNotFound.AsError()
		}
		return DictDataResponse{}, err
	}

	oldType, err := s.repo.GetTypeByID(ctx, item.TypeID)
	if err != nil {
		if err == ErrDictTypeNotFound {
			return DictDataResponse{}, errcode.ErrDictTypeNotFound.AsError()
		}
		return DictDataResponse{}, err
	}
	newType, err := s.repo.GetTypeByID(ctx, req.TypeID)
	if err != nil {
		if err == ErrDictTypeNotFound {
			return DictDataResponse{}, errcode.ErrDictTypeNotFound.AsError()
		}
		return DictDataResponse{}, err
	}

	item.TypeID = req.TypeID
	item.Label = strings.TrimSpace(req.Label)
	item.Value = strings.TrimSpace(req.Value)
	item.Sort = req.Sort
	item.Status = req.Status
	item.Remark = strings.TrimSpace(req.Remark)
	item.UpdatedAt = time.Now().UTC()

	updated, err := s.repo.UpdateData(ctx, item)
	if err != nil {
		if err == ErrDictDuplicated {
			return DictDataResponse{}, errcode.ErrDictDataLabelExists.AsError()
		}
		if err == ErrDictTypeNotFound {
			return DictDataResponse{}, errcode.ErrDictTypeNotFound.AsError()
		}
		return DictDataResponse{}, err
	}
	s.invalidateCache(ctx, oldType.Code)
	s.invalidateCache(ctx, newType.Code)
	return DictDataResponse{
		ID:        updated.ID,
		TypeID:    updated.TypeID,
		TypeCode:  newType.Code,
		Label:     updated.Label,
		Value:     updated.Value,
		Sort:      updated.Sort,
		Status:    updated.Status,
		Remark:    updated.Remark,
		CreatedAt: updated.CreatedAt,
	}, nil
}

func (s *service) DeleteData(ctx context.Context, id uint) error {
	item, err := s.repo.GetDataByID(ctx, id)
	if err != nil {
		if err == ErrDictDataNotFound {
			return errcode.ErrDictDataNotFound.AsError()
		}
		return err
	}
	dictType, err := s.repo.GetTypeByID(ctx, item.TypeID)
	if err != nil {
		if err == ErrDictTypeNotFound {
			return errcode.ErrDictTypeNotFound.AsError()
		}
		return err
	}
	if err := s.repo.DeleteData(ctx, id); err != nil {
		if err == ErrDictDataNotFound {
			return errcode.ErrDictDataNotFound.AsError()
		}
		return err
	}
	s.invalidateCache(ctx, dictType.Code)
	return nil
}

func (s *service) LookupByTypeCode(ctx context.Context, typeCode string) ([]DictDataResponse, error) {
	typeCode = strings.TrimSpace(typeCode)
	if typeCode == "" {
		return nil, errcode.ErrInvalidParam.AsError()
	}

	if s.redis != nil {
		raw, err := s.redis.Get(ctx, cacheKeyPrefix+typeCode).Result()
		if err == nil {
			var cached []DictDataResponse
			if json.Unmarshal([]byte(raw), &cached) == nil {
				return cached, nil
			}
		}
	}

	if _, err := s.repo.GetTypeByCode(ctx, typeCode); err != nil {
		if err == ErrDictTypeNotFound {
			return nil, errcode.ErrDictTypeNotFound.AsError()
		}
		return nil, err
	}

	items, err := s.repo.ListDataByTypeCode(ctx, typeCode, true)
	if err != nil {
		return nil, err
	}
	result := toDataResponses(items)

	if s.redis != nil {
		if payload, err := json.Marshal(result); err == nil {
			_ = s.redis.Set(ctx, cacheKeyPrefix+typeCode, payload, time.Hour).Err()
		}
	}
	return result, nil
}

func (s *service) invalidateCache(ctx context.Context, typeCode string) {
	if s.redis == nil || strings.TrimSpace(typeCode) == "" {
		return
	}
	_ = s.redis.Del(ctx, cacheKeyPrefix+typeCode).Err()
}

func toDataResponses(items []DictDataWithType) []DictDataResponse {
	result := make([]DictDataResponse, 0, len(items))
	for _, item := range items {
		result = append(result, DictDataResponse{
			ID:        item.Data.ID,
			TypeID:    item.Data.TypeID,
			TypeCode:  item.TypeCode,
			Label:     item.Data.Label,
			Value:     item.Data.Value,
			Sort:      item.Data.Sort,
			Status:    item.Data.Status,
			Remark:    item.Data.Remark,
			CreatedAt: item.Data.CreatedAt,
		})
	}
	return result
}
