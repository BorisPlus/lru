package lru_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/BorisPlus/lru"
)

func generateData(count int) []lru.KeyValue {
	keyValues := []lru.KeyValue{}
	for key, value := range rand.Perm(count) {
		lruKey := lru.Key(fmt.Sprint(key))
		keyValues = append(keyValues, *lru.NewKeyValuePair(lruKey, value))
	}
	return keyValues
}

var TestDatasets = []struct {
	data []lru.KeyValue
}{
	{data: generateData(10)},
	{data: generateData(100)},
	{data: generateData(10000)},
}

// sum - сумма.
func sum(arr []float64) float64 {
	var sum float64
	for _, value := range arr {
		sum += value
	}
	return sum
}

// average - среднее значение (времени проведения операции Set/Get).
func average(data []float64) float64 {
	return sum(data) / float64(len(data))
}

func BenchmarkSet(b *testing.B) {
	for _, testDataset := range TestDatasets {
		b.Run(fmt.Sprintf("%d", len(testDataset.data)), func(b *testing.B) {
			// чтоб не было операций вымещения capasity=valuesCount
			cache := lru.NewCache(len(testDataset.data))
			durationsSet, durationsGet := []float64{}, []float64{}
			b.ResetTimer()
			// Собираем данные для статистики в отношении метода Set
			for _, keyValue := range testDataset.data {
				start := time.Now()
				b.StartTimer()
				cache.Set(keyValue.Key(), keyValue.Value())
				b.StopTimer()
				duration := time.Since(start)
				durationsSet = append(durationsSet, float64(duration.Microseconds()))
			}
			// Собираем данные для статистики в отношении метода Get
			for _, keyValue := range testDataset.data {
				start := time.Now()
				b.StartTimer()
				cache.Get(keyValue.Key())
				b.StopTimer()
				duration := time.Since(start)
				// durationsGet = append(durationsGet, float64(duration.Microseconds()))
				durationsGet = append(durationsGet, float64(duration.Nanoseconds()))
			}
			// Среднее времени добавления в LRU
			b.ReportMetric(average(durationsSet), "avg_t/set")
			// Среднее времени взятия из LRU
			b.ReportMetric(average(durationsGet), "avg_t/get")
		})
	}
}
