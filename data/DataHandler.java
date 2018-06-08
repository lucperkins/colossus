package colossus;

import colossus.data.Data;
import colossus.data.DataServiceGrpc;
import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.stub.StreamObserver;
import io.prometheus.client.Counter;
import io.prometheus.client.exporter.HTTPServer;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.logging.Logger;
import java.util.stream.IntStream;

public class DataHandler {
    private static final Logger LOG = Logger.getLogger(DataHandler.class.getName());
    private static final int PORT = 1111;
    private static final Counter syncRequests = Counter.build()
            .name("data_svc_sync_requests")
            .help("Sync requests to the data service")
            .labelNames("request_key")
            .register();

    private static final Counter streamingRequests = Counter.build()
            .name("data_svc_streaming_requests")
            .help("Streaming requests to the data service")
            .register();

    private Server grpcServer;
    private static HTTPServer prometheusHttpServer;

    static class StreamingResponder implements StreamObserver<Data.DataRequest> {
        private StreamObserver<Data.DataResponse> observer;
        private List<String> items = new ArrayList<>();

        StreamingResponder(StreamObserver<Data.DataResponse> observer) {
            this.observer = observer;
        }

        @Override
        public void onNext(Data.DataRequest req) {
            items.add(req.getRequest().replace("f", "9").toUpperCase());
        }

        @Override
        public void onError(Throwable t) {
            observer.onError(t);
        }

        @Override
        public void onCompleted() {
            Data.DataResponse res = Data.DataResponse.newBuilder()
                    .setValue(items.toString())
                    .build();

            this.observer.onNext(res);

            this.observer.onCompleted();
        }
    }

    static class DataImpl extends DataServiceGrpc.DataServiceImplBase {
        private static final Logger LOG = Logger.getLogger(DataImpl.class.getName());

        @Override
        public void get(Data.DataRequest req, StreamObserver<Data.DataResponse> resObserver) {
            String request = req.getRequest();
            LOG.info(String.format("Request received for the string: \"%s\"", request));
            String computedValue = request.toUpperCase();
            LOG.info(String.format("Computed value: \"%s\"", computedValue));
            Data.DataResponse res = Data.DataResponse.newBuilder()
                    .setValue(computedValue)
                    .build();

            syncRequests.labels(request).inc();

            resObserver.onNext(res);
            resObserver.onCompleted();
        }

        @Override
        public void streamingGet(Data.EmptyRequest req, StreamObserver<Data.DataResponse> resObserver) {
            LOG.info("Request received for streaming data");

            Data.DataResponse.Builder resBldr = Data.DataResponse.newBuilder();

            IntStream.range(0, 10).forEach(i -> {
                String value = String.format("Response %d", i);

                streamingRequests.inc();

                resObserver.onNext(resBldr.setValue(value).build());
            });

            resObserver.onCompleted();
        }

        @Override
        public StreamObserver<Data.DataRequest> streamingPut(final StreamObserver<Data.DataResponse> resObserver) {
            return new StreamingResponder(resObserver);
        }

    }

    private void blockUntilShutdown() throws InterruptedException {
        if (grpcServer != null) {
            grpcServer.awaitTermination();
        }
    }

    private void stop() {
        if (grpcServer != null) grpcServer.shutdown();
    }

    private void shutdownHook() {
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            LOG.info("Shutting down gRPC data grpcServer due to JVM shutdown");
            DataHandler.this.stop();
            LOG.info("Server successfully shut down");
        }));
    }

    private void start() throws IOException {
        grpcServer = ServerBuilder.forPort(PORT)
            .addService(new DataImpl())
            .build()
            .start();
        LOG.info(String.format("Server successfully started on port %d", PORT));
        shutdownHook();
    }

    public static void main(String[] args) throws InterruptedException, IOException {
        LOG.info(String.format("Starting up gRPC data grpcServer on port %d", PORT));
        final DataHandler handler = new DataHandler();

        try {
            LOG.info("Starting Prometheus HTTP server");
            prometheusHttpServer = new HTTPServer(9092);
            LOG.info("Successfully started Prometheus HTTP server");
        } catch (IOException e) {
            LOG.severe("Could not start Prometheus HTTP server");
            System.exit(1);
        }

        handler.start();
        handler.blockUntilShutdown();
    }
}