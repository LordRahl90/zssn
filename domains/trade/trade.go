package trade

import (
	"context"
	"fmt"
	"time"

	"zssn/domains/entities"
	"zssn/domains/inventory"
	"zssn/domains/trade/store"
	"zssn/domains/users"
)

// TradeService to implement ITradeService
type TradeService struct {
	Storage          store.ITradeStorage
	UserService      users.IUserService
	InventoryService inventory.IInventoryService
}

// New returns an implementation of ITradeService
func New(storage store.ITradeStorage, usr users.IUserService, inv inventory.IInventoryService) ITradeService {
	return &TradeService{
		Storage:          storage,
		UserService:      usr,
		InventoryService: inv,
	}
}

// Execute implements ITradeService
func (ts *TradeService) Execute(ctx context.Context, seller, buyer *entities.TradeItems) error {
	if err := ts.VerifyTransaction(ctx, seller, buyer); err != nil {
		return err
	}
	b := buyer.ToDBTradeItemEntities()
	s := seller.ToDBTradeItemEntities()
	if err := ts.Storage.Execute(ctx, s, b); err != nil {
		return err
	}
	seller.Reference = s.Reference
	buyer.Reference = s.Reference

	// reduce the balance from seller
	balances, err := ts.InventoryService.FindMultipleInventory(ctx, seller.UserID, buyer.UserID)
	if err != nil {
		return err
	}
	sellerBalance := balances[seller.UserID]
	buyerBalance := balances[buyer.UserID]
	for _, v := range s.Items {
		dd := sellerBalance[v.Item]
		newBalance := dd.Balance - v.Quantity
		if err := ts.InventoryService.UpdateBalance(ctx, seller.UserID, v.Item, newBalance); err != nil {
			return err
		}

		// get the buyer's balance of this same item
		bb := buyerBalance[v.Item]
		newBalance = bb.Balance + v.Quantity
		fmt.Printf("BB Item: %s\t Qty: %d", v.Item.String(), newBalance)
		// increase the buyer's balance here
		if err := ts.InventoryService.UpdateBalance(ctx, buyer.UserID, v.Item, newBalance); err != nil {
			return err
		}
	}

	for _, v := range b.Items {
		dd := buyerBalance[v.Item]
		newBalance := dd.Balance - v.Quantity
		if err := ts.InventoryService.UpdateBalance(ctx, buyer.UserID, v.Item, newBalance); err != nil {
			return err
		}

		// get the buyer's balance of this same item
		sb := sellerBalance[v.Item]
		newBalance = sb.Balance + v.Quantity
		// increase the buyer's balance here
		if err := ts.InventoryService.UpdateBalance(ctx, seller.UserID, v.Item, newBalance); err != nil {
			return err
		}
	}

	return nil
}

// History implements ITradeService
func (ts *TradeService) History(ctx context.Context, id string, startDate time.Time, endDate time.Time) ([]*entities.Transaction, error) {
	var result []*entities.Transaction
	res, err := ts.Storage.History(ctx, id, startDate, endDate)
	if err != nil {
		return nil, err
	}
	for _, v := range res {
		result = append(result, entities.FromDBTransactionEntity(v))
	}

	return result, nil
}

func (ts *TradeService) IsTransactionAmountEqual(sellerItem *entities.TradeItems, buyerItem *entities.TradeItems) error {
	if sellerItem.Calculate() != buyerItem.Calculate() {
		return fmt.Errorf("value of the trade doesn't match")
	}
	return nil
}

// AnyParticipantInfected confirms if any of the participants have been infected
func (ts *TradeService) AnyParticipantInfected(users ...*entities.User) error {
	if len(users) == 0 {
		return fmt.Errorf("invalid users provided")
	}
	for _, v := range users {
		if v == nil {
			return fmt.Errorf("one of the participants is invalid")
		}
		if v.Infected {
			return fmt.Errorf("participant %s is infected, cannot proceed with transaction", v.Name)
		}
	}
	return nil
}

// EnoughStock confirms if there is enough stock to fulfill trade
func (ts *TradeService) EnoughStock(stock entities.Stock, item *entities.TradeItems) error {
	if len(stock) == 0 {
		return fmt.Errorf("invalid stock provided")
	}
	if len(item.Items) == 0 {
		return fmt.Errorf("invalid items in trade items")
	}
	// confirm that the items in stock can fulfil trade
	for _, v := range item.Items {
		is, ok := stock[v.Item]
		if !ok {
			return fmt.Errorf("user doesn't have the item in stock " + v.Item.String())
		}
		if is.Quantity < v.Quantity {
			return fmt.Errorf("user doesn't have enough to fulfill transaction")
		}
	}
	return nil
}

// VerifyTransaction implements ITradeService
func (ts *TradeService) VerifyTransaction(ctx context.Context, sellerItem *entities.TradeItems, buyerItem *entities.TradeItems) error {
	// verify that the values are the same
	if err := ts.IsTransactionAmountEqual(sellerItem, buyerItem); err != nil {
		return err
	}
	// verify if both users can transact
	users, err := ts.UserService.FindUsers(ctx, sellerItem.UserID, buyerItem.UserID)
	if err != nil {
		return err
	}
	if len(users) < 2 {
		// means one of the participants is not a user anymore
		return fmt.Errorf("a participant has been removed or is invalid")
	}

	// confirm none of the participants have been infected
	if err := ts.AnyParticipantInfected(users[sellerItem.UserID], users[buyerItem.UserID]); err != nil {
		return err
	}

	// confirm the inventory of the participants
	invs, err := ts.InventoryService.FindMultipleInventory(ctx, sellerItem.UserID, buyerItem.UserID)
	if err != nil {
		return err
	}

	if err := ts.EnoughStock(invs[sellerItem.UserID], sellerItem); err != nil {
		return err
	}

	if err := ts.EnoughStock(invs[buyerItem.UserID], buyerItem); err != nil {
		return err
	}

	return nil
}
