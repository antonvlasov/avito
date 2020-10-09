package handler

import (
	"errors"
	"fmt"
	"time"

	"github.com/antonvlasov/avito/netmanager"

	"github.com/antonvlasov/avito/dbmanager"
)

//Run starts handler which refreshes information every updateInterval seconds
func Run(updateInterval int) {
	const dbpath = "../resources/"
	defer dbmanager.Connect(dbpath)
	go func() {
		defer dbmanager.Close()
		tick := time.Tick(time.Duration(updateInterval) * time.Second)
		for {
			if shouldStop {
				break
			}
			<-tick
			err := UpdateAll()
			if err != nil {
				fmt.Println(err)
				dbmanager.Close()
				dbmanager.Connect(dbpath)
			}
		}
	}()
}

var shouldStop bool = false

func Stop() {
	shouldStop = true
}

func enqueueParse(items []dbmanager.Item, parseOutputChan chan netmanager.ParseResult, chanelCount int) {
	parseInputChan := make(chan dbmanager.Item, len(items))
	for i := range items {
		parseInputChan <- items[i]
	}
	close(parseInputChan)
	for i := 0; i < chanelCount; i++ {
		go netmanager.GetPrice(parseInputChan, parseOutputChan)
	}
}
func updateItems(inputChan <-chan dbmanager.Item, count int, errChan chan<- error) {
	for i := 0; i < count; i++ {
		item := <-inputChan
		item.ChangeStatus = 0
		err := dbmanager.UpdateOrInsertItem(item)
		if err != nil {
			errChan <- err
		}
	}
	select {
	case <-inputChan:
		errChan <- errors.New("too many items")
	default:
		errChan <- nil
	}
}
func updateStatusAndNotify(inputChan <-chan dbmanager.Item, count int, errChan chan<- error) {
	for i := 0; i < count; i++ {
		item := <-inputChan
		subscriptions, err := dbmanager.ReadSubscriptionsForItem(item.Id)
		if err != nil {
			errChan <- err
		}
		for i := range subscriptions {
			if (!subscriptions[i].IsNew && item.ChangeStatus == 1) || item.ChangeStatus == 2 {
				go netmanager.Notify(subscriptions[i].Mail, item.Url, item.Price)
			}
			if item.ChangeStatus > 2 {
				errChan <- errors.New("unknown change status")
			}
			err := dbmanager.MakeSubscriptionOld(subscriptions[i])
			if err != nil {
				errChan <- err
			}
		}
	}
	select {
	case <-inputChan:
		errChan <- errors.New("too many items")
	default:
		errChan <- nil
	}
}
func UpdateAll() error {
	items, err := dbmanager.ReadAllItems()
	if err != nil {
		return err
	}
	chanelCount := 1

	parseOutputChan := make(chan netmanager.ParseResult, chanelCount)
	enqueueParse(items, parseOutputChan, chanelCount)

	updatePriceChan := make(chan dbmanager.Item)
	updateStatusAndNotifyChan := make(chan dbmanager.Item)
	errChan1 := make(chan error)
	errChan2 := make(chan error)
	go updateItems(updatePriceChan, len(items), errChan1)
	go updateStatusAndNotify(updateStatusAndNotifyChan, len(items), errChan2)
	for range items {
		parseResult := <-parseOutputChan
		if parseResult.Err != nil {
			return parseResult.Err
		}
		if parseResult.NewPrice != parseResult.Price {
			parseResult.ChangeStatus++
			updatePriceChan <- parseResult.Item
		}
		updateStatusAndNotifyChan <- parseResult.Item
	}
	requirment := 0
	for {
		select {
		case err := <-errChan1:
			if err != nil {
				return err
			}
			requirment++
			if requirment == 2 {
				return nil
			}

		case err := <-errChan2:
			if err != nil {
				return err
			}
			requirment++
			if requirment == 2 {
				return nil
			}

		}
	}
	return nil
}
