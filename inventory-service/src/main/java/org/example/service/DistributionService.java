package org.example.service;

import jakarta.transaction.Transactional;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.example.entity.*;
import org.example.grpc.DistributionPlan;
import org.example.grpc.Move;
import org.example.repository.*;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;

@Slf4j
@Service
@RequiredArgsConstructor
public class DistributionService {

    private final ShipmentRepository shipmentRepository;
    private final WarehouseRepository warehouseRepository;
    private final ProductRepository productRepository;
    private final ShipmentItemRepository shipmentItemRepository;
    private final SupplyRepository supplyRepository;

    @Transactional
    public void applyDistributionPlan(DistributionPlan plan) {
        log.info("Processing plan for Request ID: {}", plan.getRequestId());

        long rawSourceId = plan.getSourceId();
        final long sourceId = (rawSourceId == 0) ? 1L : rawSourceId;

        if (rawSourceId == 0) {
            log.warn("Warning: Source ID is 0! Using fallback ID = 1.");
        }

        Warehouse sourceWarehouse = warehouseRepository.findById(sourceId)
                .orElseThrow(() -> new RuntimeException("Source warehouse not found: " + sourceId));

        for (Move move : plan.getMovesList()) {
            Long destId = Long.parseLong(move.getWarehouseId());
            Long prodId = Long.parseLong(move.getProductId());

            Warehouse destWarehouse = warehouseRepository.findById(destId)
                    .orElseThrow(() -> new RuntimeException("Destination warehouse not found: " + destId));

            Product product = productRepository.findById(prodId)
                    .orElseThrow(() -> new RuntimeException("Product not found: " + prodId));

            Shipment shipment = new Shipment();
            shipment.setSource(sourceWarehouse);
            shipment.setDestination(destWarehouse);
            shipment.setStatus("PLANNED");
            shipment.setCreatedBy("system_algo");
            shipment.setCreatedAt(LocalDateTime.now());
            shipment.setLastModifiedBy("system_algo");

            shipment = shipmentRepository.save(shipment);

            ShipmentItem item = new ShipmentItem();
            item.setShipment(shipment);
            item.setProduct(product);
            item.setQuantity(move.getQuantity());

            shipmentItemRepository.save(item);

            log.info("Saved Shipment #{} ({} -> {})", shipment.getId(), sourceId, destId);
        }
        if (plan.getUnallocatedItemsCount() > 0) {
            log.warn("⚠️ ALARM: Some items could not be allocated!");

            for (org.example.grpc.UnallocatedItem unallocated : plan.getUnallocatedItemsList()) {
                log.error("❌ Product ID {} (Volume: {}) failed. Reason: {}",
                        unallocated.getProductId(),
                        unallocated.getVolumeM3(),
                        unallocated.getReason()
                );

                // todo save to DB table for unallocated items
            }
        } else {
            log.info("✅ Perfect! All items were allocated successfully.");
        }
        Supply supply = supplyRepository.findById(plan.getSourceId()).orElse(null);
        if (supply != null) {
            supply.setStatus(SupplyStatus.valueOf("PROCESSED"));
            supplyRepository.save(supply);
            log.info("Supply #{} status updated to PROCESSED", supply.getId());
        }
    }
}