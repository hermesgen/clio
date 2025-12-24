---
title: 'Go Concurrency Patterns: Goroutines and Channels'
slug: go-concurrency-patterns:-goroutines-and-channels-42665ace7e98
permalink: ""
tags:
- ""
layout: Tech
draft: false
featured: false
excerpt: ""
summary: ""
description: ""
image: /static/images/tech/go-concurrency-patterns-goroutines-and-channels-42665ace7e98/go-concurrency-patterns.png
social-image: /static/images/tech/go-concurrency-patterns-goroutines-and-channels-42665ace7e98/go-concurrency-patterns.png
published-at: 2025-09-21T15:00:00Z
created-at: 2025-12-25T11:52:41.177327429+01:00
updated-at: 2025-12-25T11:52:41.177327429+01:00
robots: ""
keywords: ""
canonical-url: ""
sitemap: ""
table-of-contents: false
comments: false
share: false
locale: ""
---
# Go Concurrency Patterns: Goroutines and Channels

Go's built-in concurrency primitives, goroutines and channels, offer a powerful and elegant way to write concurrent programs. Unlike traditional thread-based concurrency, goroutines are lightweight and managed by the Go runtime, while channels provide a safe and synchronized way for goroutines to communicate. This post explores fundamental concurrency patterns in Go.

## Goroutines: Lightweight Threads

A goroutine is a lightweight thread of execution. You can launch a goroutine by simply prefixing a function call with the `go` keyword:

```go
func doSomething() {
    // ...
}

func main() {
    go doSomething()
    // ...
}
```

## Channels: Communicating Sequential Processes

Channels are the conduits through which goroutines communicate. They allow you to send and receive values with a guarantee of synchronization.

```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Printf("Worker %d started job %d\n", id, j)
        time.Sleep(time.Second) // Simulate work
        fmt.Printf("Worker %d finished job %d\n", id, j)
        results <- j * 2
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    for j := 1; j <= 5; j++ {
        jobs <- j
    }
    close(jobs)

    for a := 1; a <= 5; a++ {
        <-results
    }
}
```

This example demonstrates how multiple workers (goroutines) can process jobs sent through a channel and send results back through another.

## Common Patterns

- **Worker Pools**: Distribute tasks among a fixed number of goroutines.
- **Fan-out/Fan-in**: Distribute work to multiple goroutines and collect their results.
- **Select Statement**: Handle multiple channel operations simultaneously.

Go's concurrency model simplifies the development of highly performant and scalable applications.
