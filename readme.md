# Price Watcher Rest Service

price watcher is an application that allows users to pass a url to an item on ebay/amazon and get notified when the price
changes.
The backend api is written in Go with the [Gin](https://github.com/gin-gonic/gin) and [Gorm](https://github.com/go-gorm/gorm) libraries

Go was choosen for this project as 
- it's lightweight(compared to using c# + asp.net as i originally planned on doing) 
- easy to distribute with it's single binary builds
