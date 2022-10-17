%{
package pattern
%}

%union{
    num     float64
    str     string

    pattern Pattern
    arrdef  ArrayDefinition
    arrdefl []ArrayDefinition
    objdef  ObjectDefinition
    objdefl []ObjectDefinition
    opid OptionalIdentifier
    opidl []OptionalIdentifier    
    ass Assignment    
    ref Reference
    ind Index
    key Key
    bnd Binding
}

%token NULL TRUE FALSE
%token <num> NUMBER
%token <str> STRING IDENTIFIER

%type <pattern> array object pattern
%type <ass> assignment literal
%type <arrdef> array_definition
%type <arrdefl> array_definition_list
%type <objdef> object_definition
%type <objdefl> object_definition_list
%type <ind> index
%type <key> key
%type <bnd> binding
%type <ref> reference
%type <opid> optional_identifier
%type <opidl> optional_identifier_list


%start pattern

%%

pattern
    : array     { yylex.(*lex).out = $1 }
    | object    { yylex.(*lex).out = $1 }

array
    : '[' ']'                       { $$ = Array{} }
    | '[' array_definition_list ']' { $$ = Array{$2} }

array_definition_list
    : array_definition                              { $$ = []ArrayDefinition{$1} }
    | array_definition_list ',' array_definition    { $$ = append($1, $3) }

array_definition
    : index ':' assignment      { $$ = ArrayDefinition{Index: $1, Optional: false, Assignment: $3} }
    | index '?' ':' assignment  { $$ = ArrayDefinition{Index: $1, Optional: true, Assignment: $4} }

object
    : '{' '}'                           { $$ = Object{} }
    | '{' object_definition_list '}'    { $$ = Object{$2} }

object_definition_list
    : object_definition                             { $$ = []ObjectDefinition{$1} }
    | object_definition_list ',' object_definition  { $$ = append($1, $3) }

object_definition
    : key ':' assignment       { $$ = ObjectDefinition{Key: $1, Optional: false, Assignment: $3} }
    | key '?' ':' assignment   { $$ = ObjectDefinition{Key: $1, Optional: true, Assignment: $4} }

index  
    : NUMBER    { $$ = Number($1) }
    /* | reference { $$ = $1 } */

key
    : STRING    { $$ = String($1) }
    /* | reference { $$ = $1 } */

assignment
    : binding           { $$ = $1 }
    | literal           { $$ = $1 }
    | binding literal   { $$ = BoundLiteral{$1, Assignment($2)} }
    
literal
    : NULL      { $$ = Null{} }
    | TRUE      { $$ = Boolean(true) }
    | FALSE     { $$ = Boolean(false) }
    | STRING    { $$ = String($1) }
    | NUMBER    { $$ = Number($1) }
    | reference { $$ = $1 }
    | pattern   { $$ = $1 }

binding
    : '<' '=' IDENTIFIER '>'  { $$ = Binding($3) }

reference
    : '<' optional_identifier_list '>'  { $$ = Reference($2) }

optional_identifier_list
    : optional_identifier                               { $$ = []OptionalIdentifier{$1} }
    | optional_identifier_list '.' optional_identifier  { $$ = append($1, $3) }


optional_identifier
    : IDENTIFIER        { $$ = OptionalIdentifier{Identifier: Identifier($1), Optional: false} }
    | IDENTIFIER '?'    { $$ = OptionalIdentifier{Identifier: Identifier($1), Optional: true} }