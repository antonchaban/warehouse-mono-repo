package org.example.controller;

import lombok.RequiredArgsConstructor;
import org.example.entity.StockLevel;
import org.example.repository.StockLevelRepository;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/inventory")
@RequiredArgsConstructor
public class InventoryController {

    private final StockLevelRepository stockLevelRepository;

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
}