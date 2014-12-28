package main

import (
	"fmt"
)
func emitPrologue() {
fmt.Print("%locations\n")
fmt.Print("%define api.pure full\n")
fmt.Print("%param {yyscan_t scanner}\n")
fmt.Print("%{\n")
fmt.Print("\t#include <stdlib.h>\n")
fmt.Print("\t#include <stdio.h>\n")
fmt.Print("\t#include <stdarg.h>\n")
fmt.Print("\t#include <ctype.h>\n")
fmt.Print("\t#include <string.h>\n")
fmt.Print("\t#include <assert.h>\n")
fmt.Print("%}\n")
fmt.Print("%code requires {\n")
fmt.Print("\t\n")
fmt.Print("\ttypedef struct sexp {\n")
fmt.Print("\t\tchar *string;\n")
fmt.Print("\t\tstruct sexp **list;\n")
fmt.Print("\t} *sexp_t;\n")
fmt.Print("\t\n")
fmt.Print("\textern sexp_t SExString(const char *str);\n")
fmt.Print("\textern sexp_t SExList(int count, ...);\n")
fmt.Print("\textern sexp_t SExPrepend(sexp_t, sexp_t);\n")
fmt.Print("\textern sexp_t SExAppend(sexp_t, sexp_t);\n")
fmt.Print("\textern void PrintSExpression(FILE *f, sexp_t ex, int indent);\n")
fmt.Print("\textern const char *UnescapeString(const char *s);\n")
fmt.Print("\textern sexp_t Unnamed();\n")
fmt.Print("\t#define YYSTYPE sexp_t\n")
fmt.Print("\t\n")
fmt.Print("\ttypedef struct YYLTYPE\n")
fmt.Print("\t{\n")
fmt.Print("\t  int first_line;\n")
fmt.Print("\t  int first_column;\n")
fmt.Print("\t  int last_line;\n")
fmt.Print("\t  int last_column;\n")
fmt.Print("\t} YYLTYPE;\n")
fmt.Print("\t#define YYLTYPE_IS_DECLARED 1\n")
fmt.Print("\t\n")
fmt.Print("\ttypedef void* yyscan_t;\n")
fmt.Print("\tvoid yyerror (YYLTYPE *locp, yyscan_t scanner, char const *msg);\n")
fmt.Print("\tint yylex (YYSTYPE *lvalp, YYLTYPE *llocp, yyscan_t scanner);\n")
fmt.Print("}\n")
fmt.Print("%error-verbose\n")
fmt.Print("%define parse.lac full\n")
for _, val := range allTokens {
fmt.Print("%token ")
fmt.Printf("%v",  val )
fmt.Print(" \n")
}
fmt.Print("%%\n")
fmt.Print("Main:\n")
fmt.Print("  File { PrintSExpression(stdout, $1, 0); }\n")
fmt.Print(";\n")
}
// Here come the grammar rules.
func emitEpilogue() {
fmt.Print("%%\n")
fmt.Print("void yyerror(YYLTYPE *locp, yyscan_t scanner, const char *s) {\n")
fmt.Print("\tfprintf(stderr, \"[%d:%d--%d:%d] %s\\n\", locp->first_line, locp->first_column, locp->last_line, locp->last_column, s);\n")
fmt.Print("}\n")
fmt.Print("extern void yylex_init(yyscan_t *);\n")
fmt.Print("extern void yylex_destroy(yyscan_t);\n")
fmt.Print("int main(void) {\n")
fmt.Print("\tyydebug = 0;\n")
fmt.Print("\tyyscan_t scanner;\n")
fmt.Print("\tyylex_init(&scanner);\n")
fmt.Print("\tint rc = yyparse(scanner);\n")
fmt.Print("\tyylex_destroy(scanner);\n")
fmt.Print("\treturn rc;\n")
fmt.Print("}\n")
fmt.Print("sexp_t SExString(const char *str) {\n")
fmt.Print("\tsexp_t ex = calloc(sizeof(struct sexp), 1);\n")
fmt.Print("\tex->string = strdup(str);\n")
fmt.Print("\tex->list = NULL;\n")
fmt.Print("\treturn ex;\n")
fmt.Print("}\n")
fmt.Print("sexp_t SExList(int count, ...) {\n")
///* Create the array. */
fmt.Print("\tsexp_t *list = calloc(sizeof(sexp_t), count + 1);\n")
fmt.Print("\t\n")
fmt.Print("\tva_list ap;\n")
fmt.Print("\tva_start(ap, count);\n")
fmt.Print("\tfor (int i = 0; i < count; i++) {\n")
fmt.Print("\t\tlist[i] = va_arg(ap, sexp_t);\n")
fmt.Print("\t}\n")
fmt.Print("\tlist[count] = NULL;\n")
fmt.Print("\tva_end(ap);\n")
fmt.Print("\t\n")
///* Create the structure. */
fmt.Print("\tsexp_t ex = calloc(sizeof(struct sexp), 1);\n")
fmt.Print("\tex->string = NULL;\n")
fmt.Print("\tex->list = list;\n")
fmt.Print("\t\n")
fmt.Print("\treturn ex;\n")
fmt.Print("}\n")
fmt.Print("sexp_t SExAppend(sexp_t ex1, sexp_t ex2) {\n")
fmt.Print("\tassert(ex1->list != NULL);\n")
fmt.Print("\t\n")
fmt.Print("\tint list_len = 0;\n")
fmt.Print("\twhile (ex1->list[list_len] != NULL) {\n")
fmt.Print("\t\tlist_len++;\n")
fmt.Print("\t}\n")
fmt.Print("\t\n")
fmt.Print("\tsexp_t *new_list = calloc(sizeof(sexp_t), list_len + 2);\n")
fmt.Print("\tfor (int i = 0; i < list_len; i++) {\n")
fmt.Print("\t\tnew_list[i] = ex1->list[i];\n")
fmt.Print("\t}\n")
fmt.Print("\tnew_list[list_len] = ex2;\n")
fmt.Print("\tnew_list[list_len+1] = NULL;\n")
fmt.Print("\t\n")
fmt.Print("\tfree(ex1->list);\n")
fmt.Print("\tex1->list = new_list;\n")
fmt.Print("\treturn ex1;\n")
fmt.Print("}\n")
fmt.Print("sexp_t SExPrepend(sexp_t ex1, sexp_t ex2) {\n")
fmt.Print("\tint list_len = 0;\n")
fmt.Print("\twhile (ex1->list[list_len] != NULL) {\n")
fmt.Print("\t\tlist_len++;\n")
fmt.Print("\t}\n")
fmt.Print("\t\n")
fmt.Print("\tsexp_t *new_list = calloc(sizeof(sexp_t), list_len + 2);\n")
fmt.Print("\tnew_list[0] = ex2;\n")
fmt.Print("\tfor (int i = 0; i < list_len; i++) {\n")
fmt.Print("\t\tnew_list[i+1] = ex1->list[i];\n")
fmt.Print("\t}\n")
fmt.Print("\tnew_list[list_len+1] = NULL;\n")
fmt.Print("\t\n")
fmt.Print("\tfree(ex1->list);\n")
fmt.Print("\tex1->list = new_list;\n")
fmt.Print("\treturn ex1;\n")
fmt.Print("}\n")
fmt.Print("void printString(FILE *f, const char *s) {\n")
fmt.Print("\tfputc('\"', f);\n")
fmt.Print("\twhile (*s != 0) {\n")
fmt.Print("\t\tif (!isascii(*s) || !isprint(*s) || *s == '\"' || *s == '\\\\' || *s == '\\n') {\n")
fmt.Print("\t\t\tfprintf(f, \"\\\\x%02x\", (unsigned)*s);\n")
fmt.Print("\t\t} else {\n")
fmt.Print("\t\t\tfputc(*s, f);\n")
fmt.Print("\t\t}\n")
fmt.Print("\t\ts++;\n")
fmt.Print("\t}\n")
fmt.Print("\tfputc('\"', f);\n")
fmt.Print("}\n")
fmt.Print("void printIndent(FILE *f, int indent) {\n")
fmt.Print("\tfor (; indent > 0; indent--) {\n")
fmt.Print("\t\tfprintf(f, \"\\t\");\n")
fmt.Print("\t}\n")
fmt.Print("}\n")
fmt.Print("void PrintSExpression(FILE *f, sexp_t ex, int indent) {\n")
fmt.Print("\tif (ex->string != NULL) {\n")
fmt.Print("\t\tprintIndent(f, indent);\n")
fmt.Print("\t\tprintString(f, ex->string);\n")
fmt.Print("\t\tfprintf(f, \"\\n\");\n")
fmt.Print("\t\treturn;\n")
fmt.Print("\t}\n")
fmt.Print("\t\n")
fmt.Print("\tif (ex->list[0] == NULL) {\n")
fmt.Print("\t\tprintIndent(f, indent);\n")
fmt.Print("\t\tfprintf(f, \"()\\n\");\n")
fmt.Print("\t\treturn;\n")
fmt.Print("\t}\n")
fmt.Print("\t\n")
fmt.Print("\tprintIndent(f, indent);\n")
fmt.Print("\tfprintf(f, \"(\\n\");\n")
fmt.Print("\t\n")
fmt.Print("\tfor (int i = 0; ex->list[i] != NULL; i++) {\n")
fmt.Print("\t\tPrintSExpression(f, ex->list[i], indent + 1);\n")
fmt.Print("\t}\n")
fmt.Print("\t\n")
fmt.Print("\tprintIndent(f, indent);\n")
fmt.Print("\tfprintf(f, \")\\n\");\n")
fmt.Print("}\n")
fmt.Print("const char *UnescapeString(const char *s) {\n")
fmt.Print("\tstatic int buffer_len = 0;\n")
fmt.Print("\tstatic char *buffer = NULL;\n")
fmt.Print("\t\n")
fmt.Print("\tint buffer_used = 0;\n")
fmt.Print("\t\n")
fmt.Print("\tint len = strlen(s);\n")
fmt.Print("\t\n")
fmt.Print("\tif (len > buffer_len || buffer == NULL) {\n")
fmt.Print("\t\tbuffer_len = len;\n")
fmt.Print("\t\tfree(buffer);\n")
fmt.Print("\t\tbuffer = malloc(buffer_len + 1);\n")
fmt.Print("\t\tif (buffer == NULL) {\n")
fmt.Print("\t\t\tfprintf(stderr, \"Out of memory.\");\n")
fmt.Print("\t\t\texit(1);\n")
fmt.Print("\t\t}\n")
fmt.Print("\t}\n")
fmt.Print("\t\n")
fmt.Print("\tint i;\n")
fmt.Print("\tfor (i = 1; i < len-1; i++) {\n")
fmt.Print("\t\tif (s[i] != '\\\\') {\n")
fmt.Print("\t\t\tbuffer[buffer_used++] = s[i];\n")
fmt.Print("\t\t\tcontinue;\n")
fmt.Print("\t\t}\n")
fmt.Print("\t\t\n")
fmt.Print("\t\ti++;\n")
fmt.Print("\t\t\n")
fmt.Print("\t\tswitch (s[i]) {\n")
fmt.Print("\t\tcase 'a':\n")
fmt.Print("\t\t\tbuffer[buffer_used++] = 0x07;\n")
fmt.Print("\t\t\tbreak;\n")
fmt.Print("\t\tcase 'b':\n")
fmt.Print("\t\t\tbuffer[buffer_used++] = 0x08;\n")
fmt.Print("\t\t\tbreak;\n")
fmt.Print("\t\tcase 'f':\n")
fmt.Print("\t\t\tbuffer[buffer_used++] = 0x0C;\n")
fmt.Print("\t\t\tbreak;\n")
fmt.Print("\t\tcase 'n':\n")
fmt.Print("\t\t\tbuffer[buffer_used++] = 0x0A;\n")
fmt.Print("\t\t\tbreak;\n")
fmt.Print("\t\tcase 'r':\n")
fmt.Print("\t\t\tbuffer[buffer_used++] = 0x0D;\n")
fmt.Print("\t\t\tbreak;\n")
fmt.Print("\t\tcase 't':\n")
fmt.Print("\t\t\tbuffer[buffer_used++] = 0x09;\n")
fmt.Print("\t\t\tbreak;\n")
fmt.Print("\t\tcase 'v':\n")
fmt.Print("\t\t\tbuffer[buffer_used++] = 0x0B;\n")
fmt.Print("\t\t\tbreak;\n")
fmt.Print("\t\tcase '\\\\':\n")
fmt.Print("\t\t\tbuffer[buffer_used++] = '\\\\';\n")
fmt.Print("\t\t\tbreak;\n")
fmt.Print("\t\tcase '\\'':\n")
fmt.Print("\t\t\tbuffer[buffer_used++] = '\\'';\n")
fmt.Print("\t\t\tbreak;\n")
fmt.Print("\t\tcase '\"':\n")
fmt.Print("\t\t\tbuffer[buffer_used++] = '\"';\n")
fmt.Print("\t\t\tbreak;\n")
fmt.Print("\t\tcase '0':\n")
fmt.Print("\t\tcase '1':\n")
fmt.Print("\t\tcase '2':\n")
fmt.Print("\t\tcase '3':\n")
//// TODO
fmt.Print("\t\t\tfprintf(stderr, \"Octal escapes unimplemented.\\n\");\n")
fmt.Print("\t\t\texit(1);\n")
fmt.Print("\t\tcase 'x':\n")
//// TODO
fmt.Print("\t\t\tfprintf(stderr, \"Hex escapes unimplemented.\\n\");\n")
fmt.Print("\t\t\texit(1);\n")
fmt.Print("\t\tcase 'u':\n")
fmt.Print("\t\tcase 'U':\n")
//// TODO
fmt.Print("\t\t\tfprintf(stderr, \"Unicode escapes unimplemented.\\n\");\n")
fmt.Print("\t\t\texit(1);\n")
fmt.Print("\t\tdefault:\n")
fmt.Print("\t\t\tfprintf(stderr, \"Invalid escape sequence unimplemented.\\n\");\n")
fmt.Print("\t\t\texit(1);\n")
fmt.Print("\t\t}\n")
fmt.Print("\t}\n")
fmt.Print("\t\n")
fmt.Print("\tbuffer[buffer_used] = 0;\n")
fmt.Print("\t\n")
fmt.Print("\treturn buffer;\n")
fmt.Print("}\n")
fmt.Print("sexp_t Unnamed() {\n")
fmt.Print("\tstatic int index = 0;\n")
fmt.Print("\t\n")
fmt.Print("\tchar buf[256];\n")
fmt.Print("\tsnprintf(buf, 255, \"__unnamed_%d\", index);\n")
fmt.Print("\tindex++;\n")
fmt.Print("\t\n")
fmt.Print("\treturn SExString(buf); \n")
fmt.Print("}\n")
}
