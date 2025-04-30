package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v3"
)

func main() {
	headerFmt := color.New(color.FgCyan, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	apiKey := os.Getenv("API_KEY")
	bridge := os.Getenv("BRIDGE_URL")
	if apiKey == "" && bridge == "" {
		log.Fatal("API_KEY and BRIDGE_URL environment variables are required for the tool to work")
	}
	cmd := &cli.Command{
		Name:  "tcli",
		Usage: "CLI application for the Tedee ecosystem",
		Commands: []*cli.Command{

			{

				Name:    "bridge-state",
				Aliases: []string{"bs"},
				Usage:   "Check the bridge status",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					endpoint := "v1.0/bridge"
					outJsonByte := getHandler([]byte(""), tokenGenerator(apiKey), bridge, endpoint)
					var returnMessage map[string]any
					if err := json.Unmarshal(outJsonByte, &returnMessage); err != nil {
						log.Printf("Failed to unmarshal JSON: %v", err)
					}
					tbl := table.New("Name", "S/N", "SSID", "Firmware version", "Wi-Fi version")
					tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
					tbl.AddRow(returnMessage["name"], returnMessage["serialNumber"], returnMessage["ssid"], returnMessage["version"], returnMessage["wifiVersion"])
					tbl.Print()
					return nil
				},
			},
			{
				Name:    "lock",
				Aliases: []string{"l"},
				Usage:   "Lock operation",
				Commands: []*cli.Command{
					{
						Name:  "list",
						Usage: "List all locks",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							endpoint := "v1.0/lock"
							outJsonByte := getHandler([]byte(""), tokenGenerator(apiKey), bridge, endpoint)
							var returnMessage []map[string]any
							if err := json.Unmarshal(outJsonByte, &returnMessage); err != nil {
								log.Printf("Failed to unmarshal JSON: %v", err)
							}
							tbl := table.New("Name", "Type", "ID", "Battery")
							tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
							for _, lock := range returnMessage {
								tbl.AddRow(lock["name"], lock["type"], lock["id"], lock["batteryLevel"])
							}
							tbl.Print()
							return nil
						},
					},
					{
						Name:  "details",
						Usage: "Retrive information about a specific lock using its ID",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							endpoint := fmt.Sprintf("v1.0/lock/%s,", cmd.Args().Get(0))
							outJsonByte := getHandler([]byte(""), tokenGenerator(apiKey), bridge, endpoint)
							var returnMessage map[string]any
							if err := json.Unmarshal(outJsonByte, &returnMessage); err != nil {
								log.Printf("Failed to unmarshal JSON: %v", err)
							}
							tbl := table.New("Name", "S/N", "Type", "Firmware version", "ID", "Battery", "Connection Status", "RSSI")
							tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
							tbl.AddRow(returnMessage["name"], returnMessage["serialNumber"], returnMessage["type"], returnMessage["version"],
								returnMessage["id"], returnMessage["batteryLevel"], returnMessage["isConnected"], returnMessage["rssi"])
							tbl.Print()
							return nil
						},
					},
					{
						Name:  "lock",
						Usage: "Lock specific lock using its ID or use fav",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							fav := os.Getenv("FAV")
							if fav == "" {
								endpoint := fmt.Sprintf("v1.0/lock/%s/lock", cmd.Args().Get(0))
								ouStatus := postHandler([]byte(""), tokenGenerator(apiKey), bridge, endpoint)
								retCodeHandler(ouStatus)
							} else {
								endpoint := fmt.Sprintf("v1.0/lock/%s/lock", fav)
								ouStatus := postHandler([]byte(""), tokenGenerator(apiKey), bridge, endpoint)
								retCodeHandler(ouStatus)
							}

							return nil
						},
					},
					{
						Name:  "unlock",
						Usage: "Unlock specific lock using its ID or use fav",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							fav := os.Getenv("FAV")
							if fav == "" {
								endpoint := fmt.Sprintf("v1.0/lock/%s/unlock", cmd.Args().Get(0))
								ouStatus := postHandler([]byte(""), tokenGenerator(apiKey), bridge, endpoint)
								retCodeHandler(ouStatus)
							} else {
								endpoint := fmt.Sprintf("v1.0/lock/%s/unlock", fav)
								ouStatus := postHandler([]byte(""), tokenGenerator(apiKey), bridge, endpoint)
								retCodeHandler(ouStatus)
							}

							return nil
						},
					},
					{
						Name:  "pull",
						Usage: "pull the latch specific lock using its ID",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							endpoint := fmt.Sprintf("v1.0/lock/%s/pull", cmd.Args().Get(0))
							ouStatus := postHandler([]byte(""), tokenGenerator(apiKey), bridge, endpoint)
							retCodeHandler(ouStatus)
							return nil
						},
					},
				},
			},

			{
				Name:    "callback",
				Aliases: []string{"c"},
				Usage:   "Callback operation",
				Commands: []*cli.Command{
					{
						Name:  "list",
						Usage: "List all callbacks",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							endpoint := "v1.0/callback"
							outJsonByte := getHandler([]byte(""), tokenGenerator(apiKey), bridge, endpoint)
							var returnMessage []map[string]any
							if err := json.Unmarshal(outJsonByte, &returnMessage); err != nil {
								log.Printf("Failed to unmarshal JSON: %v", err)
							}
							tbl := table.New("ID", "URL", "METHOD")
							tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
							for _, callback := range returnMessage {
								tbl.AddRow(callback["id"], callback["url"], callback["method"])
							}
							tbl.Print()
							return nil
						},
					},
					{
						Name:  "add",
						Usage: "Add a callback endpoint",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							endpoint := "v1.0/callback"
							entry := CallbackEntry{
								URL:     cmd.Args().Get(0),
								Method:  "POST",
								Headers: []string{"{}"}}
							jsonData, _ := json.Marshal(entry)
							ouStatus := postHandler([]byte(jsonData), tokenGenerator(apiKey), bridge, endpoint)
							retCodeHandler(ouStatus)
							return nil
						},
					},
					{
						Name:  "delete",
						Usage: "Delete a callback endpoint",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							endpoint := fmt.Sprintf("v1.0/callback/%s", cmd.Args().Get(0))
							ouStatus := deleteHandler([]byte(""), tokenGenerator(apiKey), bridge, endpoint)
							retCodeHandler(ouStatus)
							return nil
						},
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
