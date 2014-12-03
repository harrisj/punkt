# Punkt sentence tokenizer

This code is a Go port of the [ruby 1.9.x port](https://github.com/lfcipriani/punkt-segmenter) of the Punkt sentence tokenizer algorithm implemented by the NLTK Project ([http://www.nltk.org/]). As the Ruby port describes it:

> Punkt is a **language-independent**, unsupervised approach to **sentence boundary detection**. It is based on the assumption that a large
> number of ambiguities in the determination of sentence boundaries can be eliminated once abbreviations have been identiﬁed.

The full description of the algorithm is presented in the following academic paper:

> Kiss, Tibor and Strunk, Jan (2006): Unsupervised Multilingual Sentence Boundary Detection.  
> Computational Linguistics 32: 485-525.  
> [Download paper]

Original Python implementation in NLTK:

- Willy (willy@csse.unimelb.edu.au) (original Python port)
- Steven Bird (sb@csse.unimelb.edu.au) (additions)
- Edward Loper (edloper@gradient.cis.upenn.edu) (rewrite)
- Joel Nothman (jnothman@student.usyd.edu.au) (almost rewrite)

Ruby Port by:
- [Luis Cipriani](https://github.com/lfcipriani)

# Caveats

This is my first project in Go to learn how to better work in the language. That said, it is very likely that some code will be nonidiomatic and inefficient, although I am trying to use tests to at least ensure it is correct. Also, general structure of the implementation will follow the OOP structure of the Python and Ruby originals, even though Go is not an OOP language.

You may have noticed that I am porting a port instead of the original. This is because my understanding of Ruby is much more solid than my understanding of Python and I appreciated the decent test coverage provided in the Ruby port. That said, if there are errors introduced in the Ruby port, my code will probably have them too (but to be honest, it's more likely the bugs will come from me porting the Ruby code)

# Basic Use

Use is as easy as this:

```
str := "When Mr. Gregor Samsa woke up one morning from unsettling dreams, he found himself changed in his bed into a monstrous vermin. He was lying on his back as hard as armor plate, and when he lifted his head a little, he saw his vaulted brown belly, sectioned by arch-shaped ribs, to whose dome the cover, about to slide off completely, could barely cling. His many legs, pitifully thin compared with the size of the rest of him, were waving helplessly before his eyes."

t := new(Tokenizer)
t.SetLanguage("english")

fmt.Printf("%#v\n", t.SentencesFromText(str))
```

Since this is a port directly from the NLTK, I have added the option to load precompiled settings for various languages extracted from the pickle files provided with the NLTK. Here is [the full list of languages supported](https://github.com/harrisj/punkt/tree/master/data). This is currently being loaded via HTTP from Github, since I don't know how to load and package JSON within Go, but you can also run it to load any files locally instead.

# Training

You can also train it with your own corpus. Note that I am still porting this code, so it might not work entirely, but it's a start.
