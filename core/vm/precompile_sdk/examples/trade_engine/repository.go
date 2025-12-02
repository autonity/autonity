package trade_engine

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm/precompile_sdk/storage"
)

type Repository interface {
	SetBook(common.Hash, OrderBook) error
	GetBook(common.Hash) (OrderBook, error)
	SetOrder(common.Hash, Order) error
	GetOrder(common.Hash) (Order, error)
	InsertToOrderBook(common.Hash, Order) error
	RemoveFromOrderBook(pairHash common.Hash, orderID common.Hash) error
}

type PrecompileRepository struct {
	st *storage.Storage
}

func NewPrecompileRepository(st *storage.Storage) *PrecompileRepository {
	return &PrecompileRepository{st: st}
}

func (ps *PrecompileRepository) GetOrder(id common.Hash) (Order, error) {
	var order Order
	orderPath := ps.st.Field("Orders").Map(id)
	if orderPath.Error() != nil {
		return order, orderPath.Error()
	}
	// ID
	id, err := storage.Get[common.Hash](orderPath.Field("ID"))
	if err != nil {
		return order, err
	}
	// User
	user, err := storage.Get[common.Address](orderPath.Field("User"))
	if err != nil {
		return order, err
	}

	// Side
	sideUint, err := storage.Get[uint8](orderPath.Field("Side"))
	if err != nil {
		return order, err
	}

	// Price
	price, err := storage.Get[storage.Uint256](orderPath.Field("Price"))
	if err != nil {
		return order, err
	}

	// Qty
	qty, err := storage.Get[storage.Uint256](orderPath.Field("Qty"))
	if err != nil {
		return order, err
	}

	// Filled
	filled, err := storage.Get[storage.Uint256](orderPath.Field("Filled"))
	if err != nil {
		return order, err
	}
	// Status
	status, err := storage.Get[uint8](orderPath.Field("Status"))
	if err != nil {
		return order, err
	}
	// timestamp
	timestamp, err := storage.Get[int64](orderPath.Field("Timestamp"))
	if err != nil {
		return order, err
	}

	order.ID = id
	order.User = user
	order.Side = sideUint
	order.Price = price
	order.Qty = qty
	order.Filled = filled
	order.Status = status
	order.Timestamp = timestamp

	return order, nil
}

func (ps *PrecompileRepository) SetOrder(id common.Hash, order Order) error {
	orderPath := ps.st.Field("Orders").Map(id)
	if orderPath.Error() != nil {
		return orderPath.Error()
	}
	// ID
	err := storage.Set[common.Hash](orderPath.Field("ID"), order.ID)
	if err != nil {
		return err
	}
	// User
	err = storage.Set[common.Address](orderPath.Field("User"), order.User)
	if err != nil {
		return err
	}

	// Side
	err = storage.Set[uint8](orderPath.Field("Side"), uint8(order.Side))
	if err != nil {
		return err
	}

	// Price
	err = storage.Set[storage.Uint256](orderPath.Field("Price"), order.Price)
	if err != nil {
		return err
	}
	// Qty
	err = storage.Set[storage.Uint256](orderPath.Field("Qty"), order.Qty)
	if err != nil {
		return err
	}
	// Filled
	err = storage.Set[storage.Uint256](orderPath.Field("Filled"), order.Filled)
	if err != nil {
		return err
	}
	// Status
	err = storage.Set[uint8](orderPath.Field("Status"), order.Status)
	if err != nil {
		return err
	}

	// Timestamp
	err = storage.Set[int64](orderPath.Field("Timestamp"), order.Timestamp)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PrecompileRepository) getLevel(path *storage.Path) (Level, error) {
	var level Level
	if path.Error() != nil {
		return level, path.Error()
	}
	// Price
	price, err := storage.Get[storage.Uint256](path.Field("Price"))
	if err != nil {
		return level, err
	}
	// TotalQty
	totalQty, err := storage.Get[storage.Uint256](path.Field("TotalQty"))
	if err != nil {
		return level, err
	}
	level.Price = price
	level.TotalQty = totalQty

	// get order IDs array path:
	var orderIDs []common.Hash
	orderIDPath := path.Field("OrderIDs")
	orderIDArr, err := storage.NewArray[common.Hash](orderIDPath)
	if err != nil {
		return level, err
	}
	orderLength, err := orderIDArr.Len()
	if err != nil {
		return level, err
	}
	for orderIndex := range orderLength {
		if orderIndex >= uint64(len(level.OrderIDs)) {
			break
		}
		orderID, err := orderIDArr.ValueAt(orderIndex)
		if err != nil {
			return level, err
		}
		orderIDs = append(orderIDs, orderID)
	}
	level.OrderIDs = orderIDs
	return level, nil
}

func (ps *PrecompileRepository) setLevel(path *storage.Path, level Level) error {
	// set price
	if err := storage.Set[storage.Uint256](path.Field("Price"), level.Price); err != nil {
		return err
	}

	if err := storage.Set[storage.Uint256](path.Field("TotalQty"), level.TotalQty); err != nil {
		return err
	}

	// get order IDs array path:
	orderIDPath := path.Field("OrderIDs")
	orderIDArr, err := storage.NewArray[common.Hash](orderIDPath)
	if err != nil {
		return err
	}
	orderLength, err := orderIDArr.Len()
	if err != nil {
		return err
	}
	for orderID := range orderLength {
		if orderID >= uint64(len(level.OrderIDs)) {
			break
		}
		err := orderIDArr.SetValueAt(orderID, level.OrderIDs[orderID])
		if err != nil {
			return err
		}
	}
	return nil
}

