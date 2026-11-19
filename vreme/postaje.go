// Package vreme omogoča beleženje temperaturnih meritev vremenskih postaj
// in izpis osnovne statistike (najnižja, najvišja in povprečna temperatura).
//
// Postaje hranimo v slovarju, katerega ključ je oznaka postaje. Polja
// strukture [Postaja] so neizvožena, zato novo postajo ustvarimo s funkcijo
// [DodajPostajo].
//
// # Primer uporabe
//
//	postaje := make(map[string]vreme.Postaja)
//
//	vreme.DodajPostajo(postaje, "MB01", "Maribor Tabor", "Maribor")
//	vreme.DodajMeritev(postaje, "MB01", 9.6, -50, 50)
//	vreme.DodajMeritev(postaje, "MB01", 3.4, -50, 50)
//	vreme.DodajMeritev(postaje, "MB01", 2.7, -50, 50)
//
//	vreme.IzpisMeritev(postaje)
//	vreme.IzpisPorocila(postaje, 3)
package vreme

import (
	f "fmt"
	"sort"
)

// Postaja predstavlja vremensko postajo z imenom, krajem in seznamom meritev.
// Polja so neizvožena, zato postajo ustvarimo s funkcijo [DodajPostajo].
type Postaja struct {
	ime     string
	kraj    string
	meritve []float64
}

// DodajMeritev postaji z dano oznako doda novo meritev temperature.
// Meritev je veljavna, če leži med minTemp in maxTemp (vključno).
// Če postaja ne obstaja ali meritev ni veljavna, funkcija izpiše
// obvestilo in meritve ne doda.
func DodajMeritev(postaje map[string]Postaja, oznaka string, temperatura float64, minTemp, maxTemp float64) {
	p, obstaja := postaje[oznaka] // v p se vrne value in v obstaja bool

	if !obstaja {
		f.Printf("Postaja z oznako %s ne obstaja\n", oznaka)
		return
	}
	if temperatura > maxTemp || temperatura < minTemp {
		f.Printf("Meritev %.1f ni veljavna\n", temperatura)
		return
	}

	// sicer gre za veljavno meritev in jo dodamo k podani postaji
	p.meritve = append(p.meritve, temperatura)
	postaje[oznaka] = p
}

// statistika vrne najnižjo, najvišjo in povprečno temperaturo postaje.
// Če postaja ne obstaja ali ima manj kot stMeritev meritev, vrne ok == false.
func statistika(postaje map[string]Postaja, oznaka string, stMeritev int) (najnizja, najvisja, povprecna float64, ok bool) {
	p, obstaja := postaje[oznaka]

	if !obstaja || len(p.meritve) < stMeritev {
		return 0, 0, 0, false
	}

	najnizja = p.meritve[0]
	najvisja = p.meritve[0]
	sum := 0.0

	for _, m := range p.meritve {
		if m < najnizja {
			najnizja = m
		}
		if m > najvisja {
			najvisja = m
		}
		sum += m
	}
	povprecna = sum / float64(len(p.meritve))
	return najnizja, najvisja, povprecna, true
}

// IzpisMeritev izpiše oznake, imena, kraje in meritve vseh postaj,
// urejene po oznaki postaje.
func IzpisMeritev(postaje map[string]Postaja) {
	// izluscimo vse kljuce
	oznake := make([]string, 0, len(postaje)) // <-- prazna rezina s kapaciteto dolzine slovarja

	for o := range postaje {
		oznake = append(oznake, o)
	}

	// jih sortiramo po abecedi
	sort.Strings(oznake)

	// izpisemo
	f.Println("MERITVE:")
	for _, o := range oznake {
		p := postaje[o]
		f.Printf("%s - %s (%s): %v\n", o, p.ime, p.kraj, p.meritve)
	}
}

// IzpisPorocila za vsako postajo izpiše najnižjo, najvišjo in povprečno
// temperaturo ter opis glede na povprečje (vroče, toplo, hladno ali mrzlo).
// Postaje z manj kot stMeritev meritvami so izpisane z obvestilom
// "premalo meritev". Postaje so urejene po oznaki.
func IzpisPorocila(postaje map[string]Postaja, stMeritev int) {
	// sortiramo zato da bo izpis vedno enak
	oznake := make([]string, 0, len(postaje))
	for o := range postaje {
		oznake = append(oznake, o)
	}
	sort.Strings(oznake)

	for _, o := range oznake {
		p := postaje[o]
		lo, hi, avg, ok := statistika(postaje, o, stMeritev)
		if !ok {
			f.Printf("%s: premalo meritev\n", p.ime)
			continue
		}

		var nominalno string
		switch {
		case avg >= 25:
			nominalno = "Vroče"
		case avg >= 15:
			nominalno = "Toplo"
		case avg >= 5:
			nominalno = "Hladno"
		default:
			nominalno = "Mrzlo"
		}
		f.Printf("%s: min %.1f °C, max %.1f °C, povprečje %.1f °C -> %s\n", p.ime, lo, hi, avg, nominalno)
	}
}

// DodajPostajo v slovar postaj doda novo postajo z dano oznako, imenom in
// krajem. Nova postaja nima meritev. Če postaja s to oznako že obstaja,
// funkcija izpiše obvestilo in slovarja ne spremeni.
func DodajPostajo(postaje map[string]Postaja, oznaka, ime, kraj string) {
	_, obstaja := postaje[oznaka]
	if obstaja {
		f.Printf("Postaja z oznako %s že obstaja\n", oznaka)
		return
	}
	postaje[oznaka] = Postaja{ime, kraj, []float64{}} // []float64{} ustvari prazno rezino, ki bo kasneje vsebovala meritve
}
