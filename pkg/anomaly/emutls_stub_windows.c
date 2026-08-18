#ifdef _WIN32
__asm__(
".globl __emutls_v._ZSt11__once_call\n"
".globl __emutls_v._ZSt15__once_callable\n"
".data\n"
"__emutls_v._ZSt11__once_call:\n"
"    .quad 0, 0, 0, 0\n"
"__emutls_v._ZSt15__once_callable:\n"
"    .quad 0, 0, 0, 0\n"
);
#endif
