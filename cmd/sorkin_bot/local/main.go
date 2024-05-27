package main

import (
	"context"
	"log"
	"sorkin_bot/internal"
)

func main() {
	ctx := context.Background()
	log.Fatal(internal.NewApp(ctx).Run(ctx))
}
