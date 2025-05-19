package pdf

import (
	"fmt"
	"health/internal/domain/models"
	"strconv"

	"github.com/jung-kurt/gofpdf"
)

func ExportToPDF(stats []*models.Statistics, adminName, adminSurname, depTitle string) *gofpdf.Fpdf {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	// Шрифты
	pdf.AddUTF8Font("timesnrcyrmt", "", "lib/pdf/timesnrcyrmt.ttf")
	pdf.AddUTF8Font("timesnrcyrmt", "B", "lib/pdf/timesnrcyrmt_bold.ttf")

	depName := fmt.Sprintf("Статистика выполнения гимнастики. Отдел: %v", depTitle)
	adminInfo := fmt.Sprintf("Руководитель отдела: %v %v", adminName, adminSurname)
	// Заголовок отчета
	pdf.SetFont("timesnrcyrmt", "B", 16)
	pdf.Cell(0, 10, depName)
	pdf.Ln(12)

	// Руководитель отдела
	pdf.SetFont("timesnrcyrmt", "", 12)
	pdf.Cell(0, 10, adminInfo)
	pdf.Ln(15)

	// Настройки таблицы
	header := []string{"Сотрудник", "Статус", "Просмотрено (%)", "Видео", "Дата выполнения"}
	columnWidths := []float64{50, 40, 60, 50, 50} // Сумма ~250mm (формат A4 landscape - 297mm)

	// Цвет фона заголовков
	pdf.SetFillColor(200, 200, 200)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetDrawColor(128, 128, 128)
	pdf.SetLineWidth(.3)
	pdf.SetFont("timesnrcyrmt", "B", 10)

	// Заголовки таблицы
	for i, str := range header {
		pdf.CellFormat(columnWidths[i], 7, str, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Данные таблицы
	pdf.SetFont("timesnrcyrmt", "", 10)
	pdf.SetFillColor(255, 255, 255)
	pdf.SetTextColor(0, 0, 0)

	for _, stat := range stats {
		percentView := strconv.Itoa(int(stat.PercentageView)) + "%"

		// Сотрудник
		employee := stat.EmployeeName + " " + stat.EmployeeSurname
		pdf.CellFormat(columnWidths[0], 6, employee, "1", 0, "L", false, 0, "")

		// Статус
		pdf.CellFormat(columnWidths[1], 6, stat.Progress, "1", 0, "L", false, 0, "")

		// Процент просмотра
		pdf.CellFormat(columnWidths[2], 6, percentView, "1", 0, "C", false, 0, "")

		// Видео
		pdf.CellFormat(columnWidths[3], 6, stat.MediaTitle, "1", 0, "L", false, 0, "")

		// Дата выполнения
		pdf.CellFormat(columnWidths[4], 6, stat.CreatedAt.Format("02.01.2006 15:04"), "1", 0, "L", false, 0, "")

		pdf.Ln(-1)
	}

	return pdf
}
