package org.example.controller;

import lombok.RequiredArgsConstructor;
import org.example.entity.Warehouse;
import org.example.repository.WarehouseRepository;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/admin")
@RequiredArgsConstructor
public class AdminController {

    private final WarehouseRepository warehouseRepository;

    @PostMapping("/warehouses")
    @PreAuthorize("hasRole('ADMIN')")
    public ResponseEntity<Warehouse> createWarehouse(@RequestBody Warehouse warehouse) {
        return ResponseEntity.ok(warehouseRepository.save(warehouse));
    }

    @GetMapping("/warehouses")
    @PreAuthorize("hasRole('ADMIN')")
    public ResponseEntity<List<Warehouse>> getAllWarehouses() {
        return ResponseEntity.ok(warehouseRepository.findAll());
    }
}