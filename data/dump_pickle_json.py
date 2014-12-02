#!/usr/local/bin/python

import json
import pickle

langs = ('czech', 'danish', 'dutch', 'english', 'estonian', 'finnish', 'french', 'german', 'greek', 'italian', 'norwegian', 'polish', 'portuguese', 'slovene', 'spanish', 'swedish', 'turkish')

for l in langs:
  print l
  src_file = "/Users/197200/nltk_data/tokenizers/punkt/" + l + ".pickle"
  dest_file = "/Users/197200/code/gocode/src/github.com/harrisj/punkt/data/" + l + ".json"
  p = pickle.load(open(src_file,"rb"))

  data = {"sentence_starters": list(p._params.sent_starters), "collocations": list(p._params.collocations), "abbrev_types": list(p._params.abbrev_types), "ortho_context": p._params.ortho_context}

  with open(dest_file, 'w') as fp:
    json.dump(data, fp)
