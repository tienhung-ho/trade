package categorybusiness

import (
	"client/internal/common/apperrors"
	categorymodel "client/internal/model/mysql/category"
	imagemodel "client/internal/model/mysql/image"
	"context"
	"gorm.io/gorm"
)

type UpdateCategoryStorage interface {
	GetCategory(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*categorymodel.Category, error)
	UpdateCategory(ctx context.Context,
		db *gorm.DB,
		cond map[string]interface{}, data *categorymodel.UpdateCategory, morekeys ...string) (*categorymodel.Category, error)
}

type UpdateCategoryCache interface {
	DeleteCategory(ctx context.Context, morekeys ...string) error
	DeleteListCache(ctx context.Context, entityName string) error
}

type updateCategoryBusiness struct {
	store UpdateCategoryStorage
	cache UpdateCategoryCache
	image UpdateImage
	db    *gorm.DB
}

func NewUpdateCategoryBiz(store UpdateCategoryStorage, cache UpdateCategoryCache, image UpdateImage, db *gorm.DB) *updateCategoryBusiness {
	return &updateCategoryBusiness{store, cache, image, db}
}

func (biz *updateCategoryBusiness) UpdateCategory(ctx context.Context, cond map[string]interface{}, data *categorymodel.UpdateCategory, morekeys ...string) (*categorymodel.Category, error) {
	// 1. Lấy category cũ
	record, err := biz.store.GetCategory(ctx, cond, morekeys...)
	if err != nil {
		return nil, apperrors.ErrNotFoundEntity(categorymodel.EntityName, err)
	}

	// 2. Bắt đầu transaction
	db := biz.db // HOẶC bạn có biz.db nếu đã truyền vào
	tx := db.Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// rollback/commit an toàn
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 3. Thực hiện UpdateCategory trong transaction
	updatedRecord, err := biz.store.UpdateCategory(ctx, tx, map[string]interface{}{"category_id": record.CategoryID}, data)
	if err != nil {
		tx.Rollback()
		return nil, apperrors.ErrCannotUpdateEntity(categorymodel.EntityName, err)
	}

	// 4. Xử lý images, nếu người dùng gửi data.Images
	if data.Images != nil {
		// 4a. Lấy danh sách ảnh cũ
		oldImages, err := biz.image.ListItem(ctx, map[string]interface{}{
			"resource_id": record.CategoryID,
		})
		if err != nil {
			tx.Rollback()
			return nil, apperrors.ErrCannotListEntity(imagemodel.EntityName, err)
		}

		oldIDs := make(map[uint64]struct{})
		for _, img := range oldImages {
			oldIDs[img.ImageID] = struct{}{}
		}

		// 4b. Tạo map các ID ảnh mới
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

		// 4c. Bulk update: set resource_id = categoryID cho ảnh toAdd
		if len(toAdd) > 0 {
			err := biz.image.BulkUpdateResourceID(ctx, tx, toAdd, &record.CategoryID)
			if err != nil {
				tx.Rollback()
				return nil, apperrors.ErrCannotUpdateEntity("Image", err)
			}
		}

		// 4d. Bulk update: set resource_id = null (hoặc xóa hẳn) cho ảnh toRemove
		if len(toRemove) > 0 {
			// TH1: Xóa liên kết => set resource_id = NULL
			err := biz.image.BulkUpdateResourceID(ctx, tx, toRemove, nil)
			// TH2: Xóa luôn record => biz.image.BulkDeleteImages(ctx, toRemove)

			if err != nil {
				tx.Rollback()
				return nil, apperrors.ErrCannotUpdateEntity("Image", err)
			}
		}
	}

	// 5. Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return updatedRecord, nil
}
