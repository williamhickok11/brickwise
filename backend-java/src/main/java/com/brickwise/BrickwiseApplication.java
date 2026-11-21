package com.brickwise;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

// Host process. Same role as backend/cmd/server/main.go and an ASP.NET Program.cs.
@SpringBootApplication
public class BrickwiseApplication {

    public static void main(String[] args) {
        SpringApplication.run(BrickwiseApplication.class, args);
    }
}
