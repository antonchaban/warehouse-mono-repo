package org.example.controller;

import org.example.dto.CreateRequest;
import org.example.entity.*;
import org.example.repository.*;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/admin")
@PreAuthorize("hasRole('ADMIN')") // Тільки адмін може створювати
public class AdminController {

    @Autowired private WarehouseRepository warehouseRepository;
    @Autowired private ProductRepository productRepository;

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

    private String getCurrentUsername() {
        Authentication auth = SecurityContextHolder.getContext().getAuthentication();
        return auth != null ? auth.getName() : "system";
    }
}