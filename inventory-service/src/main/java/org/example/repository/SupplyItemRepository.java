package org.example.repository;

import org.example.entity.SupplyItem;
import org.springframework.data.jpa.repository.JpaRepository;

public interface SupplyItemRepository extends JpaRepository<SupplyItem, Long> {
}
