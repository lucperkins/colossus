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
    private static final int PORT = 1111;

    private Server server;

    static class DataImpl extends DataServiceGrpc.DataServiceImplBase {
        @Override
        public void get(DataRequest req, StreamObserver<DataResponse> resObserver) {
            String request = req.getRequest();
            LOG.info(String.format("Request received for the string: \"%s\"", request));
            String computedValue = request.toUpperCase();
            LOG.info(String.format("Computed value: \"%s\"", computedValue));
            DataResponse res = DataResponse.newBuilder()
                    .setValue(computedValue)
                    .build();
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
        server = ServerBuilder.forPort(PORT)
            .addService(new DataImpl())
            .build()
            .start();
        LOG.info(String.format("Server successfully started on port $d", PORT));
        Runtime.getRuntime().addShutdownHook(new Thread() {
            @Override
            public void run() {
                System.err.println("Shutting down gRPC data server due to JVM shutdown");
                DataHandler.this.stop();
                System.err.println("Server successfully shut down");
            }
        });
    }

    public static void main(String[] args) throws InterruptedException, IOException {
        LOG.info(String.format("Starting up gRPC data server on port $d", PORT));
        final DataHandler server = new DataHandler();
        server.start();
        server.blockUntilShutdown();
    }
}