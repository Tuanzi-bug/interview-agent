package model

var DimensionDetailDao _DimensionDetail

// DimensionDetail 维度评分详情模型
type (
	_DimensionDetail struct {
	}
	DimensionDetail struct {
		ID           uint64 `json:"id" gorm:"primaryKey;autoIncrement;comment:维度评分唯一ID"`
		InterviewID  uint64 `json:"interview_id" gorm:"not null;index:idx_interview_id;comment:面试ID (外键)"`
		Title        string `json:"title" gorm:"size:255;not null;comment:维度标题，如'技术基础与实践能力'"`
		Content      string `json:"content" gorm:"type:text;not null;comment:维度的详细评价内容"`
		Score        uint8  `json:"score" gorm:"not null;comment:该维度的得分 (0-100)"`
		DisplayOrder uint8  `json:"display_order" gorm:"not null;default:0;comment:显示顺序，用于排序"`
	}
)

// TableName 指定表名
func (d *DimensionDetail) TableName() string {
	return "dimension_details"
}

// CreateDimensionDetail 创建维度评分详情
func (d *_DimensionDetail) CreateDimensionDetail(detail *DimensionDetail) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	return getDB().Create(detail).Error
}

// GetDimensionDetailsByInterviewID 根据面试ID查询维度评分详情列表
func (d *_DimensionDetail) GetDimensionDetailsByInterviewID(interviewID uint64) ([]*DimensionDetail, error) {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	var details []*DimensionDetail
	err := getDB().Where("interview_id = ?", interviewID).
		Order("display_order ASC").
		Find(&details).Error
	if err != nil {
		return nil, err
	}
	return details, nil
}

// UpdateDimensionDetail 更新维度评分详情
func (d *_DimensionDetail) UpdateDimensionDetail(detail *DimensionDetail) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	return getDB().Model(&DimensionDetail{}).
		Where("id = ?", detail.ID).
		Updates(detail).Error
}

// DeleteDimensionDetail 删除维度评分详情
func (d *_DimensionDetail) DeleteDimensionDetail(id uint64) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	return getDB().Where("id = ?", id).Delete(&DimensionDetail{}).Error
}
