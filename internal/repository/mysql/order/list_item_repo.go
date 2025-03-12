package orderrepo

import (
	"client/internal/common/apperrors"
	"client/internal/common/filter"
	"client/internal/common/paging"
	ordermodel "client/internal/model/mysql/order"
	"context"

	"gorm.io/gorm"
)

// ListItem liệt kê các đơn hàng với các điều kiện lọc và phân trang
func (r *mysqlOrder) ListItem(ctx context.Context, cond map[string]interface{},
	paging *paging.Paging,
	filter *filter.Filter, morekeys ...string) ([]ordermodel.Order, error) {

	db := r.db.WithContext(ctx)

	// Build base query
	query := r.buildBaseQuery(db, cond)

	// Apply filters
	query = r.applyFilters(query, filter)

	// Count total records trước khi áp dụng phân trang và sắp xếp
	if err := r.countTotalRecords(query, paging); err != nil {
		return nil, err
	}

	// Apply pagination and sorting sau khi đã đếm
	query, err := r.addPaging(query, paging)
	if err != nil {
		return nil, err
	}

	// Execute main query with preload
	orders, err := r.executeMainQuery(query, paging, filter)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

// buildBaseQuery xây dựng truy vấn cơ bản với các điều kiện
func (r *mysqlOrder) buildBaseQuery(db *gorm.DB, cond map[string]interface{}) *gorm.DB {
	return db.Table(ordermodel.Order{}.TableName()).Where(cond)
}

// applyFilters áp dụng các bộ lọc lên truy vấn chính
func (r *mysqlOrder) applyFilters(query *gorm.DB, filter *filter.Filter) *gorm.DB {
	if filter == nil {
		return query
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// Áp dụng lọc theo khoảng giá
	if filter.MinPrice > 0 {
		query = query.Where("total_amount >= ?", filter.MinPrice)
	}

	if filter.MaxPrice > 0 {
		query = query.Where("total_amount <= ?", filter.MaxPrice)
	}

	// Áp dụng lọc theo ngày
	if filter.StartDate != nil && filter.EndDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", filter.StartDate, filter.EndDate)
	} else if filter.StartDate != nil {
		query = query.Where("created_at >= ?", filter.StartDate)
	} else if filter.EndDate != nil {
		query = query.Where("created_at <= ?", filter.EndDate)
	}

	// Áp dụng bộ lọc ProductID
	if filter.ProductID != 0 {
		subQuery := r.db.Model(&ordermodel.OrderItem{}).
			Select("order_id").
			Where("product_id = ?", filter.ProductID).
			Where("deleted_at IS NULL")

		query = query.Where("order_id IN (?)", subQuery)
	}

	// Các bộ lọc thêm có thể được thêm vào đây

	return query
}

// countTotalRecords đếm tổng số bản ghi sau khi áp dụng bộ lọc
func (r *mysqlOrder) countTotalRecords(query *gorm.DB, paging *paging.Paging) error {
	var total int64
	// Tạo một bản sao của truy vấn mà không bao gồm phân trang và sắp xếp
	countQuery := query.Session(&gorm.Session{}).Distinct("order_id").Count(&total)
	if countQuery.Error != nil {
		return apperrors.ErrDB(countQuery.Error)
	}
	paging.Total = total
	return nil
}

// addPaging thêm phân trang và sắp xếp vào truy vấn
func (r *mysqlOrder) addPaging(db *gorm.DB, paging *paging.Paging) (*gorm.DB, error) {
	// Định nghĩa các trường cho phép sắp xếp
	allowedSortFields := map[string]bool{
		"order_id":     true,
		"total_amount": true,
		"created_at":   true,
		"updated_at":   true,
		"status":       true,
	}

	// Parse và validate các trường sắp xếp
	sortFields, err := paging.ParseSortFields(paging.Sort, allowedSortFields)
	if err != nil {
		return nil, apperrors.NewErrorResponse(err, "Invalid sort parameters", err.Error(), "InvalidSort")
	}

	// Áp dụng sắp xếp vào truy vấn
	if len(sortFields) > 0 {
		for _, sortField := range sortFields {
			db = db.Order(sortField)
		}
	} else {
		// Sắp xếp mặc định nếu không có tham số sắp xếp
		db = db.Order("order_id desc")
	}

	// Áp dụng phân trang
	offset := (paging.Page - 1) * paging.Limit
	db = db.Offset(offset).Limit(paging.Limit)

	return db, nil
}

// executeMainQuery thực hiện truy vấn chính với preload
func (r *mysqlOrder) executeMainQuery(
	query *gorm.DB,
	paging *paging.Paging,
	filter *filter.Filter,
) ([]ordermodel.Order, error) {
	var orders []ordermodel.Order

	// Áp dụng các preload cần thiết
	mainQuery := query.
		Preload("OrderItems", r.orderItemsPreloadCondition(filter)).
		Preload("OrderItems.Product").
		Preload("User")

	// Thêm preload cho Images nếu cần
	//	if contains(morekeys, "images") {
	//		mainQuery = mainQuery.Preload("OrderItems.Product.Images")
	//	}

	// Thêm preload cho Category nếu cần
	//	if contains(morekeys, "category") {
	//		mainQuery = mainQuery.Preload("OrderItems.Product.Category")
	//	}

	err := mainQuery.Find(&orders).Error
	if err != nil {
		return nil, apperrors.ErrDB(err)
	}

	return orders, nil
}

// orderItemsPreloadCondition tạo điều kiện preload cho OrderItems
func (r *mysqlOrder) orderItemsPreloadCondition(filter *filter.Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		}

		// Áp dụng các bộ lọc cho OrderItems
		if filter.MinPrice > 0 {
			db = db.Where("price >= ?", filter.MinPrice)
		}

		if filter.MaxPrice > 0 {
			db = db.Where("price <= ?", filter.MaxPrice)
		}

		// Có thể thêm các bộ lọc khác cho OrderItems

		return db
	}
}

// Hàm tiện ích kiểm tra xem một chuỗi có trong slice hay không
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
