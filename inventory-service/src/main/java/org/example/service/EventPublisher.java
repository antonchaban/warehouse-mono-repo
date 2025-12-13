package org.example.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.example.config.CustomUserDetails;
import org.example.config.RabbitMQConfig;
import org.example.dto.DistributionEvent;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Service;

@Slf4j
@Service
@RequiredArgsConstructor
public class EventPublisher {

    private final RabbitTemplate rabbitTemplate;

    public void sendCalculationRequest(DistributionEvent.DistributionEventBuilder eventBuilder) {
        var authentication = SecurityContextHolder.getContext().getAuthentication();
        var principal = (CustomUserDetails) authentication.getPrincipal();

        DistributionEvent event = eventBuilder
                .initiatedByUserId(principal.getId())
                .initiatedByUsername(principal.getUsername())
                .build();

        log.info("Publishing event to RabbitMQ by user {}: {}", principal.getUsername(), event);

        rabbitTemplate.convertAndSend(
                RabbitMQConfig.EXCHANGE_NAME,
                RabbitMQConfig.ROUTING_KEY,
                event
        );
    }
}