package org.example.repository;

import org.example.entity.StockLevel;
import org.springframework.data.jpa.repository.JpaRepository;

public interface StockLevelRepository extends JpaRepository<StockLevel, Long> {
}
