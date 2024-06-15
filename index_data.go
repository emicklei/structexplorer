package structexplorer

type (
	indexData struct {
		Rows []tableRow
	}
	tableRow struct {
		Cells []fieldList
	}
	fieldList struct {
		Label      string
		Path       string
		Row        int
		Column     int
		Type       string
		Access     string
		Fields     []fieldAccess
		SelectSize int
		SelectID   string
	}
)
