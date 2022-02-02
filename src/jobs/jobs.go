package jobs

import (
	"RestService/src/logger"
	"RestService/src/model"
	"RestService/src/web"
	"sync"
)

// check all sites in database and update prices if changed
func CheckSites() {
	var items []model.Product
	// grab all items from database
	result := model.DB.Find(&items)

	if result.Error != nil {
		logger.Log.Error().
			Err(result.Error).
			Msg("Failed to fetch items from the db")

		return
	}

	logger.Log.Info().
		Int("number of products to check", len(items)).
		Msg("Checking all items")

	if len(items) > 0 {
		var group sync.WaitGroup
		// iterate over all products
		// spawn a goroutine and add it to the wait group
		for _, i := range items {
			group.Add(1)
			go checkPrice(i, &group)
		}

		// wait for all routines to finish
		group.Wait()
		logger.Log.Info().Msg("Finished checking all items")
	}
}

// check the price of an item and update the db if the price has changed
func checkPrice(product model.Product, wg *sync.WaitGroup) {
	defer wg.Done()
	price, _ := web.Fetch(product.Url)
	if price != product.CurrentPrice {
		product.CurrentPrice = price
		result := model.DB.Save(&product)
		if result.Error != nil {
			logger.Log.Error().Err(result.Error).
				Uint("item key", product.ID).
				Msg("failed to update model in database!")
		}
	}
}
