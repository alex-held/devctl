package main

import (
	"context"
	"fmt"
	"io"
	"os"
)

var (
	//Out is the io.Writer to log standard messages
	// defaults to os.Stdout
	Out io.Writer = os.Stdout
)

func SetStdout(w io.Writer) error {
	Out = w
	return nil
}

func PluginName() string {
	fmt.Fprintf(Out, "sdk-01")
	return "sdk-01"
}

func Install(ctx context.Context, args []string) error {
	fmt.Fprintf(Out, "[%s]\tInstall(%T, %v)", PluginName(), ctx, args)
	return nil
}

func Download(ctx context.Context, args []string) error {
	fmt.Printf("\n%s.Download called. args=%v\n", PluginName(), args)
	return nil
}

func List(ctx context.Context, args []string) error {
	fmt.Printf("\n%s.List called. args=%v\n", PluginName(), args)
	return nil
}

func Current(ctx context.Context, args []string) error {
	fmt.Println(args[0])
	return nil
}

func Use(ctx context.Context, args []string) error {
	fmt.Printf("\n%s.Use called. args=%v\n", PluginName(), args)
	return nil
}
