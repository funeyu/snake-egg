package cron

import (
	"os"
	"snake/store"
)

func Remove() {
	os.RemoveAll("../badger")
}

func refresh() {
	store := store.InitBadger("./badger")

}