package domain

type Student struct {
	ID    string `json:"id"`
	Nome  string `json:"nome"`
	Plano int    `json:"plano"`
	Email string `json:"email"`
	Curso string `json:"curso"`
}

type GetStudentsResponse struct {
	Curso  string    `json:"id"`
	Alunos []Student `json:"alunos"`
}
