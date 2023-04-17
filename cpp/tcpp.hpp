#ifndef T_H
#define T_H

// __cplusplus gets defined when a C++ compiler processes the file
#ifdef __cplusplus
// extern "C" is needed so the C++ compiler exports the symbols without name
// manging.
extern "C" {
#endif


const char* a(char* a);

#ifdef __cplusplus
}
#endif


#endif