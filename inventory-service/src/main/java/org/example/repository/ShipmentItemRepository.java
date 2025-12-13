package org.example.repository;

import org.example.entity.ShipmentItem;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

public interface ShipmentItemRepository extends JpaRepository<ShipmentItem,Long> {
    @Query("SELECT COALESCE(SUM(si.quantity * p.volumeM3), 0) " +
            "FROM ShipmentItem si " +
            "JOIN si.shipment s " +
            "JOIN si.product p " +
            "WHERE s.destination.id = :warehouseId")
    Double calculateUsedVolumeByWarehouse(@Param("warehouseId") Long warehouseId);
}
