package service

import (
	"context"
	"errors"

	"cinefinder/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Interface (IMPORTANTE para testes)
type LoanServiceInterface interface {
	Create(loan model.Loan) (*model.Loan, error)
	List() ([]model.Loan, error)
	GetByID(id int) (*model.Loan, error)
}

// Implementação real
type LoanService struct {
	db *pgxpool.Pool
}

func NewLoanService(db *pgxpool.Pool) *LoanService {
	return &LoanService{db: db}
}

func (s *LoanService) Create(loan model.Loan) (*model.Loan, error) {
	var unreturnedCount int
	
	checkQuery := "SELECT COUNT(*) FROM loans WHERE user_id = $1 AND returned = false"
	
	err := s.db.QueryRow(context.Background(), checkQuery, loan.User.ID).Scan(&unreturnedCount)
	if err != nil {
		return nil, err
	}
	
	if unreturnedCount > 0 {
		return nil, errors.New("Usuário possui empréstimo em aberto")
	}

	query := `
	INSERT INTO loans (user_id, movie_id, loan_date, return_date, price, returned)
	VALUES ($1, $2, NOW(), $3::timestamptz, $4, $5)
	RETURNING id, loan_date, price, returned;
	`

	err = s.db.QueryRow(context.Background(), query,
		loan.User.ID,
		loan.Movie.ID,
		loan.ReturnDate,
		loan.Price,
		loan.Returned,
	).Scan(&loan.ID, &loan.LoanDate, &loan.Price, &loan.Returned)

	if err != nil {
		return nil, err
	}

	return &loan, nil
}
func (s *LoanService) List() ([]model.Loan, error) {
	rows, err := s.db.Query(context.Background(),
		"SELECT id, user_id, movie_id, loan_date, return_date, price FROM loans",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	loans := []model.Loan{}

	for rows.Next() {
		var l model.Loan
		err := rows.Scan(&l.ID, &l.User.ID, &l.Movie.ID, &l.LoanDate, &l.ReturnDate, &l.Price)
		if err != nil {
			return nil, err
		}
		loans = append(loans, l)
	}

	return loans, nil
}

func (s *LoanService) GetByID(id int) (*model.Loan, error) {
	query := `
	SELECT id, user_id, movie_id, loan_date, return_date, price
	FROM loans
	WHERE id = $1;
	`

	var l model.Loan

	err := s.db.QueryRow(context.Background(), query, id).
		Scan(&l.ID, &l.User.ID, &l.Movie.ID, &l.LoanDate, &l.ReturnDate, &l.Price)

	if err != nil {
		return nil, err
	}

	return &l, nil
}
