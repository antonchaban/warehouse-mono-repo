package org.example.repository;

import org.example.entity.SupplyItem;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

public interface SupplyItemRepository extends JpaRepository<SupplyItem, Long> {
    @Query(value = """
        SELECT COALESCE(SUM(si.quantity * p.volume_m3), 0)
        FROM supply_items si
        JOIN supplies s ON si.supply_id = s.id
        JOIN products p ON si.product_id = p.id
        WHERE s.warehouse_id = :warehouseId 
          AND s.status = 'RECEIVED'
    """, nativeQuery = true)
    Double calculatePendingVolumeByWarehouse(@Param("warehouseId") Long warehouseId);
}
