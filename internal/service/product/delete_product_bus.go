package productbus

import (
	"client/internal/common/apperrors"
	imagemodel "client/internal/model/mysql/image"
	productmodel "client/internal/model/mysql/product"
	"context"

	"gorm.io/gorm"
)

type DeleteProductInterface interface {
	DeleteProduct(ctx context.Context, db *gorm.DB,
		cond map[string]interface{}, morekyes ...string) error
}

type DeleteProductBusiness struct {
	store      DeleteProductInterface
	imageStore UpdateImage
	db         *gorm.DB
}

func NewDeleteProducBiz(store DeleteProductInterface,
	imageStore UpdateImage, db *gorm.DB) *DeleteProductBusiness {
	return &DeleteProductBusiness{
		store:      store,
		imageStore: imageStore,
		db:         db,
	}
}

func (biz *DeleteProductBusiness) DeleteProduct(ctx context.Context, db *gorm.DB,
	cond map[string]interface{}, morekyes ...string) error {

	tx := biz.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}

	defer func() {
		// Dùng recover nếu cần
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 2. Tạo Category (dùng store có hỗ trợ transaction)
	err := biz.store.DeleteProduct(ctx, tx, cond)
	// HOẶC: biz.store.CreateCategory() -> nhưng ta cần nó chạy trong transaction =>
	// bạn cần pass tx xuống store, tùy thiết kế
	if err != nil {
		tx.Rollback()
		return apperrors.ErrCannotDeleteEntity(productmodel.EntityName, err)
	}

	// 3. Lấy ảnh cũ
	oldImages, err := biz.imageStore.ListItem(ctx, map[string]interface{}{"resource_id": cond["product_id"]})
	if err != nil {
		tx.Rollback()
		return apperrors.ErrCannotListEntity(imagemodel.EntityName, err)
	}

	var toRemove []uint64
	for _, item := range oldImages {
		toRemove = append(toRemove, item.ImageID)
	}

	// 6. Xoá liên kết (hoặc xoá hẳn) cho "toRemove"
	if len(toRemove) > 0 {
		// Nếu chỉ muốn xóa liên kết => set resource_id = NULL
		err := biz.imageStore.BulkUpdateResourceID(ctx, tx, toRemove, nil, nil)

		// Nếu muốn xoá luôn record =>
		// err := biz.imageStore.BulkDeleteImages(ctx, toRemove)

		if err != nil {
			tx.Rollback()
			return apperrors.ErrCannotUpdateEntity("image", err)
		}
	}

	// 7. Commit Transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
