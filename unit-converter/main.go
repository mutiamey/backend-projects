package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type PageData struct {
	ActiveTab string
	Value     string
	FromUnit  string
	ToUnit    string
	Result    string
	HasResult bool
	Error     string
	Units     []UnitOption
}

type UnitOption struct {
	Value string
	Label string
}

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Unit Converter</title>
    <style>
        body { font-family: system-ui, -apple-system, sans-serif; max-width: 480px; margin: 40px auto; padding: 24px; border: 2px solid #222; border-radius: 12px; background: #fff; }
        h1 { margin-top: 0; font-size: 28px; }
        .tabs { display: flex; gap: 16px; border-bottom: 2px solid #ddd; padding-bottom: 8px; margin-bottom: 24px; }
        .tab { text-decoration: none; color: #666; font-weight: bold; font-size: 18px; padding-bottom: 4px; }
        .tab.active { color: #2563eb; border-bottom: 3px solid #2563eb; }
        .form-group { margin-bottom: 16px; }
        label { display: block; margin-bottom: 6px; font-weight: 600; color: #333; }
        input, select { width: 100%; padding: 10px; border: 2px solid #333; border-radius: 6px; font-size: 16px; box-sizing: border-box; }
        .btn { width: auto; padding: 10px 24px; background: #fff; border: 2px solid #333; border-radius: 6px; font-weight: bold; cursor: pointer; font-size: 16px; margin-top: 8px; }
        .btn:hover { background: #f0f0f0; }
        .result-box { margin-top: 20px; }
        .result-box h3 { color: #444; font-weight: normal; margin-bottom: 8px; }
        .result-text { font-size: 28px; font-weight: bold; margin: 12px 0 20px 0; }
        .error { color: #dc2626; margin-top: 12px; font-weight: 500; }
    </style>
</head>
<body>
    <h1>Unit Converter</h1>
    
    <div class="tabs">
        <a href="/length" class="tab {{if eq .ActiveTab "length"}}active{{end}}">Length</a>
        <a href="/weight" class="tab {{if eq .ActiveTab "weight"}}active{{end}}">Weight</a>
        <a href="/temperature" class="tab {{if eq .ActiveTab "temperature"}}active{{end}}">Temperature</a>
    </div>

    {{if .HasResult}}
        <div class="result-box">
            <h3>Result of your calculation</h3>
            <div class="result-text">{{.Result}}</div>
            <a href="/{{.ActiveTab}}"><button class="btn">Reset</button></a>
        </div>
    {{else}}
        <form action="/{{.ActiveTab}}" method="POST">
            <div class="form-group">
                <label>Enter the {{.ActiveTab}} to convert</label>
                <input type="number" step="any" name="value" required value="{{.Value}}">
            </div>
            
            <div class="form-group">
                <label>Unit to Convert from</label>
                <select name="from">
                    {{range .Units}}
                        <option value="{{.Value}}" {{if eq $.FromUnit .Value}}selected{{end}}>{{.Label}}</option>
                    {{end}}
                </select>
            </div>

            <div class="form-group">
                <label>Unit to Convert to</label>
                <select name="to">
                    {{range .Units}}
                        <option value="{{.Value}}" {{if eq $.ToUnit .Value}}selected{{end}}>{{.Label}}</option>
                    {{end}}
                </select>
            </div>

            <button type="submit" class="btn">Convert</button>
        </form>
    {{end}}

    {{if .Error}}
        <div class="error">{{.Error}}</div>
    {{end}}
</body>
</html>
`

func getUnits(category string) []UnitOption {
	switch category {
	case "length":
		return []UnitOption{
			{"mm", "Millimeter (mm)"},
			{"cm", "Centimeter (cm)"},
			{"m", "Meter (m)"},
			{"km", "Kilometer (km)"},
			{"in", "Inch (in)"},
			{"ft", "Foot (ft)"},
			{"yd", "Yard (yd)"},
			{"mi", "Mile (mi)"},
		}
	case "weight":
		return []UnitOption{
			{"mg", "Milligram (mg)"},
			{"g", "Gram (g)"},
			{"kg", "Kilogram (kg)"},
			{"oz", "Ounce (oz)"},
			{"lb", "Pound (lb)"},
		}
	case "temperature":
		return []UnitOption{
			{"C", "Celsius (°C)"},
			{"F", "Fahrenheit (°F)"},
			{"K", "Kelvin (K)"},
		}
	default:
		return nil
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/length", http.StatusSeeOther)
	})

	http.HandleFunc("/length", func(w http.ResponseWriter, r *http.Request) { handleCategory(w, r, "length") })
	http.HandleFunc("/weight", func(w http.ResponseWriter, r *http.Request) { handleCategory(w, r, "weight") })
	http.HandleFunc("/temperature", func(w http.ResponseWriter, r *http.Request) { handleCategory(w, r, "temperature") })

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func handleCategory(w http.ResponseWriter, r *http.Request, category string) {
	units := getUnits(category)
	tmpl := template.Must(template.New("index").Parse(htmlTemplate))

	if r.Method == http.MethodGet {
		tmpl.Execute(w, PageData{
			ActiveTab: category,
			Units:     units,
		})
		return
	}

	if r.Method == http.MethodPost {
		valStr := r.FormValue("value")
		from := r.FormValue("from")
		to := r.FormValue("to")

		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			tmpl.Execute(w, PageData{
				ActiveTab: category,
				Units:     units,
				Error:     "Please enter a valid number.",
			})
			return
		}

		resVal, err := convert(val, category, from, to)
		if err != nil {
			tmpl.Execute(w, PageData{
				ActiveTab: category,
				Units:     units,
				Error:     err.Error(),
			})
			return
		}

		// Format output e.g., "20 ft = 609.6 cm"
		resFormatted := fmt.Sprintf("%g %s = %g %s", val, from, math.Round(resVal*10000)/10000, to)

		tmpl.Execute(w, PageData{
			ActiveTab: category,
			Value:     valStr,
			FromUnit:  from,
			ToUnit:    to,
			Result:    resFormatted,
			HasResult: true,
			Units:     units,
		})
	}
}

func convert(val float64, category, from, to string) (float64, error) {
	switch category {
	case "length":
		return convertLength(val, from, to)
	case "weight":
		return convertWeight(val, from, to)
	case "temperature":
		return convertTemperature(val, from, to)
	default:
		return 0, fmt.Errorf("invalid category")
	}
}

func convertLength(val float64, from, to string) (float64, error) {
	rates := map[string]float64{
		"mm": 0.001,
		"cm": 0.01,
		"m":  1.0,
		"km": 1000.0,
		"in": 0.0254,
		"ft": 0.3048,
		"yd": 0.9144,
		"mi": 1609.344,
	}

	fromRate, ok1 := rates[from]
	toRate, ok2 := rates[to]
	if !ok1 || !ok2 {
		return 0, fmt.Errorf("invalid length unit")
	}

	return (val * fromRate) / toRate, nil
}

func convertWeight(val float64, from, to string) (float64, error) {
	rates := map[string]float64{
		"mg": 0.000001,
		"g":  0.001,
		"kg": 1.0,
		"oz": 0.028349523125,
		"lb": 0.45359237,
	}

	fromRate, ok1 := rates[from]
	toRate, ok2 := rates[to]
	if !ok1 || !ok2 {
		return 0, fmt.Errorf("invalid weight unit")
	}

	return (val * fromRate) / toRate, nil
}

func convertTemperature(val float64, from, to string) (float64, error) {
	var celsius float64

	switch from {
	case "C":
		celsius = val
	case "F":
		celsius = (val - 32) * 5 / 9
	case "K":
		celsius = val - 273.15
	default:
		return 0, fmt.Errorf("invalid temperature unit")
	}

	switch to {
	case "C":
		return celsius, nil // Tambahkan ', nil'
	case "F":
		return (celsius * 9 / 5) + 32, nil // Tambahkan ', nil'
	case "K":
		return celsius + 273.15, nil // Tambahkan ', nil'
	default:
		return 0, fmt.Errorf("invalid temperature unit")
	}
}
