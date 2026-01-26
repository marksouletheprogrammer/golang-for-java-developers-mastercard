package main

import (
	"context"
	"testing"
	"time"
)

func BenchmarkSequential(b *testing.B) {
	transactions := GenerateTestData(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ProcessTransactionsSequential(transactions)
	}
}

func BenchmarkWorkerPool1(b *testing.B) {
	transactions := GenerateTestData(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ProcessTransactionsWithWorkerPool(transactions, 1)
	}
}

func BenchmarkWorkerPool5(b *testing.B) {
	transactions := GenerateTestData(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ProcessTransactionsWithWorkerPool(transactions, 5)
	}
}

func BenchmarkWorkerPool10(b *testing.B) {
	transactions := GenerateTestData(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ProcessTransactionsWithWorkerPool(transactions, 10)
	}
}

func BenchmarkWorkerPool50(b *testing.B) {
	transactions := GenerateTestData(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ProcessTransactionsWithWorkerPool(transactions, 50)
	}
}

func BenchmarkFanOutFanIn(b *testing.B) {
	transactions := GenerateTestData(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ProcessTransactionsWithFanOutFanIn(transactions, 10)
	}
}

func BenchmarkContext(b *testing.B) {
	transactions := GenerateTestData(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		ProcessTransactionsWithContext(ctx, transactions, 10)
		cancel()
	}
}

func BenchmarkSemaphore5(b *testing.B) {
	transactions := GenerateTestData(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ProcessTransactionsWithSemaphore(transactions, 10, 5)
	}
}

func BenchmarkSemaphore20(b *testing.B) {
	transactions := GenerateTestData(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ProcessTransactionsWithSemaphore(transactions, 10, 20)
	}
}
