package org.example.controller;

import lombok.RequiredArgsConstructor;
import org.example.config.JwtTokenProvider;
import org.example.dto.AuthRequest;
import org.example.dto.AuthResponse;
import org.example.dto.CreateRequest;
import org.example.entity.Role;
import org.example.entity.User;
import org.example.repository.RoleRepository;
import org.example.repository.UserRepository;
import org.springframework.http.ResponseEntity;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.web.bind.annotation.*;

import java.time.LocalDateTime;
import java.util.HashSet;
import java.util.Set;

@RestController
@RequestMapping("/api/auth")
@RequiredArgsConstructor
public class AuthController {

    private final AuthenticationManager authenticationManager;
    private final JwtTokenProvider tokenProvider;
    private final UserRepository userRepository;
    private final RoleRepository roleRepository;
    private final PasswordEncoder passwordEncoder;

    @PostMapping("/login")
    public ResponseEntity<AuthResponse> login(@RequestBody AuthRequest request) {
        Authentication authentication = authenticationManager.authenticate(
                new UsernamePasswordAuthenticationToken(request.getUsername(), request.getPassword())
        );

        SecurityContextHolder.getContext().setAuthentication(authentication);
        String token = tokenProvider.generateToken(authentication);

        return ResponseEntity.ok(new AuthResponse(token));
    }

    @PostMapping("/register")
    public ResponseEntity<?> register(@RequestBody CreateRequest.UserRegister request) {
        if (userRepository.findByUsername(request.username()).isPresent()) {
            return ResponseEntity.badRequest().body("Username is already taken!");
        }

        User user = new User();
        user.setUsername(request.username());
        // ВИПРАВЛЕНО: поле в класі називається password, хоча колонка password_hash
        user.setPassword(passwordEncoder.encode(request.password()));
        user.setEmail(request.email());
        user.setCreatedAt(LocalDateTime.now());

        // ВИПРАВЛЕНО: для примітива boolean lombok робить setActive
        user.setActive(true);

        // Отримуємо роль
        Role role = roleRepository.findByName("ROLE_LOGISTICIAN")
                .orElseThrow(() -> new RuntimeException("Error: Role is not found."));

        // Ініціалізуємо сет ролей (бо він може бути null) і додаємо роль
        Set<Role> roles = new HashSet<>();
        roles.add(role);
        user.setRoles(roles);

        userRepository.save(user);

        return ResponseEntity.ok("User registered successfully");
    }
}