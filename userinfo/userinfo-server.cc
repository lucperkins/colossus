#include <iostream>

#include "userinfo-server.h"
#include "userinfo.grpc.pb.h"

using grpc::Server;

class Service;

void RunServer() {
    std::string server_address("0.0.0.0:7777");
    std::cout << "Starting up the server on " << server_address << std::endl;
}

int main(int argc, char** argv) {
    std::cout << "Welcome to the userinfo C++ server!" << std::endl;
    RunServer();
    return 0;
}