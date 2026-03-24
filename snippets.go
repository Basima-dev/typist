package main

import "math/rand"

type Snippet struct {
	Code     string
	Language string
}

var snippets = map[string][]Snippet{
	"go": {
		{Language: "go", Code: `func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}`},
		{Language: "go", Code: `func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}`},
		{Language: "go", Code: `type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}`},
		{Language: "go", Code: `func binarySearch(nums []int, target int) int {
	lo, hi := 0, len(nums)-1
	for lo <= hi {
		mid := (lo + hi) / 2
		if nums[mid] == target {
			return mid
		} else if nums[mid] < target {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return -1
}`},
		{Language: "go", Code: `func mergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	mid := len(arr) / 2
	left := mergeSort(arr[:mid])
	right := mergeSort(arr[mid:])
	return merge(left, right)
}`},
		{Language: "go", Code: `func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ch := make(chan Result, 1)
	go func() {
		ch <- doWork(ctx)
	}()

	select {
	case result := <-ch:
		fmt.Println(result)
	case <-ctx.Done():
		fmt.Println("timed out")
	}
}`},
		{Language: "go", Code: `func wordCount(s string) map[string]int {
	counts := make(map[string]int)
	for _, word := range strings.Fields(s) {
		counts[strings.ToLower(word)]++
	}
	return counts
}`},
		{Language: "go", Code: `type Node struct {
	Val  int
	Next *Node
}

func reverseList(head *Node) *Node {
	var prev *Node
	curr := head
	for curr != nil {
		next := curr.Next
		curr.Next = prev
		prev = curr
		curr = next
	}
	return prev
}`},
	},

	"js": {
		{Language: "js", Code: `function debounce(fn, delay) {
  let timer;
  return function (...args) {
    clearTimeout(timer);
    timer = setTimeout(() => fn.apply(this, args), delay);
  };
}`},
		{Language: "js", Code: `const deepClone = (obj) => {
  if (obj === null || typeof obj !== 'object') return obj;
  if (Array.isArray(obj)) return obj.map(deepClone);
  return Object.fromEntries(
    Object.entries(obj).map(([k, v]) => [k, deepClone(v)])
  );
};`},
		{Language: "js", Code: `async function fetchWithRetry(url, retries = 3) {
  for (let i = 0; i < retries; i++) {
    try {
      const res = await fetch(url);
      if (!res.ok) throw new Error(res.statusText);
      return await res.json();
    } catch (err) {
      if (i === retries - 1) throw err;
      await new Promise(r => setTimeout(r, 1000 * 2 ** i));
    }
  }
}`},
		{Language: "js", Code: `function memoize(fn) {
  const cache = new Map();
  return function (...args) {
    const key = JSON.stringify(args);
    if (cache.has(key)) return cache.get(key);
    const result = fn.apply(this, args);
    cache.set(key, result);
    return result;
  };
}`},
		{Language: "js", Code: `const pipe = (...fns) => (x) => fns.reduce((v, f) => f(v), x);

const transform = pipe(
  (s) => s.trim(),
  (s) => s.toLowerCase(),
  (s) => s.replace(/\s+/g, '-')
);`},
		{Language: "js", Code: `class EventEmitter {
  constructor() { this.events = {}; }

  on(event, listener) {
    (this.events[event] ??= []).push(listener);
    return () => this.off(event, listener);
  }

  emit(event, ...args) {
    this.events[event]?.forEach(fn => fn(...args));
  }
}`},
	},

	"python": {
		{Language: "python", Code: `def quicksort(arr):
    if len(arr) <= 1:
        return arr
    pivot = arr[len(arr) // 2]
    left = [x for x in arr if x < pivot]
    mid = [x for x in arr if x == pivot]
    right = [x for x in arr if x > pivot]
    return quicksort(left) + mid + quicksort(right)`},
		{Language: "python", Code: `from functools import wraps
import time

def retry(times=3, delay=1.0):
    def decorator(fn):
        @wraps(fn)
        def wrapper(*args, **kwargs):
            for attempt in range(times):
                try:
                    return fn(*args, **kwargs)
                except Exception as e:
                    if attempt == times - 1:
                        raise
                    time.sleep(delay * 2 ** attempt)
        return wrapper
    return decorator`},
		{Language: "python", Code: `def flatten(nested, depth=None):
    result = []
    for item in nested:
        if isinstance(item, list) and depth != 0:
            result.extend(flatten(item, None if depth is None else depth - 1))
        else:
            result.append(item)
    return result`},
		{Language: "python", Code: `class LRUCache:
    def __init__(self, capacity):
        self.cache = {}
        self.capacity = capacity
        self.order = []

    def get(self, key):
        if key not in self.cache:
            return -1
        self.order.remove(key)
        self.order.append(key)
        return self.cache[key]

    def put(self, key, value):
        if key in self.cache:
            self.order.remove(key)
        elif len(self.cache) >= self.capacity:
            del self.cache[self.order.pop(0)]
        self.cache[key] = value
        self.order.append(key)`},
		{Language: "python", Code: `def chunk(lst, size):
    return [lst[i:i + size] for i in range(0, len(lst), size)]

def sliding_window(lst, k):
    return [lst[i:i + k] for i in range(len(lst) - k + 1)]`},
	},

	"rust": {
		{Language: "rust", Code: `fn fibonacci(n: u64) -> u64 {
    match n {
        0 => 0,
        1 => 1,
        _ => {
            let (mut a, mut b) = (0u64, 1u64);
            for _ in 2..=n {
                (a, b) = (b, a + b);
            }
            b
        }
    }
}`},
		{Language: "rust", Code: `use std::collections::HashMap;

fn word_count(text: &str) -> HashMap<&str, usize> {
    let mut counts = HashMap::new();
    for word in text.split_whitespace() {
        *counts.entry(word).or_insert(0) += 1;
    }
    counts
}`},
		{Language: "rust", Code: `fn binary_search<T: Ord>(arr: &[T], target: &T) -> Option<usize> {
    let (mut lo, mut hi) = (0, arr.len());
    while lo < hi {
        let mid = lo + (hi - lo) / 2;
        match arr[mid].cmp(target) {
            std::cmp::Ordering::Equal => return Some(mid),
            std::cmp::Ordering::Less => lo = mid + 1,
            std::cmp::Ordering::Greater => hi = mid,
        }
    }
    None
}`},
		{Language: "rust", Code: `#[derive(Debug)]
enum Tree<T> {
    Leaf(T),
    Node(Box<Tree<T>>, Box<Tree<T>>),
}

impl<T: std::fmt::Display> Tree<T> {
    fn depth(&self) -> usize {
        match self {
            Tree::Leaf(_) => 0,
            Tree::Node(l, r) => 1 + l.depth().max(r.depth()),
        }
    }
}`},
		{Language: "rust", Code: `fn two_sum(nums: &[i32], target: i32) -> Option<(usize, usize)> {
    let mut seen = std::collections::HashMap::new();
    for (i, &n) in nums.iter().enumerate() {
        if let Some(&j) = seen.get(&(target - n)) {
            return Some((j, i));
        }
        seen.insert(n, i);
    }
    None
}`},
	},
}

var langKeys = []string{"go", "js", "python", "rust"}

func randomSnippet(lang string) Snippet {
	pool := snippets[lang]
	return pool[rand.Intn(len(pool))]
}
