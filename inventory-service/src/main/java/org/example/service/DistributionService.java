package org.example.service;

import jakarta.transaction.Transactional;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.example.entity.*;
import org.example.grpc.DistributionPlan;
import org.example.grpc.Move;
import org.example.repository.ProductRepository;
import org.example.repository.ShipmentRepository;
import org.example.repository.SupplyRepository;
import org.example.repository.WarehouseRepository;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

@Slf4j
@Service
@RequiredArgsConstructor
public class DistributionService {

    private final ShipmentRepository shipmentRepository;
    private final WarehouseRepository warehouseRepository;
    private final ProductRepository productRepository;

    @Transactional
    public void applyDistributionPlan(DistributionPlan plan) {
        Warehouse sourceWarehouse = warehouseRepository.findById(1L)
                .orElseThrow(() -> new RuntimeException("Source warehouse not found"));

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
            shipment.setLastModifiedBy("system_algo");

            shipment = shipmentRepository.save(shipment);

            ShipmentItem item = new ShipmentItem();
            item.setShipment(shipment);
            item.setProduct(product);
            item.setQuantity(move.getQuantity());
        }
    }
}