### UNDER DEVELOPMENT ###

# Zer0 Search
Zer0 Search is a platform that allows users to store documents and efficiently search through them using keywords.

## Feature
- Store documents
- Keyword-based document search

## To-do
- Add authentication
- Role based authentication
- Public and private documents
- Fix concurrent issues

**Learning Material :**

Projects Focused on Inverted Index Concepts
1. https://github.com/dwayhs/go-search-engine
 Though it uses in-memory maps, it clearly shows how to structure indexing, querying, ranking — a great reference to layer mmap + binary into.

2. https://github.com/el10savio/Inverted-Index-Generator
Focuses on core inverted-index building logic — tokenization, storing word→doc mapping — which you can adapt to a mmap-backed system with fixed struct entries.

**Resources**

1. https://dev.to/ravikishan/building-a-high-performance-full-text-search-engine-in-go-m4b#:~:text=This%20blog%20will%20walk%20you,performance%20applications
2. https://artem.krylysov.com/blog/2020/07/28/lets-build-a-full-text-search-engine/