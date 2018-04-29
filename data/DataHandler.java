package colossus;

import colossus.data.DataProto.DataRequest;
import colossus.data.DataProto.DataResponse;
import colossus.data.DataServiceGrpc;
import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.stub.StreamObserver;

import java.io.IOException;
import java.util.logging.Logger;

public class DataHandler {
    private static final Logger LOG = Logger.getLogger(DataHandler.class.getName());

    private Server server;

    static class DataImpl extends DataServiceGrpc.DataServiceImplBase {
        @Override
        public void get(DataRequest req, StreamObserver<DataResponse> resObserver) {
            String key = req.getKey();
            String value = key.toUpperCase();
            DataResponse res = DataResponse.newBuilder().setValue(value).build();
            resObserver.onNext(res);
            resObserver.onCompleted();
        }
    }

    private void blockUntilShutdown() throws InterruptedException {
        if (server != null) {
            server.awaitTermination();
        }
    }

    private void stop() {
        if (server != null) server.shutdown();
    }

    private void start() throws IOException {
        int port = 1111;
        server = ServerBuilder.forPort(port)
            .addService(new DataImpl())
            .build()
            .start();
        LOG.info("Server successfully started on port 1111");
        Runtime.getRuntime().addShutdownHook(new Thread() {
            @Override
            public void run() {
                System.err.println("Shutting down gRPC data server due to JVM shut down");
                DataHandler.this.stop();
                System.err.println("Server successfully shut down");
            }
        });
    }

    public static void main(String[] args) throws InterruptedException, IOException {
        LOG.info("Starting up gRPC data server on port 1111");
        final DataHandler server = new DataHandler();
        server.start();
        server.blockUntilShutdown();
    }
}