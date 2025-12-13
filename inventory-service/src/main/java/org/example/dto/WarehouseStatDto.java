package org.example.dto;

import lombok.Builder;
import lombok.Data;

@Data
@Builder
public class WarehouseStatDto {
    private String name;
    private double totalCapacity;
    private double usedCapacity;
    private double freeCapacity;
    private double utilizationPercentage;
}