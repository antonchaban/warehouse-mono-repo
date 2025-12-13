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
    @Autowired private SupplyRepository supplyRepository;
    @Autowired private SupplyItemRepository supplyItemRepository;

    @PostMapping("/supplies")
    public ResponseEntity<?> createSupply(@RequestBody CreateRequest.Supply request) {
        // 1. Створюємо "шапку" поставки
        Supply supply = new Supply();
        supply.setWarehouseId(request.warehouseId());
        // supply.setStatus(SupplyStatus.RECEIVED); // Краще використовувати Enum, якщо він є
        supply.setStatus(SupplyStatus.valueOf("RECEIVED")); // Або String, як у вас було в попередніх версіях
        supply.setArrivalDate(java.time.LocalDateTime.now());
        supply.setCreatedBy(getCurrentUsername());

        supply = supplyRepository.save(supply); // Отримуємо ID (наприклад, 557)

        // 2. Додаємо товари в цю поставку
        SupplyItem item = new SupplyItem();

        // --- ВИПРАВЛЕННЯ ТУТ ---
        // item.setId(supply.getId()); <--- БУЛО (Неправильно: це перезаписує ID самого товару)

        item.setSupply(supply); // <--- СТАЛО (Правильно: ми пов'язуємо товар з поставкою)
        // -----------------------

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

    private String getCurrentUsername() {
        Authentication auth = SecurityContextHolder.getContext().getAuthentication();
        return auth != null ? auth.getName() : "system";
    }
}