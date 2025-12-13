package org.example.entity;

import jakarta.persistence.*;
import lombok.Data;
import java.time.LocalDateTime;
import java.util.List;

@Entity
@Table(name = "shipments")
@Data
public class Shipment {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @ManyToOne
    @JoinColumn(name = "source_id")
    private Warehouse source;

    @ManyToOne
    @JoinColumn(name = "destination_id")
    private Warehouse destination;

    @Column(nullable = false)
    private String status;

    @OneToMany(mappedBy = "shipment", cascade = CascadeType.ALL, fetch = FetchType.EAGER)
    private List<ShipmentItem> items;

    private String createdBy;
    private LocalDateTime createdAt = LocalDateTime.now();
    private String lastModifiedBy;
    private LocalDateTime lastModifiedAt = LocalDateTime.now();
}