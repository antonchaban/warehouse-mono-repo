package org.example.dto;

import java.time.LocalDateTime;
import java.util.List;

public record ShipmentResponse(
        Long id,
        Long sourceId,
        Long destinationId,
        String status,
        LocalDateTime createdAt,
        List<ItemResponse> items
) {}