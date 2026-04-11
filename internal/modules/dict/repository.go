package dict

import (
	"context"
	"errors"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	ErrDictTypeNotFound = errors.New("dict type not found")
	ErrDictDataNotFound = errors.New("dict data not found")
	ErrDictDuplicated   = errors.New("dict duplicated")
)

type DictDataWithType struct {
	Data     DictData
	TypeCode string
}

type Repository interface {
	CreateType(ctx context.Context, item DictType) (DictType, error)
	GetTypeByID(ctx context.Context, id uint) (DictType, error)
	GetTypeByCode(ctx context.Context, code string) (DictType, error)
	ListTypes(ctx context.Context, q ListDictTypesQuery) ([]DictType, error)
	UpdateType(ctx context.Context, item DictType) (DictType, error)
	DeleteType(ctx context.Context, id uint) error
	CountTypeItems(ctx context.Context, typeID uint) (int64, error)

	CreateData(ctx context.Context, item DictData) (DictData, error)
	GetDataByID(ctx context.Context, id uint) (DictData, error)
	ListData(ctx context.Context, q ListDictDataQuery) ([]DictDataWithType, error)
	ListDataByTypeCode(ctx context.Context, typeCode string, onlyActive bool) ([]DictDataWithType, error)
	UpdateData(ctx context.Context, item DictData) (DictData, error)
	DeleteData(ctx context.Context, id uint) error
}

type inMemoryRepository struct {
	mu          sync.RWMutex
	types       map[uint]DictType
	data        map[uint]DictData
	typeCounter uint64
	dataCounter uint64
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		types: make(map[uint]DictType),
		data:  make(map[uint]DictData),
	}
}

func (r *inMemoryRepository) nextTypeID() uint {
	return uint(atomic.AddUint64(&r.typeCounter, 1))
}

func (r *inMemoryRepository) nextDataID() uint {
	return uint(atomic.AddUint64(&r.dataCounter, 1))
}

func (r *inMemoryRepository) CreateType(_ context.Context, item DictType) (DictType, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.types {
		if strings.EqualFold(existing.Code, item.Code) {
			return DictType{}, ErrDictDuplicated
		}
	}
	item.ID = r.nextTypeID()
	r.types[item.ID] = item
	return item, nil
}

func (r *inMemoryRepository) GetTypeByID(_ context.Context, id uint) (DictType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.types[id]
	if !ok {
		return DictType{}, ErrDictTypeNotFound
	}
	return item, nil
}

func (r *inMemoryRepository) GetTypeByCode(_ context.Context, code string) (DictType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, item := range r.types {
		if strings.EqualFold(item.Code, code) {
			return item, nil
		}
	}
	return DictType{}, ErrDictTypeNotFound
}

func (r *inMemoryRepository) ListTypes(_ context.Context, q ListDictTypesQuery) ([]DictType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]DictType, 0, len(r.types))
	for _, item := range r.types {
		if q.Name != "" && !strings.Contains(strings.ToLower(item.Name), strings.ToLower(q.Name)) {
			continue
		}
		if q.Code != "" && !strings.Contains(strings.ToLower(item.Code), strings.ToLower(q.Code)) {
			continue
		}
		items = append(items, item)
	}
	slices.SortFunc(items, func(a, b DictType) int { return int(a.ID) - int(b.ID) })
	return items, nil
}

func (r *inMemoryRepository) UpdateType(_ context.Context, item DictType) (DictType, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.types[item.ID]; !ok {
		return DictType{}, ErrDictTypeNotFound
	}
	for id, existing := range r.types {
		if id != item.ID && strings.EqualFold(existing.Code, item.Code) {
			return DictType{}, ErrDictDuplicated
		}
	}
	r.types[item.ID] = item
	return item, nil
}

func (r *inMemoryRepository) DeleteType(_ context.Context, id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.types[id]; !ok {
		return ErrDictTypeNotFound
	}
	delete(r.types, id)
	return nil
}

func (r *inMemoryRepository) CountTypeItems(_ context.Context, typeID uint) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var total int64
	for _, item := range r.data {
		if item.TypeID == typeID {
			total++
		}
	}
	return total, nil
}

func (r *inMemoryRepository) CreateData(_ context.Context, item DictData) (DictData, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.types[item.TypeID]; !ok {
		return DictData{}, ErrDictTypeNotFound
	}
	for _, existing := range r.data {
		if existing.TypeID == item.TypeID && strings.EqualFold(existing.Label, item.Label) {
			return DictData{}, ErrDictDuplicated
		}
	}
	item.ID = r.nextDataID()
	r.data[item.ID] = item
	return item, nil
}

func (r *inMemoryRepository) GetDataByID(_ context.Context, id uint) (DictData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.data[id]
	if !ok {
		return DictData{}, ErrDictDataNotFound
	}
	return item, nil
}

func (r *inMemoryRepository) ListData(_ context.Context, q ListDictDataQuery) ([]DictDataWithType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]DictDataWithType, 0, len(r.data))
	for _, item := range r.data {
		typ, ok := r.types[item.TypeID]
		if !ok {
			continue
		}
		if q.TypeID != nil && item.TypeID != *q.TypeID {
			continue
		}
		if q.TypeCode != "" && !strings.EqualFold(typ.Code, q.TypeCode) {
			continue
		}
		if q.Label != "" && !strings.Contains(strings.ToLower(item.Label), strings.ToLower(q.Label)) {
			continue
		}
		if q.Status != nil && item.Status != *q.Status {
			continue
		}
		result = append(result, DictDataWithType{Data: item, TypeCode: typ.Code})
	}
	slices.SortFunc(result, func(a, b DictDataWithType) int {
		if a.Data.Sort != b.Data.Sort {
			return a.Data.Sort - b.Data.Sort
		}
		return int(a.Data.ID) - int(b.Data.ID)
	})
	return result, nil
}

func (r *inMemoryRepository) ListDataByTypeCode(ctx context.Context, typeCode string, onlyActive bool) ([]DictDataWithType, error) {
	items, err := r.ListData(ctx, ListDictDataQuery{TypeCode: typeCode})
	if err != nil {
		return nil, err
	}
	if !onlyActive {
		return items, nil
	}
	filtered := make([]DictDataWithType, 0, len(items))
	for _, item := range items {
		if item.Data.Status == 1 {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

func (r *inMemoryRepository) UpdateData(_ context.Context, item DictData) (DictData, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[item.ID]; !ok {
		return DictData{}, ErrDictDataNotFound
	}
	if _, ok := r.types[item.TypeID]; !ok {
		return DictData{}, ErrDictTypeNotFound
	}
	for id, existing := range r.data {
		if id != item.ID && existing.TypeID == item.TypeID && strings.EqualFold(existing.Label, item.Label) {
			return DictData{}, ErrDictDuplicated
		}
	}
	r.data[item.ID] = item
	return item, nil
}

func (r *inMemoryRepository) DeleteData(_ context.Context, id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return ErrDictDataNotFound
	}
	delete(r.data, id)
	return nil
}
