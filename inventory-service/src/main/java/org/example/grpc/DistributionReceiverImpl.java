package org.example.grpc;

import io.grpc.stub.StreamObserver;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import net.devh.boot.grpc.server.service.GrpcService;
import org.example.service.DistributionService;

@Slf4j
@GrpcService
@RequiredArgsConstructor
public class DistributionReceiverImpl extends DistributionResultReceiverGrpc.DistributionResultReceiverImplBase {

    private final DistributionService distributionService;

    @Override
    public void processPlan(DistributionPlan request, StreamObserver<Empty> responseObserver) {
        try {
            log.info("Received plan from Go. RequestID: {}", request.getRequestId());
            long supplyId = request.getSupplyId();
            Long finalSupplyId = (supplyId == 0) ? null : supplyId;

            distributionService.applyDistributionPlan(request, finalSupplyId);

            responseObserver.onNext(Empty.newBuilder().build());
            responseObserver.onCompleted();
        } catch (Exception e) {
            log.error("Failed to process distribution plan", e);
            responseObserver.onError(e);
        }
    }
}