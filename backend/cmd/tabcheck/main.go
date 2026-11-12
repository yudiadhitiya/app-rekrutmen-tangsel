// Command tabcheck lists the tabs of the configured spreadsheet.
//
// Temporary diagnostic tool; delete after use.
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"rekrutmen-tangsel/backend/internal/store/sheets"
)

func main() {
	st, err := sheets.New(sheets.Config{
		SpreadsheetID: os.Getenv("GOOGLE_SHEET_ID"),
		KeyFile:       os.Getenv("GOOGLE_KEY_FILE"),
	})
	if err != nil {
		fmt.Println("new:", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	mode := ""
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	switch mode {
	case "clear":
		tab := os.Args[2]
		row := 0
		fmt.Sscanf(os.Args[3], "%d", &row)
		if err := st.ClearRow(tab, row); err != nil {
			fmt.Println("clear:", err)
			return
		}
		fmt.Println("cleared:", tab, row)
	case "find":
		tab := os.Args[2]
		needle := os.Args[3]
		h, rows, err := st.ReadTab(ctx, tab)
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		for i, r := range rows {
			for j, c := range r {
				if strings.EqualFold(strings.TrimSpace(c), needle) {
					name := fmt.Sprintf("col%d", j)
					if j < len(h) {
						name = h[j]
					}
					fmt.Printf("row %d (%s):\n", i+2, name)
					for k, v := range r {
						col := fmt.Sprintf("col%d", k)
						if k < len(h) {
							col = h[k]
						}
						if len(v) > 16 {
							v = v[:16] + "…"
						}
						fmt.Printf("  %-20s = %q\n", col, v)
					}
					return
				}
			}
		}
		fmt.Println("tidak ditemukan:", needle)
	case "tabs":
		tabs, err := st.Tabs(ctx)
		if err != nil {
			fmt.Println("tabs:", err)
			return
		}
		for _, t := range tabs {
			fmt.Println("tab:", t)
		}
	case "row":
		tab := os.Args[2]
		h, rows, err := st.ReadTab(ctx, tab)
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		fmt.Println("headers:", h)
		for i, r := range rows {
			if i >= 2 {
				break
			}
			fmt.Printf("--- row %d (len=%d)\n", i+1, len(r))
			for j, c := range r {
				name := fmt.Sprintf("col%d", j)
				if j < len(h) {
					name = h[j]
				}
				if len(c) > 12 {
					c = c[:12] + "…"
				}
				fmt.Printf("  %-20s = %q\n", name, c)
			}
		}
		fmt.Println("total rows:", len(rows))
	case "users":
		h, rows, err := st.ReadTab(ctx, "USERS")
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		fmt.Println("headers:", h)
		idx := func(name string) int {
			for i, c := range h {
				if strings.EqualFold(strings.TrimSpace(c), name) {
					return i
				}
			}
			return -1
		}
		iu, ir, ia := idx("USERNAME"), idx("ROLE"), idx("AKTIF")
		n := 0
		for _, r := range rows {
			get := func(i int) string {
				if i >= 0 && i < len(r) {
					return r[i]
				}
				return ""
			}
			fmt.Printf("  %-28s role=%-18s aktif=%s\n", get(iu), get(ir), get(ia))
			n++
			if n >= 15 {
				fmt.Println("  ...")
				break
			}
		}
		fmt.Println("total user:", len(rows))
	}
}
