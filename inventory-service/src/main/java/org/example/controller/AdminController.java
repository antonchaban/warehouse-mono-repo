package org.example.controller;

import org.example.dto.CreateRequest;
import org.example.dto.WarehouseStatDto;
import org.example.entity.*;
import org.example.repository.*;

import java.util.ArrayList;
import java.util.List;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/admin")
@PreAuthorize("hasRole('ADMIN')")
public class AdminController {

    @Autowired private WarehouseRepository warehouseRepository;
    @Autowired private ProductRepository productRepository;
    @Autowired private SupplyRepository supplyRepository;
    @Autowired private SupplyItemRepository supplyItemRepository;
    @Autowired private ShipmentItemRepository shipmentItemRepository;

    @GetMapping("/warehouses/stats")
    public ResponseEntity<List<WarehouseStatDto>> getWarehouseStats() {
        List<Warehouse> warehouses = warehouseRepository.findAll();
        List<WarehouseStatDto> stats = new ArrayList<>();

        for (Warehouse w : warehouses) {
            double allocatedVolume = shipmentItemRepository.calculateUsedVolumeByWarehouse(w.getId());

            double pendingVolume = supplyItemRepository.calculatePendingVolumeByWarehouse(w.getId());

            double totalUsed = allocatedVolume + pendingVolume;

            double total = w.getTotalCapacity();
            double free = total - totalUsed;
            if (free < 0) free = 0;

            double percentage = (total > 0) ? (totalUsed / total) * 100 : 0;

            stats.add(WarehouseStatDto.builder()
                    .name("WH-" + w.getId())
                    .totalCapacity(total)
                    .usedCapacity(totalUsed)
                    .freeCapacity(free)
                    .utilizationPercentage(percentage)
                    .build());
        }

        return ResponseEntity.ok(stats);
    }

    @PostMapping("/supplies")
    public ResponseEntity<?> createSupply(@RequestBody CreateRequest.Supply request) {
        Supply supply = new Supply();
        supply.setWarehouseId(request.warehouseId());
        supply.setStatus(SupplyStatus.valueOf("RECEIVED"));
        supply.setArrivalDate(java.time.LocalDateTime.now());
        supply.setCreatedBy(getCurrentUsername());

        supply = supplyRepository.save(supply);

        SupplyItem item = new SupplyItem();
        item.setSupply(supply);

        item.setProductId(request.productId());
        item.setQuantity(request.quantity());

        supplyItemRepository.save(item);

        return ResponseEntity.ok("Supply created with ID: " + supply.getId());
    }

    @PostMapping("/warehouses")
    public ResponseEntity<?> createWarehouse(@RequestBody CreateRequest.Warehouse request) {
        Warehouse w = new Warehouse();
        w.setTotalCapacity(request.capacity());
        w.setCreatedBy(getCurrentUsername());
        warehouseRepository.save(w);
        return ResponseEntity.ok("Warehouse created with ID: " + w.getId());
    }

    @PostMapping("/products")
    public ResponseEntity<?> createProduct(@RequestBody CreateRequest.Product request) {
        Product p = new Product();
        p.setVolumeM3(request.volume());
        p.setCreatedBy(getCurrentUsername());
        productRepository.save(p);
        return ResponseEntity.ok("Product created with ID: " + p.getId());
    }

    @GetMapping("/warehouses")
    public ResponseEntity<List<Warehouse>> getAllWarehouses() {
        return ResponseEntity.ok(warehouseRepository.findAll());
    }

    @GetMapping("/products")
    public ResponseEntity<List<Product>> getAllProducts() {
        return ResponseEntity.ok(productRepository.findAll());
    }

    @GetMapping("/supplies")
    public ResponseEntity<List<Supply>> getAllSupplies() {
        List<Supply> supplies = supplyRepository.findAll();
        supplies.sort((a, b) -> b.getId().compareTo(a.getId()));
        return ResponseEntity.ok(supplies);
    }

    private String getCurrentUsername() {
        Authentication auth = SecurityContextHolder.getContext().getAuthentication();
        return auth != null ? auth.getName() : "system";
    }
}