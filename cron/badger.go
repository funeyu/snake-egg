package cron

import (
	"os"
)

func Remove() {
	os.RemoveAll("../badger")
}

func refresh() {
	//store := store.InitBadger("./badger")

}