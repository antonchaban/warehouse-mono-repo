package org.example.controller;

import lombok.RequiredArgsConstructor;
import org.example.dto.CalculateRequest;
import org.example.dto.DistributionEvent;
import org.example.entity.Supply;
import org.example.repository.SupplyRepository;
import org.example.service.EventPublisher;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

import java.util.Map;
import java.util.UUID;

@RestController
@RequestMapping("/api/v1/distribution")
@RequiredArgsConstructor
public class DistributionController {

    private final SupplyRepository supplyRepository;
    private final EventPublisher eventPublisher;

    @PostMapping("/calculate")
    @PreAuthorize("hasRole('LOGISTICIAN')")
    public ResponseEntity<?> triggerCalculation(@RequestBody CalculateRequest request) {

        Supply supply = supplyRepository.findById(request.getSupplyId())
                .orElseThrow(() -> new RuntimeException("Supply not found"));

        String requestId = UUID.randomUUID().toString();

        var eventBuilder = DistributionEvent.builder()
                .requestId(requestId)
                .supplyId(supply.getId())
                .sourceWarehouseId(supply.getWarehouseId());

        eventPublisher.sendCalculationRequest(eventBuilder);

        return ResponseEntity.accepted().body(Map.of(
                "message", "Calculation triggered successfully",
                "request_id", requestId
        ));
    }
}