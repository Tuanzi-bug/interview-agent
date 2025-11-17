package impl

import (
	questionbankapi "ai-eino-interview-agent/api/model/questionbank"
	"ai-eino-interview-agent/internal/model"
	"ai-eino-interview-agent/internal/repository"
	"context"

	"gorm.io/gorm"
)

// QuestionBankServiceImpl 试题库服务实现
type QuestionBankServiceImpl struct{}

// NewQuestionBankServiceImpl 创建试题库服务实例
func NewQuestionBankServiceImpl() *QuestionBankServiceImpl {
	return &QuestionBankServiceImpl{}
}

// GetMainCategories 获取所有大类（1类）标签
func (s *QuestionBankServiceImpl) GetMainCategories(ctx context.Context) (*questionbankapi.GetMainCategoriesResponse, error) {
	db := repository.GetDB()
	var categories []model.QuestionBankEntrance

	// 查询所有 CategoryType=1 的分类（大类）
	categoryType := int8(1)
	err := db.Where("category_type = ?", categoryType).
		Where("(is_deleted IS NULL OR is_deleted = 0)").
		Order("created_time DESC").
		Find(&categories).Error
	if err != nil {
		return nil, err
	}

	// 转换为响应项
	items := make([]*questionbankapi.MainCategoryItem, 0, len(categories))
	for _, cat := range categories {
		item := &questionbankapi.MainCategoryItem{
			ID:           int64(cat.ID),
			CategoryName: cat.CategoryName,
			CreatedTime:  cat.CreatedTime.Unix(),
			UpdatedTime:  cat.UpdatedTime.Unix(),
		}
		if cat.ImageURL != "" {
			item.ImageURL = &cat.ImageURL
		}
		items = append(items, item)
	}

	return &questionbankapi.GetMainCategoriesResponse{
		Categories: items,
	}, nil
}

// GetSubCategories 根据大类ID获取所有小分支（2类）
func (s *QuestionBankServiceImpl) GetSubCategories(ctx context.Context, parentID int64) (*questionbankapi.GetSubCategoriesResponse, error) {
	db := repository.GetDB()
	var subCategories []model.QuestionBankEntrance
	var parentCategory model.QuestionBankEntrance

	// 先查询父级分类信息
	err := db.Where("id = ?", uint(parentID)).
		Where("category_type = ?", int8(1)).
		Where("(is_deleted IS NULL OR is_deleted = 0)").
		First(&parentCategory).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	// 查询所有 CategoryType=2 且 ParentID=parentID 的分类（小分支）
	categoryType := int8(2)
	parentIDUint := uint(parentID)
	err = db.Where("category_type = ?", categoryType).
		Where("parent_id = ?", parentIDUint).
		Where("(is_deleted IS NULL OR is_deleted = 0)").
		Order("created_time DESC").
		Find(&subCategories).Error
	if err != nil {
		return nil, err
	}

	// 转换为响应项
	items := make([]*questionbankapi.SubCategoryItem, 0, len(subCategories))
	for _, cat := range subCategories {
		item := &questionbankapi.SubCategoryItem{
			ID:           int64(cat.ID),
			CategoryName: cat.CategoryName,
			ParentID:     int64(*cat.ParentID),
			CreatedTime:  cat.CreatedTime.Unix(),
			UpdatedTime:  cat.UpdatedTime.Unix(),
		}
		if cat.ImageURL != "" {
			item.ImageURL = &cat.ImageURL
		}
		items = append(items, item)
	}

	return &questionbankapi.GetSubCategoriesResponse{
		ParentID:   int64(parentCategory.ID),
		ParentName: parentCategory.CategoryName,
		Categories: items,
	}, nil
}
