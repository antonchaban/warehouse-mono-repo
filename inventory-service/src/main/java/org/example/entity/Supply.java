package org.example.entity;

import jakarta.persistence.*;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.time.LocalDateTime;
import java.util.List;

@Data
@EqualsAndHashCode(callSuper = true)
@Entity
@Table(name = "supplies")
public class Supply extends AuditableEntity{
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(name = "warehouse_id", nullable = false)
    private Long warehouseId; // Source Hub

    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private SupplyStatus status;

    @Column(name = "arrival_date")
    private LocalDateTime arrivalDate;

    // One-to-Many relationship for items
    @OneToMany(mappedBy = "supply", cascade = CascadeType.ALL, fetch = FetchType.LAZY)
    private List<SupplyItem> items;
}