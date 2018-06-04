#include <iostream>

#include <grpcpp/grpcpp.h>

#include "userinfo-server.h"
#include "proto/userinfo/userinfo.grpc.pb.h"

using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerAsyncResponseWriter;
using grpc::ServerContext;
using grpc::Status;

using userinfo::UserInfoRequest;
using userinfo::UserInfoResponse;
using userinfo::UserInfo;

class UserInfoServiceImpl final : public UserInfo::Service {
    Status GetUserInfo(ServerContext* context, const UserInfoRequest* request, UserInfoResponse* response) {
        std::string user_info(request->username());

        response->set_user_info(user_info);

        return Status::OK;
    }
};

void RunServer() {
    std::string server_address("0.0.0.0:7777");
    std::cout << "Starting up the server on " << server_address << std::endl;

    UserInfoServiceImpl service;
    ServerBuilder builder;
    builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
    builder.RegisterService(&service);
    std::unique_ptr<Server> server(builder.BuildAndStart());
    std::cout << "Server now listening on " << server_address << std::endl;

    server->Wait();
}

int main(int argc, char** argv) {
    std::cout << "Welcome to the userinfo C++ server!" << std::endl;
    RunServer();
    return 0;
}