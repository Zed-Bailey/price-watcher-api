package jobs

import (
	"RestService/src/model"
	"sync"

	"github.com/rs/zerolog/log"
)

// check all sites in database and update prices if changed
func CheckSites() {
	var items []model.Product
	// grab all items from database
	result := model.DB.Find(&items)

	if result.Error != nil {
		log.Error().
			Str("ERROR", result.Error.Error())
		return
	}

	log.Info().
		Str("message", "Checking all items").
		Int("number of products to check", len(items))

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
		log.Info().
			Str("message", "Finished checking all items")
	}
}

func checkPrice(product model.Product, wg *sync.WaitGroup) {
	defer wg.Done()

}
