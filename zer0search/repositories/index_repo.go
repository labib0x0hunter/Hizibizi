package repositories

import (
	memorymapper "searchengine/memory_mapper"
)

type IndexRepo struct {
	dict *memorymapper.Dictionary
	post *memorymapper.Posting
}

func NewIndexRepo(dict *memorymapper.Dictionary, post *memorymapper.Posting) *IndexRepo {
	return &IndexRepo{
		dict: dict,
		post: post,
	}
}

// search word in dictionary.index
// true : get docIds from posting.index, append new docId and update dictionary.index
// false : append docId in posting.index and append postingOffset in dictionary.index
func (i *IndexRepo) Update(wordHash uint64, docId int64) error {
	found, offset, postingOffset, postingLen, err := i.dict.Search(wordHash)
	if err != nil {
		return err
	}
	if !found {
		postingOffset, err = i.post.Append(uint64(docId), true)
		if err != nil {
			return err
		}
		err = i.dict.Append(wordHash, postingOffset, uint64(1))
		if err != nil {
			return err
		}
		return nil
	}
	postingOffset, err = i.post.Update(postingOffset, postingLen, uint64(docId))
	if err != nil {
		return err
	}
	err = i.dict.Update(offset, postingOffset, postingLen+1)
	if err != nil {
		return err
	}
	return nil
}

// search word in dictionary.index
// get docIds from posting.index
func (i *IndexRepo) GetDocIds(wordHash uint64) ([]uint64, error) {
	found, _, postingOffset, postingLen, err := i.dict.Search(wordHash)
	if err != nil {
		return []uint64{}, err
	}
	if !found {
		return []uint64{}, nil
	}
	return i.post.Search(postingOffset, postingLen)
}
