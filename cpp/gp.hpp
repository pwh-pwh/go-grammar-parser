#ifndef G_H
#define G_H

// __cplusplus gets defined when a C++ compiler processes the file
#ifdef __cplusplus
// extern "C" is needed so the C++ compiler exports the symbols without name
// manging.
extern "C" {
#endif


const char* b(char* a);
const char * getRR(char *st);
const char * getRRAL(char *st);
const char * getFirstF(char * st);
const char * getFollowF(char * st);
const char * getTableF(char * st);
const char * getTreeF(char * st,char * gr);




#ifdef __cplusplus
}
#endif


#endif