package generator

type Options struct {
	Database   string
	Insert     int
	Update     int
	Delete     int
	BatchSize  int
	GoRoutines int
	InsertType string
}
