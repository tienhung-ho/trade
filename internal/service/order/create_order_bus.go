package orderbus

import (
	"client/internal/common/apperrors"
	ordermodel "client/internal/model/mysql/order"
	productmodel "client/internal/model/mysql/product"
	usermodel "client/internal/model/mysql/user"
	cosmosservice "client/internal/service/cosmos"
	bech32util "client/internal/util/bech32"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"gorm.io/gorm"
)

type CreateOrderInterface interface {
	CreateOrder(ctx context.Context, db *gorm.DB,
		data *ordermodel.CreateOrder, morekeys ...string) (uint64, error)
}

type CreateOrderItemInterface interface {
	CreatePatchOrderItem(ctx context.Context,
		db *gorm.DB, data []ordermodel.CreateOrderItem, morekeys ...string) error
}

type ProductInterface interface {
	ListItemByIDs(ctx context.Context, cond []uint64) ([]productmodel.Product, error)
	BulkUpdateProductQuantity(ctx context.Context, db *gorm.DB,
		updates []productmodel.ProductQuantityUpdate) error
}

type UserInterface interface {
	GetUsers(ctx context.Context, IDs []uint64,
		morekeys ...string) ([]usermodel.User, error)
	GetUser(ctx context.Context, cond map[string]interface{},
		morekeys ...string) (*usermodel.User, error)
}

type CreateOrderBusiness struct {
	store          CreateOrderInterface
	productStore   ProductInterface
	orderItemStore CreateOrderItemInterface
	userStore      UserInterface
	cosmosStore    cosmosservice.CosmosInterface
	db             *gorm.DB
}

func NewCreateOrder(store CreateOrderInterface,
	productStore ProductInterface,
	orderItemStorage CreateOrderItemInterface,
	userStore UserInterface,
	cosmos cosmosservice.CosmosInterface,
	db *gorm.DB) *CreateOrderBusiness {
	return &CreateOrderBusiness{
		store:          store,
		productStore:   productStore,
		orderItemStore: orderItemStorage,
		userStore:      userStore,
		cosmosStore:    cosmos,
		db:             db,
	}
}

type PrepareData struct {
	data map[uint64][]ordermodel.CreateOrderItem
}

func (biz *CreateOrderBusiness) preparingData(ctx context.Context, createOrderItem []ordermodel.CreateOrderItem) (*PrepareData, error) {

	var existsProductIDs []uint64
	var quantityMap = make(map[uint64]uint) // Sử dụng map để theo dõi số lượng theo ProductID

	for _, item := range createOrderItem {
		if item.Quantity <= 0 {
			return nil, apperrors.ErrInvalidRequest(errors.New("product quantity must be greater than 0"))
		}

		if item.ProductID == 0 {
			return nil, apperrors.ErrInvalidRequest(errors.New("product ID cannot be 0"))
		}

		existsProductIDs = append(existsProductIDs, item.ProductID)
		quantityMap[item.ProductID] = item.Quantity
	}

	products, err := biz.productStore.ListItemByIDs(ctx, existsProductIDs)
	if err != nil {
		return nil, apperrors.ErrNotFoundEntity(productmodel.EntityName, err)
	}

	productMap := make(map[uint64]productmodel.Product)
	for _, item := range products {
		productMap[item.ProductID] = item
	}

	orderItemMap := make(map[uint64][]ordermodel.CreateOrderItem)
	for _, item := range createOrderItem {
		product := productMap[item.ProductID]

		if product.Stock < int(quantityMap[product.ProductID]) {

			return nil, apperrors.NewErrorResponse(
				nil,
				fmt.Sprintf("product %s (ID: %d) has insufficient stock (%d requested, %d available)",
					product.Name, product.ProductID, quantityMap[product.ProductID], product.Stock),
				"Product",
				"InsufficientStock",
			)
		}

		itemTotal := float64(item.Quantity) * product.Price

		orderItem := ordermodel.CreateOrderItem{
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			UnitPrice:  product.Price,
			TotalPrice: itemTotal,
		}

		orderItemMap[product.UserID] = append(orderItemMap[product.UserID], orderItem)

	}

	return &PrepareData{
		data: orderItemMap,
	}, nil
}

