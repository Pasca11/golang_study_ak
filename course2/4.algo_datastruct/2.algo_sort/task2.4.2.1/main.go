package main

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"sort"
	"time"
)

type Product struct {
	Name      string
	Price     float64
	CreatedAt time.Time
	Count     int
}

func (p Product) String() string {
	return fmt.Sprintf("Name: %s, Price: %f, Count: %v", p.Name, p.Price, p.Count)
}
func generateProducts(n int) []Product {
	gofakeit.Seed(time.Now().UnixNano())
	products := make([]Product, n)
	for i := range products {
		products[i] = Product{
			Name:      gofakeit.Word(),
			Price:     gofakeit.Price(1.0, 100.0),
			CreatedAt: gofakeit.Date(),
			Count:     gofakeit.Number(1, 100),
		}
	}
	return products
}

type ByPrice []Product

func (a ByPrice) Len() int           { return len(a) }
func (a ByPrice) Less(i, j int) bool { return a[i].Price < a[j].Price }
func (a ByPrice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByCreatedAt []Product

func (a ByCreatedAt) Len() int           { return len(a) }
func (a ByCreatedAt) Less(i, j int) bool { return a[i].CreatedAt.Before(a[j].CreatedAt) }
func (a ByCreatedAt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByCount []Product

func (a ByCount) Len() int           { return len(a) }
func (a ByCount) Less(i, j int) bool { return a[i].Count < a[j].Count }
func (a ByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func main() {
	products := generateProducts(10)
	fmt.Println("Исходный список:")
	fmt.Println(products)
	// Сортировка продуктов по цене
	sort.Sort(ByPrice(products))
	fmt.Println("\nОтсортировано по цене:")
	fmt.Println(products)
	// Сортировка продуктов по дате создания
	sort.Sort(ByCreatedAt(products))
	fmt.Println("\nОтсортировано по дате создания:")
	fmt.Println(products)
	// Сортировка продуктов по количеству
	sort.Sort(ByCount(products))
	fmt.Println("\nОтсортировано по количеству:")
	fmt.Println(products)
}
