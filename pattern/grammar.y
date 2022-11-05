%{
package pattern
%}

%union{
    num     float64
    str     string

    pattern ValidatedPattern
    obj Object
    arr Array
    arrdef  Element
    arrdefl []Element
    objdef  Field
    objdefl []Field
    opid OptionalIdentifier
    opidl []OptionalIdentifier    
    val Value    
    ref Reference
    ind Index
    key Key
    bnd Binding
}

%token NULL TRUE FALSE
%token <num> NUMBER
%token <str> STRING IDENTIFIER

%type <pattern> pattern
%type <obj> object
%type <arr> array
%type <val> value binding_or_value 
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
    : array_definition                              { $$ = []Element{$1} }
    | array_definition_list ',' array_definition    { $$ = append($1, $3) }

array_definition
    : index ':' binding_or_value      { $$ = Element{Index: $1, Optional: false, Value: $3} }
    | index '?' ':' binding_or_value  { $$ = Element{Index: $1, Optional: true, Value: $4} }

object
    : '{' '}'                           { $$ = Object{} }
    | '{' object_definition_list '}'    { $$ = Object{$2} }

object_definition_list
    : object_definition                             { $$ = []Field{$1} }
    | object_definition_list ',' object_definition  { $$ = append($1, $3) }

object_definition
    : key ':' binding_or_value       { $$ = Field{Key: $1, Optional: false, Value: $3} }
    | key '?' ':' binding_or_value   { $$ = Field{Key: $1, Optional: true, Value: $4} }

index  
    : NUMBER    { $$ = Number($1) }
    /* | reference { $$ = $1 } */

key
    : STRING    { $$ = String($1) }
    /* | reference { $$ = $1 } */

binding_or_value
    : binding         { $$ = $1 }
    | value           { $$ = $1 }
    | binding value   { $$ = BoundLiteral{Name: $1, Value: $2} }
    
value
    : NULL      { $$ = Null{} }
    | TRUE      { $$ = Boolean(true) }
    | FALSE     { $$ = Boolean(false) }
    | STRING    { $$ = String($1) }
    | NUMBER    { $$ = Number($1) }
    | array     { $$ = $1 }
    | object    { $$ = $1 }
    | reference { $$ = $1 }

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