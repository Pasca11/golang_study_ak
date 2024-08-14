package main

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncMain(t *testing.T) {
	tempFile, err := os.CreateTemp("", "valid_commits_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	jsonContent := `[
			{"Date": "2023-01-01", "Message": "Initial commit"},
			{"Date": "2023-01-02", "Message": "Second commit"}
		]`

	if _, err := tempFile.Write([]byte(jsonContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	tempFile.Close()

	dll := &DoubleLinkedList{}
	err = dll.LoadData(tempFile.Name())
	//fmt.Println(dll.head.next.data.Message)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if dll.head == nil || dll.head.data.Date != "2023-01-01" {
		t.Errorf("Expected first commit date to be '2023-01-01', got %v", dll.head.data.Date)
	}

	if dll.Len() != 2 {
		t.Errorf("Expected length 2, got %v", dll.Len())
	}

	dll.SetCurrent(0)

	if dll.Current().data.Message != "Initial commit" {
		t.Errorf("Expected Initial commit, got %v", dll.Current().data.Message)
	}

	dll.Next()

	if dll.Current().data.Message != "Second commit" {
		t.Errorf("Expected Second commit, got %v", dll.Current().data.Message)
	}

	dll.Prev()

	if dll.Current().data.Message != "Initial commit" {
		t.Errorf("Expected Initial commit, got %v", dll.Current().data.Message)
	}

	newCommit := Commit{
		Message: "Hello world",
		Date:    "2023-01-03",
		UUID:    "3",
	}

	dll.Insert(2, newCommit)

	if dll.Len() != 3 {
		t.Errorf("Expected 3, got %v", dll.Len())
	}

	dll.SetCurrent(2)
	dll.DeleteCurrent()

	if dll.Len() != 2 {
		t.Errorf("Expected 2, got %v", dll.Len())
	}

	if dll.Current().data.Message != "Second commit" {
		t.Errorf("Expected Second commit, got %v", dll.Current().data.Message)
	}

	dll.Push(newCommit)

	if dll.Len() != 3 {
		t.Errorf("Expected 3, got %v", dll.Len())
	}

	dll.Delete(2)

	if dll.Len() != 2 {
		t.Errorf("Expected 2, got %v", dll.Len())
	}

	index, err := dll.Index()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if index != 1 {
		t.Errorf("Expected index 1, got %v", index)
	}

	secondNode, err := dll.GetByIndex(1)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if secondNode.data.Message != "Second commit" {
		t.Errorf("Expected Second commit, got %v", dll.Current().data.Message)
	}

	dll.Prev()
	dll.Push(newCommit)
	dll.Pop()

	if dll.Len() != 2 {
		t.Errorf("Expected 2, got %v", dll.Len())
	}

	dll.Shift()

	if dll.Len() != 1 {
		t.Errorf("Expected 1, got %v", dll.Len())
	}

	currentNode := dll.Current()

	if currentNode.data.Message != "Second commit" {
		t.Errorf("Expected Second commit, got %v", dll.Current().data.Message)
	}

	dll.Push(newCommit)
	uuidNode := dll.SearchUUID("3")

	if uuidNode.data.Message != "Hello world" {
		t.Errorf("Expected message Hello World, got %v", uuidNode.data.Message)
	}

	messageSearchNode := dll.Search("Hello world")

	if messageSearchNode.data.UUID != "3" {
		t.Errorf("Expected UUID 3, got %v", messageSearchNode.data.UUID)
	}

	dll.Reverse()

	if dll.head.data.UUID != "3" {
		t.Errorf("Expected UUID3, got %v", dll.head.data.UUID)
	}

	dll.Pop()
	dll.Pop()

	_, err = dll.GetByIndex(0)

	expectedError := fmt.Errorf("link is empty")
	if errors.Is(err, expectedError) {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	_, err = dll.GetByIndex(-4)
	expectedError = fmt.Errorf("index out of range")
	if errors.Is(err, expectedError) {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	err = dll.Insert(-4, newCommit)
	if errors.Is(err, expectedError) {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	err = dll.Insert(0, newCommit)
	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	dll.Pop()
	err = dll.Push(newCommit)
	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

}

func TestInit(t *testing.T) {
	commits := GenerateData(5)
	dll := &DoubleLinkedList{}
	dll.Init(commits)

	assert.Equal(t, 5, dll.Len())
	assert.NotNil(t, dll.head)
	assert.NotNil(t, dll.tail)
	assert.Equal(t, commits[0], *dll.head.data)
	assert.Equal(t, commits[4], *dll.tail.data)
}

func TestQuickSort(t *testing.T) {
	commits := []Commit{
		{Message: "Commit 3", UUID: "uuid3", Date: "2022-01-01"},
		{Message: "Commit 1", UUID: "uuid1", Date: "2020-01-01"},
		{Message: "Commit 2", UUID: "uuid2", Date: "2021-01-01"},
	}

	QuickSort(commits)

	expectedDates := []string{"2020-01-01", "2021-01-01", "2022-01-01"}
	for i, commit := range commits {
		if commit.Date != expectedDates[i] {
			t.Errorf("Expected date %s but got %s", expectedDates[i], commit.Date)
		}
	}
}
