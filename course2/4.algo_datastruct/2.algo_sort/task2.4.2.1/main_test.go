package main

import (
	"sort"
	"testing"
	"time"
)

func TestByPrice(t *testing.T) {
	products := []Product{
		{
			Name:      "product 1",
			Price:     20,
			CreatedAt: time.Now(),
			Count:     20,
		},
		{
			Name:      "product 2",
			Price:     10,
			CreatedAt: time.Now(),
			Count:     10,
		},
		{
			Name:      "product 3",
			Price:     30,
			CreatedAt: time.Now(),
			Count:     30,
		},
	}

	sort.Sort(ByPrice(products))

	if !sort.IsSorted(ByPrice(products)) {
		t.Errorf("Not sorted by price %v", products)
	}
}

func TestByCount(t *testing.T) {
	products := []Product{
		{
			Name:      "product 1",
			Price:     20,
			CreatedAt: time.Now(),
			Count:     20,
		},
		{
			Name:      "product 2",
			Price:     10,
			CreatedAt: time.Now(),
			Count:     10,
		},
		{
			Name:      "product 3",
			Price:     30,
			CreatedAt: time.Now(),
			Count:     30,
		},
	}

	sort.Sort(ByCount(products))

	if !sort.IsSorted(ByPrice(products)) {
		t.Errorf("Not sorted by count %v", products)
	}
}

func TestByCreatedAt(t *testing.T) {
	products := []Product{
		{
			Name:      "product 1",
			Price:     20,
			CreatedAt: time.Now().Add(time.Hour * 1),
			Count:     20,
		},
		{
			Name:      "product 2",
			Price:     10,
			CreatedAt: time.Now().Add(time.Hour * -1),
			Count:     10,
		},
		{
			Name:      "product 3",
			Price:     30,
			CreatedAt: time.Now(),
			Count:     30,
		},
	}

	sort.Sort(ByCreatedAt(products))

	if !sort.IsSorted(ByCreatedAt(products)) {
		t.Errorf("Not sorted by created_at %v", products)
	}
}
