package dao

type BookTable struct {
	ID       int
	Title    string
	AuthorID int
}

type AuthorTable struct {
	ID   int
	Name string
}
