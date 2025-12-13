package org.example.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Builder;
import lombok.Data;

@Data
@Builder
public class DistributionEvent {
    @JsonProperty("request_id")
    private String requestId;

    @JsonProperty("supply_id")
    private Long supplyId;

    @JsonProperty("source_warehouse_id")
    private Long sourceWarehouseId;

    @JsonProperty("initiated_by_user_id")
    private Long initiatedByUserId;

    @JsonProperty("initiated_by_username")
    private String initiatedByUsername;
}