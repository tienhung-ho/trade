package productbus

import (
	"client/internal/common/apperrors"
	imagemodel "client/internal/model/mysql/image"
	productmodel "client/internal/model/mysql/product"
	"context"

	"gorm.io/gorm"
)

type UpdateProductInterface interface {
	UpdateProduct(ctx context.Context, db *gorm.DB, cond map[string]interface{},
		data *productmodel.UpdateProduct, morekeys ...string) (*productmodel.Product, error)
}

type UpdateProductBusiness struct {
	store      UpdateProductInterface
	imageStore UpdateImage
	db         *gorm.DB
}

func NewUpdateProductBiz(store UpdateProductInterface,
	imageStore UpdateImage,
	db *gorm.DB) *UpdateProductBusiness {
	return &UpdateProductBusiness{
		store:      store,
		imageStore: imageStore,
		db:         db,
	}
}

func (biz *UpdateProductBusiness) UpdateProduct(ctx context.Context, db *gorm.DB, cond map[string]interface{},
	data *productmodel.UpdateProduct, morekeys ...string) (*productmodel.Product, error) {

	tx := biz.db.Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}

	defer func() {
		// Dùng recover nếu cần
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 2. Tạo Category (dùng store có hỗ trợ transaction)
	record, err := biz.store.UpdateProduct(ctx, tx, cond, data)
	// HOẶC: biz.store.CreateCategory() -> nhưng ta cần nó chạy trong transaction =>
	// bạn cần pass tx xuống store, tùy thiết kế
	if err != nil {
		tx.Rollback()
		return nil, apperrors.ErrCannotUpdateEntity(productmodel.EntityName, err)
	}

	// 3. Lấy ảnh cũ
	oldImages, err := biz.imageStore.ListItem(ctx, map[string]interface{}{"resource_id": record.ProductID})
	if err != nil {
		tx.Rollback()
		return nil, apperrors.ErrCannotListEntity(imagemodel.EntityName, err)
	}

	oldIDs := make(map[uint64]struct{})
	for _, img := range oldImages {
		oldIDs[img.ImageID] = struct{}{}
	}

	// 4. Ảnh mới user đẩy lên
	newIDs := make(map[uint64]struct{})
	for _, img := range data.Images {
		newIDs[img.ImageID] = struct{}{}
	}

	var toAdd []uint64
	var toRemove []uint64

	for id := range newIDs {
		if _, existed := oldIDs[id]; !existed {
			toAdd = append(toAdd, id)
		}
	}
	for id := range oldIDs {
		if _, keep := newIDs[id]; !keep {
			toRemove = append(toRemove, id)
		}
	}

	entityName := productmodel.EntityName

	// 5. Thêm liên kết (UPDATE resource_id = recordID) cho các ảnh "toAdd"
	if len(toAdd) > 0 {
		// gọi bulk update
		err := biz.imageStore.BulkUpdateResourceID(ctx, tx, toAdd, &record.ProductID, &entityName)
		if err != nil {
			tx.Rollback()
			return nil, apperrors.ErrCannotUpdateEntity("image", err)
		}
	}

	// 6. Xoá liên kết (hoặc xoá hẳn) cho "toRemove"
	if len(toRemove) > 0 {
		// Nếu chỉ muốn xóa liên kết => set resource_id = NULL
		err := biz.imageStore.BulkUpdateResourceID(ctx, tx, toRemove, nil, nil)

		// Nếu muốn xoá luôn record =>
		// err := biz.imageStore.BulkDeleteImages(ctx, toRemove)

		if err != nil {
			tx.Rollback()
			return nil, apperrors.ErrCannotUpdateEntity("image", err)
		}
	}

	// 7. Commit Transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return record, nil
}