func (ps *PrecompileRepository) GetBook(pair common.Hash) (OrderBook, error) {
	bookPath := ps.st.Field("Books").Map(pair)
	if bookPath.Error() != nil {
		return OrderBook{}, bookPath.Error()
	}
	var orderBook OrderBook
	bids := bookPath.Field("Bids")
	bidsArr, err := storage.NewArray[common.Hash](bids)
	if err != nil {
		return orderBook, err
	}
	bidsLen, err := bidsArr.Len()
	var bidLevels []Level
	for index := range bidsLen {
		if index >= uint64(len(orderBook.Bids)) {
			break
		}
		level, err := ps.getLevel(bids.Index(index))
		if err != nil {
			return OrderBook{}, err
		}
		bidLevels = append(bidLevels, level)
	}

	asks := bookPath.Field("Asks")
	asksArr, err := storage.NewArray[common.Hash](asks)
	if err != nil {
		return orderBook, err
	}
	asksLen, err := asksArr.Len()
	var askLevels []Level
	for index := range asksLen {
		if index >= uint64(len(orderBook.Asks)) {
			break
		}
		level, err := ps.getLevel(asks.Index(index))
		if err != nil {
			return OrderBook{}, err
		}
		askLevels = append(askLevels, level)
	}
	nextId, err := storage.Get[storage.Uint256](bookPath.Field("NextID"))
	if err != nil {
		return orderBook, err
	}
	orderBook.Bids = bidLevels
	orderBook.Asks = askLevels
	orderBook.NextID = nextId
	return orderBook, nil
}

func (ps *PrecompileRepository) SetBook(pair common.Hash, book OrderBook) error {
	bookPath := ps.st.Field("Books").Map(pair)
	if bookPath.Error() != nil {
		return bookPath.Error()
	}
	// NextID
	if err := storage.Set[storage.Uint256](bookPath.Field("NextID"), book.NextID); err != nil {
		return err
	}
	// Levels
	bidsPath := bookPath.Field("Bids")
	if bidsPath.Error() != nil {
		return bidsPath.Error()
	}

	bidsArr, err := storage.NewArray[common.Hash](bidsPath)
	if err != nil {
		return err
	}
	bidsLength, err := bidsArr.Len()
	if err != nil {
		return err
	}
	for bidIdx := range bidsLength {
		if bidIdx >= uint64(len(book.Bids)) {
			break
		}
		// get level path reference
		levelPath := bidsArr.ReferenceAt(bidIdx)
		if levelPath.Error() != nil {
			return levelPath.Error()
		}
		if err := ps.setLevel(levelPath, book.Bids[bidIdx]); err != nil {
			return err
		}
	}
	return nil
}

func (ps *PrecompileRepository) InsertToOrderBook(pair common.Hash, order Order) error {
	bookPath := ps.st.Field("Books").Map(pair)
	if bookPath.Error() != nil {
		return bookPath.Error()
	}
	var sidePath *storage.Path
	if order.Side == 0 { // bid
		sidePath = bookPath.Field("Bids")
	} else {
		sidePath = bookPath.Field("Asks")
	}
	if sidePath.Error() != nil {
		return sidePath.Error()
	}

	// todo: hardcode to index 0 for now
	levelPath := sidePath.Index(0)
	if levelPath.Error() != nil {
		return levelPath.Error()
	}
	orderIDPath := levelPath.Field("OrderIDs")
	if orderIDPath.Error() != nil {
		return orderIDPath.Error()
	}
	orderIDArr, err := storage.NewArray[common.Hash](orderIDPath)
	if err != nil {
		return err
	}
	newElemPath, err := orderIDArr.Grow()
	if err != nil {
		return err
	}
	err = storage.Set[common.Hash](newElemPath, order.ID)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PrecompileRepository) RemoveFromOrderBook(pair common.Hash, orderID common.Hash) error {
	bookPath := ps.st.Field("Books").Map(pair)
	if bookPath.Error() != nil {
		return bookPath.Error()
	}

	if bookPath.Error() != nil {
		return bookPath.Error()
	}
	var sidePath *storage.Path
	order, err := ps.GetOrder(orderID)
	if err != nil {
		return err
	}
	if order.Side == 0 { // bid
		sidePath = bookPath.Field("Bids")
	} else {
		sidePath = bookPath.Field("Asks")
	}
	if sidePath.Error() != nil {
		return sidePath.Error()
	}

	// todo: assume on 0 level
	levelPath := sidePath.Index(0)
	if levelPath.Error() != nil {
		return levelPath.Error()
	}
	orderIDPath := levelPath.Field("OrderIDs")
	if orderIDPath.Error() != nil {
		return orderIDPath.Error()
	}
	orderIDArr, err := storage.NewArray[common.Hash](orderIDPath)
	if err != nil {
		return err
	}
	err = orderIDArr.Shrink()
	if err != nil {
		return err
	}
	return nil
}
