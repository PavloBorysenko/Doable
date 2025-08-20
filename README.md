# Dobble Card Generator

A Windows GUI application that generates optimal Dobble (Spot It!) cards using mathematical projective plane principles.

## Features

-   **Simple Input**: Enter only the number of symbols per card
-   **Automatic Calculation**: Optimal symbol count and card generation
-   **Interactive GUI**: Built with Fyne framework
-   **Card Validation**: Ensures exactly one common symbol between any two cards
-   **Processing Tracker**: Checkbox system to mark processed cards

## How It Works

For **n** symbols per card (where **q = n-1**):

-   **Total symbols** = q² + q + 1
-   **Total cards** = q² + q + 1
-   Each symbol appears exactly **n** times

## Supported Values

| Symbols per card | Total symbols | Total cards | Note                 |
| ---------------- | ------------- | ----------- | -------------------- |
| 2                | 3             | 3           | Basic case           |
| 3                | 7             | 7           | Classic Dobble       |
| 4                | 13            | 13          | Optimal              |
| 8                | 57            | 57          | Original Dobble size |
| 9                | 73            | 73          | Optimal              |

## Usage

1. Run `dobble-generator.exe`
2. Enter number of symbols per card (2-15)
3. Click "Generate Cards"
4. View generated cards and mark as processed

## Example

**Input**: 3 symbols per card  
**Output**: 7 symbols, 7 cards  
**Cards**: [1,2,3], [1,4,5], [1,6,7], [2,4,6], [2,5,7], [3,4,7], [3,5,6]

## Build

```bash
go mod tidy
go build -o dobble-generator.exe
```

## Requirements

-   Windows 10+
-   Go 1.21+ (for building)
