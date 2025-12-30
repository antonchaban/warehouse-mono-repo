package org.example.controller;

import lombok.RequiredArgsConstructor;
import org.example.entity.Shipment;
import org.example.entity.ShipmentStatus;
import org.example.entity.StockLevel;
import org.example.repository.ShipmentRepository;
import org.example.repository.StockLevelRepository;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

import java.time.LocalDateTime;
import java.util.List;

@RestController
@RequestMapping("/api/inventory")
@RequiredArgsConstructor
public class InventoryController {

    private final StockLevelRepository stockLevelRepository;
    private final ShipmentRepository shipmentRepository;

    @GetMapping("/stocks")
    @PreAuthorize("isAuthenticated()")
    public ResponseEntity<List<StockLevel>> getCurrentStocks() {
        return ResponseEntity.ok(stockLevelRepository.findAll());
    }

    @PostMapping("/inbound")
    @PreAuthorize("hasRole('STOREKEEPER')")
    public ResponseEntity<String> registerInbound() {
        // TODO: Implement inbound registration logic
        return ResponseEntity.ok("Inbound registered successfully");
    }

    @PutMapping("/shipments/{id}/status")
    @PreAuthorize("hasAnyRole('STOREKEEPER', 'ADMIN', 'LOGISTICIAN')")
    public ResponseEntity<?> updateShipmentStatus(
            @PathVariable Long id,
            @RequestParam String status
    ) {
        Shipment shipment = shipmentRepository.findById(id)
                .orElseThrow(() -> new RuntimeException("Shipment not found"));

        try {
            ShipmentStatus newStatus = ShipmentStatus.valueOf(status);

            shipment.setStatus(newStatus.name());
            shipment.setLastModifiedAt(LocalDateTime.now());

            // if (newStatus == ShipmentStatus.DELIVERED) { ... upd StockLevel ... }

            shipmentRepository.save(shipment);

            return ResponseEntity.ok("Shipment status updated to " + status);
        } catch (IllegalArgumentException e) {
            return ResponseEntity.badRequest().body("Invalid status provided");
        }
    }
}