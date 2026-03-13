package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/normify-ru/normify-go"
)

func main() {
	client := normify.NewClient("ваш_api_ключ_здесь")

	req := &normify.ProcessRequest{
		Entity: "company_name",
		Data: []normify.DataItem{
			{ID: "1", Value: "ООО Рога и Копыта"},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Process(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Success: %v\n", resp.Success)
	fmt.Printf("Entity: %s\n", resp.Data.Entity)
	for _, out := range resp.Data.Result.Output {
		fmt.Printf("ID: %s, Value: %v, Metadata: %v\n", out.ID, out.Value, out.Metadata)
	}
}
