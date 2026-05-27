package tag

type Output struct {
	ID   uint
	Name string
}
type CreateTagInput struct {
	Name string
}

type UpdateTagInput struct {
	TagID   uint
	NewName string
}
