package repository

import (
	"fmt"

	"invoice-backend/services/shared/pkg/database"
	"invoice-backend/services/shared/pkg/types"
)

type CustomerRepository struct {
	db *database.Client
}

func NewCustomerRepository(db *database.Client) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// GetAll returns all customers
func (r *CustomerRepository) GetAll() ([]types.Customer, error) {
	var customers []types.Customer
	_, err := r.db.Supabase.From("customers").
		Select("*", "", false).
		ExecuteTo(&customers)
	return customers, err
}

// GetByID returns a customer by ID
func (r *CustomerRepository) GetByID(id string) (*types.Customer, error) {
	var customers []types.Customer
	_, err := r.db.Supabase.From("customers").
		Select("*", "", false).
		Eq("id", id).
		ExecuteTo(&customers)
	
	if err != nil {
		return nil, err
	}
	if len(customers) == 0 {
		return nil, fmt.Errorf("customer not found")
	}
	return &customers[0], nil
}

// Create creates a new customer
func (r *CustomerRepository) Create(customer types.Customer) (*types.Customer, error) {
	var result []types.Customer
	_, err := r.db.Supabase.From("customers").
		Insert(customer, false, "", "", "").
		ExecuteTo(&result)
	
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no customer created")
	}
	return &result[0], nil
}

// Update updates a customer
func (r *CustomerRepository) Update(id string, customer types.Customer) (*types.Customer, error) {
	var result []types.Customer
	_, err := r.db.Supabase.From("customers").
		Update(customer, "", "").
		Eq("id", id).
		ExecuteTo(&result)
	
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("customer not found")
	}
	return &result[0], nil
}

// Delete deletes a customer
func (r *CustomerRepository) Delete(id string) error {
	_, _, err := r.db.Supabase.From("customers").
		Delete("", "").
		Eq("id", id).
		Execute()
	return err
}
