# Алгоритм LRU-кэширования

[https://goreportcard.com/badge/github.com/BorisPlus/lru](![Go Report Card](https://goreportcard.com/badge/github.com/BorisPlus/lru))

Переработано [решение](https://github.com/BorisPlus/OTUS-Go-2023-03/blob/master/hw04_lru_cache/REPORT.md) из [домашнего задания](https://github.com/BorisPlus/OTUS-Go-2023-03/blob/master/hw04_lru_cache/README.md).

## Документация

<details>
<summary>см. "go doc -all ."</summary>

```text
package lru // import "github.com/BorisPlus/lru"


TYPES

type Cacher interface {
    Set(key Key, value interface{}) bool
    Get(key Key) (interface{}, bool)
    Clear()
}
    Cacher - интерфейс хранения кэша.

func NewCache(capacity int) Cacher
    NewCache - функция-конструктор кэша.

type Key string

type KeyValue struct {
    // Has unexported fields.
}
    KeyValue - в хранилище будет учтена пара.

    Пара пригодится при извлесении элемента из списка и необходимостью поиска в
    карте, в частности, при очистке абсолютно заполненного кэша.

func NewKeyValuePair(key Key, value interface{}) *KeyValue

func (kv *KeyValue) Key() Key
    Key - получить значение KeyValue.key.

func (kv KeyValue) String() string
    String - наглядное представление значения KeyValue-структуры.

func (kv *KeyValue) Value() interface{}
    Value - получить значение KeyValue.value.

type List struct {
    // Has unexported fields.
}
    List - структура двусвязного списка.

func (list *List) Back() *ListItem
    Back() - получить последний элемент двусвязного списка.

func (list *List) Front() *ListItem
    Front() - получить первый элемент двусвязного списка.

func (list *List) Len() int
    Len() - получить длину двусвязного списка.

func (list *List) MoveToFront(i *ListItem)
    MoveToFront() - переместить элемент в начало двусвязного списка.

func (list *List) PushBack(data interface{}) *ListItem
    PushBack() - добавить значение в конец двусвязного списка.

func (list *List) PushFront(data interface{}) *ListItem
    PushFront() - добавить значение в начало двусвязного списка.

func (list *List) Remove(i *ListItem)
    Remove() - удалить элемент из двусвязного списка.

func (list *List) String() string
    String - наглядное представление всего двусвязного списка.

    Например,

    - пустой список:

        (nil:0x0)
            |
            V
        (nil:0x0)

    - список из двух элементов:

            (nil:0x0)
                |
                V
        -------------------
        Item: 0xc00002e3a0 <--------┐
        -------------------         |
        Data: 2                     |
        Prev: 0x0                   |
        Next: 0xc00002e380  >>>-----|---┐ Next 0xc00002e380
        -------------------         |   | ссылается на
                |                   |   | блок 0xc00002e380
                V                   |   |
        -------------------         |   |
        Item: 0xc00002e380  <-----------┘
        -------------------         | Prev 0xc00002e3a0
        Data: 1                     | ссылается на
        Prev: 0xc00002e3a0  >>>-----┘ блок 0xc00002e3a0
        Next: 0x0
        -------------------
                |
                V
            (nil:0x0)

type ListItem struct {
    Data interface{}
    Prev *ListItem
    Next *ListItem
}
    ListItem - элемент двусвязного списка.

func (listItem *ListItem) String() string
    String - наглядное представление значения элемента двусвязного списка.

    Например,

        -------------------             -------------------
        Item: 0xc00002e400              Item: 0xc00002e400
        -------------------             -------------------
        Data: 30                или     Data: 30
        Prev: 0xc00002e3c0              Prev: 0x0
        Next: 0xc00002e440              Next: 0x0
        -------------------             -------------------

type Lister interface {
    Len() int
    Front() *ListItem
    Back() *ListItem
    PushFront(v interface{}) *ListItem
    PushBack(v interface{}) *ListItem
    Remove(i *ListItem)
    MoveToFront(i *ListItem)
}
    Lister - интерфейс двусвязного списка.

func NewList() Lister

type LruCache struct {
    // Has unexported fields.
}
    LruCache - структура кэша.

func (cache *LruCache) Clear()
    Clear - "очистка" кэша.

func (cache *LruCache) Get(key Key) (interface{}, bool)
    Get - получение элемента из кэша.

func (cache *LruCache) Set(key Key, value interface{}) bool
    Set - уставновка элемента в кэш.

```

</details>

## Результаты тестирование

### Результаты тестирование двусвязного списка

```shell
go clean -testcache && go test ./tests/ -coverpkg=./...
ok      github.com/BorisPlus/lru/tests  12.109s coverage: 96.2% of statements in ./...
...
```

Подробные тесты

```shell
go clean -testcache && go test -v ./tests/ -run TestList -coverpkg=./...
...
PPASS
        github.com/BorisPlus/lru        coverage: 67.9% of statements in ./...
ok      github.com/BorisPlus/lru/tests  0.024s  coverage: 67.9% of statements in ./...

go clean -testcache && go test -v ./tests/ -run TestListComplex -coverpkg=./...
...
PASS
        github.com/BorisPlus/lru        coverage: 64.1% of statements in ./...
ok      github.com/BorisPlus/lru/tests  0.019s  coverage: 64.1% of statements in ./...

go clean -testcache && go test -v ./tests/ -run TestCache -coverpkg=./...
...
PASS
        github.com/BorisPlus/lru        coverage: 65.4% of statements in ./...
ok      github.com/BorisPlus/lru/tests  0.024s  coverage: 65.4% of statements in ./...

go clean -testcache && go test -v ./tests/ -run TestCacheGoroutined -coverpkg=./...
...
PASS
        github.com/BorisPlus/lru        coverage: 65.4% of statements in ./...
ok      github.com/BorisPlus/lru/tests  0.024s  coverage: 65.4% of statements in ./...

```

```shell
go clean -testcache && go test ./tests/ -run TestBenchmark 

=== RUN   TestBenchmark
--------------------------------------------------------
Operation - Set - with dataset 10 values
--------------------------------------------------------
Number of run: 1000000000
Memory allocations: 10
Memory allocations (AVERAGE): 1.000000
Number of bytes allocated: 0
Number of bytes allocated (AVERAGE): 0.000000
Time taken: 28.949µs
Time taken (AVERAGE, nanosecs.): 2894.900000  


--------------------------------------------------------
Operation - Get - with dataset 10 values
--------------------------------------------------------
Number of run: 1000000000
Memory allocations: 0
Memory allocations (AVERAGE): 0.000000
Number of bytes allocated: 0
Number of bytes allocated (AVERAGE): 0.000000
Time taken: 11.938µs
Time taken (AVERAGE, nanosecs.): 1193.800000  


--------------------------------------------------------
Operation - Set - with dataset 100 values
--------------------------------------------------------
Number of run: 1000000000
Memory allocations: 100
Memory allocations (AVERAGE): 1.000000
Number of bytes allocated: 0
Number of bytes allocated (AVERAGE): 0.000000
Time taken: 230.149µs
Time taken (AVERAGE, nanosecs.): 2301.490000  


--------------------------------------------------------
Operation - Get - with dataset 100 values
--------------------------------------------------------
Number of run: 1000000000
Memory allocations: 0
Memory allocations (AVERAGE): 0.000000
Number of bytes allocated: 0
Number of bytes allocated (AVERAGE): 0.000000
Time taken: 80.545µs
Time taken (AVERAGE, nanosecs.): 805.450000  


--------------------------------------------------------
Operation - Set - with dataset 10000 values
--------------------------------------------------------
Number of run: 1000000000
Memory allocations: 10000
Memory allocations (AVERAGE): 1.000000
Number of bytes allocated: 0
Number of bytes allocated (AVERAGE): 0.000000
Time taken: 25.673289ms
Time taken (AVERAGE, nanosecs.): 2567.328900  


--------------------------------------------------------
Operation - Get - with dataset 10000 values
--------------------------------------------------------
Number of run: 1000000000
Memory allocations: 0
Memory allocations (AVERAGE): 0.000000
Number of bytes allocated: 0
Number of bytes allocated (AVERAGE): 0.000000
Time taken: 11.383995ms
Time taken (AVERAGE, nanosecs.): 1138.399500  


--- PASS: TestBenchmark (12.69s)
PASS
ok      github.com/BorisPlus/lru/tests  12.713s

```

```shell
go clean -testcache && go test ./tests/ -count=5 -benchmem -bench BenchmarkSet 
...

BenchmarkSet/10-4               1000000000               0.0000445 ns/op             90705 avg_t/get            95.40 avg_t/set
BenchmarkSet/10-4               1000000000               0.0000431 ns/op            209442 avg_t/get            96.40 avg_t/set
BenchmarkSet/10-4               1000000000               0.0000492 ns/op             89525 avg_t/get           110.1 avg_t/set
BenchmarkSet/10-4               1000000000               0.0000528 ns/op            137579 avg_t/get           240.3 avg_t/set
BenchmarkSet/10-4               1000000000               0.0000469 ns/op            180261 avg_t/get           105.1 avg_t/set
BenchmarkSet/100-4              1000000000               0.0003379 ns/op             73088 avg_t/get            92.17 avg_t/set
BenchmarkSet/100-4              1000000000               0.0004058 ns/op             97987 avg_t/get           102.2 avg_t/set
BenchmarkSet/100-4              1000000000               0.0003361 ns/op            117718 avg_t/get            90.22 avg_t/set
BenchmarkSet/100-4              1000000000               0.0003697 ns/op             71953 avg_t/get            96.02 avg_t/set
BenchmarkSet/100-4              1000000000               0.0003680 ns/op             89030 avg_t/get           101.6 avg_t/set
BenchmarkSet/10000-4            1000000000               0.04182 ns/op       86699 avg_t/get            92.53 avg_t/set
BenchmarkSet/10000-4            1000000000               0.03844 ns/op       90012 avg_t/get            81.65 avg_t/set
BenchmarkSet/10000-4            1000000000               0.04439 ns/op       95220 avg_t/get            97.03 avg_t/set
BenchmarkSet/10000-4            1000000000               0.03942 ns/op       85455 avg_t/get            92.38 avg_t/set
BenchmarkSet/10000-4            1000000000               0.04166 ns/op      100288 avg_t/get           101.2 avg_t/set

```

| BenchmarkSet         | 2          |           ns/op | Averaget of Get() | Averaget of Set() |
| :------------------- | ---------- | --------------: | ----------------: | ----------------- |
| BenchmarkSet/10-4    | 1000000000 | 0.0000445 ns/op |   90705 avg_t/get | 95.40 avg_t/set   |
| BenchmarkSet/10-4    | 1000000000 | 0.0000431 ns/op |  209442 avg_t/get | 96.40 avg_t/set   |
| BenchmarkSet/10-4    | 1000000000 | 0.0000492 ns/op |   89525 avg_t/get | 110.1 avg_t/set   |
| BenchmarkSet/10-4    | 1000000000 | 0.0000528 ns/op |  137579 avg_t/get | 240.3 avg_t/set   |
| BenchmarkSet/10-4    | 1000000000 | 0.0000469 ns/op |  180261 avg_t/get | 105.1 avg_t/set   |
| BenchmarkSet/100-4   | 1000000000 | 0.0003379 ns/op |   73088 avg_t/get | 92.17 avg_t/set   |
| BenchmarkSet/100-4   | 1000000000 | 0.0004058 ns/op |   97987 avg_t/get | 102.2 avg_t/set   |
| BenchmarkSet/100-4   | 1000000000 | 0.0003361 ns/op |  117718 avg_t/get | 90.22 avg_t/set   |
| BenchmarkSet/100-4   | 1000000000 | 0.0003697 ns/op |   71953 avg_t/get | 96.02 avg_t/set   |
| BenchmarkSet/100-4   | 1000000000 | 0.0003680 ns/op |   89030 avg_t/get | 101.6 avg_t/set   |
| BenchmarkSet/10000-4 | 1000000000 |   0.04182 ns/op |   86699 avg_t/get | 92.53 avg_t/set   |
| BenchmarkSet/10000-4 | 1000000000 |   0.03844 ns/op |   90012 avg_t/get | 81.65 avg_t/set   |
| BenchmarkSet/10000-4 | 1000000000 |   0.04439 ns/op |   95220 avg_t/get | 97.03 avg_t/set   |
| BenchmarkSet/10000-4 | 1000000000 |   0.03942 ns/op |   85455 avg_t/get | 92.38 avg_t/set   |
| BenchmarkSet/10000-4 | 1000000000 |   0.04166 ns/op |  100288 avg_t/get | 101.2 avg_t/set   |

```shell
go clean -testcache && go test -v ./tests/ -race -run TestGoroutinedCache 
=== RUN   TestGoroutinedCache
--- PASS: TestGoroutinedCache (61.75s)
PASS
ok      github.com/BorisPlus/lru/tests  61.896s

go clean -testcache && go test ./tests/ -race
ok      github.com/BorisPlus/lru/tests  43.311s

```
