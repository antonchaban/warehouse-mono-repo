package org.example.config;

import io.jsonwebtoken.*;
import io.jsonwebtoken.security.Keys;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Component;

import java.security.Key;
import java.util.Date;

@Component
public class JwtTokenProvider {

    // Секретний ключ (в реальності краще тримати в application.yaml)
    // Має бути довгим (мінімум 32 символи) для безпеки
    private static final String JWT_SECRET = "diplomaSecretKeyForGoodsDistributionSystem2025";

    // Час життя токена (наприклад, 1 день в мілісекундах)
    private static final long JWT_EXPIRATION = 86400000L;

    private final Key key = Keys.hmacShaKeyFor(JWT_SECRET.getBytes());

    // Генерація токена
    public String generateToken(Authentication authentication) {
        String username = authentication.getName();
        Date now = new Date();
        Date expiryDate = new Date(now.getTime() + JWT_EXPIRATION);

        return Jwts.builder()
                .setSubject(username)
                .setIssuedAt(new Date())
                .setExpiration(expiryDate)
                .signWith(key, SignatureAlgorithm.HS256)
                .compact();
    }

    // Отримання username з токена
    public String getUsernameFromJWT(String token) {
        Claims claims = Jwts.parserBuilder()
                .setSigningKey(key)
                .build()
                .parseClaimsJws(token)
                .getBody();
        return claims.getSubject();
    }

    // Валідація токена
    public boolean validateToken(String authToken) {
        try {
            Jwts.parserBuilder().setSigningKey(key).build().parseClaimsJws(authToken);
            return true;
        } catch (JwtException | IllegalArgumentException ex) {
            // Тут можна логувати помилки (Expired, Malformed тощо)
            System.err.println("Invalid JWT token: " + ex.getMessage());
        }
        return false;
    }
}