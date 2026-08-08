# Unit Converter

A web-based Unit Converter application built with Go standard library (`net/http` and `html/template`) to convert measurements for Length, Weight, and Temperature.

This project is built to complete the **[Unit Converter Challenge](https://roadmap.sh/projects/unit-converter)** from **roadmap.sh**.

## Features

- **Length Conversion**: Millimeter (mm), Centimeter (cm), Meter (m), Kilometer (km), Inch (in), Foot (ft), Yard (yd), Mile (mi).
- **Weight Conversion**: Milligram (mg), Gram (g), Kilogram (kg), Ounce (oz), Pound (lb).
- **Temperature Conversion**: Celsius (°C), Fahrenheit (°F), Kelvin (K).
- Interactive tabbed interface (`Length`, `Weight`, `Temperature`).
- Clear display of calculation results with a Reset button.
- Built strictly with the **Go Standard Library** (zero external dependencies).

## How to Run

Navigate into this folder and run:
```bash
go run main.go
```