astgen
======

A simple parser description language, and a library that uses it.

This language is intended as a simple declarative alternative to parser generators like Flex+Bison, etc. It should mostly be useful when prototyping a new language. The description defines the data types that the AST representation will consist of, and grammar productions corresponding to those types. Software built on top of the library can then generate an AST definition and a parser for the language in arbitrary implementation language. Currently, a generator for a Flex+Bison parser is included that generates a simpler AST description, and generators for libraries that can parse the simple description into AST tree and then process the tree using the visitor paradigm.


Language
--------

There are three data type kinds.


One-of-selection type:

	type TypeName = PossibleType1 | PossibleType2 | PossibleType3 | ... | PossibleTypeN
	
This type can hold any of a selection of other types.

Example: A type Expression that can be either PlusExpression, MinusExpression or MultiplyExpression.

	type Expression = PlusExpression | MinusExpression



Structured type:

This in a type that encodes a non-trivial AST node.

	// [production1]
	// [production2]
	type TypeName = struct {
		member1: TypeOfMember1
		
		// annotation
		member2: TypeOfMember2
	}


... TODO ...
