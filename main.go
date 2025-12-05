package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
)

func main() {
	color.NoColor = false

	var buffer bytes.Buffer

	wr := io.MultiWriter(&buffer, os.Stdout)
	fmt.Fprintln(wr, "splitters")

	if err := execute(wr); err != nil {
		fmt.Fprintln(wr, err.Error())
	}

	if err := saveHTML("splitter.html", buffer.String()); err != nil {
		fmt.Fprintf(os.Stderr, "failed to save HTML output: %s", err)
	}
}

func execute(w io.Writer) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	if len(os.Args) == 2 {
		dir = os.Args[1]
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %q: %w", dir, err)
	}

	_, _ = w.Write([]byte(dir))

	for _, file := range files {
		if err := processFile(w, dir, file); err != nil {
			fmt.Fprintf(w, "processing file %q failed: %s", file.Name(), err)
		}
	}

	return nil
}

func processFile(w io.Writer, dir string, file os.DirEntry) error {
	if file.IsDir() {
		// папки пропускаем
		return nil
	} else if !strings.HasSuffix(file.Name(), ".xlsx") {
		// файлы с другими расширениями тоже пропускаем
		return nil
	}

	// Складываем имя файла и имя папки в один путь.
	path := filepath.Join(dir, file.Name())

	// Пишем загадочную умную надпись
	if _, err := fmt.Fprintln(w, "Processing", path); err != nil {
		return err
	}

	// Считываем исходные данные из найденного файла.
	floors, splitters, err := parseFile(path)
	if err != nil {
		return err
	}

	// Распределяем квартиры по разделителям.
	splitters = calculate(floors, splitters)

	// Отобразить план этажей в терминале
	if err := displayFloors(w, floors, splitters); err != nil {
		return err
	}

	// Отобразить разделители с соответствующими квартирами
	if err := displayDividers(w, splitters); err != nil {
		return err
	}

	// Записываем результаты в файлек
	if err := writeResults(path, splitters); err != nil {
		return err
	}

	return nil
}

func saveHTML(path, content string) error {
	table, err := buildHTMLTable(content)
	if err != nil {
		return err
	}

	const htmlTemplate = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>Splitters output</title>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/ansi_up/5.2.1/ansi_up.min.js" integrity="sha512-24F1cCih5Hp8CZXlOKfF/QUfPnouPnK9C+lCW29FEJTgx6x8V7POvloWTaBI15ruBeJYMi/0046BX6XjwDv/2w==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
    <style>
        body { font-family: "Fira Code", "JetBrains Mono", Menlo, Consolas, monospace; background-color: #0f0f0f; color: #e8e8e8; padding: 24px; }
        h1 { color: #f5f5f5; }
        .log-table { border-collapse: collapse; width: 100%; }
        .log-table th, .log-table td { background: #1c1c1c; border: 1px solid #2a2a2a; padding: 8px 12px; vertical-align: top; }
        .log-table th { text-align: left; color: #f5f5f5; }
        .log-table tr:nth-child(even) td { background: #161616; }
    </style>
</head>
<body>
    <h1>Splitters output</h1>
    <div id="log-table">{{.Table}}</div>
    <script>
        const ansiUp = new AnsiUp();
        document.querySelectorAll('#log-table td').forEach((cell) => {
            const raw = cell.textContent;
            cell.innerHTML = ansiUp.ansi_to_html(raw);
        });
    </script>
</body>
</html>`

	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse HTML template: %w", err)
	}

	var output bytes.Buffer
	err = tmpl.Execute(&output, struct{ Table template.HTML }{Table: template.HTML(table)})
	if err != nil {
		return fmt.Errorf("failed to build HTML: %w", err)
	}

	return os.WriteFile(path, output.Bytes(), 0o644)
}

func buildHTMLTable(content string) (string, error) {
	tableBuffer := &bytes.Buffer{}
	htmlRenderer := renderer.NewHTML(renderer.HTMLConfig{
		TableCSSClass: "log-table",
		EscapeContent: true,
	})

	table := tablewriter.NewWriter(tableBuffer, tablewriter.WithRenderer(htmlRenderer))
	table.SetHeader([]string{"Output"})
	table.SetAutoWrapText(false)
	table.SetBorder(false)
	table.SetHeaderLine(true)
	table.SetRowLine(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, line := range strings.Split(strings.TrimSuffix(content, "\n"), "\n") {
		table.Append([]string{line})
	}

	table.Render()
	return tableBuffer.String(), nil
}
