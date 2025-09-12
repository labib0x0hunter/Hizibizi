package services

import (
	"errors"
	"log/slog"
	"searchengine/repositories"
	"searchengine/tokenizer"
	"searchengine/utils"
)

// business logic, user repo

type EngineService struct {
	indexRepo *repositories.IndexRepo
	docRepo   *repositories.DocumentRepo
	hasher    *utils.Hash
	docId     int64
}

func NewEngineService(indexRepo *repositories.IndexRepo, docRepo *repositories.DocumentRepo, hasher *utils.Hash) *EngineService {
	return &EngineService{
		indexRepo: indexRepo,
		docRepo:   docRepo,
		hasher:    hasher,
		docId:     1,
	}
}

/***
1. Assign docId
2. Tokenize the document
3. for each word ::
		- search in dict.index
		- if presernt
			- append docId storing at offset in post.index
		- else
			- append docId in post.index (at the last), store offset in dict.index
4. Insert document to mysql database, (Todo -> store document in docs.dat)
**/

func (e *EngineService) IndexDocument(document string) error {
	tokens := tokenizer.GetTokens(document)
	var insertedFlag bool = false
	for _, tok := range tokens.Tokens {
		tokenHash := e.getHash(tok)
		if err := e.indexRepo.Update(tokenHash, e.docId); err != nil {
			slog.Error("[engine_service.go]		[IndexDocument()]	", err)
			continue
		}
		insertedFlag = true
	}
	if !insertedFlag {
		slog.Error("[engine_service.go]		[IndexDocument()]	document not inserted")
		return errors.New("document not inserted")
	}

	id, err := e.docRepo.Insert(document)
	if err != nil {
		slog.Error("[engine_service.go]		[IndexDocument()]	", err)
		// (To-Do) Document is not stored in database, rollback to previous state
		return err
	}
	e.docId = id + 1
	return nil
}

/**
1. Tokenize the document
2. For each word ::
	- search doct.index and get offset
	- read post.index and get docId slice
	- intersect docIds
3. Retrive documents from mysql database

To-do:
	- ranking mechanism
**/

func (e *EngineService) SearchDocument(document string) []string {
	tokens := tokenizer.GetTokens(document)
	foundDocIds := make(map[uint64]struct{})
	for _, tok := range tokens.Tokens {
		tokenHash := e.getHash(tok)
		tempSlice, err := e.indexRepo.GetDocIds(tokenHash)
		if err != nil {
			slog.Error("[engine_service.go]		[SearchDocument()]	", err)
			continue
		}
		// AND operation (intersection) on tempSlice and foundDocIds
		if len(tempSlice) == 0 {
			continue
		}

		if len(foundDocIds) == 0 {
			for _, docId := range tempSlice {
				foundDocIds[docId] = struct{}{}
			}
			continue
		}
		commonDocId := make(map[uint64]struct{}, 0)
		for _, docId := range tempSlice {
			if _, ok := foundDocIds[docId]; ok {
				commonDocId[docId] = struct{}{}
			}
		}
		foundDocIds, commonDocId = commonDocId, foundDocIds
	}

	result := make([]string, 0, len(foundDocIds))
	for docId := range foundDocIds {
		// Search from MySql
		document, err := e.docRepo.Query(int(docId))
		if err != nil {
			slog.Error("[engine_service.go]		[SearchDocument()]	", err)
			continue
		}
		result = append(result, document)
	}
	return result
}

func (e *EngineService) getHash(word string) uint64 {
	e.hasher.WriteString(word)
	wordHash := e.hasher.Sum()
	e.hasher.Reset()
	return wordHash
}
