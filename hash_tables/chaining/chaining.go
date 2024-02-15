package chaining

import (
	"fmt"
)

// djb2 hash function. See http://www.cse.yorku.ca/~oz/hash.html.
func hash(value string) int {
	hash := 5381
	for _, ch := range value {
		hash = ((hash << 5) + hash) + int(ch)
	}

	// Make sure the result is non-negative.
	if hash < 0 {
		hash = -hash
	}
	return hash
}

type Employee struct {
	name  string
	phone string
}

type Bucket struct {
	items []*Employee
}

type ChainingHashTable struct {
	numBuckets int
	buckets    []*Bucket
}

// Initialize a ChainingHashTable and return a pointer to it.
func NewChainingHashTable(numBuckets int) *ChainingHashTable {
	hashTable := &ChainingHashTable{numBuckets: numBuckets, buckets: make([]*Bucket, numBuckets)}
	for i := 0; i < numBuckets; i++ {
		hashTable.buckets[i] = &Bucket{items: make([]*Employee, 0)}
	}
	return hashTable
}

// Display the hash table's contents.
func (hashTable *ChainingHashTable) dump() {
	for i, bucket := range hashTable.buckets {
		fmt.Printf("Bucket %d:\n", i)
		for _, employee := range bucket.items {
			fmt.Printf("    %s: %s\n", employee.name, employee.phone)
		}
	}
}

// Find the bucket and Employee holding this key.
// Return the bucket number and Employee number in the bucket.
// If the key is not present, return the bucket number and -1.
func (hashTable *ChainingHashTable) find(name string) (int, int) {
	key := hash(name)
	bucket := key % hashTable.numBuckets
	if hashTable.buckets != nil && hashTable.buckets[bucket] != nil {
		for i := range hashTable.buckets[bucket].items {
			if hashTable.buckets[bucket].items[i].name == name {
				return bucket, i
			}
		}
	}
	return bucket, -1
}

// Add an item to the hash table.
func (hashTable *ChainingHashTable) set(name string, phone string) {
	bucket, employee := hashTable.find(name)
	if employee > -1 {
		hashTable.buckets[bucket].items[employee].phone = phone
		return
	}
	hashTable.buckets[bucket].items = append(hashTable.buckets[bucket].items, &Employee{name: name, phone: phone})
}

// Return an item from the hash table.
func (hashTable *ChainingHashTable) get(name string) string {
	bucket, employee := hashTable.find(name)
	if employee > -1 && hashTable.buckets[bucket].items[employee] != nil {
		return hashTable.buckets[bucket].items[employee].phone
	}
	return "NOT FOUND"
}

// Return true if the person is in the hash table.
func (hashTable *ChainingHashTable) contains(name string) bool {
	bucket, employee := hashTable.find(name)
	if employee > -1 && hashTable.buckets[bucket].items[employee] != nil && hashTable.buckets[bucket].items[employee].name == name {
		return true
	}
	return false
}

// Delete this key's entry.
func (hashTable *ChainingHashTable) delete(name string) {
	bucket, employee := hashTable.find(name)
	if employee > -1 && hashTable.buckets[bucket].items[employee] != nil {
		hashTable.buckets[bucket].items = append(hashTable.buckets[bucket].items[:employee], hashTable.buckets[bucket].items[employee+1:]...)
	}
}

func ChainingRun() {
	fmt.Println("Running ChainingRun()")

	// Make some names.
	employees := []Employee{
		{"Ann Archer", "202-555-0101"},
		{"Bob Baker", "202-555-0102"},
		{"Cindy Cant", "202-555-0103"},
		{"Dan Deever", "202-555-0104"},
		{"Edwina Eager", "202-555-0105"},
		{"Fred Franklin", "202-555-0106"},
		{"Gina Gable", "202-555-0107"},
		{"Herb Henshaw", "202-555-0108"},
		{"Ida Iverson", "202-555-0109"},
		{"Jeb Jacobs", "202-555-0110"},
	}

	hashTable := NewChainingHashTable(10)
	for _, employee := range employees {
		hashTable.set(employee.name, employee.phone)
	}
	hashTable.dump()

	fmt.Printf("Table contains Sally Owens: %t\n", hashTable.contains("Sally Owens"))
	fmt.Printf("Table contains Dan Deever: %t\n", hashTable.contains("Dan Deever"))
	fmt.Println("Deleting Dan Deever")
	hashTable.delete("Dan Deever")
	fmt.Printf("Table contains Dan Deever: %t\n", hashTable.contains("Dan Deever"))
	fmt.Printf("Sally Owens: %s\n", hashTable.get("Sally Owens"))
	fmt.Printf("Fred Franklin: %s\n", hashTable.get("Fred Franklin"))
	fmt.Println("Changing Fred Franklin")
	hashTable.set("Fred Franklin", "202-555-0100")
	fmt.Printf("Fred Franklin: %s\n", hashTable.get("Fred Franklin"))
}
