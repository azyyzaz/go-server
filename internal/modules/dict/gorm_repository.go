package dict

import (
	"context"
	"strings"

	"gorm.io/gorm"
)

func isDuplicateError(err error) bool {
	return strings.Contains(err.Error(), "Duplicate entry") ||
		strings.Contains(err.Error(), "1062") ||
		strings.Contains(strings.ToLower(err.Error()), "unique constraint")
}

type gormRepository struct {
	db *gorm.DB
}

func NewGORMRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateType(ctx context.Context, item DictType) (DictType, error) {
	if err := r.db.WithContext(ctx).Create(&item).Error; err != nil {
		if isDuplicateError(err) {
			return DictType{}, ErrDictDuplicated
		}
		return DictType{}, err
	}
	return item, nil
}

func (r *gormRepository) GetTypeByID(ctx context.Context, id uint) (DictType, error) {
	var item DictType
	err := r.db.WithContext(ctx).First(&item, id).Error
	if err == gorm.ErrRecordNotFound {
		return DictType{}, ErrDictTypeNotFound
	}
	return item, err
}

func (r *gormRepository) GetTypeByCode(ctx context.Context, code string) (DictType, error) {
	var item DictType
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&item).Error
	if err == gorm.ErrRecordNotFound {
		return DictType{}, ErrDictTypeNotFound
	}
	return item, err
}

func (r *gormRepository) ListTypes(ctx context.Context, q ListDictTypesQuery) ([]DictType, error) {
	var items []DictType
	db := r.db.WithContext(ctx).Model(&DictType{})
	if q.Name != "" {
		db = db.Where("name LIKE ?", "%"+q.Name+"%")
	}
	if q.Code != "" {
		db = db.Where("code LIKE ?", "%"+q.Code+"%")
	}
	err := db.Order("id ASC").Find(&items).Error
	return items, err
}

func (r *gormRepository) UpdateType(ctx context.Context, item DictType) (DictType, error) {
	if err := r.db.WithContext(ctx).Save(&item).Error; err != nil {
		if isDuplicateError(err) {
			return DictType{}, ErrDictDuplicated
		}
		return DictType{}, err
	}
	return item, nil
}

func (r *gormRepository) DeleteType(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&DictType{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrDictTypeNotFound
	}
	return nil
}

func (r *gormRepository) CountTypeItems(ctx context.Context, typeID uint) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&DictData{}).Where("type_id = ?", typeID).Count(&total).Error
	return total, err
}

func (r *gormRepository) CreateData(ctx context.Context, item DictData) (DictData, error) {
	if err := r.db.WithContext(ctx).Create(&item).Error; err != nil {
		if isDuplicateError(err) {
			return DictData{}, ErrDictDuplicated
		}
		return DictData{}, err
	}
	return item, nil
}

func (r *gormRepository) GetDataByID(ctx context.Context, id uint) (DictData, error) {
	var item DictData
	err := r.db.WithContext(ctx).First(&item, id).Error
	if err == gorm.ErrRecordNotFound {
		return DictData{}, ErrDictDataNotFound
	}
	return item, err
}

func (r *gormRepository) ListData(ctx context.Context, q ListDictDataQuery) ([]DictDataWithType, error) {
	type row struct {
		ID        uint
		TypeID    uint
		Label     string
		Value     string
		Sort      int
		Status    int8
		Remark    string
		CreatedAt int64
		UpdatedAt int64
		TypeCode  string
	}

	var rows []struct {
		ID        uint
		TypeID    uint
		Label     string
		Value     string
		Sort      int
		Status    int8
		Remark    string
		CreatedAt string
		UpdatedAt string
		TypeCode  string
	}
	db := r.db.WithContext(ctx).Table("dict_data AS d").
		Select("d.id, d.type_id, d.label, d.value, d.sort, d.status, d.remark, d.created_at, d.updated_at, t.code AS type_code").
		Joins("JOIN dict_types t ON t.id = d.type_id")

	if q.TypeID != nil {
		db = db.Where("d.type_id = ?", *q.TypeID)
	}
	if q.TypeCode != "" {
		db = db.Where("t.code = ?", q.TypeCode)
	}
	if q.Label != "" {
		db = db.Where("d.label LIKE ?", "%"+q.Label+"%")
	}
	if q.Status != nil {
		db = db.Where("d.status = ?", *q.Status)
	}
	if err := db.Order("d.sort ASC, d.id ASC").Scan(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]DictDataWithType, 0, len(rows))
	for _, row := range rows {
		items = append(items, DictDataWithType{
			Data: DictData{
				ID:     row.ID,
				TypeID: row.TypeID,
				Label:  row.Label,
				Value:  row.Value,
				Sort:   row.Sort,
				Status: row.Status,
				Remark: row.Remark,
			},
			TypeCode: row.TypeCode,
		})
	}
	return items, nil
}

func (r *gormRepository) ListDataByTypeCode(ctx context.Context, typeCode string, onlyActive bool) ([]DictDataWithType, error) {
	q := ListDictDataQuery{TypeCode: typeCode}
	if onlyActive {
		active := int8(1)
		q.Status = &active
	}
	return r.ListData(ctx, q)
}

func (r *gormRepository) UpdateData(ctx context.Context, item DictData) (DictData, error) {
	if err := r.db.WithContext(ctx).Save(&item).Error; err != nil {
		if isDuplicateError(err) {
			return DictData{}, ErrDictDuplicated
		}
		return DictData{}, err
	}
	return item, nil
}

func (r *gormRepository) DeleteData(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&DictData{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrDictDataNotFound
	}
	return nil
}
