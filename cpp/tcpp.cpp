#include "tcpp.hpp"
#include <iostream>
#include <string>
using namespace std;


const char* a(char* a) {
    string s = a;
    std::cout <<a<<std::endl;
    return s.c_str();
}
