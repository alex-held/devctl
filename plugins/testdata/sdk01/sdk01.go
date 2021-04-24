package main

import (
	"context"
	"fmt"
)

const Name = "sdk01"

func Install(ctx context.Context, args []string) error {
	fmt.Printf("\n%s.Install called. args=%v\n", Name, args)
	return nil
}


func Download(ctx context.Context, args []string) error {
	fmt.Printf("\n%s.Download called. args=%v\n", Name, args)
	return nil
}

func List(ctx context.Context, args []string) error {
	fmt.Printf("\n%s.List called. args=%v\n", Name, args)
	return nil
}
