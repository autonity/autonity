package trade_engine

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm/pdk/storage"
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
	err := storage.Load(orderPath, &order)
	if err != nil {
		return Order{}, err
	}
	return order, nil
}

func (ps *PrecompileRepository) SetOrder(id common.Hash, order Order) error {
	orderPath := ps.st.Field("Orders").Map(id)
	if orderPath.Error() != nil {
		return orderPath.Error()
	}
	return storage.Save(orderPath, &order)
}

func (ps *PrecompileRepository) getLevel(path *storage.Path) (Level, error) {
	var level Level
	if path.Error() != nil {
		return level, path.Error()
	}

	err := storage.Load(path, &level)
	if err != nil {
		return Level{}, err
	}
	// load the dynamic members manually
	orderIDPath := path.Field("OrderIDs")
	orderSlice := storage.NewSlice[common.Hash](orderIDPath)
	var orderIDs []common.Hash
	orderLength, err := orderSlice.Len()
	if err != nil {
		return level, err
	}
	for orderIndex := range orderLength {
		if orderIndex >= uint64(len(level.OrderIDs)) {
			break
		}
		orderID, err := orderSlice.Get(orderIndex)
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
	orderIDArr := storage.NewSlice[common.Hash](orderIDPath)
	orderLength, err := orderIDArr.Len()
	if err != nil {
		return err
	}
	targetLen := uint64(len(level.OrderIDs))
	for i := range targetLen {
		if i >= orderLength {
			err = orderIDArr.Append(level.OrderIDs[i])
			if err != nil {
				return err
			}
		} else {
			err = orderIDArr.Set(i, level.OrderIDs[i])
			if err != nil {
				return err
			}
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
	// load the primitive members
	if err := storage.Load(bookPath, &orderBook); err != nil {
		return OrderBook{}, err
	}

	bids := bookPath.Field("Bids")
	bidsWrapper := storage.NewSlice[Level](bids)
	bidsLen, err := bidsWrapper.Len()
	if err != nil {
		return OrderBook{}, err
	}

	orderBook.Bids = make([]Level, bidsLen)
	for index := range bidsLen {
		level, err := ps.getLevel(bids.Index(index))
		if err != nil {
			return OrderBook{}, err
		}
		orderBook.Bids[index] = level
	}

	asks := bookPath.Field("Asks")
	asksArr := storage.NewSlice[Level](asks)
	asksLen, err := asksArr.Len()
	orderBook.Asks = make([]Level, asksLen)
	for index := range asksLen {
		if index >= uint64(len(orderBook.Asks)) {
			break
		}
		level, err := ps.getLevel(asks.Index(index))
		if err != nil {
			return OrderBook{}, err
		}
		orderBook.Asks[index] = level
	}
	return orderBook, nil
}

func (ps *PrecompileRepository) SetBook(pair common.Hash, book OrderBook) error {
	bookPath := ps.st.Field("Books").Map(pair)
	if bookPath.Error() != nil {
		return bookPath.Error()
	}

	// save primitive members
	err := storage.Save(bookPath, &book)
	if err != nil {
		return err
	}
	// Bids
	bidsPath := bookPath.Field("Bids")
	if bidsPath.Error() != nil {
		return bidsPath.Error()
	}

	bidsArr := storage.NewSlice[Level](bidsPath)
	bidsLength, err := bidsArr.Len()
	if err != nil {
		return err
	}
	for i := uint64(0); i < uint64(len(book.Bids)); i++ {
		if i < bidsLength {
			err = ps.setLevel(bidsPath.Index(i), book.Bids[i])
			if err != nil {
				return err
			}
		} else {
			_, err := bidsArr.Grow()
			if err != nil {
				return err
			}
			err = ps.setLevel(bidsPath.Index(i), book.Bids[i])
			if err != nil {
				return err
			}
		}
	}

	for i := uint64(len(book.Bids)); i < bidsLength; i++ {
		err = bidsArr.Pop()
		if err != nil {
			return err
		}
	}

	asksPath := bookPath.Field("Asks")
	asksSlice := storage.NewSlice[Level](asksPath)
	asksLength, err := asksSlice.Len()
	if err != nil {
		return err
	}

	for i := uint64(0); i < uint64(len(book.Asks)); i++ {
		if i < asksLength {
			err = ps.setLevel(asksPath.Index(i), book.Asks[i])
			if err != nil {
				return err
			}
		} else {
			_, err := asksSlice.Grow()
			if err != nil {
				return err
			}
			err = ps.setLevel(asksPath.Index(i), book.Asks[i])
			if err != nil {
				return err
			}
		}
	}
	for i := uint64(len(book.Asks)); i < asksLength; i++ {
		err = asksSlice.Pop()
		if err != nil {
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

	levelSlice := storage.NewSlice[Level](sidePath)
	levelLen, err := levelSlice.Len()
	if err != nil {
		return err
	}

	if levelLen == 0 {
		// need to create the first level
		_, err := levelSlice.Grow()
		if err != nil {
			return err
		}
		err = ps.setLevel(sidePath.Index(0), Level{Price: order.Price, TotalQty: order.Qty, OrderIDs: []common.Hash{order.ID}})
		if err != nil {
			return err
		}
	}
	levelPath := sidePath.Index(0)
	if levelPath.Error() != nil {
		return levelPath.Error()
	}

	orderIDPath := levelPath.Field("OrderIDs")
	if orderIDPath.Error() != nil {
		return orderIDPath.Error()
	}

	orderIDSlice := storage.NewSlice[common.Hash](orderIDPath)
	err = orderIDSlice.Append(order.ID)
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
	orderIDSlice := storage.NewSlice[Level](orderIDPath)
	err = orderIDSlice.Pop()
	if err != nil {
		return err
	}
	return nil
}