func (biz *CreateOrderBusiness) CreateOrder(ctx context.Context,
	data *ordermodel.CreateOrder, morekeys ...string) (*ordermodel.CreateOrderResponse, error) {
	// Kiểm tra dữ liệu đầu vào
	denom := os.Getenv("COIN_NAME")
	adminName := os.Getenv("ADMIN_NAME")
	if adminName == "" {
		return nil, errors.New("ADMIN_NAME is empty")
	}
	if data == nil || len(data.CreateOrderItems) == 0 {
		return nil, apperrors.ErrInvalidRequest(errors.New("order must have at least one item"))
	}

	// Bắt đầu transaction
	tx := biz.db.Begin()
	if err := tx.Error; err != nil {
		return nil, apperrors.ErrDB(err)
	}

	// Đảm bảo rollback nếu có lỗi
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Lọc và chuẩn bị dữ liệu sản phẩm
	prepareData, err := biz.preparingData(ctx, data.CreateOrderItems)

	if err != nil {
		return nil, apperrors.ErrCannotListEntity(productmodel.EntityName, err)
	}
	var orderIDs []uint64
	var sellerIDs []uint64
	finalAmountMap := make(map[uint64]float64)
	for sellerID, items := range prepareData.data {
		// Tính tổng tiền của nhóm seller này
		var sumPrice float64
		for _, it := range items {
			sumPrice += it.TotalPrice
		}

		finalAmountMap[sellerID] = sumPrice + data.ShippingFee - data.DiscountAmount
		newOrder := ordermodel.CreateOrder{
			UserID:          data.UserID, // người mua
			TotalAmount:     sumPrice,
			ShippingFee:     data.ShippingFee,    // tuỳ business
			DiscountAmount:  data.DiscountAmount, // tuỳ business
			FinalAmount:     finalAmountMap[sellerID],
			ShippingAddress: data.ShippingAddress,
			RecipientName:   data.RecipientName,
			RecipientPhone:  data.RecipientPhone,
			SellerID:        sellerID,
		}

		// Tạo đơn hàng
		orderID, err := biz.store.CreateOrder(ctx, tx, &newOrder)
		orderIDs = append(orderIDs, orderID)
		if err != nil {
			tx.Rollback()
			return nil, apperrors.ErrCannotCreateEntity(ordermodel.EntityName, err)
		}

		// Cập nhật OrderID cho tất cả các orderItem
		for i := range items {
			items[i].OrderID = orderID
		}
		// Tạo các orderItem
		if err := biz.orderItemStore.CreatePatchOrderItem(ctx, tx, items); err != nil {
			tx.Rollback()
			return nil, apperrors.ErrCannotCreateEntity(ordermodel.ItemEntityName, err)
		}

		sellerIDs = append(sellerIDs, sellerID)
	}

	var productUpdates []productmodel.ProductQuantityUpdate
	for _, item := range data.CreateOrderItems {
		productUpdates = append(productUpdates, productmodel.ProductQuantityUpdate{
			ProductID:  item.ProductID,
			Adjustment: -1 * int(item.Quantity), // Reduce stock by quantity ordered
		})
	}
	if err := biz.productStore.BulkUpdateProductQuantity(ctx, tx, productUpdates); err != nil {
		tx.Rollback()
		return nil, apperrors.ErrCannotUpdateEntity(productmodel.EntityName, err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, apperrors.ErrDB(err)
	}

	user, err := biz.userStore.GetUser(ctx, map[string]interface{}{"user_id": data.UserID})

	if err != nil {
		return nil, apperrors.ErrCannotGetEntity(usermodel.EntityName, err)
	}

	sellers, err := biz.userStore.GetUsers(ctx, sellerIDs)

	if err != nil {
		return nil, apperrors.ErrCannotGetEntity(usermodel.EntityName, err)
	}
	var unsignedTxs []string
	for _, seller := range sellers {
		sellerWallAddrs := seller.Wallets.WalletAddress
		amount := int(finalAmountMap[seller.ID])

		// Gọi hàm build Tx chưa ký
		unsignedTxBz, err := biz.BuildUnsignedSendTokenTx(ctx, user, denom, sellerWallAddrs, amount)
		if err != nil {
			return nil, fmt.Errorf("failed to build unsigned tx for seller %d: %w", seller.ID, err)
		}

		// Encode base64
		encodedTx := base64.StdEncoding.EncodeToString(unsignedTxBz)
		unsignedTxs = append(unsignedTxs, encodedTx)
	}

	response := &ordermodel.CreateOrderResponse{
		OrderIDs:    orderIDs,
		UnsignedTxs: unsignedTxs,
	}

	return response, nil
}

func (biz *CreateOrderBusiness) BuildUnsignedSendTokenTx(
	ctx context.Context,
	user *usermodel.User,
	denom string,
	receiverWallAddr string,
	numberOfToken int,
) ([]byte, error) {

	alias := bech32util.NormalizeBech32Address(user.Fullname)

	// 1) Tạo message
	msg := biz.cosmosStore.SendTokens(user.Wallets.WalletAddress, denom, int32(numberOfToken), receiverWallAddr)
	if msg == nil {
		return nil, fmt.Errorf("msg is nil")
	}

	// 2) Lấy & cập nhật accountNumber, sequence từ on-chain (nếu cần)
	if err := biz.cosmosStore.UpdateTxFactoryAccountSequence(alias); err != nil {
		log.Printf("Không thể cập nhật sequence cho %s: %v", alias, err)
		// Tùy bạn có muốn trả lỗi hay không
	}

	// 3) Dùng BuildUnsignedTx thay vì BuildTx
	accNum, seq, err := biz.cosmosStore.GetTxFactoryAccAndSeq()
	if err != nil {
		return nil, err
	}
	unsignedTxBz, err := biz.cosmosStore.BuildUnsignedTx(
		msg,
		"citcoin", // denom fee
		"500",     // feeAmount
		300000,    // gasLimit
		accNum,
		seq,
		"my offline tx", // memo (tùy ý)
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build unsigned tx: %w", err)
	}

	return unsignedTxBz, nil
}

func (biz *CreateOrderBusiness) SendToken(ctx context.Context, user *usermodel.User, denom,
	receiverWallAddr string, numberOfToken int) error {

	alias := bech32util.NormalizeBech32Address(user.Fullname)

	// Số lần thử lại tối đa
	maxRetries := 3
	var err error
	var res sdk.TxResponse

	for attempt := 0; attempt < maxRetries; attempt++ {
		// Cập nhật sequence trực tiếp trước mỗi lần gửi
		if err := biz.cosmosStore.UpdateTxFactoryAccountSequence(alias); err != nil {
			log.Printf("Không thể cập nhật sequence cho %s: %v", alias, err)
			// Tiếp tục thử với sequence hiện tại
		}

		msg := biz.cosmosStore.SendTokens(user.Wallets.WalletAddress, denom, int32(numberOfToken), receiverWallAddr)

		txBuilder, err := biz.cosmosStore.BuildTx(msg, "citcoin", "500", 300000)
		if err != nil {
			return apperrors.ErrInternal(fmt.Errorf("failed to build tx: %w", err))
		}

		if err := biz.cosmosStore.SignTx(ctx, alias, txBuilder); err != nil {
			return apperrors.ErrInternal(fmt.Errorf("failed to sign tx: %w", err))
		}

		txBytes, err := biz.cosmosStore.EncodeTxBytes(txBuilder)
		if err != nil {
			return fmt.Errorf("cannot encode tx: %w", err)
		}

		res, err = biz.cosmosStore.BroadcastTx(txBytes)

		// Kiểm tra kết quả
		if err == nil {
			// Giao dịch đã được gửi thành công
			log.Printf("Giao dịch thành công sau %d lần thử. TxHash: %s", attempt+1, res.TxHash)
			return nil
		}

		// Kiểm tra nếu lỗi là sequence mismatch
		if strings.Contains(err.Error(), "account sequence mismatch") {
			log.Printf("Sequence mismatch, đang thử lại (lần %d/%d)...", attempt+1, maxRetries)

			// Đợi một khoảng thời gian trước khi thử lại để đảm bảo blockchain đã cập nhật state
			time.Sleep(500 * time.Millisecond)
			continue
		}

		// Nếu gặp lỗi khác, không thử lại
		return fmt.Errorf("cannot broadcast tx: %w", err)
	}

	// Nếu đã thử hết số lần mà vẫn không thành công
	return fmt.Errorf("failed to send token after %d attempts, last error: %w", maxRetries, err)
}
