//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

// +build example1

package main

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzers/custom"
	"github.com/blevesearch/bleve/analysis/analyzers/keyword"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/analysis/tokenizers/unicode"
	"github.com/blevesearch/bleve/analysis/tokens/lowercase"
	"github.com/blevesearch/bleve/analysis/tokens/porter"
	"github.com/blevesearch/bleve/analysis/tokens/truncate"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/blevex/detectlang"
)

const textFieldAnalyzer = "en"

func buildIndexMapping() (mapping.IndexMapping, error) {

	// a custom field definition that uses our custom analyzer
	notTooLongFieldMapping := bleve.NewTextFieldMapping()
	notTooLongFieldMapping.Analyzer = "enNotTooLong"

	// a generic reusable mapping for english text
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	// a generic reusable mapping for keyword text
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name

	// a specific mapping to index the description fields
	// detected language
	descriptionLangFieldMapping := bleve.NewTextFieldMapping()
	descriptionLangFieldMapping.Name = "descriptionLang"
	descriptionLangFieldMapping.Analyzer = detectlang.AnalyzerName
	descriptionLangFieldMapping.Store = false
	descriptionLangFieldMapping.IncludeTermVectors = false
	descriptionLangFieldMapping.IncludeInAll = false

	beerMapping := bleve.NewDocumentMapping()

	// name
	beerMapping.AddFieldMappingsAt("name", englishTextFieldMapping)

	// description
	beerMapping.AddFieldMappingsAt("description",
		notTooLongFieldMapping,
		descriptionLangFieldMapping)

	beerMapping.AddFieldMappingsAt("type", keywordFieldMapping)
	beerMapping.AddFieldMappingsAt("style", keywordFieldMapping)
	beerMapping.AddFieldMappingsAt("category", keywordFieldMapping)

	breweryMapping := bleve.NewDocumentMapping()
	breweryMapping.AddFieldMappingsAt("name", englishTextFieldMapping)
	breweryMapping.AddFieldMappingsAt("description", englishTextFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("beer", beerMapping)
	indexMapping.AddDocumentMapping("brewery", breweryMapping)

	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = textFieldAnalyzer

	err := indexMapping.AddCustomTokenFilter("notTooLong",
		map[string]interface{}{
			"type":   truncate.Name,
			"length": 5.0,
		})
	if err != nil {
		return nil, err
	}

	err = indexMapping.AddCustomAnalyzer("enNotTooLong",
		map[string]interface{}{
			"type":      custom.Name,
			"tokenizer": unicode.Name,
			"token_filters": []string{
				"notTooLong",
				en.PossessiveName,
				lowercase.Name,
				en.StopName,
				porter.Name,
			},
		})
	if err != nil {
		return nil, err
	}

	return indexMapping, nil
}
