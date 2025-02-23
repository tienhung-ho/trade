package categorybusiness

import (
	"client/internal/common/apperrors"
	categorymodel "client/internal/model/mysql/category"
	imagemodel "client/internal/model/mysql/image"
	"context"
	"sync"
)

type DeleteCategoryStorage interface {
	GetCategory(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*categorymodel.Category, error)
	DeleteCategory(ctx context.Context, cond map[string]interface{}, morekeys ...string) error
}

type DeleteCategoryCache interface {
	DeleteCategory(ctx context.Context, morekeys ...string) error
	DeleteListCache(ctx context.Context, entityName string) error
}

type DeleteImage interface {
	DeleteImage(ctx context.Context, cond map[string]interface{}, morekeys ...string) error
	UpdateImage(ctx context.Context, cond map[string]interface{}, data *imagemodel.UpdateImage) error
	ListItem(ctx context.Context, cond map[string]interface{}) ([]imagemodel.Image, error)
}

type deleteCategoryBusiness struct {
	store DeleteCategoryStorage
	cache DeleteCategoryCache
	image DeleteImage
}

func NewDeleteCategoryBiz(store DeleteCategoryStorage, cache DeleteCategoryCache, image DeleteImage) *deleteCategoryBusiness {
	return &deleteCategoryBusiness{store, cache, image}
}

func (biz *deleteCategoryBusiness) DeleteCategory(ctx context.Context, cond map[string]interface{}, morekeys ...string) error {

	record, err := biz.store.GetCategory(ctx, cond)

	if err != nil {
		return apperrors.ErrNotFoundEntity(categorymodel.EntityName, err)
	}

	if err := biz.store.DeleteCategory(ctx, map[string]interface{}{"category_id": record.CategoryID}, morekeys...); err != nil {
		return apperrors.ErrCannotUpdateEntity(categorymodel.EntityName, err)
	}

	// 3. Lấy danh sách hình ảnh cũ liên kết với sản phẩm
	imageCond := map[string]interface{}{
		"resource_id": record.CategoryID,
	}

	oldImages, err := biz.image.ListItem(ctx, imageCond)

	if err != nil {
		return apperrors.ErrCannotListEntity(imagemodel.EntityName, err)
	}

	var imagesToRemove []uint64
	for _, imgID := range oldImages {
		imagesToRemove = append(imagesToRemove, imgID.ImageID)
	}

	var wg sync.WaitGroup
	var deleteErr error
	var mu sync.Mutex

	// Giới hạn số lượng goroutines chạy đồng thời để tránh làm quá tải hệ thống
	sem := make(chan struct{}, 10) // Giới hạn 10 goroutines

	for _, imgID := range imagesToRemove {
		wg.Add(1)
		sem <- struct{}{}
		go func(id uint64) {
			defer wg.Done()
			defer func() { <-sem }()
			err := biz.image.DeleteImage(ctx, map[string]interface{}{"image_id": id})
			if err != nil {
				mu.Lock()
				deleteErr = err
				mu.Unlock()
			}
		}(imgID)
	}

	wg.Wait()

	if deleteErr != nil {
		// Xử lý lỗi (rollback transaction nếu cần)
		return apperrors.ErrCannotUpdateEntity("Image", deleteErr)
	}

	return nil
}
