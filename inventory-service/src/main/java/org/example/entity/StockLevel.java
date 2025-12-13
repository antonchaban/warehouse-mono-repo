package org.example.entity;

import jakarta.persistence.*;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.io.Serializable;

@Data
@Entity
@EqualsAndHashCode(callSuper = true)
@Table(name = "stock_levels")
public class StockLevel extends AuditableEntity {

    @EmbeddedId
    private StockLevelId id;

    @Column(nullable = false)
    private Integer quantity;

    // --- Composite Key Class ---
    @Data
    @Embeddable
    public static class StockLevelId implements Serializable {
        @Column(name = "warehouse_id")
        private Long warehouseId;

        @Column(name = "product_id")
        private Long productId;
    }
}