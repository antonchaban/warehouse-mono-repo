package org.example.dto;

public class CreateRequest {
    public record Warehouse(Double capacity) {}
    public record Product(Double volume) {}
    public record UserRegister(String username, String password, String email) {}
}