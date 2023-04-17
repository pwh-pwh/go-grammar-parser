#ifndef G_H
#define G_H

// __cplusplus gets defined when a C++ compiler processes the file
#ifdef __cplusplus
// extern "C" is needed so the C++ compiler exports the symbols without name
// manging.
extern "C" {
#endif


void b(char* a);

#ifdef __cplusplus
}
#endif


#endif