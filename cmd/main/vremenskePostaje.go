package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/vamp3k/dn5/vreme"
)

func main() {
	cmd := &cli.Command{
		Name:  "vremenske-postaje",
		Usage: "Belezenje temperaturnih meritev in izpis statistike vremenskih postaj.",
		Flags: []cli.Flag{
			&cli.FloatFlag{
				Name:  "minTemp",
				Value: -50,
				Usage: "Najnizja dovoljena temperatura",
			},
			&cli.FloatFlag{
				Name:  "maxTemp",
				Value: 50,
				Usage: "Najvisja dovoljena temperatura",
			},
			&cli.IntFlag{
				Name:  "stMeritev",
				Value: 3,
				Usage: "Najmanjse stevilo meritev za izracun statistike",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			minTemp := cmd.Float("minTemp")
			maxTemp := cmd.Float("maxTemp")
			stMeritev := cmd.Int("stMeritev")

			if minTemp >= maxTemp {
				return fmt.Errorf("minTemp mora biti manjsi od maxTemp")
			}

			if stMeritev < 1 {
				return fmt.Errorf("stMeritev mora biti vsaj 1")
			}

			return izvediProgram(minTemp, maxTemp, stMeritev)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
func izvediProgram(minTemp, maxTemp float64, stMeritev int) error {
	// Ustvarimo slovar postaj
	postaje := make(map[string]vreme.Postaja)

	// Začetne postaje
	vreme.DodajPostajo(postaje, "CE01", "Celje center", "Celje")
	vreme.DodajPostajo(postaje, "MB01", "Maribor Tabor", "Maribor")
	vreme.DodajPostajo(postaje, "KP01", "Koper Luka", "Koper")
	vreme.DodajPostajo(postaje, "LJ01", "Ljubljana Bežigrad", "Ljubljana")

	// Začetne meritve
	vreme.DodajMeritev(postaje, "CE01", 15.7, minTemp, maxTemp)
	vreme.DodajMeritev(postaje, "CE01", 14.2, minTemp, maxTemp)
	vreme.DodajMeritev(postaje, "CE01", 16.3, minTemp, maxTemp)

	vreme.DodajMeritev(postaje, "MB01", 9.6, minTemp, maxTemp)
	vreme.DodajMeritev(postaje, "MB01", 3.4, minTemp, maxTemp)
	vreme.DodajMeritev(postaje, "MB01", 2.7, minTemp, maxTemp)

	vreme.DodajMeritev(postaje, "KP01", 24.5, minTemp, maxTemp)
	vreme.DodajMeritev(postaje, "KP01", 27, minTemp, maxTemp)
	vreme.DodajMeritev(postaje, "KP01", 26.1, minTemp, maxTemp)
	vreme.DodajMeritev(postaje, "KP01", 28.4, minTemp, maxTemp)

	vreme.DodajMeritev(postaje, "LJ01", 2, minTemp, maxTemp)
	vreme.DodajMeritev(postaje, "LJ01", -1.5, minTemp, maxTemp)
	vreme.DodajMeritev(postaje, "LJ01", 3.5, minTemp, maxTemp)

	// 1. IZPIS NA ZAČETKU
	fmt.Println("\nIZPIS NA ZAČETKU")
	fmt.Println("-----------------------------------------------------------")
	vreme.IzpisMeritev(postaje)
	fmt.Println()
	vreme.IzpisPorocila(postaje, stMeritev)

	// 2. IZPIS PO DODAJANJU MERITEV
	fmt.Println("\nIZPIS PO DODAJANJU MERITEV")
	fmt.Println("-----------------------------------------------------------")

	vreme.DodajMeritev(postaje, "KR01", 2, minTemp, maxTemp)
	vreme.DodajMeritev(postaje, "LJ01", -99.7, minTemp, maxTemp)
	vreme.DodajMeritev(postaje, "KP01", 50.8, minTemp, maxTemp)
	vreme.DodajMeritev(postaje, "KP01", -10, minTemp, maxTemp)

	fmt.Println()
	vreme.IzpisMeritev(postaje)
	fmt.Println()
	vreme.IzpisPorocila(postaje, stMeritev)

	// 3. IZPIS PO DODAJANJU POSTAJE
	fmt.Println("\nIZPIS PO DODAJANJU POSTAJE")
	fmt.Println("-----------------------------------------------------------")

	vreme.DodajPostajo(postaje, "NM01", "Novo mesto", "Novo mesto")
	vreme.DodajMeritev(postaje, "NM01", 11.1, minTemp, maxTemp)

	vreme.IzpisMeritev(postaje)
	fmt.Println()
	vreme.IzpisPorocila(postaje, stMeritev)

	return nil
}
