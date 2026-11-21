package com.brickwise;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.autoconfigure.jdbc.DataSourceAutoConfiguration;

/**
 * Entry point for the Brickwise API.
 *
 * <p><b>C# /.NET analogy:</b> This is like {@code Program.cs} / {@code WebApplication.CreateBuilder}
 * in ASP.NET Core. {@code @SpringBootApplication} turns on component scanning, auto-configuration,
 * and marks this as a configuration class — similar to the minimal hosting model plus
 * {@code [ApiController]} discovery.
 *
 * <p>We exclude {@link DataSourceAutoConfiguration} because {@code DatabaseConfig} builds the
 * DataSource manually (SQLite vs Postgres from env), the same idea as custom connection-string
 * wiring in ASP.NET.
 */
@SpringBootApplication(exclude = DataSourceAutoConfiguration.class)
public class BrickwiseApplication {

    public static void main(String[] args) {
        SpringApplication.run(BrickwiseApplication.class, args);
    }
}
