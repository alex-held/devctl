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
	fmt.Fprintf(Out, "Install + %v", args[0])
	return nil
}

func Download(ctx context.Context, args []string) error {
	fmt.Fprintf(Out, "Download + %v", args[0])
	return nil
}

func List(ctx context.Context, args []string) error {
	fmt.Fprintf(Out, "List + %v", args[0])
	return nil
}

func Current(ctx context.Context, args []string) error {
	fmt.Fprintf(Out, "Current + %v", args[0])
	return nil
}

func Use(ctx context.Context, args []string) error {
	fmt.Fprintf(Out, "Use + %v", args[0])
	return nil
}
